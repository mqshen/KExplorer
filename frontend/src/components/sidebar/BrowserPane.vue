<script setup>
import { isEmpty } from 'lodash'
import { computed, reactive, ref, onMounted } from 'vue'
import useBrowserStore from 'stores/browser'
import BrowserTree from './BrowserTree.vue'

const browserStore = useBrowserStore()

const browserTreeRef = ref(null)
const inCheckState = ref(false)
const loading = ref(false)
const fullyLoaded = ref(false)

const props = defineProps({
    server: String,
})

// const loadProgress = computed(() => {
//     const hasPattern = !isEmpty(filterForm.pattern)
//     if (hasPattern) {
//         return 100
//     }

//     const db = browserStore.getDatabase(props.server)
//     if (db.maxKeys <= 0) {
//         return 100
//     }
//     return (db.keyCount * 100) / Math.max(db.keyCount, db.maxKeys)
// })

const filterForm = reactive({
    type: '',
    exact: false,
    pattern: '',
    filter: '',
})

const onReload = async () => {
    try {
        loading.value = true
        // tabStore.setSelectedKeys(props.server)
        // browserStore.closeConnection(props.server)

        // let matchType = unref(filterForm.type)
        // if (!types.hasOwnProperty(matchType)) {
        //     matchType = ''
        // }
        // browserStore.setKeyFilter(props.server, {
        //     type: matchType,
        //     pattern: unref(filterForm.pattern),
        //     exact: unref(filterForm.exact) === true,
        // })
        console.log(" test for test", props.server)
        await browserStore.getKafkaMetaData(props.server)
        // fullyLoaded.value = await browserStore.loadMoreKeys(props.server)
        // $message.success(i18n.t('dialogue.reload_succ'))
    } catch (e) {
        console.warn(e)
    } finally {
        loading.value = false
    }
}

onMounted(() => onReload())

</script>
<template>

<div class="nav-pane-container flex-box-v">
    <!-- tree view -->
    <browser-tree
            ref="browserTreeRef"
            :check-mode="inCheckState"
            :full-loaded="fullyLoaded"
            :loading="loading"
            :pattern="filterForm.filter"
            :server="props.server" />
</div>
    
</template>