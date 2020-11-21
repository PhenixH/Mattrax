<template>
  <div>
    <Header :is-menu-active="menuActive" />
    <Sidebar :is-menu-active="menuActive" />
    <main :class="{ icons: !$store.state.dashboard.menuActive }">
      <div v-if="$store.state.dashboard.error">
        <h1>An Error Occured</h1>
        <p>{{ $store.state.dashboard.error }}</p>
      </div>
      <div v-else>
        <Nuxt ref="body" />
        <div
          v-if="$store.state.dashboard.editting !== null"
          class="page-footer"
        >
          <button
            @click="
              $store.state.dashboard.editting
                ? $refs.body.$children[0].savebtn === undefined
                  ? $refs.body.$children[0].$refs.body.savebtn()
                  : $refs.body.$children[0].savebtn()
                : $store.commit('dashboard/setEditting', true)
            "
          >
            {{ $store.state.dashboard.editting ? 'Save' : 'Edit' }}
          </button>
          <button
            v-if="$store.state.dashboard.deletable"
            :disabled="!$store.state.dashboard.editting"
            class="red"
            @click="
              $refs.body.$children[0].deletebtn === undefined
                ? $refs.body.$children[0].$refs.body.deletebtn()
                : $refs.body.$children[0].deletebtn()
            "
          >
            Delete
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  middleware: ['auth', 'administrators-only', 'tenant-required'],
  data() {
    return {
      menuActive: false,
    }
  },
  updated() {
    if (
      this.$store.state.dashboard.error !== null &&
      this.$store.state.dashboard.error.name === 'AuthError'
    ) {
      this.$store
        .dispatch('authentication/logout')
        .then(() => this.$router.push('/login'))
        .catch(console.error)
    }
  },
})
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Raleway');

:root {
  --primary-color: #0082c8;
  --primary-color-accent: #1d75b4;
  --secondary-color: #353435;
  --secondary-color-accent: #232528;
  --light-text-color: white;
}

body {
  font-family: Raleway, sans-serif;
  font-weight: 300;
  height: 100vh;
  overflow: hidden;
  background-color: #f2f2f2;
}

h1,
h2 {
  font-weight: 400;
}

h1,
h2,
h5,
p {
  margin: 0;
  padding: 5px;
}

.mttx-brand {
  font-size: 28px;
  font-weight: 400;
  color: inherit;
  letter-spacing: 1px;
  text-transform: uppercase;
  text-decoration: none;
  line-height: 50px;
}

main {
  height: 100%;
  margin: 50px 0 0 250px;
  transition: all 0.2s linear;
}

main.icons {
  margin: 50px 0 0 54px;
}

/* START OF NEW */

.page-head {
  background: white;
  padding-left: 5px;
}

.page-head h1 {
  padding-left: 15px;
}

.page-nav {
  width: 100%;

  background: white;
  margin: 0;
  box-shadow: 0 0 0 2px lightgrey;
}

.page-nav button {
  background-color: inherit;
  border: none;
  outline: none;
  padding: 14px 16px;
  border-bottom: 3px solid rgba(255, 255, 255, 0);
}

.page-nav button.active {
  border-bottom: 3px solid #353435;
}

.page-nav button:hover {
  border-bottom: 3px solid #616161;
  transition: 0.1s;
}

.page-body {
  padding: 15px;
  overflow: scroll;

  /* min-height: 100%; */
  /* height: auto !important; */
  /* height: calc(100% - 50%); */
  /* background: orange; */
  /* height: 100%; */
  /* top: 0; */
  /* bottom: 100px;
  position: relative;
  overflow: scroll; */
}

.danger {
  color: #dc3545;
  font-weight: 600;
}

.page-body input {
  display: block;
  margin: 10px;
  padding: 5px;
  width: 100%;
  max-width: 300px;
}

.page-body select {
  display: block;
  margin: 10px;
  padding: 5px;
  width: 100%;
  max-width: 300px;
}

.page-footer {
  width: 100%;
  background: #fff;
  border-top: 3px solid #353435;
  padding: 10px;
  position: absolute;
  bottom: 0;
}

.page-footer button {
  background-color: var(--primary-color-accent);
  border: none;
  outline: none;
  color: white;
  padding: 10px 32px;
  font-size: 16px;
  margin: 10px;
  width: 100px;
}

.page-footer button.red {
  background-color: #dc3545;
}

.page-footer button.red:disabled {
  background-color: #e35d6a;
}

.loading {
  margin: 10px;
}

.breadcrumb {
  font-size: 0.7em;
}

ul.breadcrumb {
  margin: 0;
  padding: 5px 10px;
  list-style: none;
}

ul.breadcrumb li {
  display: inline;
  font-size: 0.9em;
}

ul.breadcrumb li + li:before {
  color: black;
  content: '/\00a0';
}

ul.breadcrumb li a {
  color: var(--primary-color);
  text-decoration: none;
}

ul.breadcrumb li a:hover {
  color: var(--secondary-color-accent);
  text-decoration: underline;
}
</style>
