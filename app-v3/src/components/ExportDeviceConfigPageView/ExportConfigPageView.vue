<template>
  <el-dialog v-model="dialogVisible" title="Export Configuraion" :close-on-click-modal="false"
    :close-on-press-escape="false" :before-close="handleCloseDialog" @open="handleOpen" width="75%" top="10vh"
    style="min-height: 80vh;">
    <template v-if="activeStep !== 100">
      <el-steps :active="activeStep" align-center style="margin-bottom: 16px;">
        <el-step title="Select Page" />
        <el-step title="Get Parameter Values" />
        <el-step title="Select Parameters" />
        <el-step title="Result" />
      </el-steps>
      <template v-if="activeStep === 1">
        <el-row justify="center">
          <el-col :span="16">
            <el-scrollbar style="height:calc(100vh - 15vh - 15vh - 54px - 66px - 60px);">
              <el-tree ref="selectTreeRef" :data="pageTreeData" :props="defaultProps" node-key="key"
                :default-expanded-keys="defaultExpandedKeys" show-checkbox @check="handleCheck"
                :filter-node-method="filterNode" />
            </el-scrollbar>
          </el-col>
        </el-row>
        <el-row justify="center">
          <el-button type="primary" plain :disabled="selectNodes.length < 1" @click="handleSelectNext">Next</el-button>
        </el-row>
      </template>
      <template v-if="activeStep === 2">
        <el-row justify="center">
          <el-col :span="16">
            <el-progress :percentage="queryProgress" :status="queryProgress === 100 ? 'success' : ''" />
            <el-scrollbar style="height:calc(80vh - 54px - 66px - 60px - 16px - 36px);" ref="scrollbarRef">
              <template v-for="(item, index) in selectNodes">
                <p v-if="item.status === 'success'" :key="index + 'success'">
                  <el-icon style="margin-right: 32px;">
                    <SuccessFilled color="#67C23A" />
                  </el-icon>
                  {{ item.fullLabel }}
                </p>
                <p v-else-if="item.status === 'error'" :key="index + 'error'">
                  <el-icon style="margin-right: 32px;">
                    <CircleCloseFilled color="#f56c6c" />
                  </el-icon>
                  {{ item.fullLabel }}
                </p>
                <p v-else :key="index + 'wait'">
                  <el-icon style="margin-right: 32px;">
                    <Loading />
                  </el-icon>{{ item.fullLabel }}
                </p>
              </template>
            </el-scrollbar>
          </el-col>
        </el-row>
        <el-row justify="center" style="min-height: 36px;">
          <el-button type="primary" plain v-if="updateParameterIsFail" @click="handleSelectRetry">Retry</el-button>
        </el-row>
      </template>
      <template v-if="activeStep === 3">
        <el-row justify="center">
          <el-col :span="20" style="margin-top: 10px 0;text-align: center;" v-loading="tableLoading">
            <div style="text-align: left;">
              <el-checkbox v-model="selectedAll" label="Select All" style="margin-left: 4px;" />
            </div>
            <el-scrollbar style="height:calc(80vh - 54px - 66px - 60px - 16px - 30px);">
              <el-collapse accordion>
                <template v-for="(pageItem) in paramList ">
                  <template v-for="(paramItem) in pageItem.paramOids ">
                    <el-collapse-item v-if="paramItem.rwOids && (paramItem.rwOids.length > 0)"
                      :key="pageItem.key + paramItem.tabKey" :name="pageItem.key + paramItem.tabKey">
                      <template #title>
                        <el-checkbox @click.stop="" v-model="paramItem.checkAll"
                          :indeterminate="(selected[pageItem.key + paramItem.tabKey]?.length < paramItem.rwOids.length) && (selected[pageItem.key + paramItem.tabKey]?.length > 0)"
                          @change="handleCheckAllChange($event, pageItem.key + paramItem.tabKey)"
                          style="margin-right: 8px;margin-left: 4px;">
                        </el-checkbox>
                        {{ pageItem.fullLabel + ((pageItem.paramOids.length > 1) ? ('/' +
                          paramItem.Name) : '') }}
                        {{ `(${selected[pageItem.key + paramItem.tabKey]?.length ?? 0} /${paramItem.rwOids.length})` }}
                      </template>
                      <template #default>
                        <param-table-view :oids="paramItem.rwOids" :defaultValues="paramItem.defaultValues"
                          :ref="(el) => setParamTableRef(el, pageItem.key + paramItem.tabKey, paramItem)"
                          @change-select="handleParamSelect($event, pageItem.key + paramItem.tabKey)" />
                      </template>
                    </el-collapse-item>
                  </template>
                </template>
              </el-collapse>
            </el-scrollbar>
          </el-col>
        </el-row>
        <el-row justify="center">
          <el-button type="primary" plain :disabled="Object.keys(selected).length == 0" @click="nextToResultStep"
            style="margin: 4px 0;">Next</el-button>
        </el-row>
      </template>
      <template v-if="activeStep === 4">
        <el-row justify="center">
          <el-result icon="success" title="" sub-title="Configuration file is ready">
            <template #extra>
              <el-button type="primary" @click="exportConfigFile">Export</el-button>
            </template>
          </el-result>
        </el-row>
      </template>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, reactive, nextTick, defineProps, defineEmits } from 'vue'
