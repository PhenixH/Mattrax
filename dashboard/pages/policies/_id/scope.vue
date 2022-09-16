<template>
  <div v-if="loading" class="loading">Loading Policy Scope...</div>
  <div v-else>
    <TableView :headings="['Groups', 'Actions']">
      <tr v-for="group in groups" :key="group.id">
        <td>
          <NuxtLink :to="'/groups/' + group.id" exact>{{
            group.name
          }}</NuxtLink>
        </td>
        <td>
          <button @click="unscopeGroup(group.id)">Unscope</button>
        </td>
      </tr>
    </TableView>
    <select ref="groupSelector">
      <option
        v-for="group in tenantGroups.filter((g) => g.visible)"
        :key="group.id"
        :value="group.id"
      >
        {{ group.name }}
      </option>
    </select>
    <button @click="addGroup()">Add Group</button>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      loading: true,
      groups: [],
      tenantGroups: [],
    }
  },
  async created() {
    await Promise.all([
      this.$store
        .dispatch('policies/getScopeByID', this.$route.params.id)
        .then((groups) => (this.groups = groups)),
      this.$store.dispatch('groups/getAll').then((groups) => {
        this.tenantGroups = groups
      }),
    ]).catch((err) => this.$store.commit('dashboard/setError', err))

    this.refreshGroupSelectorSection()
    this.loading = false
  },
  methods: {
    refreshGroupSelectorSection() {
      this.tenantGroups.forEach((group) => {
        if (this.groups.findIndex((g) => g.id === group.id) === -1) {
          group.visible = true
        } else {
          group.visible = false
        }
      })
    },
    addGroup() {
      this.$store
        .dispatch('policies/addToGroup', {
          gid: this.$refs.groupSelector.value,
          pid: this.$route.params.id,
        })
        .then(() => {
          this.groups.push({
            id: this.$refs.groupSelector.value,
            name: this.$refs.groupSelector.options[
              this.$refs.groupSelector.selectedIndex
            ].text,
          })
          this.refreshGroupSelectorSection()
        })
        .catch((err) => this.$store.commit('dashboard/setError', err))
    },
    unscopeGroup(policyID: string) {
      this.$store
        .dispatch('policies/removeFromGroup', {
          gid: policyID,
          pid: this.$route.params.id,
        })
        .then(() => {
          this.groups = this.groups.filter((group) => group.id !== policyID)
          this.refreshGroupSelectorSection()
        })
        .catch((err) => this.$store.commit('dashboard/setError', err))
    },
  },
})
</script>

<style scoped></style>
