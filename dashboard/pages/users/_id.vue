<template>
  <div v-if="loading" class="loading">Loading User...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/users">Users</NuxtLink></li>
      </ul>
      <h1>{{ user.fullname }}</h1>
    </PageHead>
    <PageNav>
      <button
        :class="{
          active:
            this.$route.path.replace('/users/' + this.$route.params.id, '') ==
            '',
        }"
        @click="navigate('')"
      >
        Overview
      </button>
      <button
        :class="{
          active:
            this.$route.path.replace('/users/' + this.$route.params.id, '') ==
            '/groups',
        }"
        @click="navigate('/groups')"
      >
        Groups
      </button>
    </PageNav>
    <NuxtChild ref="body" :user="user" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      user: {},
    }
  },
  created() {
    this.$store.commit('dashboard/setDeletable', true)
    this.$store
      .dispatch('users/getByID', this.$route.params.id)
      .then((user) => {
        this.user = user
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/users/' + this.$route.params.id + pathSuffix)
    },
    async delete(): Promise<string> {
      await this.$store.dispatch('users/deleteUser', this.$route.params.id)
      return '/users'
    },
  },
})
</script>

<style scoped></style>
