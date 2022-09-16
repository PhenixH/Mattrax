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

    <p class="field-title">Ownership:</p>
    <input :value="titleCaseStr(device.ownership)" type="text" disabled />

    <h2>Management Details</h2>

    <p class="field-title">Protocol:</p>
    <input :value="titleCaseStr(device.protocol)" type="text" disabled />

    <p class="field-title">Scope:</p>
    <input :value="device.scope" type="text" disabled />

    <p class="field-title">State:</p>
    <input :value="titleCaseStr(device.state)" type="text" disabled />

    <div v-if="device.azure_did !== null">
      <p class="field-title">Azure Device ID:</p>
      <input :value="device.azure_did" type="text" disabled />
    </div>

    <p class="field-title">Last Seen At:</p>
    <input :value="device.lastseen" type="text" disabled />

    <p class="field-title">Enrolled At:</p>
    <input :value="device.enrolled_at" type="text" disabled />

    <h2>Operating System</h2>

    <p class="field-title">Version Major:</p>
    <input :value="device.os_major" type="text" disabled />

    <p class="field-title">Version Minor:</p>
    <input :value="device.os_minor" type="text" disabled />

    <h2>Hardware</h2>

    <p class="field-title">Serial Number:</p>
    <input :value="device.serial_number" type="text" disabled />

    <p class="field-title">Manufacturer:</p>
    <input
      :value="titleCaseStr(device.model_manufacturer)"
      type="text"
      disabled
    />

    <p class="field-title">Model:</p>
    <input :value="device.model" type="text" disabled />
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
