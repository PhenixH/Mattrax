<template>
  <div v-if="loading" class="loading">Creating Application...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/applications">Applications</NuxtLink></li>
      </ul>
      <h1>Add New Application</h1>
    </PageHead>
    <div class="page-body">
      <form class="create-form" @submit.prevent="createPolicy">
        <p class="field-title">Name:</p>
        <input
          v-model="policy.name"
          type="text"
          placeholder="Baseline Restrictions"
          required
        />

        <button type="submit" class="submit">Create Policy</button>
      </form>
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
      passwordVisible: false,
      policy: {
        name: '',
        type: 'custom',
      },
    }
  },
  methods: {
    createPolicy() {
      this.loading = true
      this.$store
        .dispatch('policies/createPolicy', this.policy)
        .then((pid) => this.$router.push('/policies/' + pid))
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

.password-field {
  width: 300px;
  position: relative;
}

.password-field span {
  position: absolute;
  top: 0;
  right: 0px;
  z-index: 2;
  top: 2.5px;
  cursor: pointer;
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
