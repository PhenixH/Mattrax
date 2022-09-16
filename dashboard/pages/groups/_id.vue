<template>
  <div v-if="loading" class="loading">Loading Group...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/groups">Groups</NuxtLink></li>
      </ul>
      <h1>{{ group.name }}</h1>
    </PageHead>
    <PageNav>
      <button
        :class="{
          active:
            this.$route.path.replace('/groups/' + this.$route.params.id, '') ==
            '',
        }"
        @click="navigate('')"
      >
        Overview
      </button>
    </PageNav>
    <NuxtChild ref="body" :group="group" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      group: {},
    }
  },
  created() {
    this.$store.commit('dashboard/setDeletable', true)
    this.$store
      .dispatch('groups/getByID', this.$route.params.id)
      .then((group) => {
        this.group = group
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/groups/' + this.$route.params.id + pathSuffix)
    },
    async delete(): Promise<string> {
      await this.$store.dispatch('groups/deleteGroup', this.$route.params.id)
      return '/groups'
    },
  },
})
</script>

<style scoped></style>
