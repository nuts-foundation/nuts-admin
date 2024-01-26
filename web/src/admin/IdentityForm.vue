<template>
  <form class="space-y-3">
    <div>
      <label for="identity-did-qualifier-input">DID Qualifier</label>
      <input type="text" v-model="localValue.did_qualifier" id="identity-did-qualifier-input">
      <p>
        The DID Qualifier is a postfix that is added to the did:web of your Nuts node.
        It can be used as human-readable part of the DID, e.g. "hospital_x".
        If not set, a GUID will be used.
      </p>
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
