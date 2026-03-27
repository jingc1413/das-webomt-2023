<template>
  <el-dialog v-model="dialogVisible" :title="'Load Configuraion - ' + fileName" :close-on-click-modal="false"
    :close-on-press-escape="false" :before-close="handleCloseDialog" @open="handleOpen" width="75%" top="10vh"
    style="min-height: 80vh;">
    <template v-if="activeStep !== 100">
      <el-steps :active="activeStep" align-center style="margin-bottom: 16px;">
        <el-step title="Get Parameter Values" />
        <el-step title="Select Parameters" />
        <el-step title="Set Parameter Values" />
        <!-- <el-step title="Result" /> -->
      </el-steps>
      <template v-if="activeStep === 1">
        <el-row justify="center" v-loading="initLoading">
          <el-col :span="16">
            <el-progress v-if="!initSelectNodeIsFail" :percentage="queryProgress"
              :status="queryProgress === 100 ? 'success' : ''" />
            <el-scrollbar style="height:calc(80vh - 54px - 66px - 60px - 16px - 36px);" ref="scrollbarGetRef">
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
          <el-button type="primary" plain v-if="initSelectNodeIsFail" @click="initSelectNodes">Retry</el-button>
        </el-row>
      </template>
      <template v-if="activeStep === 2">
        <el-row justify="center">
          <el-col :span="20" style="margin-top: 10px 0;text-align: center;" v-loading="tableLoading">
            <div style="text-align: left;">
              <el-checkbox v-model="selectedAll" label="Select All" style="margin-left: 4px;" />
            </div>
            <el-scrollbar style="height:calc(80vh - 54px - 66px - 60px - 16px);">
              <el-collapse accordion>
                <template v-for="(pageItem) in paramList">
                  <el-collapse-item v-if="pageItem.paramOids.length > 0" :key="pageItem.key" :name="pageItem.fullLabel">
                    <template #title>
                      <el-checkbox @click.stop="" v-model="pageItem.checkAll"
                        :indeterminate="(selected[pageItem.key]?.data?.length < pageItem.paramOids.length) && (selected[pageItem.key]?.data?.length > 0)"
                        @change="handleCheckAllChange($event, pageItem.key)"
                        style="margin-right: 8px;margin-left: 4px;">
                      </el-checkbox>
                      {{ pageItem.fullLabel }}
                      {{ `(${selected[pageItem.key]?.data?.length ?? 0} /${pageItem.paramOids.length})` }}
                    </template>
                    <template #default>
                      <param-table-view :oids="pageItem.paramOids"
                        :ref="(el) => setParamTableRef(el, pageItem.key, pageItem)"
                        @change-select="handleParamSelect($event, pageItem.key, pageItem)" />
                    </template>
                  </el-collapse-item>
                </template>
              </el-collapse>
            </el-scrollbar>
          </el-col>
        </el-row>
        <el-row justify="center">
          <el-button type="primary" plain :disabled="Object.keys(selected).length == 0" @click="handleSetParamNext"
            style="margin: 4px 0;">Next</el-button>
        </el-row>
      </template>
      <template v-if="activeStep === 3">
        <el-row justify="center">
          <el-col :span="16">
            <el-progress :percentage="setProgress" :status="setProgress === 100 ? 'success' : ''" />
            <el-scrollbar style="height:calc(80vh - 54px - 66px - 60px - 16px - 36px);" ref="scrollbarSetRef">
              <template v-for="(item, index) in setNodes">
                <p v-if="item.status === 'success'" :key="index + 'success'">
                  <el-icon style="margin-right: 32px;">
                    <SuccessFilled color="#67C23A" />
                  </el-icon>
                  {{ item.name }}
                </p>
                <p v-else-if="item.status === 'error'" :key="index + 'error'">
                  <el-icon style="margin-right: 32px;">
                    <CircleCloseFilled color="#f56c6c" />
                  </el-icon>
                  {{ item.name }}
                </p>
                <p v-else-if="item.status === 'warring'" :key="index + 'warring'">
                  <el-icon style="margin-right: 32px;">
                    <InfoFilled color="#E6A23C" />
                  </el-icon>
                  {{ item.name }}
                </p>
                <p v-else :key="index + 'wait'">
                  <el-icon style="margin-right: 32px;">
                    <Loading />
                  </el-icon>{{ item.name }}
                </p>
              </template>
            </el-scrollbar>
          </el-col>
        </el-row>
        <el-row justify="center" style="min-height: 36px;">
          <el-button type="primary" plain v-if="setParameterIsFail" @click="handleSetRetry">Retry</el-button>
          <el-button type="primary" plain v-if="setAllParamSuccess" @click="dialogVisible = false">Close</el-button>
        </el-row>
      </template>
      <template v-if="activeStep === 4">
        <el-row justify="center">
          <el-result :icon="setParameterIsWarring ? 'warning' : 'success'" title=""
            :sub-title="`Load configuration file successfully ${(setParameterIsWarring ? $t('tip.ResponseWithFaultCode') : '')}`" />
        </el-row>
      </template>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, reactive, nextTick, defineProps, defineEmits } from 'vue'
