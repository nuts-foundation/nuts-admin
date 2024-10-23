<template>
  <modal-window :cancelRoute="{name: 'admin.identities'}" :confirmFn="checkForm" confirmText="Create Identity"
                title="Create a new identity" type="add">

    <p class="mb-3 text-sm">
      Here you can create new identities (subject) on the Nuts Node. Identities consist of one or more DIDs.
    </p>

    <ErrorMessage v-if="apiError" :message="apiError" title="Could not create identity"/>

    <div class="p-3 bg-red-100 rounded-md" v-if="formErrors.length">
      <b>Please correct the following error(s):</b>
      <ul>
        <li v-for="(error, idx) in formErrors" :key="`err-${idx}`">* {{ error }}</li>
      </ul>
    </div>

    <div class="mt-4">
      <identity-form mode="new" :value="identity" @input="(newIdentity)=> {identity = newIdentity}"/>
    </div>
  </modal-window>
</template>

<script>
import ModalWindow from '../components/ModalWindow.vue'
import IdentityForm from './IdentityForm.vue'
import ErrorMessage from "../components/ErrorMessage.vue";

export default {
  components: {
    ErrorMessage,
    ModalWindow,
    IdentityForm
  },
  data () {
    return {
      apiError: '',
      formErrors: [],
      identity: {
        subject: '',
      }
    }
  },
  emits: ['statusUpdate'],
  methods: {
    checkForm (e) {
      // reset the errors
      this.formErrors.length = 0
      this.apiError = ''
      return this.confirm()
    },
    confirm () {
      this.$api.post('api/id', this.identity)
        .then(response => {
          this.$emit('statusUpdate', 'Identity created')
          this.$router.push({ name: 'admin.identities' })
        })
        .catch(response => {
          this.apiError = response
        })
    }
  }
}
</script>
