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
      <!-- <tr v-for="group in groups" :key="group.id">
        <td>
          <NuxtLink :to="'/groups/' + group.id">{{ group.name }}</NuxtLink>
        </td>
        <td>
          <span v-for="policy in group.policies" :key="policy.id">
            <NuxtLink :to="'/policies/' + policy.id">{{ policy.name }}</NuxtLink
            >,
          </span>
        </td>
      </tr> -->
    </TableView>
    <TableView :headings="['Scoped Devices']">
      <!-- <tr v-for="group in groups" :key="group.id">
        <td>
          <NuxtLink :to="'/groups/' + group.id">{{ group.name }}</NuxtLink>
        </td>
        <td>
          <span v-for="policy in group.policies" :key="policy.id">
            <NuxtLink :to="'/policies/' + policy.id">{{ policy.name }}</NuxtLink
            >,
          </span>
        </td>
      </tr> -->
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
