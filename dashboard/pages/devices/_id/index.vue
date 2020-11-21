<template>
  <div class="page-body">
    <p class="field-title">Name:</p>
    <input
      name="name"
      :value="device.name"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Owner:</p>
    <input :value="device.owner" type="text" disabled />

    <p class="field-title">Protocol:</p>
    <input :value="titleCaseStr(device.protocol)" type="text" disabled />

    <p class="field-title">State:</p>
    <input :value="titleCaseStr(device.state)" type="text" disabled />

    <p class="field-title">Azure Device ID:</p>
    <input :value="device.azure_did" type="text" disabled />

    <p class="field-title">Model:</p>
    <input :value="device.model" type="text" disabled />

    <p class="field-title">Enrolled At:</p>
    <input :value="device.enrolled_at" type="text" disabled />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import resource from '@/mixins/resource'

export default Vue.extend({
  mixins: [resource],
  props: {
    device: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async save(patch: object) {
      await this.$store.dispatch('devices/patchDevice', {
        id: this.$route.params.id,
        patch,
      })

      Object.keys(patch).forEach((key) => (this.device[key] = patch[key]))
    },
  },
})
</script>

<style scoped></style>
