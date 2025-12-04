<template>
  <modal-window :cancelRoute="{name: 'admin.identityDetails', params: {subjectID: subjectID}}"
                title="Request Credential Issuance"
                type="add"
                :confirmFn="issueCredential"
                confirmText="Request">
    <ErrorMessage v-if="issueError" :message="issueError"/>
    <div>
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
        <label>Issuer</label>
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
  computed: {
    subjectID() {
      return this.$route.params.subjectID
    }
  },
  data() {
    return {
      selectedCredentialType: '',
      selectedWalletDID: '',
      issueError: undefined,
      credentialProfiles: [],
      walletDIDs: [],
    }
  },
  created() {
    // Fetch config for credential profiles
    this.$api.get('api/config')
        .then(data => {
          this.credentialProfiles = data.credential_profiles || []
        })
        .catch(response => {
          console.error('Failed to fetch config:', response)
          this.issueError = 'Failed to fetch credential profiles: ' + response
        })

    // Fetch identity details for wallet DIDs
    this.$api.get('api/id/' + this.subjectID)
        .then(data => {
          this.walletDIDs = data.did_documents ? data.did_documents.map(doc => doc.id) : []
          // Pre-select did:web DID after data is loaded
          this.preselectWalletDID()
        })
        .catch(response => {
          console.error('Failed to fetch identity:', response)
          this.issueError = 'Failed to fetch wallet DIDs: ' + response
        })
  },
  methods: {
    preselectWalletDID() {
      if (this.walletDIDs && this.walletDIDs.length > 0) {
        const didWebDID = this.walletDIDs.find(did => did.startsWith('did:web:'))
        if (didWebDID) {
          this.selectedWalletDID = didWebDID
        } else {
          this.selectedWalletDID = this.walletDIDs[0]
        }
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
            // Validate that the redirect_uri is from the same origin or a trusted external issuer
            try {
              const redirectUrl = new URL(data.redirect_uri)
              // Allow same-origin redirects or external HTTPS URLs (for OpenID4VCI flow)
              if (redirectUrl.origin === window.location.origin || redirectUrl.protocol === 'https:') {
                window.location.href = data.redirect_uri
              } else {
                this.issueError = 'Invalid redirect URL received from server'
              }
            } catch (e) {
              this.issueError = 'Invalid redirect URL format'
            }
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

