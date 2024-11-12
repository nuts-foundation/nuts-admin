<template>
  <form class="space-y-3">
    <div>
      <div v-if="error" class="text-red-500">{{ error }}</div>
      <textarea v-model="localValue" rows="10" cols="30" placeholder="Verifiable Credential as JSON or JWT"></textarea>
    </div>
  </form>
</template>
<script>
export default {
  props: {
    value: Object | String,
    mode: String
  },
  data() {
    return {
      error: undefined,
      localValue: JSON.stringify(this.value, null, 2)
    }
  },
  emits: ['input'],
  watch: {
    localValue(newValue) {
      // if starts with a {, assume JSON. Otherwise, parse as JWT.
      this.error = '';
      if (newValue.startsWith('{')) {
          try {
            let parsedInput = JSON.parse(newValue)
            this.$emit('input', parsedInput);
          } catch (e) {
            this.error = 'Invalid JSON format';
          }
        } else {
          // simplistic parsing as form of validation
          let parts = newValue.split('.');
          try {
            if (parts.length !== 3) {
              this.error = 'Invalid JWT format';
            } else {
              // first and second part should be base64 encoded JSON
              JSON.parse(atob(parts[0]));
              JSON.parse(atob(parts[1]));
            }
            this.$emit('input', newValue);
          } catch (e) {
            console.log(e);
            this.error = 'Invalid JWT format';
          }
        }
    }
  }
}
</script>