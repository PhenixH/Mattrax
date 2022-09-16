<template>
  <div ref="page" class="page-body">
    <p class="field-title">Name:</p>
    <input
      name="name"
      :value="policy.name"
      type="text"
      :disabled="!$store.state.dashboard.editting"
    />

    <p class="field-title">Type:</p>
    <select v-model="policy.type" name="type" disabled>
      <option v-for="(v, key) in payloads_json" :key="key" :value="key">
        {{ v.display_name }}
      </option>
    </select>

    <div v-if="active_payload !== undefined" ref="payloads">
      <h2 class="field-title">{{ active_payload.display_name }}</h2>
      <div v-for="(field, field_id) in active_payload.fields" :key="field_id">
        <p class="field-title">{{ field.display_name }}</p>
        <select
          v-if="field.type === 'select'"
          :name="field_id"
          data-type="select"
          :disabled="!$store.state.dashboard.editting"
        >
          <option
            v-for="(display_name, value) in field.values"
            :key="value"
            :value="value"
          >
            {{ display_name }}
          </option>
        </select>
        <input
          v-else-if="field.type === 'checkbox'"
          :name="`payload.${policy.type}.${field_id}`"
          type="checkbox"
          :checked="policy.payload[policy.type][field_id]"
          value=""
          :disabled="!$store.state.dashboard.editting"
        />
        <input
          v-else
          :name="`payload.${policy.type}.${policy.type}.${field_id}`"
          :type="field.type !== null ? field.type : 'text'"
          :value="policy.payload[policy.type][field_id]"
          :placeholder="field.placeholder"
          :disabled="!$store.state.dashboard.editting"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import resource from '@/mixins/resource'
import policiesJson from '@/policies.json'

export default Vue.extend({
  mixins: [resource],
  props: {
    policy: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      payloads_json: policiesJson.payloads,
    }
  },
  computed: {
    active_payload() {
      return this.payloads_json[this.policy.type]
    },
  },
  created() {
    if (this.policy.payload[this.policy.type] === undefined)
      this.policy.payload[this.policy.type] = {}
  },
  methods: {
    async save(patch: object) {
      await this.$store.dispatch('policies/patchPolicy', {
        id: this.$route.params.id,
        patch,
      })

      Object.keys(patch).forEach((key) => (this.policy[key] = patch[key]))
    },
  },
})
</script>

<style scoped></style>
