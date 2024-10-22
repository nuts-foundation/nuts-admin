<template>
  <div>
    <h1>Verifiable Credential</h1>
    <p v-if="fetchError" class="m-4">Error: {{ fetchError }}</p>
    <div v-if="credential">
      <section>
        <div>
          <label for="credentialID">ID</label>
          <input id="credentialID" v-model="credential.credentialSubject.id" type="text" readonly>
        </div>
        <div>
          <label for="credentialType">Type</label>
          <input id="credentialType" v-model="credentialType" type="text" readonly>
        </div>
        <div>
          <label for="credentialIssuer">Issuer</label>
          <input id="credentialIssuer" v-model="credential.issuer" type="text" readonly>
        </div>
      </section>
      <section>
        <header>Credential Subject</header>
        <table class="w-full">
          <thead>
          <tr>
            <th class="w-1/3">Key</th>
            <th>Value</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="(param, idx) in credentialSubjectProperties" :key="'credentialSubjectProperty-' + idx">
            <td>{{ param.key }}</td>
            <td>{{ param.value }}</td>
          </tr>
          </tbody>
        </table>
      </section>
    </div>

  </div>
</template>

<script>

export default {
  data() {
    return {
      fetchError: undefined,
      credential: undefined
    }
  },
  computed: {
    credentialType() {
      return this.credential.type.filter(t => t !== "VerifiableCredential").join(', ')
    },
    credentialSubjectProperties() {
      // recursively flatten credential.credentialSubject, concatenating keys with parent keys (with dots inbetween)
      const flatten = (obj, parentKey = '') => {
        return Object.keys(obj).reduce((acc, key) => {
          const newKey = parentKey ? `${parentKey}.${key}` : key
          if (typeof obj[key] === 'object') {
            return acc.concat(flatten(obj[key], newKey))
          } else {
            return acc.concat({key: newKey, value: obj[key]})
          }
        }, [])
      }
      return flatten(this.credential.credentialSubject)
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
