
import { KafkaNodeItem } from '@/objects/kafkaNodeItem'
import { NodeType } from "@/consts/kafka_node_type"
/**
 * server connection state
 */
export class KafkaServerState {
    /**
     * @typedef {Object} LoadingState
     * @property {boolean} loading indicated that is loading children now
     * @property {boolean} fullLoaded indicated that all children already loaded
     */

    name: string
    stats: Object
    loadingState: Object
    nodeMap: Map<string, KafkaNodeItem>

    /**
     * @param {string} name server name
     * @param {number} db current opened database
     * @param {number} reloadKey try to reload when changed
     * @param {{}} stats current server status info
     * @param {string|null} patternFilter pattern filter
     * @param {string|null} typeFilter redis type filter
     * @param {boolean} exactFilter exact match filter keyword
     * @param {LoadingState} loadingState all loading state in opened connections map by server and LoadingState
     * @param {KeyViewType} viewType view type selection for all opened connections group by 'server'
     * @param {Map<string, KafkaNodeItem>} nodeMap map nodes by "type#key"
     */
    constructor(
        name,
        stats = {},
        loadingState = {},
        nodeMap = new Map(),
    ) {
        this.name = name
        this.stats = stats
        this.loadingState = loadingState
        this.nodeMap = nodeMap
    }

    addNodes(brokers, nodeKey, nodeType) {
        const result = {
            success: false,
            newLayer: 0,
            newKey: 0,
            replaceKey: 0,
        }
        const root = this.getRoot()
        let nodes = brokers.sort().map((broker) => new KafkaNodeItem({
            key: broker,
            label: broker,
            name: broker,
            type: nodeType,
            isLeaf: true,
        })
        )
        let node = new KafkaNodeItem({
            key: nodeKey,
            label: nodeKey,
            name: nodeKey,
            keyCount: brokers.length,
            isLeaf: false,
            children: nodes
        })


        root.addChild(node)
        return result
    }

    /**
     * get tree root item
     * @returns {RedisNodeItem}
     */
    getRoot() {
        const rootKey = "root"
        let root = this.nodeMap.get(rootKey)
        if (root == null) {
            // create root node
            root = new KafkaNodeItem({
                key: rootKey,
                label: root,
                children: [],
            })
            this.nodeMap.set(rootKey, root)
        }
        return root
    }

    dispose() {
        this.stats = {}
        this.nodeMap.clear()
    }

}
