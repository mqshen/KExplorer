import { defineStore } from 'pinia'
import { isEmpty, get } from 'lodash'
import { CloseConnection, OpenConnection, GetKafkaMetaData } from 'wailsjs/go/services/browserService'
import useTabStore from 'stores/tab'
import { KafkaServerState } from '@/objects/kafkaServerState'
import { NodeType } from "@/consts/kafka_node_type"

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
        /**
         * open database and load all keys
         * @param server
         * @param db
         * @returns {Promise<void>}
         */
        async getKafkaMetaData(server) {
            console.log("start get kafka meta data")
            const { data, success, msg } = await GetKafkaMetaData(server)
            if (!success) {
                throw new Error(msg)
            }
            const { brokers = [], topics = [], consumers = [] } = data

            /** @type {KafkaServerState} **/
            const serverInst = this.servers[server]
            if (serverInst == null) {
                return
            }
            serverInst.loadingState.fullLoaded = true

            if (isEmpty(brokers) && isEmpty(topics) && isEmpty(consumers)) {
                serverInst.nodeMap.clear()
            } else {
                // append db node to current connection's children
                serverInst.addNodes(brokers, "Brokers", NodeType.Broker)
                serverInst.addNodes(topics, "Topics", NodeType.Topic)
                serverInst.addNodes(consumers, "Consumers", NodeType.Consumer)
            }
            // serverInst.tidyNode('', false)
        },
        /**
         * get key struct in current database
         * @param {string} server
         * @param {boolean} [includeRoot]
         * @return {RedisNodeItem[]}
         */
        getKeyStruct(server, includeRoot) {
            /** @type {KafkaServerState} **/
            const serverInst = this.servers[server]
            let rootNode = null
            if (serverInst != null) {
                rootNode = serverInst.getRoot()
            }
            if (includeRoot === true) {
                return [rootNode]
            }
            let result = get(rootNode, 'children', [])
            return result;
        },
    },
    getters: {
        anyConnectionOpened() {
            return !isEmpty(this.servers)
        },
    },
})

export default useBrowserStore