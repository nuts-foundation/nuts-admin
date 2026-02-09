<template>
  <div>
    <h1>Verifiable Credential</h1>
      <section>
        <div>
          <label>ID</label>
          <div>{{credentialSubject.id}}</div>
        </div>
        <div>
          <label>Type</label>
          <div>{{credentialType}}</div>
        </div>
        <div>
          <label>Issuer</label>
          <div>{{credential.issuer}}</div>
        </div>
        <div>
          <label>Issuance date</label>
          <div>{{credential.issuanceDate}}</div>
        </div>
        <div>
          <label>Expiration date</label>
          <div>{{credential.expirationDate}}</div>
        </div>
        <div v-if="credential.status">
          <label>Status</label>
          <div>
            <span :class="statusClass">{{credential.status}}</span>
          </div>
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

      <!-- Button bar for actions -->
      <div v-if="showRevokeButton && credential.status === 'active'" class="mt-4 pt-4 border-t border-gray-200">
        <button @click="revokeCredential"
                class="bg-red-600 hover:bg-red-700 text-white font-semibold py-2 px-4 rounded disabled:opacity-50 disabled:cursor-not-allowed"
                :disabled="revoking">
          {{ revoking ? 'Revoking...' : 'Revoke Credential' }}
        </button>
        <p v-if="revokeError" class="text-red-600 mt-2">{{ revokeError }}</p>
        <p v-if="revokeSuccess" class="text-green-600 mt-2 font-semibold">Credential revoked successfully</p>
      </div>
    </div>
</template>
<script>
import { parseApiError } from '../../lib/errors.js'

export default {
  props: {
    credential: Object,
    showRevokeButton: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      revoking: false,
      revokeError: '',
      revokeSuccess: false
    }
  },
  computed: {
    credentialSubject() {
      if (Array.isArray(this.credential.credentialSubject) && this.credential.credentialSubject.length === 1) {
        return this.credential.credentialSubject[0]
      }
      return this.credential.credentialSubject
    },
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
      if (Array.isArray(this.credential.credentialSubject) && this.credential.credentialSubject.length === 1) {
        return flatten(this.credential.credentialSubject[0])
      }
      return flatten(this.credential.credentialSubject)
    },
    statusClass() {
      switch(this.credential.status) {
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
  },
  methods: {
    async revokeCredential() {
      if (!confirm('Are you sure you want to revoke this credential? This action cannot be undone.')) {
        return
      }

      this.revoking = true
      this.revokeError = ''
      this.revokeSuccess = false

      try {
        const credentialId = this.credential.id
        if (!credentialId) {
          throw new Error('Credential ID is missing')
        }

        // Call the proxy endpoint to revoke the credential
        await this.$api.delete('api/proxy/internal/vcr/v2/issuer/vc/' + encodeURIComponent(credentialId))

        this.revokeSuccess = true
        // Update the credential status locally
        this.credential.status = 'revoked'

        // Emit event so parent component can refresh if needed
        this.$emit('revoked')
      } catch (error) {
        this.revokeError = parseApiError(error)
      } finally {
        this.revoking = false
      }
    }
  }
}
</script>
