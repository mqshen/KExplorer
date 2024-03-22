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
    actions: {
        openNewDialog() {
            this.connParam = null
            this.connType = ConnDialogType.NEW
            this.connDialogVisible = true
        },
    }
})


export default useDialogStore