<template>
  <div>
    <h1 class="mb-2">Issued Credentials</h1>
    <ErrorMessage v-if="fetchError" :message="fetchError" :title="'Could not fetch data'"/>
    <section>
      <label for="credentialTypes" class="inline">Credential types (comma-separated): </label>
      <input type="text" id="credentialTypes" v-model="credentialTypes" v-on:change="fetchData" class="inline" style="width: 50%">
      <table class="table w-full divide-y divide-gray-200 mt-4" v-if="credentials.length > 0">
        <thead>
        <tr>
          <th class="thead">Issuer</th>
          <th class="thead">Subject</th>
          <th class="thead">Type</th>
          <th class="thead">Status</th>
          <th class="thead">Issuance date</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="credential in credentials" :key="credential.id"
          @click="chosenCredential = credential" style="cursor: pointer">
          <td>{{ credential.issuer }}</td>
          <td>{{ Array.isArray(credential.credentialSubject) ? credential.credentialSubject[0].id : credential.credentialSubject.id }}</td>
          <td>{{ credential.type.filter(t => t !== 'VerifiableCredential').join(', ') }}</td>
          <td>
            <span :class="statusClass(credential.status)">{{ credential.status }}</span>
          </td>
          <td>{{ new Date(credential.issuanceDate).toLocaleString() }}</td>
        </tr>
        </tbody>
      </table>
      <p v-else>
        No credentials found.
      </p>
    </section>
  </div>
  <CredentialDetails
    :credential="chosenCredential"
    v-if="chosenCredential"
    :showRevokeButton="true"
    @revoked="fetchData"
    style="margin-top: 20px;" />
</template>

<script>
import CredentialDetails from "./CredentialDetails.vue";
import ErrorMessage from "../../components/ErrorMessage.vue";

export default {
  components: {ErrorMessage, CredentialDetails},
  data() {
    return {
      fetchError: '',
      credentials: [],
      credentialTypes: '*',
      chosenCredential: undefined,
    }
  },
  mounted() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.chosenCredential = undefined
      this.$api.get('api/issuer/vc?credentialTypes=' + encodeURIComponent(this.credentialTypes))
          .then(data => {
            this.credentials = data
          })
          .catch(response => {
            this.fetchError = response
          })
    },
    statusClass(status) {
      switch(status) {
        case 'active':
          return 'text-green-600 font-semibold'
        case 'expired':
          return 'text-gray-500 font-semibold'
        case 'revoked':
          return 'text-red-600 font-semibold'
        default:
          return ''
      }
    }
  }
}
</script>
