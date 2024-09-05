<template>
  <form class="space-y-3">
    <div>
      <label for="subject-input">Name</label>
      <input type="text" v-model="localValue.subject" id="subject-input">
      <p>
        The name is for how you refer to the identity when calling the Nuts node's APIs.
        It can be set to have a human-readable identifier, e.g. "hospital_x". Allowed characters are <code>a-z A-Z 0-9 _ -</code>.
      </p>
      <p>If not set, a random GUID will be generated.</p>
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
      localValue: this.value
    }
  },
  emits: ['input'],
  watch: {
    localValue () {
      this.$emit('input', this.localValue)
    }
  }
}
</script>
