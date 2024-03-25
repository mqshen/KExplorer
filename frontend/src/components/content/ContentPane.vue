<script setup></script>
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
    </n-tabs>
  </div>
</template>
