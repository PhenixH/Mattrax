<template>
  <div v-if="loading" class="loading">Loading Device Scope...</div>
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

    <TableView :headings="['Computed Policies']">
      <tr v-for="policy in policies" :key="policy.name">
        <td>
          <NuxtLink :to="'/policies/' + policy.id" exact>{{
            policy.name
          }}</NuxtLink>
        </td>
      </tr>
    </TableView>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  data() {
    return {
      loading: true,
      groups: [],
      policies: [],
      tenantGroups: [],
    }
  },
  async created() {
    await Promise.all([
      this.$store
        .dispatch('devices/getScopeByID', this.$route.params.id)
        .then((scope) => {
          this.groups = scope.groups
          this.policies = scope.policies
        }),
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
        .dispatch('devices/addToGroup', {
          gid: this.$refs.groupSelector.value,
          udid: this.$route.params.id,
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
        .dispatch('devices/removeFromGroup', {
          gid: policyID,
          udid: this.$route.params.id,
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
