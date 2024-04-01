<script setup>
import { h, markRaw, nextTick, reactive, ref } from "vue";
import { NIcon, NSpace, NText, useThemeVars } from "naive-ui";
import { isEmpty, includes, get, split } from "lodash";
import { useI18n } from "vue-i18n";
import Folder from "@/components/icons/Folder.vue";
import Server from "@/components/icons/Server.vue";
import Cluster from "@/components/icons/Cluster.vue";
import Connect from "@/components/icons/Connect.vue";
import Config from "@/components/icons/Config.vue";
import Delete from "@/components/icons/Delete.vue";
import Unlink from "@/components/icons/Unlink.vue";
import IconButton from "@/components/common/IconButton.vue";
import { ConnectionType } from "@/consts/connection_type";
import useConnectionStore from "stores/connections";
import usePreferencesStore from "stores/preferences";
import useBrowserStore from "stores/browser";
import useTabStore from "stores/tab";
import useDialogStore from 'stores/dialog'

const i18n = useI18n();

const prefStore = usePreferencesStore();
const connectionStore = useConnectionStore();
const browserStore = useBrowserStore();
const tabStore = useTabStore();
const dialogStore = useDialogStore()

const expandedKeys = ref([]);
const selectedKeys = ref([]);

const props = defineProps({
  filterPattern: {
    type: String,
  },
});

const contextMenuParam = reactive({
  show: false,
  x: 0,
  y: 0,
  options: null,
  currentNode: null,
});

/**
 * get mark color of server saved in preferences
 * @param name
 * @return {null|string}
 */
const getServerMarkColor = (name) => {
  const { markColor = "" } = connectionStore.serverProfile[name] || {};
  if (!isEmpty(markColor)) {
    const rgb = parseHexColor(markColor);
    const rgb2 = hexGammaCorrection(rgb, 0.75);
    return toHexColor(rgb2);
  }
  return null;
};

const renderPrefix = ({ option }) => {
  const iconTransparency = prefStore.isDark ? 0.75 : 1;
  switch (option.type) {
    case ConnectionType.Group:
      const opened = indexOf(expandedKeys.value, option.key) !== -1;
      return h(
        NIcon,
        { size: 20 },
        {
          default: () =>
            h(Folder, {
              open: opened,
              fillColor: `rgba(255,206,120,${iconTransparency})`,
            }),
        }
      );
    case ConnectionType.Server:
      const connected = browserStore.isConnected(option.name);
      const color = getServerMarkColor(option.name);
      const icon = option.cluster === true ? Cluster : Server;
      return h(
        NIcon,
        { size: 20, color: !!!connected ? color : "#dc423c" },
        {
          default: () =>
            h(icon, {
              inverse: !!connected,
              fillColor: `rgba(220,66,60,${iconTransparency})`,
            }),
        }
      );
  }
};

/**
 * Open connection
 * @param name
 * @returns {Promise<void>}
 */
const openConnection = async (name) => {
  try {
    connectingServer.value = name;
    if (!browserStore.isConnected(name)) {
      await browserStore.openConnection(name);
    }
    // check if connection already canceled before finish open
    if (!isEmpty(connectingServer.value)) {
      tabStore.upsertTab({
        server: name,
        // db: browserStore.getSelectedDB(name),
      });
    }
  } catch (e) {
    $message.error(e.message);
    // node.isLeaf = undefined
  } finally {
    connectingServer.value = "";
  }
};

const onCancelOpen = () => {
  if (!isEmpty(connectingServer.value)) {
    browserStore.closeConnection(connectingServer.value);
    connectingServer.value = "";
  }
};

const nodeProps = ({ option }) => {
  return {
    onDblclick: async () => {
      if (option.type === ConnectionType.Server) {
        openConnection(option.name).then(() => {});
      } else if (option.type === ConnectionType.Group) {
        // toggle expand
        nextTick().then(() => expandKey(option.key));
      }
    },
    onContextmenu(e) {
      e.preventDefault();
      const mop = menuOptions[option.type];
      if (mop == null) {
        return;
      }
      contextMenuParam.show = false;
      nextTick().then(() => {
        contextMenuParam.options = markRaw(mop(option));
        contextMenuParam.currentNode = option;
        contextMenuParam.x = e.clientX;
        contextMenuParam.y = e.clientY;
        contextMenuParam.show = true;
        selectedKeys.value = [option.key];
      });
    },
  };
};

const handleDrop = ({ node, dragNode, dropPosition }) => {};

const renderLabel = ({ option }) => {
  if (option.type === ConnectionType.Server) {
    const color = getServerMarkColor(option.name);
    if (color != null) {
      return h(
        NText,
        {
          style: {
            color,
            fontWeight: "450",
          },
        },
        () => option.label
      );
    }
  }
  return option.label;
};

