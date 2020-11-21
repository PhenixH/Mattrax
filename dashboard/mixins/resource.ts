function generateFormPatch(el: any) {
  let patch: any = null
  el.querySelectorAll('input, select, checkbox, textarea').forEach(
    (node: HTMLInputElement) => {
      if (node.value !== node.defaultValue) {
        if (patch === null) patch = {}
        if (node.getAttribute('data-type') === 'bool') {
          patch[node.name] = node.value === 'true'
        } else {
          patch[node.name] = node.value
        }
      }
    }
  )
  return patch
}

export default {
  created() {
    if (this.save === undefined) {
      console.error(
        "Error mounting resource mixin without 'save' function defined!"
      )
      return
    }

    this.$store.commit('dashboard/setEditting', false)
  },
  destroyed() {
    this.$store.commit('dashboard/setEditting', null)
    this.$store.commit('dashboard/setDeletable', false)
  },
  updated() {
    this.storeDefaultValues()
  },
  methods: {
    storeDefaultValues() {
      this.$el
        .querySelectorAll('input, select, checkbox, textarea')
        .forEach((node: any) => (node.defaultValue = node.value))
    },
    savebtn() {
      const patch = generateFormPatch(this.$el)
      if (patch === null) {
        this.$store.commit('dashboard/setEditting', false)
        return
      }

      this.save(patch)
        .then(() => {
          this.storeDefaultValues()
          this.$store.commit('dashboard/setEditting', false)
        })
        .catch((err: Error) => this.$store.commit('dashboard/setError', err)) // TODO: Warning that saving failed
    },
    deletebtn() {
      ;(this.delete !== undefined ? this.delete : this.$parent.delete)()
        .then((dest: string) => {
          this.$store.commit('dashboard/setEditting', false)
          this.$router.push(dest)
        })
        .catch((err: Error) => this.$store.commit('dashboard/setError', err)) // TODO: Warning that saving failed
    },
    titleCaseStr(str: string): string {
      if (str === '') return ''
      return str.charAt(0).toUpperCase() + str.slice(1)
    },
  },
  beforeRouteLeave(_to: any, _from: any, next: any) {
    if (this.$store.state.dashboard.editting === false) {
      next()
      return
    } else if (generateFormPatch(this.$el) === null) {
      next()
      return
    }
    const answer = window.confirm(
      'Do you really want to leave? you have unsaved changes!'
    )
    if (answer) {
      next()
    } else {
      next(false)
    }
  },
}
