<template>
  <div class="page-body">
    <h2>Tenant Details</h2>
    <p class="field-title">Tenant Name:</p>
    <input
      name="display_name"
      :value="tenant_settings.display_name"
      :test="tenant_settings.display_name"
      type="text"
      :disabled="!editting"
    />

    <p class="field-title">Tenant Email:</p>
    <input
      name="email"
      :value="tenant_settings.email"
      :test="tenant_settings.email"
      type="email"
      :disabled="!editting"
    />

    <p class="field-title">Tenant Phone:</p>
    <input
      name="phone"
      :value="tenant_settings.phone"
      :test="tenant_settings.phone"
      type="tel"
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
      tenant_settings: {},
    }
  },
  created() {
    this.$store
      .dispatch('settings/getForCurrentTenant')
      .then((settings) => {
        this.tenant_settings = settings
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    async save(patch: object) {
      await this.$store.dispatch('settings/updateForCurrentTenant', {
        id: this.$route.params.id,
        patch,
      })

      Object.keys(patch).forEach(
        (key) => (this.tenant_settings[key] = patch[key])
      )

      if (
        this.tenant_settings.display_name !==
        this.$store.state.tenants.tenant.display_name
      ) {
        const tenant = Object.assign({}, this.$store.state.tenants.tenant, {
          display_name: this.tenant_settings.display_name,
        })
        this.$store.commit('tenants/set', tenant)
        this.$store.commit('tenants/clearTenants')
      }
    },
  },
})
</script>

<style scoped></style>
