<template>
  <modal-window :cancelRoute="{name: 'admin.identityDetails', params: {subjectID: this.$route.params.subjectID}}" :confirmFn="confirm" confirmText="Upload Credential"
                title="Upload a Credential JSON" type="add">

    <p class="mb-3 text-sm">
      Here you can add a credential issued by another party to be stored in your wallet. Paste the JSON object or JWT below.
    </p>

    <ErrorMessage v-if="apiError" :message="apiError" title="Could not upload credential"/>

    <div class="mt-4">
      <upload-credential-form mode="new" :value="credential" @input="(credentialPaste)=> {credential = credentialPaste}"/>
    </div>
  </modal-window>
</template>

<script>
import ModalWindow from '../../components/ModalWindow.vue'
import UploadCredentialForm from "./UploadCredentialForm.vue";
import ErrorMessage from "../../components/ErrorMessage.vue";

export default {
  components: {
    ErrorMessage,
    ModalWindow,
    UploadCredentialForm
  },
  created() {
    this.subjectID = this.$route.params.subjectID;
  },
  data () {
    return {
      apiError: '',
      subjectID: '',
      credential: undefined,
    }
  },
  methods: {
    confirm () {
      this.$api.post(`api/proxy/internal/vcr/v2/holder/${this.subjectID}/vc`, this.credential)
          .then(() => {
            this.$emit('uploadCredential', 'Verifiable Credential issued, and loaded into wallet')
            this.$router.push({name: 'admin.identityDetails', params: {subjectID: this.subjectID}})
          })
          .catch(reason => {
            this.apiError = reason
          })
    }
  }
}
</script>
