import { isEmpty } from 'lodash'
import { defineStore } from 'pinia'
import { SaveTopic, ListTopic } from 'wailsjs/go/services/topicService'

class Topic {
    keySerializer: number
    valueSerializer: number
    constructor(keySerializer: number, valueSerializer: number) {
        this.keySerializer = keySerializer;
        this.valueSerializer = valueSerializer;
    }
}

const useTopicStore = defineStore('topics', {
    state: () => ({
        topics: new Map<string, Topic>()
    }),

    actions: {
        async initTopics(force) {
            if (!force && !isEmpty(this.connections)) {
                return
            }
            const topics: Map<string, Topic> = new Map<string, Topic>()
            const { data = [{ server: '', topic: '', keySerializer: 0, valueSerializer: 0 }] } = await ListTopic()
            for (const topic of data) {
                topics.set(`${topic.server}_${topic.topic}`, new Topic(topic.keySerializer, topic.valueSerializer));
            }
            this.topics = topics;
        }, 
        async saveTopic(server, topic, param) {
            console.log(server, topic, param)
            const { success, msg } = await SaveTopic(server, topic, param)
            if (!success) {
                return { success: false, msg }
            }

            // reload connection list
            await this.initTopics(true)
            return { success: true }
        },
    }
})



export default useTopicStore