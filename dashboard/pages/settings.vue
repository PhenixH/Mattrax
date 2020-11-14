<template>
  <div v-if="loading" class="loading">Loading Settings...</div>
  <div v-else>
    <div class="page-head">
      <ul class="breadcrumb">
        <li><NuxtLink to="/">Dashboard</NuxtLink></li>
      </ul>
      <h1>Settings</h1>
    </div>
    <div class="page-nav">
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
    </div>
    <div class="page-body">
      <NuxtChild :editting="editting" />
    </div>
    <div class="page-footer">
      <button @click="editting ? save() : (editting = true)">
        {{ editting ? 'Save' : 'Edit' }}
      </button>
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
      editting: false,
    }
  },
  methods: {
    navigate(pathSuffix: string) {
      this.$router.push('/settings' + pathSuffix)
    },
    save() {
      if (this.editting === false) return

      let patch = null
      this.$children[1].$el
        .querySelectorAll('input, select, checkbox, textarea')
        .forEach((node) => {
          if (node.value !== node.defaultValue) {
            if (patch === null) patch = {}
            if (node.getAttribute('data-type') === 'bool') {
              patch[node.name] = node.value === 'true'
            } else {
              patch[node.name] = node.value
            }
          }
        })
      if (patch === null) {
        this.editting = false
        return
      }

      this.$children[1]
        .save(patch)
        .then(() => (this.editting = false))
        .catch((err) => this.$store.commit('dashboard/setError', err)) // TODO: Warning that saving failed
    },
  },
})
</script>

<style scoped></style>
