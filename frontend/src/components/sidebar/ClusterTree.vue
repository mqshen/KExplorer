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
import CopyLink from '@/components/icons/CopyLink.vue'
import IconButton from "@/components/common/IconButton.vue";
import { ClusterType } from "@/consts/cluster_type";
import useClusterStore from "stores/cluster";
import usePreferencesStore from "stores/preferences";
import useBrowserStore from "stores/browser";
import useTabStore from "stores/tab";
import useDialogStore from "stores/dialog";
import { useRender } from '@/utils/render'

const i18n = useI18n();
const render = useRender()

const prefStore = usePreferencesStore();
const clusterStore = useClusterStore();
const browserStore = useBrowserStore();
const tabStore = useTabStore();
const dialogStore = useDialogStore();

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

const menuOptions = {
    [ClusterType.Group]: ({ opened }) => [
        {
            key: 'group_rename',
            label: 'interface.rename_conn_group',
            icon: Edit,
        },
        {
            key: 'group_delete',
            label: 'interface.remove_conn_group',
            icon: Delete,
        },
    ],
    [ClusterType.Server]: ({ name }) => {
        const connected = browserStore.isConnected(name)
        if (connected) {
            return [
                {
                    key: 'server_close',
                    label: 'interface.disconnect',
                    icon: Unlink,
                },
                {
                    key: 'server_edit',
                    label: 'interface.edit_conn',
                    icon: Config,
                },
                {
                    key: 'server_dup',
                    label: 'interface.dup_conn',
                    icon: CopyLink,
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'server_remove',
                    label: 'interface.remove_conn',
                    icon: Delete,
                },
            ]
        } else {
            return [
                {
                    key: 'server_open',
                    label: 'interface.open_connection',
                    icon: Connect,
                },
                {
                    key: 'server_edit',
                    label: 'interface.edit_conn',
                    icon: Config,
                },
                {
                    key: 'server_dup',
                    label: 'interface.dup_conn',
                    icon: CopyLink,
                },
                {
                    type: 'divider',
                    key: 'd1',
                },
                {
                    key: 'server_remove',
                    label: 'interface.remove_conn',
                    icon: Delete,
                },
            ]
        }
    },
}
/**
 * get mark color of server saved in preferences
 * @param name
 * @return {null|string}
 */
const getServerMarkColor = (name) => {
  const { markColor = "" } = clusterStore.clusterProfile[name] || {};
  if (!isEmpty(markColor)) {
    const rgb = parseHexColor(markColor);
    const rgb2 = hexGammaCorrection(rgb, 0.75);
    return toHexColor(rgb2);
  }
  return null;
};

const renderPrefix = ({ option }) => {
  const iconTransparency = prefStore.isDark ? 0.75 : 1;
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
      tabStore.upsertTab(name);
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
      openConnection(option.name).then(() => {});
    },
    onContextmenu(e) {
      e.preventDefault();
      const mop = menuOptions[ClusterType.Server];
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
        selectedKeys.value = [option.name];
      });
    },
  };
};

const handleDrop = ({ node, dragNode, dropPosition }) => {};

const renderLabel = ({ option }) => {
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
      () => option.name
    );
  }
  return option.name;
};

const connectingServer = ref("");

const removeConnection = (name) => {
  $dialog.warning(
    i18n.t("dialogue.remove_tip", {
      type: i18n.t("dialogue.connection.conn_name"),
      name,
    }),
    async () => {
      clusterStore.deleteCluster(name).then(({ success, msg }) => {
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
  console.log(keys)
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
  if (includes(selectedKeys.value, option.name)) {
    const connected = browserStore.isConnected(option.name);
    return renderIconMenu(getServerMenu(connected));
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
      v-if="isEmpty(clusterStore.clusters)"
      :description="$t('interface.empty_server_list')"
      class="empty-content"
    />
    <n-tree
      v-else
      :animated="false"
      :block-line="true"
      :block-node="true"
      :cancelable="false"
      :data="clusterStore.clusters"
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
