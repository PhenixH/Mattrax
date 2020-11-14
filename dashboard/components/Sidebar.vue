<template>
  <aside :class="{ sidebar: true, icons: !$store.state.dashboard.menuActive }">
    <NuxtLink to="/" exact>
      <HomeIcon />
      <span v-if="$store.state.dashboard.menuActive">Dashboard</span>
    </NuxtLink>
    <NuxtLink to="/devices" exact>
      <PhoneIcon />
      <span v-if="$store.state.dashboard.menuActive">Devices</span>
    </NuxtLink>
    <NuxtLink to="/users" exact>
      <UserIcon />
      <span v-if="$store.state.dashboard.menuActive">Users</span>
    </NuxtLink>
    <NuxtLink to="/policies" exact>
      <BookIcon />
      <span v-if="$store.state.dashboard.menuActive">Policies</span>
    </NuxtLink>
    <NuxtLink to="/groups" exact>
      <GridIcon />
      <span v-if="$store.state.dashboard.menuActive">Groups</span>
    </NuxtLink>
    <!-- <NuxtLink to="/reporting" exact>
      <DatabaseIcon />
      <span v-if="$store.state.dashboard.menuActive">Reporting</span>
    </NuxtLink> -->
    <NuxtLink to="/settings" exact>
      <CogIcon />
      <span v-if="$store.state.dashboard.menuActive">Settings</span>
    </NuxtLink>
  </aside>
</template>

<script lang="ts">
export default {}
</script>

<style scoped>
.sidebar {
  height: 100%;
  width: 250px;
  position: fixed;
  left: 0;
  top: 50px;
  background: var(--secondary-color);
  color: var(--light-text-color);
  transition: all 0.2s linear;
}

.sidebar.icons {
  width: 54px;
}

.sidebar a {
  display: block;
  padding: 10px 15px;
  text-decoration: none;
  color: inherit;
}

.sidebar a:hover {
  background: var(--secondary-color-accent);
}

.sidebar a svg {
  display: inline-block;
  vertical-align: middle;
}

.sidebar a span {
  padding-left: 8px;
  vertical-align: middle;
}
</style>
