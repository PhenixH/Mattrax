<template>
  <div class="page-body">
    <p class="field-title">Name:</p>
    <input
      name="name"
      :value="group.name"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Scope:</p>
    <TableView :headings="['Scoped Policies']">
      <tr v-for="policy in scopedPolicies" :key="policy.id">
        <td>
          <NuxtLink :to="'/policies/' + policy.id">{{ policy.name }}</NuxtLink>
        </td>
      </tr>
    </TableView>
    <TableView :headings="['Scoped Devices']">
      <tr v-for="device in scopedDevices" :key="device.id">
        <td>
          <NuxtLink :to="'/devices/' + device.id">{{ device.name }}</NuxtLink>
        </td>
      </tr>
    </TableView>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import resource from '@/mixins/resource'

export default Vue.extend({
  mixins: [resource],
  props: {
    group: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      scopedPolicies: [],
      scopedDevices: [],
    }
  },
  created() {
    this.$store
      .dispatch('groups/getScopedPolicies', this.$route.params.id)
      .then((policies) => (this.scopedPolicies = policies))
      .catch((err) => this.$store.commit('dashboard/setError', err))

    this.$store
      .dispatch('groups/getScopedDevices', this.$route.params.id)
      .then((devices) => (this.scopedDevices = devices))
      .catch((err) => this.$store.commit('dashboard/setError', err))
  },
  methods: {
    async save(patch: object) {
      await this.$store.dispatch('groups/patchGroup', {
        id: this.$route.params.id,
        patch,
      })

      Object.keys(patch).forEach((key) => (this.group[key] = patch[key]))
    },
  },
})
</script>

<style scoped></style>
