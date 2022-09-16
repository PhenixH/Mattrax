<template>
  <aside :class="{ sidebar: true, icons: !$store.state.dashboard.menuActive }">
    <NuxtLink to="/" exact>
      <HomeIcon />
      <transition name="slide-fade">
        <span v-if="$store.state.dashboard.menuActive">Dashboard</span>
      </transition>
    </NuxtLink>
    <NuxtLink to="/devices" exact>
      <PhoneIcon />
      <transition name="slide-fade">
        <span v-if="$store.state.dashboard.menuActive">Devices</span>
      </transition>
    </NuxtLink>
    <NuxtLink to="/users" exact>
      <UserIcon />
      <transition name="slide-fade">
        <span v-if="$store.state.dashboard.menuActive">Users</span>
      </transition>
    </NuxtLink>
    <NuxtLink to="/applications" exact>
      <BrowserIcon />
      <transition name="slide-fade">
        <span v-if="$store.state.dashboard.menuActive">Applications</span>
      </transition>
    </NuxtLink>
    <NuxtLink to="/policies" exact>
      <BookIcon />
      <transition name="slide-fade">
        <span v-if="$store.state.dashboard.menuActive">Policies</span>
      </transition>
    </NuxtLink>
    <NuxtLink to="/groups" exact>
      <GridIcon />
      <transition name="slide-fade">
        <span v-if="$store.state.dashboard.menuActive">Groups</span>
      </transition>
    </NuxtLink>
    <!-- <NuxtLink to="/reporting" exact>
      <DatabaseIcon />
      <transition name="slide-fade">
      <span v-if="$store.state.dashboard.menuActive">Reporting</span>
      </transition>
    </NuxtLink> -->
    <NuxtLink to="/settings" exact>
      <CogIcon />
      <transition name="slide-fade">
        <span v-if="$store.state.dashboard.menuActive">Settings</span>
      </transition>
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

.slide-fade-enter-active {
  opacity: 0;
  transition: all 0.3s ease;
}
.slide-fade-leave-active {
  transition: opacity 0.1s;
}
</style>
