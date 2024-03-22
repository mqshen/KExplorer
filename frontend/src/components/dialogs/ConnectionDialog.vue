<script setup>
import { isEmpty } from "lodash";
import { computed, nextTick, ref, watch } from "vue";
import useDialog, { ConnDialogType } from "stores/dialog";
import { useI18n } from "vue-i18n";
import { TestConnection } from "wailsjs/go/services/connectionService";
import useConnectionStore from "stores/connections";
const i18n = useI18n();

const editName = ref("");

const dialogStore = useDialog();
const connectionStore = useConnectionStore();

const tab = ref("general");
const testing = ref(false);
const testResult = ref(null);

const generalFormRef = ref(null);
const generalForm = ref({});

const resetForm = () => {
  generalForm.value = connectionStore.newDefaultConnection();
  generalFormRef.value?.restoreValidation();
  testing.value = false;
  testResult.value = null;
  tab.value = "general";
  loadingSentinelMaster.value = false;
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
  return {
    name: [
      { required: true, message: requiredMsg, trigger: "input" },
      {
        validator: (rule, value) => {
          return every(illegalChars, (c) => !includes(value, c));
        },
        message: i18n.t("dialogue.illegal_characters"),
        trigger: "input",
      },
    ],
    addr: { required: true, message: requiredMsg, trigger: "input" },
    defaultFilter: { required: true, message: requiredMsg, trigger: "input" },
    keySeparator: { required: true, message: requiredMsg, trigger: "input" },
  };
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
        ? $t('dialogue.connection.edit_title')
        : $t('dialogue.connection.new_title')
    "
    preset="dialog"
    style="width: 600px"
    transform-origin="center"
  >
    <n-spin :show="closingConnection">
      <n-tabs
        :value="tab"
        animated
        pane-style="min-height: 50vh;"
        placement="left"
        tab-style="justify-content: right; font-weight: 420;"
        type="line"
      >
        <n-tab-pane
          :tab="$t('dialogue.connection.general')"
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
            <n-grid :x-gap="10">
              <n-form-item-gi
                :label="$t('dialogue.connection.conn_name')"
                :span="24"
                path="name"
                required
              >
                <n-input
                  :value="generalForm.name"
                  :placeholder="$t('dialogue.connection.name_tip')"
                />
              </n-form-item-gi>
              <n-form-item-gi
                v-if="!isEditMode"
                :label="$t('dialogue.connection.group')"
                :span="24"
                required
              >
                <n-select
                  :value="generalForm.group"
                  :options="groupOptions"
                  :render-label="
                    ({ label, value }) => (value === '' ? $t(label) : label)
                  "
                />
              </n-form-item-gi>
              <n-form-item-gi
                :label="$t('dialogue.connection.addr')"
                :span="24"
                path="addr"
                required
              >
                <n-input-group>
                  <n-input
                    :value="generalForm.addr"
                    :placeholder="$t('dialogue.connection.addr_tip')"
                  />
                  <n-text style="width: 40px; text-align: center">:</n-text>
                  <n-input-number
                    :value="generalForm.port"
                    :max="65535"
                    :min="1"
                    :show-button="false"
                    placeholder="9092"
                    style="width: 200px"
                  />
                </n-input-group>
              </n-form-item-gi>
            </n-grid>
          </n-form>
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
          {{ $t("dialogue.connection.test") }}
        </n-button>
      </div>
      <div class="flex-item n-dialog__action">
        <n-button
          :disabled="closingConnection"
          :focusable="false"
          @click="pasteFromClipboard"
        >
          {{ $t("dialogue.connection.parse_url_clipboard") }}
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
