import { NodeType } from "@/consts/kafka_node_type"
export class KafkaNodeItem {
    key: string
    label: string
    name: string
    type: NodeType
    keyCount: Number = 0
    isLeaf: Boolean = false
    opened: Boolean = false
    expanded: Boolean = false
    children: KafkaNodeItem[]

    constructor({
        key,
        label,
        name,
        type = NodeType.Topic,
        keyCount = 0,
        isLeaf = false,
        opened = false,
        expanded = false,
        children = []
    }: {
        key: string,
        label: string,
        name: string,
        type: NodeType,
        keyCount: Number,
        isLeaf: Boolean,
        opened: Boolean,
        expanded: Boolean,
        children: KafkaNodeItem[]
    }) {
        this.key = key
        this.label = label
        this.name = name
        this.type = type
        this.keyCount = keyCount
        this.isLeaf = isLeaf
        this.opened = opened
        this.expanded = expanded
        this.children = children
    }

    /**
 *
 * @param {RedisNodeItem} child
 * @param {boolean} [sorted]
 */
    addChild(child) {
        this.children.push(child)
    }

}