import { useAuthStore } from "@/stores/auth";
import { useDasDevices } from "@/stores/das-devices";
import { sleep } from "@/utils/index.js";
import ParamTableView from "./ParamTableView.vue";
import { toRaw } from 'vue';

const props = defineProps(['configDialogVisible']);

const emit = defineEmits(['update:configDialogVisible', 'closePrintView']);

const dialogVisible = computed({
  get() {
    return props.configDialogVisible
  },
  set(value) {
    emit('update:configDialogVisible', value)
  }
})

let defaultProps = reactive({
  children: 'children',
  label: 'label',
});

let selectTreeRef = ref();
let appStore = useDasDevices();
let dev = useDasDevices().currentDevice;

let defaultExpandedKeys = reactive(["/"]);
let activeStep = ref(0);
let selectNodes = ref([]);
let queryProgress = ref(0);
let queryList = ref([]);

let closeViewFlag = ref(false);

onMounted(() => {
  handleOpen()
  getFileName()
})
onBeforeUnmount(() => {
  closeViewFlag.value = true;
})

function handleOpen() {
  queryList.value = [];
  activeStep.value = 1;
  queryProgress.value = 0;
  nextTick(() => {
    selectTreeRef.value && selectTreeRef.value.filter();
  })
}

let pageTreeData = computed(() => dev.layout.pageTreeData());

function filterNode(value, data) {
  const auth = useAuthStore();
  if (auth.superModeDisabled && data.superModeOnly) {
    return false;
  }
  if (data.key == '/overview/das_topo') return false;
  if (data.page && data.page.rwOids.length == 0 && data.page.woOids.length == 0) return false;
  return true;
}

function handleCheck(v) {
  selectNodes.value = selectTreeRef.value && selectTreeRef.value.getCheckedNodes(true);
}


let scrollbarRef = ref();
let updateParameterIsFail = ref(false);

function handleSelectNext() {
  selectNodes.value = JSON.parse(JSON.stringify(selectNodes.value)).filter(item => filterNode(null, item));
  handleSelectRetry();
}

async function handleSelectRetry() {
  activeStep.value = 2;

  updateParameterIsFail.value = false;
  selectNodes.value.forEach(v => {
    if (v.status !== "success") {
      v.status = "waiting";
    }
  });
  let count = queryList.value.length;
  let len = selectNodes.value.length;
  queryProgress.value = Number((count * 100 / len).toFixed(0)) || 0;
  for (let i = queryList.value.length; i < len; i++) {
    if (closeViewFlag.value) {
      return
    }
    const nodeItem = selectNodes.value[i];
    let res = true;
    if (nodeItem.page && nodeItem.page.rwOids && Object.keys(nodeItem.page.rwOids).length > 0) {
      let resultValues = await dev.params.getParameterValues({ oids: nodeItem.page.rwOids });
      if (resultValues.length == 0) {
        res = false
      }
    }
    if ((process.env.NODE_ENV === 'development') || res) {
      nodeItem.status = "success";
      queryList.value.push(nodeItem);
      count += 1;
      queryProgress.value = Number((count * 100 / len).toFixed(0)) || 0;
      scrollbarRef.value && scrollbarRef.value.setScrollTop((count - 4) * 31)
    } else {
      nodeItem.status = "error";
      updateParameterIsFail.value = true;
      return
    }
  }
  await sleep(1000);
  nextToParameterStep()
}

function nextToParameterStep() {
  if (closeViewFlag.value) {
    return
  }
  loadParamTables()
}

