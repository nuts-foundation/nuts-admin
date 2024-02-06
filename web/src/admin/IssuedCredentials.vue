<template>
  <div>

    <div class="flex justify-between mb-6">
      <h1>Issued Credentials</h1>

    </div>

    <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <p v-if="fetchError" class="m-4">Could not fetch identities: {{ fetchError }}</p>
    </div>
  </div>
</template>

<script>

export default {
  data() {
    return {
      fetchError: '',
      identities: [],
    }
  },
  mounted() {
    this.fetchData()
  },
  emits: ['statusUpdate'],
  methods: {
    updateStatus(event) {
      this.fetchData()
      this.$emit('statusUpdate', event)
    },
    fetchData() {
      this.$api.get('api/id')
          .then(data => {
            this.identities = data
          })
          .catch(response => {
            this.fetchError = response
          })
    }
  }
}
</script>
