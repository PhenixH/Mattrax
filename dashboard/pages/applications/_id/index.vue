<template>
  <div ref="page" class="page-body">
    <p class="field-title">Name:</p>
    <input
      name="name"
      :value="app.name"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import resource from '@/mixins/resource'

export default Vue.extend({
  mixins: [resource],
  props: {
    app: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async save(patch: object) {
      // await this.$store.dispatch('policies/patchPolicy', {
      //   id: this.$route.params.id,
      //   patch,
      // })

      Object.keys(patch).forEach((key) => (this.app[key] = patch[key]))
    },
  },
})
</script>

<style scoped></style>
