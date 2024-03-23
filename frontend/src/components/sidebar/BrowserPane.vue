<script setup>
import { computed, reactive } from 'vue'
import useBrowserStore from 'stores/browser'
import BrowserTree from './BrowserTree.vue'

const browserStore = useBrowserStore()

const browserTreeRef = ref(null)
const inCheckState = ref(false)
const loading = ref(false)
const fullyLoaded = ref(false)

const props = defineProps({
    server: String,
    db: {
        type: Number,
        default: 0,
    },
})

const loadProgress = computed(() => {
    const hasPattern = !isEmpty(filterForm.pattern)
    if (hasPattern) {
        return 100
    }

    const db = browserStore.getDatabase(props.server, props.db)
    if (db.maxKeys <= 0) {
        return 100
    }
    return (db.keyCount * 100) / Math.max(db.keyCount, db.maxKeys)
})

const filterForm = reactive({
    type: '',
    exact: false,
    pattern: '',
    filter: '',
})
</script>
<template>

<div class="nav-pane-container flex-box-v">
    <!-- tree view -->
    <browser-tree
            ref="browserTreeRef"
            :check-mode="inCheckState"
            :db="props.db"
            :full-loaded="fullyLoaded"
            :loading="loading && loadProgress <= 0"
            :pattern="filterForm.filter"
            :server="props.server" />
            jsjdfjsjdfjsjdfj
</div>
    
</template>