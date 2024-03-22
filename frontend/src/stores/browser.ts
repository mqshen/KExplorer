import { defineStore } from 'pinia'
import { isEmpty } from 'lodash'

const useBrowserStore = defineStore('browser', {
    state: () => ({
        servers: {},
    }),
    getters: {
        anyConnectionOpened() {
            return !isEmpty(this.servers)
        },
    },
})

export default useBrowserStore