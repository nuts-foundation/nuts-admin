<template>
  <ErrorMessage v-if="fetchError" :message="fetchError" :title="'Could not fetch credential'"/>
  <CredentialDetails :credential="credential" v-if="credential" />
</template>

<script>
import CredentialDetails from './CredentialDetails.vue'
import ErrorMessage from "../../components/ErrorMessage.vue";
export default {
  components: {
    ErrorMessage,
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
