<template>
  <div v-if="loading" class="loading">Loading Policies...</div>
  <div v-else>
    <div class="page-head">
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
      </ul>
      <h1>Policies</h1>
    </div>
    <div class="page-nav">
      <button @click="$router.push('/policies/new')">Create New Policy</button>
      <input type="text" placeholder="Search..." disabled />
    </div>
    <div class="page-body">
      <TableView :headings="['Name', 'Payloads']">
        <tr v-for="policy in policies" :key="policy.id">
          <td>
            <NuxtLink :to="'/policies/' + policy.id" exact>{{
              policy.name
            }}</NuxtLink>
          </td>
          <td>
            <!-- {{ policy.payloads.join(', ') }} -->
          </td>
        </tr>
      </TableView>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      policies: [],
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
