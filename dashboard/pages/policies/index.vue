<template>
  <div v-if="loading" class="loading">Loading Policies...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
      </ul>
      <h1>Policies</h1>
    </PageHead>
    <PageNav>
      <button @click="$router.push('/policies/new')">Create New Policy</button>
      <input type="text" placeholder="Search..." disabled />
    </PageNav>
    <div class="page-body">
      <TableView :headings="['Name', 'Type']">
        <tr v-for="policy in policies" :key="policy.id">
          <td>
            <NuxtLink :to="'/policies/' + policy.id" exact>{{
              policy.name
            }}</NuxtLink>
          </td>
          <td>
            {{ payloads_json[policy.type].display_name }}
          </td>
        </tr>
      </TableView>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import policiesJson from '@/policies.json'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      policies: [],
      payloads_json: policiesJson.payloads,
    }
  },
  created() {
    this.$store
      .dispatch('policies/getAll')
      .then((policies) => {
        this.policies = policies
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style scoped></style>
