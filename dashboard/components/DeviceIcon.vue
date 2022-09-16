<template>
  <WindowsIcon v-if="protocol == 'windows'" />
  <AgentIcon v-else-if="protocol == 'agent'" />
  <AppleIcon v-else-if="protocol == 'apple'" />
  <AndroidIcon v-else-if="protocol == 'android'" />
  <span v-else></span>
</template>

<script lang="ts">
export default {
  props: {
    protocol: {
      type: String,
      required: true,
    },
  },
}
</script>

<style scoped></style>
