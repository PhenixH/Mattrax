<template>
  <div>
    <p>
      Open <i>Access Work and School</i> then enter your organisation email
      address.
    </p>

    <button v-if="platform == 'Win32'" @click="triggerEnroll()">
      ENROLL THIS DEVICE
    </button>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  computed: {
    platform: () => window.navigator.platform,
  },
  methods: {
    triggerEnroll() {
      window.open(
        'ms-device-enrollment:?mode=mdm&username=' +
          this.$store.state.authentication.user.upn
      )
    },
  },
})
</script>

<style scoped></style>
