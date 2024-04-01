<script setup>
const props = defineProps({
  strokeWidth: {
    type: [Number, String],
    default: 3,
  },
});
</script>

<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    viewBox="0 0 32 32"
  >
    <path
      d="M28 28H4a2.002 2.002 0 0 1-2-2v-5h2v5h24v-5h2v5a2.002 2.002 0 0 1-2 2z"
      fill="currentColor"
    ></path>
    <path d="M7 21h18v2H7z" fill="currentColor"></path>
    <path d="M7 16h18v2H7z" fill="currentColor"></path>
    <path d="M7 11h18v2H7z" fill="currentColor"></path>
    <path d="M7 6h18v2H7z" fill="currentColor"></path>
  </svg>
</template>

<style lang="scss" scoped></style>
