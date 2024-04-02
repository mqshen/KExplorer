import { findIndex, isEmpty, set, get, find, size } from 'lodash'
import { defineStore } from 'pinia'
import { TabItem } from '@/objects/tabItem'
import useTopicStore from "stores/topics";

const useTabStore = defineStore('tab', {
    state: () => ({
        nav: 'server',
        asideWidth: 300,
        tabList: new Array<TabItem>(),
        activatedIndex: 0,
    }),
    getters: {
        /**
         * get current tab list item
         * @returns {TabItem[]}
         */
        tabs(): Array<TabItem> {
            // if (isEmpty(this.tabList)) {
            //     this.newBlankTab()
            // }
            return this.tabList
        },

        /**
         * get current activated tab item
         * @returns {TabItem|null}
         */
        currentTab(): TabItem {
            return get(this.tabs, this.activatedIndex)
        },
        currentTabName(): string {
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
        _setActivatedIndex(idx: number, switchNav: boolean, subTab: string = "") {
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
        upsertTab(server: string): void {
            let tabIndex = this.tabList.findIndex((tab: TabItem) => tab.name == server)
            if (tabIndex === -1) {
                const tabItem = new TabItem(server, server)
                this.tabList.push(tabItem)
                tabIndex = this.tabList.length - 1
            } else {
                const tab = this.tabList[tabIndex]
                tab.blank = false
                // tab.title = db !== undefined ? `${server}/db${db}` : `${server}`
                tab.title = server
            }
            this._setActivatedIndex(tabIndex, true)
        },/**
        * set expanded keys for server
        * @param {string} server
        * @param {string[]} keys
        */
        setExpandedKeys(server: string, keys = []) {
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
        setSelectedKeys(server: string, keys = null, node: any) {
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

        switchTab(tabIndex: number) {
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
        switchSubTab(name: string) {
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
        removeTabByName(tabName: string) {
            const idx = findIndex(this.tabs, { name: tabName })
            if (idx !== -1) {
                this.removeTab(idx)
            }
        },

        /**
         *
         * @param {number} tabIndex
         * @returns {*|null}
         */
        removeTab(tabIndex: number) {
            const len = size(this.tabs)
            // ignore remove last blank tab
            if (len === 1 && this.tabs[0].blank) {
                return null
            }

            if (tabIndex < 0 || tabIndex >= len) {
                return null
            }
            const removed = this.tabList.splice(tabIndex, 1)

            // update select index if removed index equal current selected
            this.activatedIndex -= 1
            if (this.activatedIndex < 0) {
                if (this.tabList.length > 0) {
                    this._setActivatedIndex(0, false)
                } else {
                    this._setActivatedIndex(-1, false)
                }
            } else {
                this._setActivatedIndex(this.activatedIndex, false)
            }

            return size(removed) > 0 ? removed[0] : null
        },

    }
})

export default useTabStore