function loadParamTables() {
  nextTick(() => {
    tableLoading.value = true
    activeStep.value = 3
  })
  setTimeout(async () => {
    await nextTick(() => {
      paramList.value = queryList.value.map(item => {
        let paramOids = getPageTableOrFormItems(item.page.Items);
        return {
          ...item,
          paramOids
        }
      });
    });
    await nextTick(() => {
      selectAllParamTable(true);
    })
    setTimeout(() => {
      tableLoading.value = false;
    }, 200);
  }, 1000);
}

let paramList = ref([]);
let tableLoading = ref(false);

let fileName = ref('');
function getFileName() {
  const dasDevices = useDasDevices();
  fileName.value = dasDevices.currentDeviceFileName('Configuration');
}

function closeView() {
  emit('closePrintView')
}

function handleCloseDialog(close) {
  close()
}

function nextToResultStep() {
  activeStep.value = 4
}

async function exportConfigFile() {
  let paramBody = filterSelectParam()
  const exportData = await dev.cfg.exportConfigurationData(paramBody)
  let exportBody = JSON.stringify(exportData);
  let blob = new Blob([exportBody], { type: "text/plain;charset=utf-8" });
  let url = URL.createObjectURL(blob);
  let link = document.createElement('a');
  link.href = url;
  link.download = fileName.value + '.json';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

/**
 * body:[{
     *    path : [...key],
     *    data: [...{oid:'', value}]
     * }]
 */
function filterSelectParam() {
  let keys = Object.keys(selected.value).filter(item => selected.value[item].length > 0);
  let keysString = keys.join(',')
  let cpParamList = toRaw(paramList.value).filter(item => keysString.indexOf(item.key) > -1);
  let defaultValues = {};
  cpParamList.forEach(item => {
    item.page.defaultValues.forEach(defItem => {
      defaultValues[defItem.oid] = defItem['value']
    });
  })
  return keys.map(item => {
    let keysInfo = item.split('/')
    let info = {
      path: [keysInfo[1], keysInfo[2]],
      data: selected.value[item].map(OidItem => {
        return {
          oid: OidItem.oid,
          value: defaultValues[OidItem.oid] || undefined
        }
      })
    }
    if (keysInfo[4]) {
      info.path.push(keysInfo[3]);
    }
    return info
  })
}

//----------------------------------//
function getPageTableOrFormItems(items, key = '', tabName = '') {
  if (!items || items.length == 0) {
    return []
  }
  if (items[0].Type == "Table" || items[0].Type == 'Form') {
    return items.map(item => {
      let rwOids = item.rwOids;
      if (item.woOids && Object.keys(item.woOids).length > 0) {
        rwOids = item.rwOids.concat(item.woOids);
      }
      return {
        ...item,
        rwOids: rwOids,
        tabKey: key + '/' + item.Key,
        tabName: tabName,
        checkAll: false
      }
    })
  }
  let result = []

  items.forEach(item => {
    let cpKey = key;
    let cpTabName = tabName
    if (item.Key) {
      cpKey += ('/' + item.Key);
      cpTabName += ('/' + item.Name);
    }
    let tableItems = getPageTableOrFormItems(item.Items, cpKey, cpTabName);
    result.push(...tableItems)
  })
  return result
}

let selected = ref({});
function handleParamSelect(selection, key) {
  if (!selected.value[key]) {
    selected.value[key] = []
  }
  selected.value[key] = selection

  if (selection.length == paramTableRefs.value[key].paramItem.rwOids.length) {
    paramTableRefs.value[key].paramItem.checkAll = true
  } else if (selection.length == 0) {
    paramTableRefs.value[key].paramItem.checkAll = false
  }
}


//---------ParamTable  All Select ---------------
let paramTableRefs = ref({});
function setParamTableRef(el, index, paramItem) {
  if (paramTableRefs[index]) return;
  paramTableRefs.value[index] = { el, paramItem }
}

function handleCheckAllChange(val, key) {
  let tableRef = paramTableRefs.value[key].el.tableRef;
  if (val) {
    tableRef.toggleAllSelection();
  } else {
    tableRef.clearSelection();
  }
}

let selectedAll = computed({
  get() {
    return Object.values(paramTableRefs.value).every((item) => item.paramItem.checkAll == true)
  },
  set(val) {
    selectAllParamTable(val);
  }
});

function selectAllParamTable(isSelected) {
  for (let key in toRaw(paramTableRefs.value)) {
    let el = paramTableRefs.value[key].el
    if (isSelected) {
      if (!paramTableRefs.value[key].paramItem.checkAll) el.tableRef && el.tableRef.toggleAllSelection();
    } else {
      el.tableRef && el.tableRef.clearSelection();
    }
  }
}
</script>

<style scoped lang="scss"></style>