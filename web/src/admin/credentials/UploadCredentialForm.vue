<template>
  <form class="space-y-3">
    <div>
      <div v-if="error" class="text-red-500">{{ error }}</div>
      <textarea v-model="localValue" rows="10" cols="30"></textarea>
    </div>
  </form>
</template>
<script>
export default {
  props: {
    value: Object,
    mode: String
  },
  data () {
    return {
      localValue: JSON.stringify(this.value, null, 2)
    }
  },
  emits: ['input'],
  watch: {
    localValue (newValue) {
      try {
        this.error = '';
        let parsedInput = JSON.parse(newValue)
        this.$emit('input', parsedInput);
      } catch (e) {
        this.error = 'Invalid JSON format';
      }
    }
  }
}
</script>