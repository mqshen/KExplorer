<script setup>
import { get, size, find } from "lodash";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import useTabStore from "stores/tab";
import useTopicStore from "stores/topics";
import Start from "@/components/icons/Start.vue";
import StopIcon from "@/components/icons/Stop.vue";
import { FetchMessages } from "wailsjs/go/services/browserService";
import VerticalResizeableWrapper from "@/components/common/VerticalResizeableWrapper.vue";
import ContentValueString from "@/components/content_value/ContentValueString.vue";
const tabStore = useTabStore();
const topicStore = useTopicStore();

const props = defineProps({
  server: String,
});

const i18n = useI18n();

const currentNode = computed(() => {
  const tab = find(tabStore.tabList, { name: props.server });
  console.log(tab);
  if (tab != null) {
    return tab.currentNode;
  }
  return {};
});

const messageOrder = ref(0);

const options = [
  { label: i18n.t("interface.oldest"), value: 0 },
  { label: i18n.t("interface.newest"), value: 1 },
];
const columns = [
  {
    title: "Partition",
    key: "partition",
    width: 120,
    align: "center",
  },
  {
    title: "Offset",
    key: "offset",
    width: 120,
    align: "center",
  },
  {
    title: "Key",
    key: "key",
    width: 200,
    align: "center",
  },
  {
    title: "Value",
    key: "value",
    titleAlign: "center",
    render: ({ value }, index) => {
      if (value.length > 64) {
        return value.substring(0, 64) + "...";
      }
      return value;
    },
  },
  {
    title: "Timestamp",
    key: "timestamp",
    width: 200,
    align: "center",
  },
];
const kafkaMessages = ref([]);

const fetchMessages = () => {
  const tab = find(tabStore.tabList, { name: props.server });
  console.log(tab);
  if (tab != null) {
    const currentNode = tab.currentNode;
    const { topic, keySerializer, valueSerializer } = currentNode;
    FetchMessages(props.server, topic, {
      keySerializer: keySerializer,
      valueSerializer: valueSerializer,
      messageOrder: messageOrder.value,
      Size: 30,
    }).then((res) => {
      const { data } = res;
      if (data) {
        const { messages = [] } = data;
        console.log(messages);
        kafkaMessages.value = messages;
      }
    });
  }
};
const asideWidth = ref(500);
const messageContent = ref(null)
const rowProps = (row) => {
  return {
    style: "cursor: pointer;",
    onClick: () => {
      messageContent.value = row.value
    },
  };
};
</script>
<template>
  <div class="content-wrapper flex-box-v message-container">
    <n-row :gutter="[0, 24]">
      <n-col :span="2"> </n-col>
      <n-col :span="6">
        <div style="display: flex" class="button-group">
          <n-button strong secondary circle @click="fetchMessages">
            <template #icon>
              <n-icon><start /></n-icon>
            </template>
          </n-button>
          <n-button strong secondary circle>
            <template #icon>
              <n-icon><stop-icon /></n-icon>
            </template>
          </n-button>
        </div>
      </n-col>
      <n-col :span="6"> </n-col>
      <n-col :span="9">
        <div style="display: flex; justify-content: flex-end">
          <span class="input-label">{{ $t("interface.message") }}</span>
          <n-select v-model:value="messageOrder" :options="options" />
        </div>
      </n-col>
      <n-col :span="1"> </n-col>
    </n-row>
    <vertical-resizeable-wrapper v-model:size="asideWidth" :min-size="300">
      <n-data-table
        :columns="columns"
        :data="kafkaMessages"
        :bordered="true"
        :row-props="rowProps"
        :max-height="420"
      />
    </vertical-resizeable-wrapper>
    <content-value-string 
        :value="messageContent"> </content-value-string>
  </div>
</template>
<style lang="scss" scoped>
.n-data-table {
  margin-top: 10px;
}
.input-label {
  line-height: 32px;
  margin-right: 5px;
}
.button-group {
  button {
    margin-right: 10px;
  }
}
.n-data-table-wrapper {
  overflow-y: scroll !important;
}
</style>