import { sleep } from "@/utils/index.js";
import ParamTableView from "./ConfigParamTableView.vue";
import { ElMessage } from "element-plus";
import { toRaw } from 'vue';
import { useI18n } from 'vue-i18n'
import { useDasDevices } from '@/stores/das-devices';

// import mock from "./moke2.json";

const { t } = useI18n();

const props = defineProps(['configDialogVisible', 'fileName']);

const emit = defineEmits(['update:configDialogVisible', 'closePrintView']);

const dialogVisible = computed({
  get() {
    return props.configDialogVisible
  },
  set(value) {
    emit('update:configDialogVisible', value)
  }
})

let dev = useDasDevices().currentDevice;
let activeStep = ref(0);
let selectNodes = ref([]);
let queryProgress = ref(0);
let queryList = ref([]);

let closeViewFlag = ref(false);

let initLoading = ref(false);

onMounted(() => {
  handleOpen()
})
onBeforeUnmount(() => {
  closeViewFlag.value = true;
})

let pageTreeData = computed(() => (dev.layout?.pageTreeData())[0].children);

function handleOpen() {
  nextTick(() => {
    queryList.value = [];
    activeStep.value = 1;
    queryProgress.value = 0;
    initLoading.value = true;
  })
  setTimeout(() => {
    nextTick(() => {
      initSelectNodes();
    })
  }, 500);

}

let configInfoTree = {};

let initSelectNodeIsFail = ref(false);

async function initSelectNodes() {
  initLoading.value = true;
  let content = await dev.files.getFile('ConfigFile', props.fileName);
  if (!content) {
    ElMessage.error(t("tip.LoadConfigurationFailed"));
    initSelectNodeIsFail.value = true;
    initLoading.value = false;
    return
  }
  initSelectNodeIsFail.value = false;
  let configInfo;
  try {
    configInfo = dev.cfg.importConfigurationData(content);
  } catch (error) {
    console.error(error);
    ElMessage.error(t("tip.LoadConfigurationFailed"));
    initSelectNodeIsFail.value = true;
    initLoading.value = false;
    return
  }
  let moduleKeys = configInfo.map(item => ('/' + item.path[0]));
  configInfoTree = {};
  configInfo.forEach(item => {
    let moduleKey = '/' + item.path[0];
    let pageKey = '/' + item.path[0] + '/' + item.path[1];
    if (!configInfoTree[moduleKey]) {
      configInfoTree[moduleKey] = {}
    }
    if (!configInfoTree[moduleKey][pageKey]) {
      configInfoTree[moduleKey][pageKey] = {}
    }
    if (item.path[2]) {
      configInfoTree[moduleKey][pageKey][item.path[2]] = { 'data': item.data }
    } else {
      configInfoTree[moduleKey][pageKey] = { 'data': item.data }
    }
  })
  let modules = JSON.parse(JSON.stringify(pageTreeData.value.filter(item => moduleKeys.includes(item.key))));
  let pages = [];
  modules.forEach(moduleItem => {
    let moduleKey = moduleItem.key;
    moduleItem.children = moduleItem.children.filter(pageItem => (pageItem.key in configInfoTree[moduleKey]));
    let filterPage = moduleItem.children.map(pageItem => {
      if (pageItem.children) {
        pageItem.children = pageItem.children.filter(tabItem => (tabItem.page.Key in configInfoTree[moduleKey][pageItem.key]));
      }
      return pageItem
    })
    pages.push(...filterPage)
  })
  pages.forEach(item => {
    if (item.children) {
      selectNodes.value.push(...item.children)
    } else {
      selectNodes.value.push(item)
    }
  })
  handleSelectNext();
}

let scrollbarGetRef = ref();
let scrollbarSetRef = ref();
let updateParameterIsFail = ref(false);

function handleSelectNext() {
  initLoading.value = false;
  handleSelectRetry();
}

