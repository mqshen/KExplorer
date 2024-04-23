<script setup>
import { get, isEmpty, map } from "lodash";
import { computed, nextTick, ref, watch } from "vue";
import useDialog, { ConnDialogType } from "stores/dialog";
import { useI18n } from "vue-i18n";
import { TestConnection } from "wailsjs/go/services/connectionService";
import useClusterStore from "stores/cluster";
import useBrowserStore from "stores/browser";
const i18n = useI18n();

const editName = ref("");

const dialogStore = useDialog();
const clusterStore = useClusterStore();
const browserStore = useBrowserStore();

const tab = ref("general");
const testing = ref(false);
const testResult = ref(null);

const generalFormRef = ref(null);
const generalForm = ref(null);

const resetForm = () => {
  generalForm.value = clusterStore.newDefaultCluster();
  generalFormRef.value?.restoreValidation();
  testing.value = false;
  testResult.value = null;
  tab.value = "general";
};

const isEditMode = computed(() => dialogStore.connType === ConnDialogType.EDIT);
const closingConnection = computed(() => {
  if (isEmpty(editName.value)) {
    return false;
  }
  return browserStore.isConnected(editName.value);
});

const generalFormRules = () => {
  const requiredMsg = i18n.t("dialogue.field_required");
  const illegalChars = ["/", "\\"];
  return {};
};

const onTestConnection = async () => {
  testResult.value = "";
  testing.value = true;
  let result = "";
  try {
    const { success = false, msg } = await TestConnection(generalForm.value);
    if (!success) {
      result = msg;
    }
  } catch (e) {
    result = e.message;
  } finally {
    testing.value = false;
  }

  if (!isEmpty(result)) {
    testResult.value = result;
  } else {
    testResult.value = "";
  }
};

const onClose = () => {
  dialogStore.closeConnDialog();
};

const onSaveConnection = async () => {
  // validate general form
  await generalFormRef.value?.validate((err) => {
    if (err) {
      nextTick(() => (tab.value = "general"));
    }
  });

  // store new connection
  const { success, msg } = await clusterStore.saveCluster(
    isEditMode.value ? editName.value : null,
    generalForm.value
  );
  if (!success) {
    $message.error(msg);
    return;
  }

  $message.success(i18n.t("dialogue.handle_succ"));
  onClose();
};

const pasteFromClipboard = async () => {};

watch(
  () => dialogStore.connDialogVisible,
  (visible) => {
    if (visible) {
      resetForm();
      editName.value = get(dialogStore.connParam, "name", "");
      generalForm.value =
        dialogStore.connParam || clusterStore.newDefaultCluster();
    }
  }
);

const bootstrapType = ref(0);
</script>
<template>
  <n-modal
    :show="dialogStore.connDialogVisible"
    :closable="false"
    :close-on-esc="false"
    :mask-closable="false"
    :on-after-leave="resetForm"
    :show-icon="false"
    :title="
      isEditMode
        ? $t('dialogue.cluster.edit_title')
        : $t('dialogue.cluster.new_title')
    "
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
  >
    <n-spin :show="closingConnection">
      <n-tabs
        v-model:value="tab"
        animated
        pane-style="min-height: 50vh;"
        placement="left"
        tab-style="justify-content: right; font-weight: 420;"
        type="line"
      >
        <n-tab-pane
          :tab="$t('dialogue.cluster.general')"
          display-directive="show:lazy"
          name="general"
        >
          <n-form
            ref="generalFormRef"
            :model="generalForm"
            :rules="generalFormRules()"
            :show-require-mark="false"
            label-placement="top"
          >
            <n-form-item
              :label="$t('dialogue.cluster.conn_name')"
              :span="24"
              path="name"
              required
            >
              <n-input
                v-model:value="generalForm.name"
                :placeholder="$t('dialogue.cluster.name_tip')"
              />
            </n-form-item>

            <!-- <n-form-item
              :label="$t('dialogue.cluster.bootstrap')"
              path="radioGroupValue"
            >
              <n-radio-group v-model:value="bootstrapType" name="BootstrapType">
                <n-space>
                  <n-radio :value="0"> Zookeeper </n-radio>
                  <n-radio :value="1"> Kafka Bootstrap Servers </n-radio>
                </n-space>
              </n-radio-group>
            </n-form-item> -->

            <n-form-item
              :label="$t('dialogue.cluster.bootstrap_server')"
              :span="24"
              path="bootstrap"
              required
            >
              <n-input
                v-model:value="generalForm.bootstrap"
                placeholder="localhost:9092"
              />
            </n-form-item>
          </n-form>
        </n-tab-pane>
        <n-tab-pane
          :tab="$t('dialogue.cluster.security')"
          display-directive="show:lazy"
          name="security"
        >
        </n-tab-pane>
      </n-tabs>
    </n-spin>

    <template #action>
      <div class="flex-item-expand">
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          :loading="testing"
          @click="onTestConnection"
        >
          {{ $t("dialogue.cluster.test") }}
        </n-button>
      </div>
      <div class="flex-item n-dialog__action">
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          @click="pasteFromClipboard"
        >
          {{ $t("dialogue.cluster.parse_url_clipboard") }}
        </n-button>
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          @click="onClose"
        >
          {{ $t("common.cancel") }}
        </n-button>
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          type="primary"
          @click="onSaveConnection"
        >
          {{
            isEditMode ? $t("preferences.general.update") : $t("common.confirm")
          }}
        </n-button>
      </div>
    </template>
  </n-modal>
</template>
