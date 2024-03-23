import { findIndex, isEmpty, set } from 'lodash'
import { defineStore } from 'pinia'
// import { TabItem } from '@/objects/tabItem'

const useTabStore = defineStore('tab', {
    state: () => ({
        nav: 'server',
        asideWidth: 300,
        tabList: [],
        activatedIndex: 0, 
    }),
    getters: {
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
                // const tabItem = new TabItem({
                //     name: server,
                //     title: server,
                //     server,
                // })
                // this.tabList.push(tabItem)
                tabIndex = this.tabList.length - 1
            } else {
                const tab = this.tabList[tabIndex]
                tab.blank = false
                // tab.title = db !== undefined ? `${server}/db${db}` : `${server}`
                tab.title = server
                tab.server = server
            }
            this._setActivatedIndex(tabIndex, true)
        }
    }
})



export default useTabStore