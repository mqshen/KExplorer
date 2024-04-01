<script setup>
import { computed, ref } from "vue";
import { useThemeVars } from 'naive-ui'
import useTabStore from "stores/tab";
import { BrowserTabType } from "@/consts/browser_tab_type";
import Status from '@/components/icons/Status.vue'
import Detail from '@/components/icons/Detail.vue'
import ContentTopicProperties from '@/components/content_value/ContentTopicProperties.vue'
import ContentTopicData from '@/components/content_value/ContentTopicData.vue'

const themeVars = useThemeVars()
const tabStore = useTabStore();

const props = defineProps({
    server: String,
})

const tabsRef = ref(null)
const selectedSubTab = computed(() => {
  const { subTab = "properties" } = tabStore.currentTab || {};
  return subTab;
});

</script>
<template>
  <div class="content-container flex-box-v">
    <n-tabs
      ref="tabsRef"
      :tabs-padding="5"
      :theme-overrides="{
        tabFontWeightActive: 'normal',
        tabGapSmallLine: '10px',
        tabGapMediumLine: '10px',
        tabGapLargeLine: '10px',
      }"
      :value="selectedSubTab"
      class="content-sub-tab"
      default-value="status"
      pane-class="content-sub-tab-pane"
      placement="top"
      tab-style="padding-left: 10px; padding-right: 10px;"
      type="line"
      @update:value="tabStore.switchSubTab"
    >
      <n-tab-pane
        :name="BrowserTabType.Properties.toString()"
        display-directive="show:lazy"
      >
        <template #tab>
          <n-space
            :size="5"
            :wrap-item="false"
            align="center"
            inline
            justify="center"
          >
            <n-icon size="16">
              <status
                :inverse="
                  selectedSubTab === BrowserTabType.Properties.toString()
                "
                :stroke-color="themeVars.tabColor"
                stroke-width="4"
              />
            </n-icon>
            <span>{{ $t("interface.sub_tab.properties") }}</span>
          </n-space>
        </template>
        <content-topic-properties
          :server="props.server"
        />
      </n-tab-pane>
      <n-tab-pane
        :name="BrowserTabType.Data.toString()"
        display-directive="show:lazy"
      >
        <template #tab>
          <n-space
            :size="5"
            :wrap-item="false"
            align="center"
            inline
            justify="center"
          >
            <n-icon size="16">
              <detail
                :inverse="
                  selectedSubTab === BrowserTabType.Data.toString()
                "
                :stroke-color="themeVars.tabColor"
                stroke-width="4"
              />
            </n-icon>
            <span>{{ $t("interface.sub_tab.data") }}</span>
          </n-space>
        </template>
        <content-topic-data
          :server="props.server"
        />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>


<style lang="scss" scoped>
@import '@/styles/content';

.content-container {
    //padding: 5px 5px 0;
    //padding-top: 0;
    box-sizing: border-box;
    background-color: v-bind('themeVars.tabColor');
}
</style>

<style lang="scss">
.content-sub-tab {
    background-color: v-bind('themeVars.tabColor');
    height: 100%;
}

.content-sub-tab-pane {
    padding: 0 !important;
    height: 100%;
    background-color: v-bind('themeVars.bodyColor');
    overflow: hidden;
}

.n-tabs .n-tabs-bar {
    transition: none !important;
}
</style>