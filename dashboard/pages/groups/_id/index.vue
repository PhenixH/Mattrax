<template>
  <div ref="page" class="page-body">
    <p class="field-title">Name:</p>
    <input name="name" :value="group.name" type="text" :disabled="!editting" />

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

export default Vue.extend({
  props: {
    group: {
      type: Object,
      required: true,
    },
    editting: {
      type: Boolean,
      required: true,
    },
  },
  mounted() {
    // TODO: Remove need for ref with this or something
    this.$refs.page
      .querySelectorAll('input, select, checkbox, textarea')
      .forEach((node: HTMLInputElement) => {
        node.defaultValue = node.value
      })
  },
})
</script>

<style scoped></style>
