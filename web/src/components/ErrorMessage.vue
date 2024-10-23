<template>
  <div
      class="appearance-none mb-2 mt-2 px-6 py-3.5 w-full text-sm text-red-500 bg-red-100 placeholder-red-500 outline-none border border-red-500 focus:ring-4 focus:ring-blue-200 rounded-md"
      role="alert">
    <div class="font-semibold mb-2">{{ title}}</div>
    {{ message }}
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
    }
  }
}
</script>