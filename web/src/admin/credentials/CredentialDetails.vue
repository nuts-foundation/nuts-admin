<template>
  <div>
    <h1>Verifiable Credential</h1>
      <section>
        <div>
          <label>ID</label>
          <div>{{credential.credentialSubject.id}}</div>
        </div>
        <div>
          <label>Type</label>
          <div>{{credentialType}}</div>
        </div>
        <div>
          <label>Issuer</label>
          <div>{{credential.issuer}}</div>
        </div>
      </section>
      <section>
        <header>Credential Subject</header>
        <table class="table w-full divide-y divide-gray-200">
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

</template>
<script>
export default {
  props: {
    credential: Object
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
  }
}
</script>
