<template>
  <div>
    <h1>Issue Verifiable Credential</h1>
    <p v-if="fetchError" class="m-4">Error: {{ fetchError }}</p>

    <section v-if="issuedCredential">
      <header>Issued Credential</header>
      <pre>{{ JSON.stringify(issuedCredential, null, 2) }}</pre>
    </section>

    <div v-else>
      <section>
        <header>Credential properties</header>
        <div>
          <label for="type">Type</label>
          <select v-on:change="e => selectCredentialType(e.target.value)" id="type">
            <option :value="currentType" v-for="currentType in Object.keys(templates)" :key="currentType"
                    :selected="currentType === credentialType">
              {{ currentType }}
            </option>
          </select>
        </div>

        <div>
          <label for="issuerDID">Issuer DID</label>
          <select v-on:change="selectIssuerDID" id="issuerDID">
            <option disabled selected value>choose issuer DID</option>
            <optgroup v-for="entry in subjects" :key="'subject-' + entry.subject" :label="entry.subject">
              <option :value="currentDID" v-for="currentDID in entry.dids" :key="'did-' + currentDID">
                {{ currentDID }}
              </option>
            </optgroup>
          </select>
        </div>

        <div>
          <label for="subjectDID">Wallet DID</label>
          <select v-on:change="selectSubjectDID" class="inline" style="width: 20%">
            <option disabled value="" selected>choose wallet DID</option>
            <optgroup v-for="entry in subjects" :key="'subject-' + entry.subject" :label="entry.subject">
              <option :value="entry.subject + '/' + currentDID" v-for="currentDID in entry.dids"
                      :key="'did-' + currentDID">
                {{ currentDID }}
              </option>
            </optgroup>
          </select>
          <input id="subjectDID" type="text" v-model="subjectDID" placeholder="Enter a DID" class="inline" style="width: 80%">
          <p>The credential will be loaded into the wallet if it's owned by the local Nuts node.</p>
        </div>

        <div v-if="template">
          <div>
            <label for="daysValid">Days valid</label>
            <input id="daysValid" v-model="daysValid" type="number">
            <p>This will be used to set the credentials <code>expirationDate</code> property.</p>
          </div>
          <div v-for="(field, idx) in template.fields" :key="field.name">
            <label :for="field.name">
              {{ field.name }}
              <p>{{ field.description }}</p>
            </label>
            <input :id="field.name" v-model="credentialFields[idx]" type="text">
          </div>
        </div>
      </section>

      <section>
        <input type="button" class="btn btn-primary" value="Preview Credential" @click="previewCredential">
        &nbsp;
        <input type="button" class="btn btn-primary" value="Issue Credential" @click="issueCredential">
      </section>

      <section v-if="credentialPreview">
        <header>Credential Preview</header>
        <pre>{{ credentialPreview }}</pre>
      </section>
    </div>
  </div>
</template>

<script>
import templates from "./templates";

export default {
  data() {
    return {
      fetchError: undefined,
      credentialType: undefined,
      subjectDID: undefined,
      holderSubjectID: undefined,
      issuerDID: undefined,
      subjects: [],
      templates: templates,
      template: undefined,
      credentialFields: [],
      daysValid: 365,
      credentialPreview: undefined,
      issuedCredential: undefined,
    }
  },
  mounted() {
    if (this.$route.params.credentialType) {
      this.selectCredentialType(this.$route.params.credentialType)
    } else if (Object.keys(this.templates).length > 0) {
      this.selectCredentialType(Object.keys(this.templates)[0])
    }
    this.subjectDID = this.$route.params.subjectDID
    this.fetchData()
  },
  methods: {
    updateStatus(event) {
      this.$emit('statusUpdate', event)
    },
    selectIssuerDID(event) {
      this.issuerDID = event.target.value
    },
    selectSubjectDID(event) {
      // subject is in form of subjectID/did, need to parse it to set both
      const parts = event.target.value.split('/')
      this.holderSubjectID = parts[0]
      this.subjectDID = parts[1]
      event.target.value = ""
    },
    selectCredentialType(type) {
      this.template = this.templates[type]
      this.credentialFields = []
    },
    previewCredential() {
      this.credentialPreview = JSON.stringify(this.template.render(this.issuerDID, this.subjectDID, this.credentialFields), null, 2)
    },
    issueCredential() {
      if (!this.issuerDID) {
        this.fetchError = 'Issuer DID is required'
        return
      }
      if (!this.subjectDID) {
        this.fetchError = 'Subject DID is required'
        return
      }
      let credentialToIssue = this.template.render(this.issuerDID, this.subjectDID, this.credentialFields)
      credentialToIssue['@context'] = credentialToIssue['@context'][0]
      credentialToIssue['type'] = credentialToIssue['type'].find(t => t !== "VerifiableCredential")
      credentialToIssue['expirationDate'] = new Date(new Date().getTime() + 1000 * 60 * 60 * 24 * this.daysValid).toISOString()
      // disable statusListRevocation2021 for now, causes issues on MS SQL Server
      //credentialToIssue.withStatusList2021Revocation = true
      this.fetchError = undefined
      this.$api.post('api/proxy/internal/vcr/v2/issuer/vc', credentialToIssue)
          .then(issuedCredential => {
            // Load issued VC into wallet
            this.issuedCredential = issuedCredential
            // If it's a local wallet, load it into the wallet
            if (this.holderSubjectID) {
              this.$api.post(`api/proxy/internal/vcr/v2/holder/${this.holderSubjectID}/vc`, issuedCredential)
                  .then(() => {
                    this.$emit('statusUpdate', 'Verifiable Credential issued, and loaded into wallet')
                  })
                  .catch(reason => {
                    this.fetchError = "Couldn't load credential into wallet: " + reason
                  })
            } else {
              this.$emit('statusUpdate', 'Verifiable Credential issued, NOTE: make sure to copy it to provide it to the holder!')
            }
          })
          .catch(reason => {
            this.fetchError = "Couldn't issue credential: " + reason
          })

    },
    fetchData() {
      this.$api.get('api/id')
          .then(data => {
            this.subjects = data
          })
          .catch(response => {
            this.fetchError = response
          })
    }
  }
}
</script>
