<template>
  <div>
    <p>
      Install the
      <a
        href="https://play.google.com/store/apps/details?id=com.google.android.apps.work.clouddpc"
        target="_blank"
        rel="noreferrer"
        >Android Device Policy</a
      >
      application and scan the code to enroll.
    </p>
    <img
      ref="qr"
      src="data:image/gif;base64,R0lGODlhAQABAAAAACwAAAAAAQABAAA="
      alt="Android Enrollment QR Code"
    />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  created() {
    this.$store
      .dispatch('android/getEnrollmentQR')
      .then((imgBlob) => {
        this.$refs.qr.src = URL.createObjectURL(imgBlob)
      })
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
})
</script>

<style scoped></style>
