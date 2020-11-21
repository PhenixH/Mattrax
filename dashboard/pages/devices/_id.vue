<template>
  <div v-if="loading" class="loading">Loading Device...</div>
  <div v-else>
    <PageHead>
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
        <li><NuxtLink to="/devices">Devices</NuxtLink></li>
      </ul>
      <h1>{{ device.name }}</h1>
    </PageHead>
    <PageNav>
      <button
        :class="{
          active:
            this.$route.path.replace('/devices/' + this.$route.params.id, '') ==
            '',
        }"
        @click="navigate('')"
      >
        Overview
      </button>
      <button
        :class="{
          active:
            this.$route.path.replace('/devices/' + this.$route.params.id, '') ==
            '/scope',
        }"
        @click="navigate('/scope')"
      >
        Scope
      </button>
    </PageNav>
    <NuxtChild ref="body" :device="device" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  data() {
    return {
      loading: true,
      device: {},
    }
  },
  created() {
    this.$store
      .dispatch('devices/getByID', this.$route.params.id)
      .then((device) => {
        this.device = device
        this.loading = false
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/devices/' + this.$route.params.id + pathSuffix)
    },
    async delete(): Promise<string> {
      await this.$store.dispatch('devices/deleteDevice', this.$route.params.id)
      return '/devices'
    },
  },
})
</script>

<style scoped></style>
