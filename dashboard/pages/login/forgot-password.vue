<template>
  <div v-if="loading" class="loading">Checking 2FA...</div>
  <div v-else class="page-content">
    <p class="title">Password Reset Email:</p>
    <form class="form" @submit.prevent="submit">
      <p v-if="errorTxt" class="error-msg">{{ errorTxt }}</p>
      <input
        v-model="email"
        required
        type="text"
        placeholder="chris@otbeaumont.me"
        maxlength="100"
        autocomplete="new-password"
        @input="errorTxt = null"
      />
      <button>SEND PASSWORD RESET EMAIL</button>
    </form>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      loading: false,
      errorTxt: null,
      email: '',
    }
  },
  async created() {
    if (await this.$store.dispatch('authentication/isAuthenticated')) {
      this.$router.push({ path: '/login/tenants', query: this.$route.query })
    }
  },
  methods: {
    submit() {},
  },
})
</script>

<style scoped>
.form input {
  outline: 0;
  background: #f2f2f2;
  width: 100%;
  border: 0;
  margin: 0 0 15px;
  padding: 15px;
  box-sizing: border-box;
  font-size: 14px;
}
.form .error-msg {
  margin-bottom: 5px;
  color: red;
  font-size: 13px;
}
.title {
  font-size: 1.5em;
}
.tenant-list button {
  margin: 5px;
}
.tenant-list .create-btn {
  background-color: #353435;
}
.logout {
  float: left;
}
.logout:hover {
  border-bottom: 1px solid black;
}
.logout-btn {
  position: absolute;
  top: 15px;
  right: 10px;
}
</style>
