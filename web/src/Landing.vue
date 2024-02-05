<template>
  <div class="">
    <h1 class="page-title">Welcome to Nuts Admin</h1>
    <p class="mt-4 leading-normal">
      This application lets you administer <router-link to="identities">identities</router-link> (e.g. organizations) in the form of DID documents,
      and Verifiable Credentials (VCs) for Nuts use cases.
    </p>
    <h2 class="mt-4 sub-title">Getting Started</h2>
    <p class="mt-4 leading-normal">
      For testing purposes or local development, you can start by <router-link to="identities">registering a new DID</router-link>,
      and then issue a <code>NutsOrganizationCredential</code> to it.
    </p>
    <p class="mt-4 leading-normal">
      If you're implementing a use case that provides a Discovery Service (given you have configured it on your Nuts node),
      you can activate the Discovery Service for the DID (identity) created by inspecting the identity.
      Registration on a Discovery Service requires service-specific Verifiable Credentials,
      which are either self-issued (e.g. using this application) or acquired from a trusted third-party issuer.
    </p>
    <p class="mt-4 leading-normal">
      You can also <router-link to="discovery">search through the Discovery Services</router-link> that are available on your Nuts node.
    </p>
  </div>
</template>

<script>
export default {
}
</script>
