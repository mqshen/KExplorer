/**
 * server connection state
 */
export class KafkaServerState {
    /**
     * @typedef {Object} LoadingState
     * @property {boolean} loading indicated that is loading children now
     * @property {boolean} fullLoaded indicated that all children already loaded
     */

    /**
     * @param {string} name server name
     * @param {number} db current opened database
     * @param {number} reloadKey try to reload when changed
     * @param {{}} stats current server status info
     * @param {Object.<number, RedisDatabaseItem>} databases database list
     * @param {string|null} patternFilter pattern filter
     * @param {string|null} typeFilter redis type filter
     * @param {boolean} exactFilter exact match filter keyword
     * @param {LoadingState} loadingState all loading state in opened connections map by server and LoadingState
     * @param {KeyViewType} viewType view type selection for all opened connections group by 'server'
     * @param {Map<string, RedisNodeItem>} nodeMap map nodes by "type#key"
     */
    constructor({
        name,
        stats = {},
        loadingState = {},
    }) {
        this.name = name
        this.stats = stats
        this.loadingState = loadingState
    }

    dispose() {
        this.stats = {}
    }

}
