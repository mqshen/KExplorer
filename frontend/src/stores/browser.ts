import { defineStore } from 'pinia'
import { isEmpty } from 'lodash'
import { CloseConnection, OpenConnection } from 'wailsjs/go/services/browserService'
import useTabStore from 'stores/tab'
import { KafkaServerState } from '@/objects/kafkaServerState'

const useBrowserStore = defineStore('browser', {
    state: () => ({
        servers: {},
    }),
    actions: {
        /**
             * check if connection is connected
             * @param name
             * @returns {boolean}
             */
        isConnected(name) {
            return this.servers.hasOwnProperty(name)
        },
        /**
         * open connection
         * @param {string} name
         * @param {boolean} [reload]
         * @returns {Promise<void>}
         */
        async openConnection(name, reload) {
            if (this.isConnected(name)) {
                if (reload !== true) {
                    return
                } else {
                    // reload mode, try close connection first
                    await CloseConnection(name)
                }
            }

            const { data, success, msg } = await OpenConnection(name)
            if (!success) {
                throw new Error(msg)
            }
            const serverInst = new KafkaServerState({
                name,
            })
            this.servers[name] = serverInst
        },
        /**
         * close connection
         * @param {string} name
         * @returns {Promise<boolean>}
         */
        async closeConnection(name) {
            const { success, msg } = await CloseConnection(name)
            if (!success) {
                // throw new Error(msg)
                return false
            }
            delete this.servers[name]

            const tabStore = useTabStore()
            tabStore.removeTabByName(name)
            return true
        },
    },
    getters: {
        anyConnectionOpened() {
            return !isEmpty(this.servers)
        },
    },
})

export default useBrowserStore