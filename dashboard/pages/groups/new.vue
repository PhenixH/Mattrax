<template>
  <div v-if="loading" class="loading">Creating Group...</div>
  <div v-else>
    <div class="panel">
      <div class="panel-head">
        <h1>
          <GridIcon view-box="0 0 24 24" height="40" width="40" />Create New
          Group
        </h1>
      </div>
      <div>
        <form class="create-form" @submit.prevent="createGroup">
          <p class="field-title">Name:</p>
          <input
            v-model="group.name"
            type="text"
            placeholder="Student Devices"
            required
          />

          <button type="submit" class="submit">Create Group</button>
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
      passwordVisible: false,
      group: {
        name: '',
      },
    }
  },
  methods: {
    createGroup() {
      this.loading = true
      this.$store
        .dispatch('groups/createGroup', this.group)
        .then((gid) => this.$router.push('/groups/' + gid))
        .catch((err) => this.$store.commit('dashboard/setError', err))
    },
  },
})
</script>

<style>
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
