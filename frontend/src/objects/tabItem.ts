/**
 * tab item
 */
export class TopicConfig {
    server: string;
    topic: string;
    keySerializer: number;
    valueSerializer: number;

    static createFrom(source: any = {}) {
        return new TopicConfig(source);
    }

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.server = source["server"];
        this.topic = source["topic"];
        this.keySerializer = source["keySerializer"];
        this.valueSerializer = source["valueSerializer"];
    }
}

export class TabItem {
    name: string
    title: string
    end: boolean
    loading: boolean
    blank: boolean = false
    expandedKeys = new Array<string>()
    selectedKeys = new Array<string>()
    currentNode: TopicConfig = new TopicConfig()
    subTab: string = "properties"

    constructor(
        name: string,
        title: string,
        end = false,
        loading = false,
    ) {
        this.name = name
        this.title = title
        this.end = end
        this.loading = loading
    }
}
