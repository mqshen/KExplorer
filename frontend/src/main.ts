import { createPinia } from 'pinia'
import { createApp, nextTick } from "vue"
import './styles/style.scss'
import App from "./App.vue"
import { i18n } from '@/utils/i18n.ts'
import { setupDiscreteApi } from '@/utils/discrete.ts'
import { loadEnvironment } from '@/utils/platform.ts'
import usePreferencesStore from 'stores/preferences.ts'

async function setupApp() {
  const app = createApp(App)
  app.use(i18n)
  app.use(createPinia())

  await loadEnvironment()


  const prefStore = usePreferencesStore()
  await prefStore.loadPreferences()
  await setupDiscreteApi()
  app.config.errorHandler = (err, instance, info) => {
      // TODO: add "send error message to author" later
      nextTick().then(() => {
          try {
              const content = err.toString()
              $notification.error(content, {
                  title: i18n.global.t('common.error'),
                  meta: 'Please see console output for more detail',
              })
              console.error(err)
          } catch (e) {}
      })
  }
  // app.config.warnHandler = (message) => {
  //     console.warn(message)
  // }
  app.mount('#app')

}

setupApp()
