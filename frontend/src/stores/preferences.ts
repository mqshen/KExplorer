import { defineStore } from 'pinia'
import {  isEmpty, join, map } from 'lodash'

const usePreferencesStore = defineStore('preferences', {
    state: () => ({
        behavior: {
            welcomed: false,
            asideWidth: 300,
            windowWidth: 0,
            windowHeight: 0,
            windowMaximised: false,
        },
        editor: {
            font: '',
            fontFamily: [],
            fontSize: 14,
            showLineNum: true,
            showFolding: true,
            dropText: true,
            links: true,
        },
    }),
    getters: {
        /**
         * current editor font
         * @return {{fontSize: string, fontFamily?: string}}
         */
        editorFont() {
            const fontStyle = {
                fontSize: (this.editor.fontSize || 14) + 'px',
            }
            if (!isEmpty(this.editor.fontFamily)) {
                fontStyle['fontFamily'] = join(
                    map(this.editor.fontFamily, (f) => `"${f}"`),
                    ',',
                )
            }
            // compatible with old preferences
            // if (isEmpty(fontStyle['fontFamily'])) {
            //     if (!isEmpty(this.editor.font) && this.editor.font !== 'none') {
            //         const font = find(this.fontList, { name: this.editor.font })
            //         if (font != null) {
            //             fontStyle['fontFamily'] = `${font.name}`
            //         }
            //     }
            // }
            if (isEmpty(fontStyle['fontFamily'])) {
                fontStyle['fontFamily'] = ['monaco']
            }
            return fontStyle
        },
    },
    actions: {
        /**
         * load preferences from local
         * @returns {Promise<void>}
         */
        async loadPreferences() {
        },

        /**
         * save preferences to local
         * @returns {Promise<boolean>}
         */
        async savePreferences() {
            // const pf = pick(this, ['behavior', 'general', 'editor', 'cli', 'decoder'])
            // const { success, msg } = await SetPreferences(pf)
            // return success === true
        },
    }

})


export default usePreferencesStore