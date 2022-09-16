<template>
  <div class="password-field">
    <input
      :value="value"
      :type="passwordVisible ? 'text' : 'password'"
      placeholder="password"
      minlength="8"
      maxlength="100"
      :disabled="disabled"
      @input="(e) => $emit('update:value', e.target.value)"
    />
    <span
      @mousedown="passwordVisible = true"
      @mouseup="passwordVisible = false"
    >
      <EyeIcon view-box="0 0 8 8" height="24" width="24" />
    </span>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    value: {
      required: true,
    },
    disabled: {
      type: Boolean,
    },
  },
  data() {
    return {
      passwordVisible: false,
    }
  },
}
</script>

<style scoped>
.password-field {
  width: 300px;
  position: relative;
}

.password-field span {
  position: absolute;
  top: 0;
  right: 0px;
  z-index: 2;
  top: 2.5px;
  cursor: pointer;
}
</style>
