<template>
  <div ref="page" class="page-body">
    <h1>Overview</h1>
    <p class="field-title">Name:</p>
    <input
      name="fullname"
      :value="user.fullname"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">User Principal Name:</p>
    <input
      name="upn"
      :value="user.upn"
      type="email"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Status:</p>
    <select
      name="disabled"
      data-type="select"
      :value="user.disabled"
      :disabled="!$store.state.dashboard.editting"
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
import resource from '@/mixins/resource'

export default Vue.extend({
  mixins: [resource],
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async save(patch: object) {
      const upn = await this.$store.dispatch('users/patchUser', {
        upn: this.$route.params.id,
        patch,
      })

      if (upn !== this.$route.params.id) {
        this.$router.push('/users/' + upn)
      } else {
        Object.keys(patch).forEach((key) => (this.user[key] = patch[key]))
      }
    },
  },
})
</script>

<style scoped></style>
