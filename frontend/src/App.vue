<script setup>
import ConnectionDialog from './components/dialogs/ConnectionDialog.vue'
import { h, onMounted, ref, watch } from 'vue'
import usePreferencesStore from './stores/preferences.ts'
import { darkThemeOverrides, themeOverrides } from '@/utils/theme.ts'

import AppContent from "./AppContent.vue"

const initializing = ref(true)

const prefStore = usePreferencesStore()

onMounted(async () => {
    try {
        initializing.value = true
    } finally {
        initializing.value = false
    }
})
</script>
<template>
  <n-config-provider
    :inline-theme-disabled="true"
    :locale="prefStore.themeLocale"
    :theme="prefStore.isDark ? darkTheme : undefined"
    :theme-overrides="prefStore.isDark ? darkThemeOverrides : themeOverrides"
    class="fill-height"
  >
    <n-dialog-provider>
      <app-content :loading="initializing" />

      <!-- top modal dialogs -->
      <connection-dialog />
    </n-dialog-provider>
  </n-config-provider>
</template>
