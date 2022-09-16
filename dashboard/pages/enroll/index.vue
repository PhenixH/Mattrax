<template>
  <div class="protocol-selector">
    <span @click="$router.push('/enroll/android')">
      <AndroidIcon height="60" width="60" />
    </span>
    <span @click="$router.push('/enroll/windows')">
      <WindowsIcon height="60" width="60" />
    </span>
    <span @click="$router.push('/enroll/apple')">
      <AppleIcon height="60" width="60" />
    </span>
    <span @click="$router.push('/enroll/agent')">
      <AgentIcon height="60" width="60" />
    </span>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({})
</script>

<style scoped>
.protocol-selector svg {
  width: 100px;
  border: 1px solid lightgrey;
  display: inline;
}
</style>
