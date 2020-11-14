<template>
  <div v-if="loading" class="loading">Loading User...</div>
  <div v-else>
    <div class="page-head">
      <ul class="breadcrumb">
        <li><NuxtLink to="/users">Users</NuxtLink></li>
      </ul>
      <h1>{{ user.fullname }}</h1>
    </div>
    <div class="page-nav">
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
    </div>
    <NuxtChild :user="user" :editting="editting" />
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
      loading: true,
      editting: false,
      user: {},
    }
  },
  created() {
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
    save() {
      if (this.editting === false) return

      let patch: any = null
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

      this.$store
        .dispatch('users/patchUser', {
          upn: this.$route.params.id,
          patch,
        })
        .then((upn) => {
          if (upn !== this.$route.params.id) {
            this.$router.push('/users/' + upn)
            return
          } else {
            Object.keys(patch).forEach((key) => (this.user[key] = patch[key]))
          }

          this.editting = false
        })
        .catch((err) => this.$store.commit('dashboard/setError', err)) // TODO: Warning that saving failed
    },
  },
})
</script>

<style scoped></style>
