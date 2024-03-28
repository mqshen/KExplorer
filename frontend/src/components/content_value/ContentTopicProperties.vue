<script setup>
import { get, size, find } from "lodash";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import useTabStore from "stores/tab";
const tabStore = useTabStore();

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
  { label: "Byte Array", value: "Byte" },
  { label: "String", value: "String" },
  { label: "Avro", value: "Avro" },
];

const updateTopicSerializator = () => {};
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
          v-model:value="currentNode.label"
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
          v-model:value="currentNode.keySerializator"
          :options="options"
        />
      </n-form-item>
      <n-form-item :label="$t('interface.value')" path="inputValue">
        <n-select
          v-model:value="currentNode.valueSerializator"
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
