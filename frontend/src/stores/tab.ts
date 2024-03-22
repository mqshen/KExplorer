import { defineStore } from 'pinia'

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
    }
})



export default useTabStore