const connectingServer = ref("");

const removeConnection = (name) => {
  $dialog.warning(
    i18n.t("dialogue.remove_tip", {
      type: i18n.t("dialogue.connection.conn_name"),
      name,
    }),
    async () => {
      connectionStore.deleteConnection(name).then(({ success, msg }) => {
        if (!success) {
          $message.error(msg);
        }
      });
    }
  );
};

const handleSelectContextMenu = (key) => {
  contextMenuParam.show = false;
  const selectedKey = get(selectedKeys.value, 0);
  if (selectedKey == null) {
    return;
  }
  const [group, name] = split(selectedKey, "/");
  console.log("selectedKey", selectedKey, group, name);
  if (isEmpty(group) && isEmpty(name)) {
    return;
  }
  switch (key) {
    case "server_open":
      openConnection(name).then(() => {});
      break;
    case "server_remove":
      removeConnection(name);
      break;
    case "server_edit":
      // ask for close relevant connections before edit
      if (browserStore.isConnected(name)) {
        $dialog.warning(i18n.t("dialogue.edit_close_confirm"), () => {
          browserStore.closeConnection(name);
          dialogStore.openEditDialog(name);
        });
      } else {
        dialogStore.openEditDialog(name);
      }
      break;
    default:
      console.warn("TODO: handle context menu:" + key);
  }
};

const onUpdateExpandedKeys = (keys, option) => {
  expandedKeys.value = keys;
};

const onUpdateSelectedKeys = (keys, option) => {
  selectedKeys.value = keys;
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

const getServerMenu = (connected) => {
  const btns = [];
  if (connected) {
    btns.push(
      h(IconButton, {
        tTooltip: "interface.disconnect",
        icon: Unlink,
        onClick: () => handleSelectContextMenu("server_close"),
      }),
      h(IconButton, {
        tTooltip: "interface.edit_conn",
        icon: Config,
        onClick: () => handleSelectContextMenu("server_edit"),
      })
    );
  } else {
    btns.push(
      h(IconButton, {
        tTooltip: "interface.open_connection",
        icon: Connect,
        onClick: () => handleSelectContextMenu("server_open"),
      }),
      h(IconButton, {
        tTooltip: "interface.edit_conn",
        icon: Config,
        onClick: () => handleSelectContextMenu("server_edit"),
      }),
      h(IconButton, {
        tTooltip: "interface.remove_conn",
        icon: Delete,
        onClick: () => handleSelectContextMenu("server_remove"),
      })
    );
  }
  return btns;
};

const renderSuffix = ({ option }) => {
  if (includes(selectedKeys.value, option.key)) {
    switch (option.type) {
      case ConnectionType.Server:
        const connected = browserStore.isConnected(option.name);
        return renderIconMenu(getServerMenu(connected));
      case ConnectionType.Group:
        return renderIconMenu(getGroupMenu());
    }
  }
  return null;
};
</script>

<template>
  <div
    class="connection-tree-wrapper"
    @keydown.esc="contextMenuParam.show = false"
  >
    <n-empty
      v-if="isEmpty(connectionStore.connections)"
      :description="$t('interface.empty_server_list')"
      class="empty-content"
    />
    <n-tree
      v-else
      :animated="false"
      :block-line="true"
      :block-node="true"
      :cancelable="false"
      :data="connectionStore.connections"
      :draggable="true"
      :expanded-keys="expandedKeys"
      :node-props="nodeProps"
      :pattern="props.filterPattern"
      :render-label="renderLabel"
      :render-prefix="renderPrefix"
      :render-suffix="renderSuffix"
      :selected-keys="selectedKeys"
      class="fill-height"
      virtual-scroll
      @drop="handleDrop"
      @update:selected-keys="onUpdateSelectedKeys"
      @update:expanded-keys="onUpdateExpandedKeys"
    />

    <!-- status display modal -->
    <n-modal :show="connectingServer !== ''" transform-origin="center">
      <n-card
        :bordered="false"
        :content-style="{ textAlign: 'center' }"
        aria-model="true"
        role="dialog"
        style="width: 400px"
      >
        <n-spin>
          <template #description>
            <n-space vertical>
              <n-text strong>{{ $t("dialogue.opening_connection") }}</n-text>
              <n-button
                :focusable="false"
                secondary
                size="small"
                @click="onCancelOpen"
              >
                {{ $t("dialogue.interrupt_connection") }}
              </n-button>
            </n-space>
          </template>
        </n-spin>
      </n-card>
    </n-modal>

    <!-- context menu -->
    <n-dropdown
      :keyboard="true"
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
      @clickoutside="contextMenuParam.show = false"
      @select="handleSelectContextMenu"
    />
  </div>
</template>
