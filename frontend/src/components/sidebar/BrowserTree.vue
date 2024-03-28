<script setup>
import { includes, isEmpty, get, size } from "lodash";
import { computed, reactive, ref, h } from "vue";
import { NIcon, NSpace, NText, useThemeVars } from "naive-ui";
import { useI18n } from "vue-i18n";
import useConnectionStore from "stores/connections";
import useBrowserStore from "stores/browser";
import useTabStore from "stores/tab";
import { parseHexColor } from "@/utils/rgb";
import Layer from "@/components/icons/Layer.vue";
import Key from "@/components/icons/Key.vue";

const props = defineProps({
  server: String,
  db: Number,
  keyView: String,
  loading: Boolean,
  pattern: String,
  fullLoaded: Boolean,
  checkMode: Boolean,
});

const i18n = useI18n();
const connectionStore = useConnectionStore();
const browserStore = useBrowserStore();
const tabStore = useTabStore();

const treeKey = ref(0);

const backgroundColor = computed(() => {
  const { markColor: hex = "" } =
    connectionStore.serverProfile[props.server] || {};
  if (isEmpty(hex)) {
    return "";
  }
  const { r, g, b } = parseHexColor(hex);
  return `rgba(${r}, ${g}, ${b}, 0.2)`;
});

const contextMenuParam = reactive({
  show: false,
  x: 0,
  y: 0,
  options: null,
});

/**
 *
 * @type {ComputedRef<string[]>}
 */
const expandedKeys = computed(() => {
  const tab = find(tabStore.tabList, { name: props.server });
  if (tab != null) {
    return get(tab, "expandedKeys", []);
  }
  return [];
});

/**
 *
 * @type {ComputedRef<string[]>}
 */
const selectedKeys = computed(() => {
  const tab = find(tabStore.tabList, { name: props.server });
  if (tab != null) {
    return get(tab, "selectedKeys", []);
  }
  return [];
});
const emit = defineEmits(['change'])

const onUpdateSelectedKeys = (keys, options) => {
  if (!isEmpty(keys)) {
    console.log(keys, options)
    tabStore.setSelectedKeys(props.server, keys, options[0]);
  } else {
    // default is load blank key to display server status
    // tabStore.openBlank(props.server)
  }
};

const onUpdateExpanded = (value, option, meta) => {
  tabStore.setExpandedKeys(props.server, value);
  if (!meta.node) {
    return;
  }

  // keep expand or collapse children while they own more than 1 child
  let node = meta.node;
  while (node != null && size(node.children) === 1) {
    const key = node.children[0].key;
    switch (meta.action) {
      case "expand":
        node.expanded = true;
        if (!includes(value, key)) {
          tabStore.addExpandedKey(props.server, key);
        }
        break;
      case "collapse":
        node.expanded = false;
        tabStore.removeExpandedKey(props.server, key);
        break;
    }
    node = node.children[0];
  }
};

const data = computed(() => {
  console.log("lllls");
  return browserStore.getKeyStruct(props.server, props.checkMode);
});

const handleSelectContextMenu = (action) => {};

const handleOutsideContextMenu = () => {
  contextMenuParam.show = false;
};

const nodeProps = ({ option }) => {
  return {
    onClick: () => {},
    onDblclick: () => {},
    onContextmenu(e) {
      e.preventDefault();
    },
    // onMouseover() {
    //   console.log('mouse over')
    // }
  };
};

const renderPrefix = ({ option }) => {
  if (option.isLeaf) {
    return h(NIcon, { size: 20 }, () => h(Key));
  } else {
    return h(
      NIcon,
      { size: 20 },
      {
        default: () => h(Layer),
      }
    );
  }
};

// render tree item label
const renderLabel = ({ option }) => {
  return option.label;
};

// render horizontal item
const renderIconMenu = (items) => {
  return h(
    NSpace,
    {
      align: "center",
      inline: true,
      size: 3,
      wrapItem: false,
      wrap: false,
      style: "margin-right: 5px",
    },
    () => items
  );
};

// render menu function icon
const renderSuffix = ({ option }) => {
  const selected = includes(selectedKeys.value, option.key);
  if (selected && !props.checkMode) {
    switch (option.type) {
      case ConnectionType.RedisDB:
        return renderIconMenu(
          calcDBMenu(option.opened, option.loading, option.fullLoaded)
        );
      case ConnectionType.RedisKey:
        return renderIconMenu(calcLayerMenu(option.loading));
      case ConnectionType.RedisValue:
        return renderIconMenu(calcValueMenu());
    }
  } else if (
    !selected &&
    !!option.redisKeyCode &&
    option.type === ConnectionType.RedisValue
  ) {
    // render binary icon
    return renderIconMenu(h(NIcon, { size: 20 }, () => h(Binary)));
  }
  return null;
};
</script>

<template>
  <div
    :style="{ backgroundColor }"
    class="flex-box-v browser-tree-wrapper"
    @contextmenu="(e) => e.preventDefault()"
    @keydown.esc="contextMenuParam.show = false"
  >
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
      :data="data"
      :render-label="renderLabel"
      :render-prefix="renderPrefix"
      :render-suffix="renderSuffix"
      :selected-keys="selectedKeys"
      :show-irrelevant-nodes="false"
      check-strategy="child"
      class="fill-height"
      virtual-scroll
      @update:selected-keys="onUpdateSelectedKeys"
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
