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

    <h2>MDM Protocols</h2>
    <button @click="setupAndroidOrganisation()">
      Configure Android for Work Organisation
    </button>

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

    <!-- <h2>Audit Log</h2>
    <textarea v-model="this.event_log" /> -->

    <h2>Administrators</h2>

    <h2>Identifiers</h2>
    <p class="field-title">Android For Work Organisation ID:</p>
    <input :value="tenant_settings.afw_enterprise_id" disabled />
    <p class="field-title">Mattrax Tenant ID:</p>
    <input :value="$store.state.tenants.tenant.id" disabled />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
// import { fetchEventSource } from '@microsoft/fetch-event-source'
import resource from '@/mixins/resource'

export default Vue.extend({
  mixins: [resource],
  data() {
    return {
      tenant_settings: {},
      domains: [],
      new_domain: '',
      event_log: '',
    }
  },
  created() {
    // await fetchEventSource(
    //   process.env.baseUrl +
    //     '/' +
    //     this.$store.state.tenants.tenant.id +
    //     '/settings/events',
    //   {
    //     headers: {
    //         'Authorization': 'Bearer ' + this.$store.state.authentication.authToken,
    //     },
    //     onmessage: this.onmessage
    //   }
    // )

    this.$store
      .dispatch('settings/getForCurrentTenant')
      .then((settings) => {
        this.tenant_settings = settings.tenant
        this.domains = settings.domains
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    onmessage(ev: any) {
      this.event_log += ev.data + '\n'
    },
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
    async setupAndroidOrganisation() {
      const signupURL = await this.$store.dispatch('android/newSignupURL')
      window.open(signupURL.url, '_blank')
    },
  },
})
</script>

<style scoped></style>
