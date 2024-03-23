<script setup>
import { computed, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import useConnectionStore from 'stores/connections'
import { parseHexColor } from '@/utils/rgb'

const props = defineProps({
    server: String,
    db: Number,
    keyView: String,
    loading: Boolean,
    pattern: String,
    fullLoaded: Boolean,
    checkMode: Boolean,
})

const i18n = useI18n()
const connectionStore = useConnectionStore()

const treeKey = ref(0)

const backgroundColor = computed(() => {
    const { markColor: hex = '' } = connectionStore.serverProfile[props.server] || {}
    if (isEmpty(hex)) {
        return ''
    }
    const { r, g, b } = parseHexColor(hex)
    return `rgba(${r}, ${g}, ${b}, 0.2)`
})

const contextMenuParam = reactive({
    show: false,
    x: 0,
    y: 0,
    options: null,
})
</script>

<template>
  <div
    :style="{ backgroundColor }"
    class="flex-box-v browser-tree-wrapper"
    @contextmenu="(e) => e.preventDefault()"
    @keydown.esc="contextMenuParam.show = false"
  >
  jsjjsjsjsjs
    <n-spin v-if="props.loading" class="fill-height" />
    <n-empty
      v-else-if="!props.loading && isEmpty(data)"
      class="empty-content"
    />
    <n-tree
      v-show="!props.loading && !isEmpty(data)"
      :key="treeKey"
      :animated="false"
      :block-line="true"
      :block-node="true"
      :cancelable="false"
      :cascade="true"
      :checkable="props.checkMode"
      :checked-keys="checkedKeys"
      :data="data"
      :expand-on-click="false"
      :expanded-keys="expandedKeys"
      :filter="(pattern, node) => includes(node.redisKey, pattern)"
      :node-props="nodeProps"
      :pattern="props.pattern"
      :render-label="renderLabel"
      :render-prefix="renderPrefix"
      :render-suffix="renderSuffix"
      :selected-keys="selectedKeys"
      :show-irrelevant-nodes="false"
      check-strategy="child"
      class="fill-height"
      virtual-scroll
      @update:selected-keys="onUpdateSelectedKeys"
      @update:expanded-keys="onUpdateExpanded"
      @update:checked-keys="onUpdateCheckedKeys"
    />

    <!-- context menu -->
    <n-dropdown
      :options="contextMenuParam.options"
      :render-icon="({ icon }) => render.renderIcon(icon)"
      :render-label="
        ({ label }) =>
          render.renderLabel($t(label), { class: 'context-menu-item' })
      "
      :show="contextMenuParam.show"
      :x="contextMenuParam.x"
      :y="contextMenuParam.y"
      placement="bottom-start"
      trigger="manual"
      @clickoutside="handleOutsideContextMenu"
      @select="handleSelectContextMenu"
    />
  </div>
</template>
