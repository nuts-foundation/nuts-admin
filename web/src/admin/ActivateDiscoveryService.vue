<template>
  <div>
    <h1>Activate Discovery Service</h1>
    <p>
      This page allows you to activate a discovery service for a subject.
    </p>
    <ErrorMessage v-if="fetchError" :message="fetchError"/>
    <section v-if="selectedDiscoveryService">
      <div>
        <label>Discovery Service</label>
        <div>{{selectedDiscoveryService.id}}</div>
      </div>
      <div>
        <label>Subject ID</label>
        <div>{{$route.params.subjectID}}</div>
      </div>
      <div>
        <label>Endpoint</label>
        <div>{{selectedDiscoveryService.endpoint}}</div>
      </div>
    </section>

    <section>
      <header>Parameters</header>
      <a v-on:click="registrationParameters.push({})" style="cursor: pointer">Add a parameter</a>
      <table class="w-full">
        <thead>
        <tr>
          <th class="w-1/3">Key</th>
          <th>Value</th>
          <th></th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(param, idx) in registrationParameters" :key="'search-' + idx">
          <td style="vertical-align: top">
            <input type="text" v-model="param.key" :id="'param-key-' + idx">
            <div v-if="param.label" style="padding-left: 10px;">{{ param.label }}</div>
          </td>
          <td style="vertical-align: top">
            <input type="text" v-model="param.value" :id="'param-value-' + idx" :placeholder="parameterPlaceholder(param.key)">
          </td>
          <td style="vertical-align: top; padding-top: 15px;">
            <svg @click="registrationParameters.splice(idx, 1)" style="cursor: pointer"
                 xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                 stroke="currentColor" class="w-6 h-6">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12H9m12 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
            </svg>

          </td>
        </tr>
        </tbody>
      </table>
    </section>

    <section>
      <input type="button" class="btn btn-primary" value="Activate" @click="activate">
    </section>
  </div>
</template>

<script>

import ErrorMessage from "../components/ErrorMessage.vue";

export default {
  components: {ErrorMessage},
  data() {
    return {
      fetchError: undefined,
      discoveryServices: [],
      selectedDiscoveryService: undefined,
      registrationParameters: [],
    }
  },
  mounted() {
    this.$api.get('api/proxy/internal/discovery/v1')
        .then(data => {
          this.services = data
          this.selectedDiscoveryService = this.services.filter(s => s.id === this.$route.params.discoveryServiceID)[0]
          let registrationParametersInputDescriptors = this.selectedDiscoveryService.presentation_definition.input_descriptors.filter(descriptor => descriptor.constraints.fields.filter(f => f.filter && f.filter.type === "string" && f.filter.const === "DiscoveryRegistrationCredential").length > 0)
          registrationParametersInputDescriptors.forEach(r => r.constraints.fields.filter(f => f.id && f.path)
              .forEach(f => {
                f.path.forEach(path => {
                  let label = f.id
                  if (f.purpose) {
                    label = "(" + label + ") " + f.purpose
                  }
                  // remove $.credentialSubject. from the beginning
                  const key = path.replace(/^\$\.credentialSubject\./, '')
                  this.registrationParameters.push({key: key, label: label})
                })
              })
          )
        })
        .catch(response => {
          this.fetchError = response.statusText
          this.services = []
        })
  },
  methods: {
    updateStatus(event) {
      this.$emit('statusUpdate', event)
    },
    parameterPlaceholder(paramName) {
      if (paramName === 'authServerURL') {
        return 'Leave empty to have the Nuts node auto-generate it.'
      }
      return '';
    },
    activate() {
      this.fetchError = undefined
      let params = {}
      this.registrationParameters.forEach(p => {
        params[p.key] = p.value
      })
      this.$api.post(`api/proxy/internal/discovery/v1/${this.selectedDiscoveryService.id}/${this.$route.params.subjectID}`, {registrationParameters: params})
          .then(data => {
            if (data.reason) {
              this.fetchError = data.reason
            } else {
              this.$router.push({name: 'admin.identityDetails', params: {subjectID: this.$route.params.subjectID}})
            }
          })
          .catch(response => {
            this.fetchError = response
          })
    }
  }
}
</script>
