<template>
  <div ref="page" class="page-body">
    <p class="field-title">Name:</p>
    <input
      name="name"
      :value="app.name"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Publisher:</p>
    <input
      name="publisher"
      :value="app.publisher"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Description:</p>
    <input
      name="description"
      :value="app.description"
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
      await this.$store.dispatch('applications/patch', {
        id: this.$route.params.id,
        patch,
      })

      Object.keys(patch).forEach((key) => (this.app[key] = patch[key]))
    },
  },
})
</script>

<style scoped></style>
