<template>
  <div
      class="appearance-none mb-2 mt-2 px-6 py-3.5 w-full text-sm text-red-500 bg-red-100 placeholder-red-500 outline-none border border-red-500 focus:ring-4 focus:ring-blue-200 rounded-md"
      role="alert">
    <div class="font-semibold mb-2">{{ title}}</div>
    <div>{{ message }}</div>
    <div v-if="showReauthenticateButton" class="mt-3">
      <button type="button"
              class="btn btn-danger"
              @click="reauthenticate"
      >
        {{ reauthenticateMessage }}
      </button>
    </div>
  </div>
</template>
<script>
export default {
  name: 'ErrorMessage',
  props: {
    message: String,
    title: {
      type: String,
      default: 'Error'
    },
    reauthenticateMessage: {
      type: String,
      default: 'Reauthenticate'
    },
    reauthenticateUrl: {
      type: String,
      default: './auth/openid-connect'
    }
  },
  computed: {
    showReauthenticateButton() {
      return !!this.message && this.message.toLowerCase().includes('unauthorized');
    }
  },
  methods: {
    reauthenticate() {
      window.location.href = this.reauthenticateUrl;
    }
  }
}
</script>