<template>
  <div v-if="loading" class="loading">Loading Applications...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
      </ul>
      <h1>Applications</h1>
    </PageHead>
    <PageNav>
      <button @click="$router.push('/applications/new')">Add Application</button>
      <input type="text" placeholder="Search..." disabled />
    </PageNav>
    <div class="page-body">
      <TableView :headings="['Name', 'Supported Platforms']">
        <tr v-for="app in applications" :key="app.id">
          <td>
            <NuxtLink :to="'/applications/' + app.id" exact>{{
              app.name
            }}</NuxtLink>
          </td>
          <td>
            
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
      applications: [],
    }
  },
  created() {
    this.$store
      .dispatch('policies/getAll')
      .then((applications) => {
        this.applications = applications
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style scoped></style>
