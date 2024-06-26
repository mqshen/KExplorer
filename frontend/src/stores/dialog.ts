import { defineStore } from 'pinia'
import useConnectionStore from './connections.ts'

/**
 * connection dialog type
 * @enum {number}
 */
export const ConnDialogType = {
    NEW: 0,
    EDIT: 1,
}

const useDialogStore = defineStore('dialog', {
    state: () => ({
        connDialogVisible: false,
    }),
    actions: {
        openNewDialog() {
            this.connParam = null
            this.connType = ConnDialogType.NEW
            this.connDialogVisible = true
        },
        closeConnDialog() {
            this.connDialogVisible = false
        },

        async openEditDialog(name) {
            const connStore = useConnectionStore()
            const profile = await connStore.getConnectionProfile(name)
            this.connParam = connStore.mergeConnectionProfile(connStore.newDefaultConnection(name), profile)
            this.connType = ConnDialogType.EDIT
            this.connDialogVisible = true
        },
    }
})


export default useDialogStore