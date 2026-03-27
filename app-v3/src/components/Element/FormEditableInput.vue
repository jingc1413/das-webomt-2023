<template>
  <div style="width: 100%;" :class="{ 'is_edit': isAllowOpenEdit }">
    <div class="readonly_div" v-show="editable === false">
      <el-text class="readonly_label" v-if="subID">{{ inputModel }}</el-text>
      <template v-else>
        <div v-if="data?.Style?.hidden" />
        <div v-else-if="data?.Style?.input === 'button'">
          <my-button key="input-button" :owner="owner" :data="data" />
        </div>
        <el-button-group v-else-if="data?.Style?.input === 'buttonGroup'">
          <my-button v-for="item in data.Items" :key="item.key" :owner="owner" :data="item" />
        </el-button-group>
        <el-text class="readonly_label" v-else>{{ formatValueWithUnit }}</el-text>
      </template>
      <el-button v-if="isAllowOpenEdit" text class="readonly_button" @click="enterEditMode()">
        <el-icon>
          <EditPen />
        </el-icon>
      </el-button>
    </div>
    <div v-if="rwValue" v-show="editable === true" class="edit_div">
      <el-input v-if="subID" v-model="inputModel" :maxlength="param?.ByteSize" show-word-limit
        @keyup.esc="handleCancel()" style="width: 100%;">
      </el-input>

      <div v-else ref="eventDivRef" :contenteditable="true" @keyup.esc="handleCancel()">
        <div :contenteditable="false">

          <el-tooltip :disabled="!hasTips" :trigger-keys="[]" :show-after="1000">
            <template #content>
              <span v-if="param.Tips">{{ `Tips: ${param.Tips}` }}<br></span>
              <span>{{ `Value: ${param.Value !== undefined && param.Value !== null ? param.Value : ""}` }}</span>
              <span v-if="param.UnitName">{{ `, Unit: ${param.UnitName}` }}</span>
              <span v-if="param.NumberMin != undefined || param.NumberMax != undefined">
                {{ `, Range: ${param.NumberMin != undefined ? param.NumberMin : ""}` }}
                {{ ` ~ ` }}
                {{ `${param.NumberMax != undefined ? param.NumberMax : " "}` }}
              </span>
              <!-- <span v-if="param.TextMin != undefined || param.TextMax != undefined">
                {{ `, Length: ${param.TextMin != undefined ? param.TextMin : ""}` }}
                {{ ` ~ ` }}
                {{ `${param.TextMax != undefined ? param.TextMax : " "}` }}
              </span> -->
              <p v-if="!auth.superTestDisabled && !appStore.debugTooltipDisabled">
                {{ data.OID + ": " + param.Name + ", " + param.Access + ", " + param.DataType }}<br>
                {{ param.Groups[0] || '' }}
              </p>
            </template>
            <template v-if="data?.Style?.hidden" />
            <template v-else-if="data?.Style?.readonly">
              <el-input v-model="formatValue" key="input-readonly" class="input-readonly" :style="paramStyle"
                autocomplete="new-password" readonly>
                <template v-if="param.UnitName" #append>
                  <span v-if="param.UnitName" style="margin-left: 8px;">{{ param.UnitName }}</span>
                </template>
              </el-input>
            </template>
            <div v-else-if="data?.Style?.input === 'button'">
              <my-button key="input-button" :owner="owner" :data="data" />
            </div>
            <template v-else>
              <el-switch v-if="data?.Style?.input === 'switch'" key="input-switch" :style="paramStyle"
                v-model="param.InputValue" inline-prompt :active-text="switchData.activeText"
                :active-value="switchData.activeValue" :inactive-text="switchData.inactiveText"
                :inactive-value="switchData.inactiveValue" :before-change="beforeInputChange"
                :disabled="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue" />
              <template v-else-if="data?.Style?.input === 'status' && data?.Style?.status === 'alarm'">
                <el-text v-if="paramValue === '00'" key="status-alarm-00" type="success">
                  <el-icon>
                    <CircleCheck />
                  </el-icon>
                  <span style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
                    {{ formatValue }}
                  </span>
                </el-text>
                <el-text v-else-if="paramValue === '01'" key="status-alarm-01" type="danger">
                  <el-icon>
                    <Warning />
                  </el-icon>
                  <span style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
                    {{ formatValue }}
                  </span>
                </el-text>
                <el-text v-else style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
                  {{ formatValue }}
                </el-text>
              </template>
              <template v-else-if="data?.Style?.input === 'status' && data?.Style?.status === 'sync'">
                <el-text v-if="paramValue === '01'" key="status-sync-01" type="danger">
                  <el-icon>
                    <Warning />
                  </el-icon>
                  <span style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
                    {{ formatValue }}
                  </span>
                </el-text>
                <el-text v-else-if="paramValue === '00'" key="status-sync-00" type="success">
                  <el-icon>
                    <CircleCheck />
                  </el-icon>
                  <span style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all; padding-left: 4px;">
                    {{ formatValue }}
                  </span>
                </el-text>
                <el-text v-else style="overflow: hidden;text-overflow: ellipsis; word-break: keep-all;">
                  {{ formatValue }}
                </el-text>
              </template>
              <el-tree-select v-else-if="data?.Style?.input === 'treeSelect'" key="treeSelect"
                v-model="param.InputValue" :data="treeSelectData" :render-after-expand="false" style="width: 240px"
                :multiple="param.MultipleOption" show-checkbox />
              <el-select v-else-if="data?.Style?.input === 'select'" key="select" :style="paramStyle"
                style="width:240px;" v-model="param.InputValue" :multiple="param.MultipleOption"
                :placeholder="param.MultipleOption ? 'NULL' : undefined"
                :disabled="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue">
                <el-option v-for="opt in param.SortOptions" :key="opt.k" :label="opt.v" :value="opt.k" />
              </el-select>
              <el-radio-group v-else-if="data?.Style?.input === 'radio'" key="radio" v-model="param.InputValue"
                :style="paramStyle">
                <el-radio-button v-for="opt in param.SortOptions" :key="opt.k" :label="opt.v" :value="opt.k" />
              </el-radio-group>
              <el-button-group v-else-if="data?.Style?.input === 'buttonGroup'">
                <my-button v-for="item in data.Items" :key="item.key" :owner="owner" :data="item" />
              </el-button-group>
              <template v-else-if="data?.Style?.input === 'number'">
                <div>
                  <el-input-number key="input-number" v-model="param.InputValue" :style="paramStyle"
                    :step="param.NumberStep" :max="param.NumberMax" :min="param.NumberMin" controls-position="right"
                    step-strictly :readonly="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue"
                    :controls="viewMode == provideKeys.viewModeDefaultValue" />
                  <span v-if="param.UnitName" style="margin-left: 8px;">{{ param.UnitName }}</span>
                  <span v-if="param.Min || param.Max" style="margin-left: 8px;">
                    [ {{ param.NumberMin }} ~ {{ param.NumberMax }} ]
                  </span>
                </div>
              </template>
              <template v-else-if="data?.Style?.input === 'password'">
                <el-input v-model="param.InputValue" key="input-password" :style="paramStyle"
                  autocomplete="new-password" :minlength="param.TextMin" :maxlength="param.TextMax" show-word-limit
                  :show-password="viewMode == provideKeys.viewModeDefaultValue"
                  :readonly="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue" />
              </template>
              <template v-else-if="data?.Style?.input === 'binary'">
                <el-input key="input-string" v-model="param.InputValue" :style="paramStyle" autocomplete="new-password"
                  :minlength="param.TextMin" :maxlength="param.TextMax" show-word-limit
                  :readonly="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue">
                  <template v-if="param.UnitName" #append>
                    {{ param.UnitName }}
                  </template>
                </el-input>
              </template>
              <template v-else-if="data?.Style?.input === 'ipv4'">
                <el-input key="input-ipv4addr" v-model="param.InputValue" :style="paramStyle"
                  autocomplete="new-password" show-word-limit
                  :readonly="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue" />
              </template>
              <template v-else-if="data?.Style?.input === 'datetime'">
                <el-input key="input-datetime" v-model="param.InputValue" :style="paramStyle"
                  autocomplete="new-password" show-word-limit
                  :readonly="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue">
                  <template v-if="(data.OID === 'T02.P0150') && (viewMode == provideKeys.viewModeDefaultValue)" #append>
                    <el-button type="primary" plain
                      @click="param.InputValue = dayjs().format('YYYY-MM-DD HH:mm:ss')">Now</el-button>
                  </template>
                </el-input>
              </template>
              <template v-else-if="data?.OID === 'T02.P0F3E'">
                <el-input key="input-tdd-slot" v-model="param.InputValue" :style="paramStyle"
                  autocomplete="new-password" :minlength="param.TextMin" :maxlength="param.TextMax" show-word-limit
                  :readonly="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue">
                  <template v-if="param.UnitName" #append>
                    {{ param.UnitName }}
                  </template>
                </el-input>
              </template>
              <template v-else>
                <el-input key="input-default" v-model="param.InputValue"
                  :class="data?.InputDisabled ? 'input-readonly' : undefined" :style="paramStyle"
                  autocomplete="new-password" :minlength="param.TextMin" :maxlength="param.TextMax" show-word-limit
                  :readonly="data?.InputDisabled || viewMode == provideKeys.viewModePrintValue">
                  <template v-if="param.UnitName" #append>
                    {{ param.UnitName }}
                  </template>
                </el-input>
              </template>
            </template>
          </el-tooltip>
        </div>
      </div>
      <div v-if="param && param.RespMsg" style="position: absolute; top: 22px;">
        <el-text v-if="param.RespMsg.error" type="danger">{{ param.RespMsg.error }}</el-text>
        <el-text v-else-if="param.RespMsg.warning" type="warning">{{ param.RespMsg.warning }}</el-text>
        <el-text v-else-if="param.RespMsg.info" type="info">{{ param.RespMsg.info }}</el-text>
      </div>

      <div class="edit_buttons">
        <el-button text @click="handleCancel()" style="background-color: #091e4221;">
          <el-icon>
            <CloseBold />
          </el-icon>
        </el-button>
        <el-button text @click="handleSubmit()" style="background-color: #091e4221;">
          <el-icon><Select /></el-icon>
        </el-button>
      </div>
    </div>
  </div>
