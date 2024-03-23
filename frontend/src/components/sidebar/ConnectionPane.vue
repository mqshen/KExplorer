<script setup>
import { ref } from 'vue'
import ConnectionTree from './ConnectionTree.vue'

const filterPattern = ref('')
</script>

<template>
  <div class="nav-pane-container flex-box-v">
    <connection-tree :filter-pattern="filterPattern" />
  </div>
</template>
