<template>
  <div v-if="loading" class="loading">Loading Device...</div>
  <div v-else>
    <div class="page-head">
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/devices">Devices</NuxtLink></li>
      </ul>
      <h1>{{ device.name }}</h1>
    </div>
    <div class="page-nav">
      <button
        :class="{
          active:
            this.$route.path.replace('/devices/' + this.$route.params.id, '') ==
            '',
        }"
        @click="navigate('')"
      >
        Overview
      </button>
      <button
        :class="{
          active:
            this.$route.path.replace('/devices/' + this.$route.params.id, '') ==
            '/scope',
        }"
        @click="navigate('/scope')"
      >
        Scope
      </button>
    </div>
    <NuxtChild :device="device" :editting="editting" />
    <div class="page-footer">
      <button @click="editting ? save() : (editting = true)">
        {{ editting ? 'Save' : 'Edit' }}
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      editting: false,
      device: {},
    }
  },
  created() {
    this.$store
      .dispatch('devices/getByID', this.$route.params.id)
      .then((device) => {
        this.device = device
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/devices/' + this.$route.params.id + pathSuffix)
    },
    save() {
      if (this.editting === false) return

      let patch = null
      this.$children[2].$el
        .querySelectorAll('input, select, checkbox, textarea')
        .forEach((node) => {
          if (node.value !== node.defaultValue) {
            if (patch === null) patch = {}
            if (node.getAttribute('data-type') === 'bool') {
              patch[node.name] = node.value === 'true'
            } else {
              patch[node.name] = node.value
            }
          }
        })
      if (patch === null) {
        this.editting = false
        return
      }

      this.$store
        .dispatch('devices/patchDevice', {
          id: this.$route.params.id,
          patch,
        })
        .then(() => {
          Object.keys(patch).forEach((key) => (this.device[key] = patch[key]))

          this.editting = false
        })
        .catch((err) => this.$store.commit('dashboard/setError', err)) // TODO: Warning that saving failed
    },
  },
})
</script>

<style scoped></style>
