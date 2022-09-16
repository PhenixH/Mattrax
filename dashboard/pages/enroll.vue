<template>
  <div class="page-content">
    <span
      @click="
        $route.name === 'enroll'
          ? $router.push('/devices')
          : $router.replace('./')
      "
    >
      <LogoutIcon class="back-btn" view-box="0 0 8 8" height="20" width="20" />
    </span>
    <h1 class="title">Mattrax Enrollment</h1>

    <NuxtChild />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  middleware: ['auth'],
})
</script>

<style scoped>
.back-btn {
  position: absolute;
  top: 15px;
  right: 10px;
}
</style>
