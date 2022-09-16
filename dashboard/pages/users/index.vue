<template>
  <div v-if="loading" class="loading">Loading Users...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
      </ul>
      <h1>Users</h1>
    </PageHead>
    <PageNav>
      <button @click="$router.push('/users/new')">Create New User</button>
      <input type="text" placeholder="Search..." disabled />
    </PageNav>
    <div class="page-body">
      <TableView :headings="['Name', 'UPN']">
        <tr v-for="user in users" :key="user.upn">
          <td>
            <NuxtLink :to="'/users/' + user.upn" exact>{{
              user.fullname
            }}</NuxtLink>
          </td>
          <td>
            <NuxtLink :to="'/users/' + user.upn" exact>{{ user.upn }}</NuxtLink>
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
      users: [],
    }
  },
  created() {
    this.$store
      .dispatch('users/getAll')
      .then((users) => {
        this.users = users
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style scoped></style>
