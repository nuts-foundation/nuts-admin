<template>
  <div>
    <h1>Identity: {{ this.$route.params.subjectID }}</h1>
    <ErrorMessage v-if="fetchError" :message="fetchError"/>
    <div v-if="details">
      <section>
        <header>DID Documents</header>
        <table class="min-w-full" v-if="details.did_documents.length > 0">
          <tbody>
          <tr v-for="didDocument in details.did_documents" :key="didDocument.id">
            <td>
              <pre v-on:click="showDIDDocument(didDocument.id)" style="cursor: pointer;">{{ didDocument.id }}</pre>
              <pre v-if="shownDIDDocument === didDocument.id">{{ JSON.stringify(didDocument, null, 2) }}</pre>
            </td>
          </tr>
          </tbody>
        </table>
      </section>
      <section>
        <header>Discovery Services</header>
        <table class="min-w-full divide-y divide-gray-200" v-if="details.discovery_services.length > 0">
          <thead>
          <tr>
            <th class="thead">Service</th>
            <th class="thead">Required Credential(s)</th>
            <th class="thead">Status</th>
            <th class="thead">Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="service in details.discovery_services" :key="service.id">
            <td>{{ service.id }}</td>
            <td>
              {{ discoveryServices[service.id].join(', ') }}
            </td>
            <td v-if="service.active && service.vps && service.vps.length > 0" class="whitespace-nowrap">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="green"
                   class="w-6 h-6 inline-block">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
              </svg>
              active
            </td>
            <td v-if="service.active && (!service.vps || service.vps.length == 0)" class="whitespace-nowrap">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                   stroke="DarkOrange" class="w-6 h-6 inline-block">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"/>
              </svg>
              missing credentials
            </td>
            <td v-if="!service.active" class="whitespace-nowrap">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                   stroke="currentColor" class="w-6 h-6 inline-block">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"/>
              </svg>
              not active
            </td>
            <td>
              <button v-if="service.active" class="btn btn-primary" @click="deactivateService(service.id)">
                Deactivate
              </button>
              <button v-else class="btn btn-primary" @click="$router.push({name: 'admin.activateDiscoveryService', params: {subjectID: this.$route.params.subjectID, discoveryServiceID: service.id}})">
                Activate
              </button>
            </td>
          </tr>
          </tbody>
        </table>
        <p v-else>
          No Discovery Services configured in the Nuts node.
        </p>
      </section>
      <section>
        <header>Credentials in Wallet</header>
        <table class="min-w-full divide-y divide-gray-200" v-if="details.wallet_credentials.length > 0">
          <thead>
          <tr>
            <th class="thead">Type</th>
            <th class="thead">Issuer</th>
            <th class="thead">Actions</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="credential in details.wallet_credentials" :key="credential.id"
              style="cursor: pointer">
            <td @click="$router.push({name: 'admin.credentialDetails', params: {subjectID: this.$route.params.subjectID, credentialID: credential.id}})">{{ credential.type.filter(t => t !== "VerifiableCredential").join(', ') }}</td>
            <td @click="$router.push({name: 'admin.credentialDetails', params: {subjectID: this.$route.params.subjectID, credentialID: credential.id}})">{{ credential.issuer }}</td>
            <td>
              <button class="btn btn-danger" @click="deleteCredential(credential.id)">
                Delete
              </button>
            </td>
          </tr>
          </tbody>
        </table>
        <p v-else>
          No credentials in wallet.
        </p>
        <br>
        <button
            id="upload-credential-button"
            @click="$router.push({name: 'admin.uploadCredential', params: {subjectID: this.$route.params.subjectID}})"
            class="btn btn-primary"
        >
          Upload
        </button>
      </section>
    </div>
  </div>
</template>

<script>

import DiscoveryServiceDefinition from "./DiscoveryServiceDefinition";
import ErrorMessage from "../components/ErrorMessage.vue";
import UploadCredential from "./credentials/UploadCredential.vue";

export default {
  components: {ErrorMessage, UploadCredential},
  data() {
    return {
      fetchError: undefined,
      details: undefined,
      shownDIDDocument: undefined,
      discoveryServices: {},
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    updateStatus(event) {
      this.$emit('statusUpdate', event)
    },
    fetchData() {
      this.$api.get('api/id/' + this.$route.params.subjectID)
          .then(data => {
            this.details = data
          })
          .catch(response => {
            this.fetchError = response
          })
      this.$api.get('api/proxy/internal/discovery/v1')
          .then(data => {
            // data is returned as array, convert into an object
            this.discoveryServices = data.reduce((result, curr) => {
              result[curr.id] = new DiscoveryServiceDefinition(curr).credentials()
              return result
            }, {})
          })
          .catch(response => {
            this.fetchError = this.fetchError + response.statusText
            this.discoveryServices = {}
          })
    },
    showDIDDocument(id) {
      this.shownDIDDocument = this.shownDIDDocument === id ? undefined : id
    },
    deactivateService(id) {
      this.fetchError = undefined
      this.$api.delete(`api/proxy/internal/discovery/v1/${id}/${this.details.subject}`)
          .then(data => {
            if (data.reason) {
              this.fetchError = data.reason
            }
          })
          .catch(response => {
            this.fetchError = response
          })
          .finally(() => {
            this.fetchData()
          })
    },
    deleteCredential(credentialId) {
      if (confirm("Are you sure you want to delete this credential?") !== true) {
        return
      }
      this.fetchError = undefined
      this.$api.delete(`api/proxy/internal/vcr/v2/holder/${encodeURIComponent(this.$route.params.subjectID)}/vc/${encodeURIComponent(credentialId)}`)
          .then(data => {
            if (data.reason) {
              this.fetchError = data.reason
            }
          })
          .catch(response => {
            this.fetchError = response
          })
          .finally(() => {
            this.fetchData()
          })
    }
  }
}
</script>
