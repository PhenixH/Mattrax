<template>
  <div class="page-body">
    <h1>Server Overview</h1>
    <p v-if="info.debug_mode" class="field-title danger">
      Danger: Debug Mode Enabled
    </p>
    <p class="field-title">Mattrax Server Version: {{ info.version }}</p>
    <p class="field-title">Mattrax Server Commit: {{ info.version_commit }}</p>
    <p class="field-title">Mattrax Server Date: {{ info.version_date }}</p>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      info: {},
    }
  },
  created() {
    this.$store
      .dispatch('settings/getOverview')
      .then((info) => (this.info = info))
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style scoped></style>
