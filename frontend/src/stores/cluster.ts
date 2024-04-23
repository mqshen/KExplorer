import { defineStore } from 'pinia'
import { SaveCluster, ListCluster } from 'wailsjs/go/services/ClusterService'

class Cluster {
    name: string
    bootstrap: string

    constructor(name: string, bootstrap: string) {
        this.name = name;
        this.bootstrap = bootstrap;
    }
}

const useClusterStore = defineStore('clusters', {
    state: () => ({
        clusters: [], // all connections
        clusterProfile: {}, // all server profile in flat list
    }),
    actions: {
        newDefaultCluster(name: string | null) {
            return {
                name: name || '',
                bootstrap: 'localhost:9092'
            }
        },
        async initClusters(force: boolean = false) {
            console.log(force)
            if (!force && this.clusters.length > 0) {
                return
            }
            const { success, data } = await ListCluster()
            if (success) {
                this.clusters = data.map(cluster => {
                    return {
                        key: '/' + cluster.name,
                        name: cluster.name,
                    }
                })
            }
        },
        async saveCluster(name: string, param: any) {
            const { success, msg } = await SaveCluster(name, param)
            if (!success) {
                return { success: false, msg }
            }

            // reload connection list
            await this.initClusters(true)
            return { success: true }
        },
    }

})

export default useClusterStore