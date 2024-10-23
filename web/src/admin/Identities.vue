<template>
  <div>

    <div class="flex justify-between mb-6">
      <h1>Identities</h1>

      <button
          id="new-identity-button"
          @click="$router.push({name: 'admin.newIdentity'})"
          class="float-right inline-flex items-center bg-nuts w-10 h-10 rounded-lg justify-center shadow-md"
      >
        <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#fff">
          <path d="M0 0h24v24H0V0z" fill="none"/>
          <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
        </svg>
      </button>
    </div>

    <ErrorMessage v-if="fetchError" :message="fetchError" :title="'Could not fetch identities'"/>

    <div class="mt-8 bg-white p-5 shadow-lg rounded-lg">
      <div class="m-4" v-if="identities.length === 0 && !fetchError">No identities yet, add one!</div>
      <table v-if="identities.length > 0" class="min-w-full divide-y divide-gray-200">
        <thead>
        <tr>
          <th class="thead">Name</th>
          <th class="thead">DIDs</th>
        </tr>
        </thead>
        <tbody class="tbody">
        <tr class="hover:bg-gray-100 cursor-pointer" v-for="{subject, dids} in identities" :key="subject"
            @click="$router.push({name: 'admin.identityDetails', params: {subjectID: subject}})">
          <td>{{ subject }}</td>
          <td>
            <div v-for="did in dids" :key="did">
              <pre>{{ did }}</pre>
            </div>
          </td>
        </tr>
        </tbody>
      </table>
      <router-view name="modal" @statusUpdate="updateStatus"></router-view>
    </div>
  </div>
</template>

<script>

import ErrorMessage from "../components/ErrorMessage.vue";

export default {
  components: {ErrorMessage},
  data() {
    return {
      fetchError: '',
      identities: [],
    }
  },
  mounted() {
    this.fetchData()
  },
  emits: ['statusUpdate'],
  methods: {
    updateStatus(event) {
      this.fetchData()
      this.$emit('statusUpdate', event)
    },
    fetchData() {
      this.$api.get('api/id')
          .then(data => {
            this.identities = data
          })
          .catch(response => {
            this.fetchError = response
          })
    }
  }
}
</script>
