<template>
  <div class="page-body">
    <TableView :headings="['Group', 'Group Policies']">
      <tr v-for="group in groups" :key="group.id">
        <td>
          <NuxtLink :to="'/groups/' + group.id">{{ group.name }}</NuxtLink>
        </td>
        <td>
          <span v-for="policy in group.policies" :key="policy.id">
            <NuxtLink :to="'/policies/' + policy.id">{{ policy.name }}</NuxtLink
            >,
          </span>
        </td>
      </tr>
    </TableView>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      // TODO: Get from API
      groups: [
        // {
        //   id: '9xBb3sxf',
        //   name: 'Student Restrictions',
        //   policies: [
        //     {
        //       id: 'LxsK17QR',
        //       name: 'Block Camera Restriction',
        //     },
        //     {
        //       id: 'LxsK17QB',
        //       name: 'Block Camera Restriction2',
        //     },
        //   ],
        // },
      ],
    }
  },
})
</script>

<style scoped></style>
