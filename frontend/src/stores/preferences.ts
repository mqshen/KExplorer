import { defineStore } from 'pinia'
import {  pick } from 'lodash'

const usePreferencesStore = defineStore('preferences', {
    state: () => ({
        behavior: {
            welcomed: false,
            asideWidth: 300,
            windowWidth: 0,
            windowHeight: 0,
            windowMaximised: false,
        },
    }),
    getters: {
    },
    actions: {
        /**
         * load preferences from local
         * @returns {Promise<void>}
         */
        async loadPreferences() {
        },

        /**
         * save preferences to local
         * @returns {Promise<boolean>}
         */
        async savePreferences() {
            // const pf = pick(this, ['behavior', 'general', 'editor', 'cli', 'decoder'])
            // const { success, msg } = await SetPreferences(pf)
            // return success === true
        },
    }

})


export default usePreferencesStore