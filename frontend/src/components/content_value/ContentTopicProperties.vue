<script setup>
import { get, size, find } from "lodash";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import useTabStore from "stores/tab";
import useTopicStore from "stores/topics";
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

const options = [
  { label: "Byte Array", value: 0 },
  { label: "String", value: 1 },
  { label: "Avro", value: 2 },
];

const updateTopicSerializator = async () => {
  const { topic, keySerializer, valueSerializer } = currentNode.value;
  if (!topic) {
    $message.error($t("message.empty_topic"));
  }
  const { success, msg } = await topicStore.saveTopic(props.server, topic, {
    keySerializer: keySerializer,
    valueSerializer: valueSerializer,
  });
  if (!success) {
    $message.error(msg);
    return;
  }

  $message.success(i18n.t("dialogue.handle_succ"));
};
</script>
<template>
  <div class="properties-container flex-box-v">
    <n-divider title-placement="left">
      {{ $t("interface.general") }}
    </n-divider>
    <n-form
      label-placement="left"
      require-mark-placement="right-hanging"
      label-width="120"
    >
      <n-form-item :label="$t('interface.topic_name')" path="inputValue">
        <n-input
          v-model:value="currentNode.topic"
          placeholder="Input"
          :disabled="true"
        />
      </n-form-item>
    </n-form>
    <n-divider title-placement="left">
      {{ $t("interface.content_types") }}
    </n-divider>
    <n-form
      label-placement="left"
      require-mark-placement="right-hanging"
      label-width="120"
    >
      <n-form-item :label="$t('interface.key')" path="inputValue">
        <n-select
          v-model:value="currentNode.keySerializer"
          :options="options"
        />
      </n-form-item>
      <n-form-item :label="$t('interface.value')" path="inputValue">
        <n-select
          v-model:value="currentNode.valueSerializer"
          :options="options"
        />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <div style="display: flex; justify-content: flex-end">
            <n-button type="primary" @click="updateTopicSerializator">
              {{ $t("interface.update") }}
            </n-button>
          </div>
        </n-col>
      </n-row>
    </n-form>
    <n-divider title-placement="left">
      {{ $t("interface.message") }}
    </n-divider>
  </div>
</template>
<style scoped>
.properties-container {
  padding: 0 20px;
}
</style>
