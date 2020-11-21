<template>
  <div v-if="loading" class="loading">Loading Policy...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/policies">Policies</NuxtLink></li>
      </ul>
      <h1>{{ policy.name }}</h1>
    </PageHead>
    <PageNav>
      <button
        :class="{
          active:
            this.$route.path.replace(
              '/policies/' + this.$route.params.id,
              ''
            ) == '',
        }"
        @click="navigate('')"
      >
        Overview
      </button>
    </PageNav>
    <NuxtChild ref="body" :policy="policy" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      policy: {},
    }
  },
  created() {
    this.$store.commit('dashboard/setDeletable', true)
    this.$store
      .dispatch('policies/getByID', this.$route.params.id)
      .then((policy) => {
        this.policy = policy
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/policies/' + this.$route.params.id + pathSuffix)
    },
    async delete(): Promise<string> {
      await this.$store.dispatch('policies/deletePolicy', this.$route.params.id)
      return '/policies'
    },
  },
})
</script>

<style scoped></style>
