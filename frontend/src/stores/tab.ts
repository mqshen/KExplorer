import { findIndex, isEmpty, set, get, find } from 'lodash'
import { defineStore } from 'pinia'
import { TabItem } from '@/objects/tabItem'
import useTopicStore from "stores/topics";

const useTabStore = defineStore('tab', {
    state: () => ({
        nav: 'server',
        asideWidth: 300,
        tabList: [],
        activatedIndex: 0,
    }),
    getters: {
        /**
         * get current tab list item
         * @returns {TabItem[]}
         */
        tabs() {
            // if (isEmpty(this.tabList)) {
            //     this.newBlankTab()
            // }
            return this.tabList
        },

        /**
         * get current activated tab item
         * @returns {TabItem|null}
         */
        currentTab() {
            return get(this.tabs, this.activatedIndex)
        },
        currentTabName() {
            return get(this.tabs, [this.activatedIndex, 'name'])
        },
    },
    actions: {
        /**
         *
         * @param idx
         * @param {boolean} [switchNav]
         * @param {string} [subTab]
         * @private
         */
        _setActivatedIndex(idx, switchNav, subTab) {
            console.log(idx, switchNav, subTab)
            this.activatedIndex = idx
            if (switchNav === true) {
                this.nav = idx >= 0 ? 'browser' : 'server'
                if (!isEmpty(subTab)) {
                    set(this.tabList, [idx, 'subTab'], subTab)
                }
            } else {
                if (idx < 0) {
                    this.nav = 'server'
                }
            }
        },
        /**
         * update or insert a new tab if not exists with the same name
         * @param {string} server
         */
        upsertTab({
            server,
        }) {
            let tabIndex = findIndex(this.tabList, { name: server })
            if (tabIndex === -1) {
                const tabItem = new TabItem(server, server)
                this.tabList.push(tabItem)
                tabIndex = this.tabList.length - 1
            } else {
                const tab = this.tabList[tabIndex]
                tab.blank = false
                // tab.title = db !== undefined ? `${server}/db${db}` : `${server}`
                tab.title = server
                tab.server = server
            }
            this._setActivatedIndex(tabIndex, true)
        },/**
        * set expanded keys for server
        * @param {string} server
        * @param {string[]} keys
        */
        setExpandedKeys(server, keys = []) {
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                if (isEmpty(keys)) {
                    tab.expandedKeys = []
                } else {
                    tab.expandedKeys = keys
                }
            }
        },/**
        * set selected keys for server
        * @param {string} server
        * @param {string|string[]} [keys]
        */
        setSelectedKeys(server, keys = null, node) {
            console.log(server, keys)
            /** @type TabItem**/
            let tab = find(this.tabList, { name: server })
            if (tab != null) {
                if (keys == null) {
                    // select nothing
                    tab.selectedKeys = []
                } else if (typeof keys === 'string') {
                    tab.selectedKeys = [keys]
                } else {
                    tab.selectedKeys = keys
                }
                const topic = node.label
                tab.currentNode.topic = topic
                const topicStore = useTopicStore()
                console.log(topicStore.topics)
                let topicConfig = topicStore.topics.get(`${server}_${topic}`)
                if (topicConfig != null) {
                    tab.currentNode.keySerializer = topicConfig.keySerializer
                    tab.currentNode.valueSerializer = topicConfig.valueSerializer
                } else {
                    tab.currentNode.keySerializer = 0
                    tab.currentNode.valueSerializer = 0
                }
            }
        },

        switchTab(tabIndex) {
            // const len = size(this.tabList)
            // if (tabIndex < 0 || tabIndex >= len) {
            //     tabIndex = 0
            // }
            // this.activatedIndex = tabIndex
            // const tabIndex = findIndex(this.tabList, {name})
            // if (tabIndex === -1) {
            //     return
            // }
            // this.activatedIndex = tabIndex
        },
        switchSubTab(name) {
            const tab = this.currentTab
            if (tab == null) {
                return
            }
            tab.subTab = name
        },
        /**
         *
         * @param {string} tabName
         */
        removeTabByName(tabName) {
            const idx = findIndex(this.tabs, { name: tabName })
            if (idx !== -1) {
                this.removeTab(idx)
            }
        },
    }
})



export default useTabStore