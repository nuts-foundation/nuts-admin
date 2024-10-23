<template>
  <p v-if="fetchError" class="m-4">Error: {{ fetchError }}</p>
  <CredentialDetails :credential="credential" v-if="credential" />
</template>

<script>
import CredentialDetails from './CredentialDetails.vue'
export default {
  components: {
    CredentialDetails
  },
  data() {
    return {
      fetchError: '',
      credential: undefined,
    }
  },
  mounted() {
    this.$api.get('api/issuer/vc?id=' + encodeURIComponent(this.$route.params.credentialID))
        .then(data => {
          this.credentials = data
        })
        .catch(response => {
          this.fetchError = response
        })
  }
}
</script>
