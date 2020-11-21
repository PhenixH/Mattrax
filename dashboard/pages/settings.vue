<template>
  <div>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
      </ul>
      <h1>Settings</h1>
    </PageHead>
    <PageNav>
      <button
        :class="{
          active: this.$route.path.replace('/settings', '') == '',
        }"
        @click="navigate('')"
      >
        Overview
      </button>
      <button
        :class="{
          active: this.$route.path.replace('/settings', '') == '/tenant',
        }"
        @click="navigate('/tenant')"
      >
        Tenant
      </button>
      <button
        :class="{
          active: this.$route.path.replace('/settings', '') == '/user',
        }"
        @click="navigate('/user')"
      >
        User
      </button>
    </PageNav>
    <NuxtChild ref="body" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/settings' + pathSuffix)
    },
  },
})
</script>

<style scoped></style>
