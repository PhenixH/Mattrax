<template>
  <div v-if="loading" class="loading">Loading Groups...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
      </ul>
      <h1>Groups</h1>
    </PageHead>
    <PageNav>
      <button @click="$router.push('/groups/new')">Create New Group</button>
      <input type="text" placeholder="Search..." disabled />
    </PageNav>
    <div class="page-body">
      <TableView :headings="['Name']">
        <tr v-for="group in groups" :key="group.id">
          <td>
            <NuxtLink :to="'/groups/' + group.id" exact>{{
              group.name
            }}</NuxtLink>
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
      groups: [],
    }
  },
  created() {
    this.$store
      .dispatch('groups/getAll')
      .then((groups) => {
        this.groups = groups
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style scoped></style>