async function handleSelectRetry() {
  activeStep.value = 1;

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
    // let oids = nodeItem.data?.map(item => item.oid);
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
      scrollbarGetRef.value && scrollbarGetRef.value.setScrollTop((count - 4) * 31)
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
    activeStep.value = 2
  })
  setTimeout(async () => {
    await nextTick(() => {
      paramList.value = queryList.value.map(queryItem => {
        let tempKeys = queryItem.key.split('/');
        let moduleKey = '/' + tempKeys[1];
        let pageKey = moduleKey + '/' + tempKeys[2];
        let tabsInfo = configInfoTree[moduleKey][pageKey];
        let rwOids = []
        if (tempKeys[3]) {
          rwOids = tabsInfo[tempKeys[3]].data
        } else {
          rwOids = tabsInfo.data
        }
        let paramOids = getArrayIntersection(rwOids, queryItem.page.rwOids, queryItem.page.woOids)
        return {
          ...queryItem,
          paramOids,
          status: 'waiting',
          checkAll: false
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

function getArrayIntersection(rwOidsWithValue, rwOids = [], woOids = []) {
  const set1 = new Set(rwOidsWithValue);
  const set2 = new Set(rwOids.concat(woOids));
  const intersection = [];
  for (let item of set1) {
    if (set2.has(item.oid)) {
      intersection.push(item);
    }
  }
  return intersection.map(item => {
    let isChanged = false;
    let param = dev.params.getParam(item.oid);
    if (String(item.value) !== String(param.Value)) isChanged = true;
    let isValidate = true
    if (param.validateInputValue(item.value) !== undefined) {
      isValidate = false
    }
    return {
      ...item,
      isValidate,
      isChanged,
    }
  });
}

let paramList = ref([]);
let tableLoading = ref(false);

function closeView() {
  emit('closePrintView')
}

function handleCloseDialog(close) {
  close()
}

function handleSetParamNext() {
  let setParams = Object.values(selected.value).filter((item) => item.data.length > 0);
  if (setParams.length == 0) {
    return ElMessage.warning(t("tip.selectOne"))
  }
  setNodes.value = setParams;
  setNodes.value.forEach(v => {
    v.status = "waiting";
  });
  queryList.value = [];
  activeStep.value = 3;
  setProgress.value = 0;
  handleSetRetry();
}

let setParameterIsFail = ref(false);
let setParameterIsWarring = ref(false);
let setProgress = ref(0);
let setNodes = ref([])
async function handleSetRetry() {
  activeStep.value = 3;
  setParameterIsFail.value = false;
  setNodes.value.forEach(v => {
    if (v.status !== "success" && v.status !== "warring") {
      v.status = "waiting";
    }
  });
  let count = queryList.value.length;
  let len = setNodes.value.length;
  setProgress.value = Number((count * 100 / len).toFixed(0)) || 0;
  for (let i = queryList.value.length; i < len; i++) {
    if (closeViewFlag.value) {
      return
    }
    const nodeItem = setNodes.value[i];
    let values = nodeItem.data?.map(item => { return { oid: item.oid, value: item.value } });
    let res = true;
    let faultValues = undefined;
    if (nodeItem.data && nodeItem.data.length > 0) {
      let resultValues = await dev.params.setParameterValues({ values });
      if (resultValues.length == 0) {
        res = false
      } else {
        faultValues = resultValues.find(v => v.code !== "00");
        faultValues && (setParameterIsWarring.value = true);
      }
    }
    if ((process.env.NODE_ENV === 'development') || res) {
      nodeItem.status = faultValues ? "warring" : "success";
      queryList.value.push(nodeItem);
      count += 1;
      setProgress.value = Number((count * 100 / len).toFixed(0)) || 0;
      scrollbarSetRef.value && scrollbarSetRef.value.setScrollTop((count - 4) * 31)
    } else {
      nodeItem.status = "error";
      setParameterIsFail.value = true;
      return
    }
  }
  await sleep(1000);
  nextToResultStep()

}

let setAllParamSuccess = ref(false);
function nextToResultStep() {
  setAllParamSuccess.value = true;
  // activeStep.value = 4;

}

//----------------------------------//

let selected = ref({});
function handleParamSelect(selection, key, pageItem) {
  if (!selected.value[key]) {
    selected.value[key] = {
      data: [],
      name: pageItem.fullLabel
    }
  }
  selected.value[key].data = selection
  let checkAllLength = pageItem.paramOids.filter(item => (!item.isFail && item.isValidate)).length;
  if (selection.length == checkAllLength) {
    pageItem.checkAll = true
  } else if (selection.length == 0) {
    pageItem.checkAll = false
  }
}

//---------ParamTable  All Select ---------------
let paramTableRefs = ref({});
function setParamTableRef(el, index, pageItem) {
  if (paramTableRefs[index]) return;
  paramTableRefs.value[index] = { el, pageItem }
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
    return Object.values(paramTableRefs.value).every((item) => item.pageItem.checkAll == true)
  },
  set(val) {
    selectAllParamTable(val);
  }
});

function selectAllParamTable(isSelected) {
  for (let key in toRaw(paramTableRefs.value)) {
    let el = paramTableRefs.value[key].el
    if (isSelected) {
      if (!paramTableRefs.value[key].pageItem.checkAll) el.tableRef && el.tableRef.toggleAllSelection();
    } else {
      el.tableRef && el.tableRef.clearSelection();
    }
  }
}

</script>

<style scoped lang="scss"></style>