<template>
  <div>
    <h1>Identity</h1>
    <p v-if="fetchError" class="m-4">Error: {{ fetchError }}</p>
    <div v-if="details">
      <section>
        <code class="inline">{{ details.did }}</code>
        &nbsp;
        <button class="btn btn-primary btn-tiny" v-on:click="showDIDDocument = !showDIDDocument">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
               stroke="currentColor" class="w-6 h-6 inline">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="m15.75 15.75-2.489-2.489m0 0a3.375 3.375 0 1 0-4.773-4.773 3.375 3.375 0 0 0 4.774 4.774ZM21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
          </svg>
          &nbsp;
          Show DID Document
        </button>
        <pre v-if="showDIDDocument">{{ details.did_document }}</pre>
      </section>
      <section>

      </section>
      <section>
        <header>Discovery Services</header>
        <table class="min-w-full divide-y divide-gray-200">
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
            <td v-if="service.active && service.vp" class="whitespace-nowrap">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="green"
                   class="w-6 h-6 inline-block">
                <path stroke-linecap="round" stroke-linejoin="round"
                      d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
              </svg>
              active
            </td>
            <td v-if="service.active && !service.vp" class="whitespace-nowrap">
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
              <button v-else class="btn btn-primary" @click="activateService(service.id)">
                Activate
              </button>
            </td>
          </tr>
          </tbody>
        </table>
      </section>
      <section>
        <header>Credentials in Wallet</header>
        <table class="min-w-full divide-y divide-gray-200" v-if="details.wallet_credentials.length > 0">
          <thead>
          <tr>
            <th class="thead">Type</th>
            <th class="thead">Issuer</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="credential in details.wallet_credentials" :key="credential.id">
            <td>{{ credential.type.filter(t => t !== "VerifiableCredential").join(', ') }}</td>
            <td>{{ credential.issuer }}</td>
          </tr>
          </tbody>
        </table>
        <p v-else>
          No credentials in wallet.
        </p>
      </section>
    </div>
  </div>
</template>

<script>

import DiscoveryServiceDefinition from "./DiscoveryServiceDefinition";

export default {
  data() {
    return {
      fetchError: undefined,
      details: undefined,
      discoveryServices: {},
      showDIDDocument: false,
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
      this.$api.get('api/id/' + this.$route.params.id)
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
    activateService(id) {
      this.fetchError = undefined
      this.$api.post(`api/proxy/internal/discovery/v1/${encodeURIComponent(id)}/${encodeURIComponent(this.details.did)}`)
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
    deactivateService(id) {
      this.fetchError = undefined
      this.$api.delete(`api/proxy/internal/discovery/v1/${encodeURIComponent(id)}/${encodeURIComponent(this.details.did)}`)
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
