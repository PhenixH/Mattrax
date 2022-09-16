<template>
  <div v-if="loading" class="loading">Creating User...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/users">Users</NuxtLink></li>
      </ul>
      <h1>Create New User</h1>
    </PageHead>
    <div class="page-body">
      <div>
        <form class="create-form" @submit.prevent="createUser">
          <p class="field-title">Name:</p>
          <input
            v-model="user.fullname"
            type="text"
            placeholder="Chris Green"
            required
          />
          <p class="field-title">Email:</p>
          <input
            v-model="user.upn"
            type="email"
            placeholder="chris@mattrax.app"
            required
          />
          <p class="field-title">Password:</p>
          <PasswordField :value.sync="user.password" />

          <button type="submit" class="submit">Create User</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: false,
      user: {
        fullname: '',
        upn: '',
        password: '',
      },
    }
  },
  methods: {
    createUser() {
      this.loading = true
      this.$store
        .dispatch('users/createUser', this.user)
        .then(() => this.$router.push('/users/' + this.user.upn))
        .catch((err) => this.$store.commit('dashboard/setError', err))
    },
  },
})
</script>

<style scoped>
.create-form {
  padding: 10px;
}

.create-form input {
  display: block;
  margin: 10px;
  padding: 5px;
  width: 100%;
  max-width: 300px;
}

.submit {
  background-color: var(--primary-color-accent);
  border: none;
  color: white;
  padding: 10px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  margin: 10px;
  width: 100%;
  max-width: 300px;
}
</style>
