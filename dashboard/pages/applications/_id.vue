<template>
  <div v-if="loading" class="loading">Loading Application...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/applications">Applications</NuxtLink></li>
      </ul>
      <h1>{{ app.name }}</h1>
    </PageHead>
    <PageNav>
      <button
        :class="{
          active:
            this.$route.path.replace(
              '/applications/' + this.$route.params.id,
              ''
            ) == '',
        }"
        @click="navigate('')"
      >
        Overview
      </button>
    </PageNav>
    <NuxtChild ref="body" :app="app" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      app: {},
    }
  },
  created() {
    this.$store.commit('dashboard/setDeletable', true)
    this.$store
      .dispatch('policies/getByID', this.$route.params.id)
      .then((app) => {
        this.app = app
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/applications/' + this.$route.params.id + pathSuffix)
    },
    async delete(): Promise<string> {
      // await this.$store.dispatch('policies/deletePolicy', this.$route.params.id)
      return '/applications'
    },
  },
})
</script>

<style scoped></style>
