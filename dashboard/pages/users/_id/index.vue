<template>
  <div ref="page" class="page-body">
    <h1>Overview</h1>
    <p class="field-title">Name:</p>
    <input
      name="fullname"
      :value="user.fullname"
      type="text"
      :disabled="!editting"
    />

    <p class="field-title">User Principal Name:</p>
    <input name="upn" :value="user.upn" type="email" :disabled="!editting" />

    <p class="field-title">Status:</p>
    <select
      name="disabled"
      data-type="bool"
      :value="user.disabled"
      :disabled="!editting"
    >
      <option value="false">Enabled</option>
      <option value="true">Disabled</option>
    </select>

    <p class="field-title">AzureAD Object ID:</p>
    <input v-model="user.azuread_oid" type="text" disabled />

    <p class="field-title">Permission:</p>
    <!-- TODO: Connect to data -->
    <select disabled>
      <option value="administrator">Administrator</option>
      <option value="user" selected>User</option>
    </select>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    user: {
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
    (<Element>this.$refs.page)
      .querySelectorAll('input, select, checkbox, textarea')
      .forEach((node: any) => {
        node.defaultValue = node.value
      })
  },
})
</script>

<style scoped></style>
