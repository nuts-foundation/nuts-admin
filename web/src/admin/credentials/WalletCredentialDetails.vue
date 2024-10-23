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
      fetchError: undefined,
      credential: undefined
    }
  },
  mounted() {
    this.$api.get('api/id/' + encodeURIComponent(this.$route.params.subjectID))
        .then(data => {
          this.credential = data.wallet_credentials.filter(c => c.id === this.$route.params.credentialID)[0]
        })
        .catch(response => {
          this.fetchError = response
        })
  },
}
</script>
