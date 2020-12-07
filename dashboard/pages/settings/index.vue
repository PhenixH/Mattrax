<template>
  <div class="page-body">
    <h1>Server Overview</h1>
    <p v-if="info.debug_mode" class="field-title danger">
      Danger: Debug Mode Enabled
    </p>
    <p v-if="info.cloud_mode" class="field-title info">
      Running on Mattrax Cloud
    </p>
    <p class="field-title">
      Primary Domain:
      <a :href="'https://' + info.primary_domain" target="_blank">{{
        info.primary_domain
      }}</a>
    </p>
    <br />

    <h1>Server Resources</h1>
    <p class="field-title">
      PostgreSQL Database:
      <span v-if="info.database_status" class="safe">Online</span
      ><span v-else class="danger">Offline</span>
    </p>
    <p v-if="info.zipkin_status" class="field-title">
      Zipkin Tracing: <span class="info">Enabled</span>
    </p>
    <p
      v-if="
        info.protocols !== undefined && info.protocols.android === undefined
      "
      class="field-title"
    >
      Android Management: <span class="danger">Disabled</span>
    </p>
    <br />

    <h1>Server Version</h1>
    <p class="field-title">Server Version: {{ info.version }}</p>
    <p class="field-title">Server Commit: {{ info.version_commit }}</p>
    <p class="field-title">Server Build Date: {{ info.version_date }}</p>
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
