import { defineStore } from 'pinia'
import { SaveConnection, ListConnection, DeleteConnection, GetConnection } from 'wailsjs/go/services/connectionService'
import { isEmpty, get, uniq } from 'lodash'
import { ConnectionType } from '@/consts/connection_type'
import useBrowserStore from 'stores/browser'

class Connection {
    key: string
    label: string
    name: string
    type: number

    children: Connection[]
    constructor(key: string, label: string, name: string,  type: number) {
        this.key = key;
        this.label = label;
        this.name = name;
        this.type = type;
      }
}

const useConnectionStore = defineStore('connections', {
    state: () => ({
        groups: [], // all group name set
        connections: [], // all connections
        serverProfile: {}, // all server profile in flat list
    }),

    actions: {
        /**
         * load all store connections struct from local profile
         * @param {boolean} [force]
         * @returns {Promise<void>}
         */
        async initConnections(force) {
            if (!force && !isEmpty(this.connections)) {
                return
            }
            const conns: Connection[] = []
            const groups: string[] = []
            const profiles = {}
            const { data = [{ groupName: '', connections: [], refreshInterval: 5 }] } = await ListConnection()
            for (const conn of data) {
                if (conn.type !== 'group') {
                    // top level
                    conns.push(new Connection(
                        '/' + conn.name,
                        conn.name,
                        conn.name,
                        ConnectionType.Server
                        // cluster: get(conn, 'cluster.enable', false),
                        // isLeaf: false,
                    ))
                    profiles[conn.name] = {
                        defaultFilter: conn.defaultFilter,
                        keySeparator: conn.keySeparator,
                        markColor: conn.markColor,
                        refreshInterval: conn.refreshInterval,
                    }
                } else {
                    // custom group
                    groups.push(conn.name)
                    const subConns = get(conn, 'connections', [])
                    const children: Connection[] = []
                    for (const item of subConns) {
                        const value = conn.name + '/' + item.name
                        children.push(new Connection(
                            value,
                            item.name,
                            item.name,
                            ConnectionType.Server
                            // cluster: get(item, 'cluster.enable', false),
                            // isLeaf: false,
                        ))
                        profiles[item.name] = {
                            defaultFilter: item.defaultFilter,
                            keySeparator: item.keySeparator,
                            markColor: item.markColor,
                            refreshInterval: item.refreshInterval,
                        }
                    }
                    conns.push({
                        key: conn.name + '/',
                        label: conn.name,
                        name: '',
                        type: ConnectionType.Group,
                        children,
                    })
                }
            }
            this.connections = conns
            this.serverProfile = profiles
            this.groups = uniq(groups)
        },

        /**
         * create a new default connection
         * @param {string} [name]
         * @returns {{}}
         */
        newDefaultConnection(name) {
            return {
                group: '',
                name: name || '',
                addr: '127.0.0.1',
                port: 9092,
            }
        },
        /**
         * create a new connection or update current connection profile
         * @param {string} name set null if create a new connection
         * @param {{}} param
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async saveConnection(name, param) {
            const { success, msg } = await SaveConnection(name, param)
            if (!success) {
                return { success: false, msg }
            }

            // reload connection list
            await this.initConnections(true)
            return { success: true }
        },

        /**
         * remove connection
         * @param name
         * @returns {Promise<{success: boolean, [msg]: string}>}
         */
        async deleteConnection(name) {
            // close connection first
            const browser = useBrowserStore()
            await browser.closeConnection(name)
            const { success, msg } = await DeleteConnection(name)
            if (!success) {
                return { success: false, msg }
            }
            await this.initConnections(true)
            return { success: true }
        },
        /**
         * get connection by name from local profile
         * @param name
         * @returns {Promise<ConnectionProfile|null>}
         */
        async getConnectionProfile(name) {
            try {
                const { data, success } = await GetConnection(name)
                if (success) {
                    this.serverProfile[name] = {
                        defaultFilter: data.defaultFilter,
                        keySeparator: data.keySeparator,
                        markColor: data.markColor,
                    }
                    return data
                }
            } finally {
            }
            return null
        },

        mergeConnectionProfile(dest, src) {
            const mergeObj = (destObj, srcObj) => {
                for (const k in srcObj) {
                    const t = typeof srcObj[k]
                    if (t === 'string') {
                        destObj[k] = srcObj[k] || destObj[k] || ''
                    } else if (t === 'number') {
                        destObj[k] = srcObj[k] || destObj[k] || 0
                    } else if (t === 'object') {
                        mergeObj(destObj[k], srcObj[k] || {})
                    } else {
                        destObj[k] = srcObj[k]
                    }
                }
                return destObj
            }
            return mergeObj(dest, src)
        },
    }
})

export default useConnectionStore