<template>
  <div class="page-body">
    <h2>Tenant Details</h2>
    <p class="field-title">Full Name:</p>
    <input
      name="fullname"
      :value="user_settings.fullname"
      :test="user_settings.fullname"
      type="text"
      :disabled="!editting"
    />

    <p class="field-title">Azure AD OID:</p>
    <input
      :value="user_settings.azuread_oid"
      :test="user_settings.azuread_oid"
      type="azuread_oid"
      :disabled="!editting"
    />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    editting: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      user_settings: {},
    }
  },
  created() {
    this.$store
      .dispatch('settings/getUser')
      .then((settings) => {
        this.user_settings = settings
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    async save(patch: object) {
      await this.$store.dispatch('settings/updateUser', {
        id: this.$route.params.id,
        patch,
      })

      Object.keys(patch).forEach(
        (key) => (this.user_settings[key] = patch[key])
      )

      // TODO: Update auth token
    },
  },
})
</script>

<style scoped></style>