</template>
<script>
import { useDasDevices } from '@/stores/das-devices'
import { useAuthStore } from "@/stores/auth";
import provideKeys from '@/utils/provideKeys.js'
import { dayjs } from "element-plus";
import { ref } from 'vue';
import { useAppStore } from '@/stores/app';

export default {
  name: 'editableInput',
  inject: ['viewMode'],
  props: {
    subID: String,
    value: String,
    rwValue: {
      default: true,
      type: Boolean
    },
    owner: Object,
    data: Object,
  },
  setup() {
    const dasDevices = useDasDevices()
    const auth = useAuthStore();
    const appStore = useAppStore();
    const dev = dasDevices.currentDevice;
    const tooltip_body = ref();
    let eventDivRef = ref();
    return {
      dasDevices,
      auth,
      appStore,
      dev,
      provideKeys,
      dayjs,
      tooltip_body,
      eventDivRef
    };
  },
  data() {
    let param = undefined;
    if (this.subID) {
      param = this.dasDevices.getDevice(this.subID)?.params.getParam(this.data?.OID)
    } else {
      param = this.dev.params.getParam(this.data?.OID);
    }
    return {
      editable: false,
      inputModel: this.data?.Value ?? this.value,
      param
    }
  },
  computed: {
    hasTips() {
      if (this.param !== undefined) {
        if (this.param.Tips != undefined) {
          return true;
        }
        if (this.param.Value != undefined) {
          return true;
        }
        if (this.param.UnitName != undefined) {
          return true;
        }
        if (this.param.NumberMax != undefined || this.param.NumberMin != undefined) {
          return true;
        }
        if (this.param.TextMax != undefined || this.param.TextMin != undefined) {
          return true;
        }
        if (!this.auth.superTestDisabled && !this.appStore.debugTooltipDisabled) {
          return true;
        }
      }
      return false;
    },
    paramStyle() {
      return this.param?.getShowStyle();
    },
    paramValue() {
      return this.param?.getValue({ defaultValue: this.inputModel });
    },
    formatValue() {
      return this.param?.getShowValue({ defaultValue: this.inputModel });
    },
    formatValueWithUnit() {
      return this.param?.getShowValue({ defaultValue: this.inputModel, withUnit: true });
    },
    switchData() {
      return this.param?.getSwitchData(this.data?.Style);
    },
    treeSelectData() {
      return this.param?.getTreeSelectData();
    },
    isAllowOpenEdit() {
      if (this.viewMode == provideKeys.viewModePrintValue) {
        return false;
      }
      if (!this.subID) {
        if (this.data?.Style?.hidden || this.data?.Style?.input === 'button' || this.data?.Style?.input === 'buttonGroup') {
          return false;
        }
        if (this.rwValue == false) return false;
        return true;
      }
      return this.rwValue
    }
  },
  watch: {
    value() {
      this.inputModel = this.value;
    }
  },
  methods: {
    handleCancel() {
      if (this.subID) {
        this.inputModel = this.data?.Value || this.value;
      } else {
        this.param.resetInputValue();
      }
      this.editable = false;
    },
    handleSubmit() {
      if (this.subID) {
        this.submitDeviceParams()
      }
    },
    async submitDeviceParams() {
      try {
        var values = [];
        values.push({ oid: this.param.PrivOid, value: this.inputModel });
        let res = await this.dasDevices.setDeviceParameterValues({
          sub: this.subID,
          values,
          showMessage: true,
        });
        if (res && res.length > 0) {
          this.editable = false;
        }
      } catch (e) {
        console.error(e)
      }
    },
    submitCurrentDeviceParams() {
    },
    beforeInputChange: function () {
      const self = this;
      if (self.data?.Style?.confirmMessage) {
        const resultPromise = new Promise((resolve, reject) => {
          self.appStore.openConfirmDialog({
            title: self.data?.Style?.confirmTitle || 'Confirm',
            content: self.data?.Style?.confirmMessage || '',
            callback: (ok) => {
              resolve(ok);
            }
          })
        });
        return resultPromise;
      }
      return true;
    },
    enterEditMode() {
      this.editable = true;
      this.$nextTick(() => {
        if (this.eventDivRef) {
          this.eventDivRef.focus();
        }
      })
    }
  }
}


</script>

<style scoped lang="scss">
.is_edit {
  .readonly_div {
    display: flex;
    flex-wrap: nowrap;
    flex-direction: row;
    justify-content: space-between;
    border: 1px solid transparent;
    min-height: 24px;
    width: 100%;

    &:hover {
      border: 1px solid var(--el-border-color);
      border-radius: 4px;
      box-shadow: var(--el-box-shadow-lighter);
    }

    &:hover .readonly_button {
      display: block;
    }



    .readonly_button {
      display: none;
      background-color: #091e4221;
      height: auto;

      &:hover {
        display: block;
      }
    }
  }

  .edit_div {
    width: 100%;

    .edit_buttons {
      position: absolute;
      right: 0;
      top: 100%;
      z-index: 1;
      box-shadow: 0 3px 6px rgba(111, 111, 111, 0.2);
      padding: 2px 4px;
      background: #fff;
    }

    & :deep(.el-tooltip__trigger) {
      box-shadow: var(--el-box-shadow-lighter);
    }
  }
}

.readonly_label {
  padding-left: 4px;
  padding-right: 4px;
  display: block;
  max-width: calc(100% - 36px);
}
</style>
