<template>
  <div ref="page" class="page-body">
    <p class="field-title">Targets:</p>

    <TableView :headings="['Type', 'File']">
      <tr v-for="(target, i) in app.targets" :key="i">
        <td>{{ target.msi_file !== '' ? 'MSI' : 'Unknown' }}</td>
        <td>
          <button @click="downloadFile(target.msi_file)">Download</button>
          <input
            type="file"
            :name="target.msi_file"
            :disabled="!$store.state.dashboard.editting"
          />
        </td>
      </tr>
    </TableView>
    <select ref="target-type">
      <option value="msi">MSI</option>
      <option value="msfb">Microsoft Store for Business</option>
    </select>
    <button @click="addTarget">Add Target</button>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import resource from '@/mixins/resource'

function saveBlob(blob, fileName) {
  const a = document.createElement('a')
  document.body.appendChild(a)
  a.style = 'display: none'

  const url = window.URL.createObjectURL(blob)
  a.href = url
  a.download = fileName
  a.click()
  window.URL.revokeObjectURL(url)
}

export default Vue.extend({
  mixins: [resource],
  props: {
    app: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async save(patch: object) {
      this.$el
        .querySelectorAll("input[type='file']")
        .forEach(async (node: HTMLInputElement) => {
          if (node.files.length <= 0) {
            return
          }

          await this.$store.dispatch('objects/upload', {
            id: node.name,
            file: node.files[0],
          })

          node.value = ''
        })
    },
    async downloadFile(object_id: string) {
      const { blob, filename } = await this.$store.dispatch(
        'objects/get',
        object_id
      )
      saveBlob(blob, filename)
    },
    async addTarget() {},
  },
})
</script>

<style scoped></style>
