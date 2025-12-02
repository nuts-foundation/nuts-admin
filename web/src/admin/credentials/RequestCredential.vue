<template>
  <modal-window :cancelRoute="{name: 'admin.identityDetails', params: {subjectID: subjectID}}"
                title="Request Credential Issuance"
                type="add"
                :confirmFn="handleConfirm"
                :confirmText="issuanceResult ? 'Close' : 'Issue'">
    <ErrorMessage v-if="issueError" :message="issueError"/>
    <div v-if="issuanceResult">
      <p>Credential issuance initiated successfully!</p>
      <p>
        <a :href="issuanceResult.redirect_uri" target="_blank" class="text-blue-600 hover:text-blue-800 underline">
          Click here to start the credential issuance
        </a>
      </p>
    </div>
    <div v-else>
     <div>
        <label for="credential-type">Credential Type</label>
        <select v-model="selectedCredentialType" id="credential-type">
          <option value="">Select a credential type</option>
          <option v-for="profile in credentialProfiles" :key="profile.type" :value="profile.type">
            {{ profile.type }}
          </option>
        </select>
      </div>
      <div>
        <label for="wallet-did">Wallet DID</label>
        <select v-model="selectedWalletDID" id="wallet-did">
          <option value="">Select a wallet DID</option>
          <option v-for="did in walletDIDs" :key="did" :value="did">
            {{ did }}
          </option>
        </select>
      </div>
      <div v-if="selectedCredentialType">
        <label for="issuer-did">Issuer</label>
        <p>{{ getIssuerForType(selectedCredentialType) }}</p>
      </div>
    </div>
  </modal-window>
</template>

<script>
import ErrorMessage from "../../components/ErrorMessage.vue";
import ModalWindow from "../../components/ModalWindow.vue";
import {encodeURIPath} from "../../lib/encode";

export default {
  components: {ErrorMessage, ModalWindow},
  props: {
    subjectID: {
      type: String,
      required: true
    },
    credentialProfiles: {
      type: Array,
      required: true
    },
    walletDIDs: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      selectedCredentialType: '',
      selectedWalletDID: '',
      issueError: undefined,
      issuanceResult: undefined,
    }
  },
  created() {
    // Pre-select did:web DID if available
    const didWebDID = this.walletDIDs.find(did => did.startsWith('did:web:'))
    if (didWebDID) {
      this.selectedWalletDID = didWebDID
    } else if (this.walletDIDs.length > 0) {
      this.selectedWalletDID = this.walletDIDs[0]
    }
  },
  methods: {
    handleConfirm() {
      if (this.issuanceResult) {
        // If we already have a result, navigate back to identity details
        this.$router.push({name: 'admin.identityDetails', params: {subjectID: this.subjectID}})
      } else {
        // Otherwise, issue the credential
        this.issueCredential()
      }
    },
    getIssuerForType(credentialType) {
      const profile = this.credentialProfiles.find(p => p.type === credentialType)
      return profile ? profile.issuer : ''
    },
    issueCredential() {
      this.issueError = undefined
      const issuerDID = this.getIssuerForType(this.selectedCredentialType)
      if (!issuerDID) {
        this.issueError = 'No issuer found for selected credential type'
        return
      }

      if (!this.selectedWalletDID) {
        this.issueError = 'Please select a wallet DID'
        return
      }

      const redirectUri = `${window.location.origin}${this.$router.resolve({name: 'admin.identityDetails', params: {subjectID: this.subjectID}}).href}`

      const requestBody = {
        credential_type: this.selectedCredentialType,
        issuer: issuerDID,
        wallet_did: this.selectedWalletDID,
        redirect_uri: redirectUri,
      }

      this.$api.post(`api/proxy/internal/auth/v2/${encodeURIPath(this.subjectID)}/request-credential`, requestBody)
          .then(data => {
            this.issuanceResult = data
            this.$emit('statusUpdate', 'Credential issuance initiated successfully')
          })
          .catch(response => {
            this.issueError = "Failed to initiate credential issuance: " + response
          })
    }
  }
}
</script>

<style scoped>
/* Make the modal wider to accommodate long DIDs */
:deep(.inline-block) {
  max-width: 48rem !important; /* 768px instead of default 32rem (512px) */
}

:deep(select) {
  width: 100%;
}
</style>

