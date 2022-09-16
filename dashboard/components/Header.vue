<template>
  <nav class="nav">
    <span
      :class="{
        'disable-select': true,
        'menu-toggle': true,
        right: $store.state.dashboard.menuActive,
      }"
      @click="$store.commit('dashboard/toggleMenuActive')"
    >
      <MenuIcon />
    </span>
    <NuxtLink to="/" exact class="mttx-brand">
      Mattrax -
      <NuxtLink to="/login/tenants" exact
        ><span>{{
          $store.state.tenants.tenant !== null
            ? $store.state.tenants.tenant.display_name
            : ''
        }}</span></NuxtLink
      >
    </NuxtLink>

    <div class="navRight">
      <!-- <span class="navNotifications"><NotificationIcon /></span> -->

      <div class="dropdown">
        <span class="navUser">
          {{ $store.state.authentication.user.name }}
          <CaretIcon />
        </span>
        <div class="dropdown-content">
          <NuxtLink to="/settings/user"> Edit Account </NuxtLink>
          <a href="#" @click.prevent="logout()">Logout</a>
        </div>
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
export default {
  methods: {
    logout() {
      this.$store
        .dispatch('authentication/logout')
        .then(() => {
          this.$router.push('/login')
        })
        .catch((err) => {
          console.error(err)
        })
    },
  },
}
</script>

<style scoped>
.nav {
  height: 50px;
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  background: var(--primary-color);
  color: var(--light-text-color);
  padding: 0 15px;
}

.menu-toggle svg {
  color: white;
  transition: all 0.2s linear;
}

.menu-toggle.right svg {
  transform: rotate(90deg);
  transition: all 0.2s linear;
}

.mttx-brand span {
  font-size: 14px;
}

.mttx-brand a {
  color: inherit;
  text-decoration: none;
}

.mttx-brand a:hover {
  color: #d8d8d8;
}

.navRight {
  float: right;
  height: 100%;
  line-height: 50px;
}

.navRight svg {
  display: inline-block;
  vertical-align: middle;
}

.navNotifications {
  float: left;
  padding: 0 7px;
}

.navUser {
  /* float: right; */
  padding: 0 5px 0 7px;
  cursor: default;
}

.dropdown {
  position: relative;
  display: inline-block;
}

.dropdown-content {
  display: none;
  position: absolute;
  background-color: #f9f9f9;
  min-width: 160px;
  box-shadow: 0px 8px 16px 0px rgba(0, 0, 0, 0.2);
  z-index: 1;
}

.dropdown-content a {
  color: black;
  padding: 0px 12px;
  text-decoration: none;
  display: block;
}

.dropdown-content a:hover {
  background-color: #f1f1f1;
}

.dropdown:hover .dropdown-content {
  display: block;
}

.disable-select {
  user-select: none; /* supported by Chrome and Opera */
  -webkit-user-select: none; /* Safari */
  -khtml-user-select: none; /* Konqueror HTML */
  -moz-user-select: none; /* Firefox */
  -ms-user-select: none; /* Internet Explorer/Edge */
}
</style>
