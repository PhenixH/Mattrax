<template>
  <div ref="page" class="page-body">
    <p class="field-title">Name:</p>
    <input name="name" :value="policy.name" type="text" :disabled="!editting" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    policy: {
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
