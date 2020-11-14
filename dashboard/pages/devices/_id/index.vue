<template>
  <div ref="page" class="page-body">
    <p class="field-title">Name:</p>
    <input name="name" :value="device.name" type="text" :disabled="!editting" />

    <p class="field-title">Owner:</p>
    <input :value="device.owner" type="text" disabled />

    <p class="field-title">Protocol:</p>
    <input :value="titleCase(device.protocol)" type="text" disabled />

    <p class="field-title">State:</p>
    <input :value="titleCase(device.state)" type="text" disabled />

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
import titleCase from '@/mixins/titleCase'

export default Vue.extend({
  mixins: [titleCase],
  props: {
    device: {
      type: Object,
      required: true,
    },
    editting: {
      type: Boolean,
      required: true,
    },
  },
  mounted() {
    // TODO: Remove need for ref with this or something
    this.$refs.page
      .querySelectorAll('input, select, checkbox, textarea')
      .forEach((node: HTMLInputElement) => {
        node.defaultValue = node.value
      })
  },
})
</script>

<style scoped></style>
