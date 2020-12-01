<template>
  <div class="page-body">
    <h2>Tenant Details</h2>
    <p class="field-title">Tenant Name:</p>
    <input
      name="display_name"
      :value="tenant_settings.display_name"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Tenant Email:</p>
    <input
      name="email"
      :value="tenant_settings.email"
      type="email"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Tenant Phone:</p>
    <input
      name="phone"
      :value="tenant_settings.phone"
      type="tel"
      :disabled="!$store.state.dashboard.editting"
    />

    <h2>Domain Management</h2>
    <p>
      To link a domain create a TXT DNS record at the domain with the linking
      code, then click verify.
    </p>
    <TableView :headings="['Domain', 'Verified', 'Linking Code', 'Actions']">
      <tr v-for="domain in domains" :key="domain.domain">
        <td>{{ domain.domain }}</td>
        <td :class="{ danger: !domain.verified }">
          {{ domain.verified === true ? 'Verified' : 'Awaiting Verification' }}
        </td>
        <td>mttx{{ domain.linking_code }}</td>
        <td>
          <button @click="verifyDomain(domain)">Verify</button>
          <button @click="deleteDomain(domain)">Delete</button>
        </td>
      </tr>
    </TableView>
    <input
      v-model="new_domain"
      placeholder="example.com"
      pattern="^[a-zA-Z][a-zA-Z\d-]{1,22}[a-zA-Z\d]$"
    />
    <button @click="addDomain()">Add Domain</button>

    <p>Tenant ID: {{ $store.state.tenants.tenant.id }}</p>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import resource from '@/mixins/resource'

export default Vue.extend({
  mixins: [resource],
  data() {
    return {
      tenant_settings: {},
      domains: [],
      new_domain: '',
    }
  },
  created() {
    this.$store
      .dispatch('settings/getForCurrentTenant')
      .then((settings) => {
        this.tenant_settings = settings.tenant
        this.domains = settings.domains
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
    async verifyDomain(domain: any) {
      const verified = await this.$store.dispatch(
        'tenants/verifyDomain',
        domain.domain
      )
      this.domains.find((d) => d.domain === domain.domain).verified = verified
    },
    async deleteDomain(domain: any) {
      await this.$store.dispatch('tenants/deleteDomain', domain.domain)
      this.domains = this.domains.filter((d) => d.domain !== domain.domain)
    },
    async addDomain() {
      if (this.new_domain === '') {
        return
      }

      const domain = await this.$store.dispatch(
        'tenants/addDomain',
        this.new_domain
      )
      this.domains.push(domain)
      this.new_domain = ''
    },
  },
})
</script>

<style scoped></style>
