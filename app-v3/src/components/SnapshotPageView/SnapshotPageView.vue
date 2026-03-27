<template>
  <el-dialog v-model="dialogVisible" title="Snapshot" @open="handleOpen" :close-on-click-modal="false"
    :close-on-press-escape="false" :before-close="handleCloseDialog" :fullscreen="activeStep == 100" id="SnapshotDialog">
    <template v-if="activeStep !== 100">
      <el-steps :active="activeStep" align-center style="margin-bottom: 16px;">
        <el-step title="Select Page" />
        <el-step title="Get Parameter Values" />
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
          <el-button type="primary" plain @click="handleSelectNext" :disabled="selectNodes.length < 1">Next</el-button>
        </el-row>
      </template>
      <template v-if="activeStep === 2">
        <el-row justify="center">
          <el-col :span="16">
            <el-progress :percentage="queryProgress" :status="queryProgress === 100 ? 'success' : ''" />
            <el-scrollbar style="height:calc(100vh - 15vh - 15vh - 54px - 66px - 60px - 16px - 36px);" ref="scrollbarRef">
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
          <el-result icon="success" title="" sub-title="Get parameter values successfully, result is ready">
            <template #extra>
              <el-button type="primary" plain @click="handleView">View</el-button>
            </template>
          </el-result>
        </el-row>
      </template>
    </template>
    <template v-else>
      <el-row>
        <el-col style="display: flex;justify-content: space-between;">
          <span class="el-dialog__title">{{ fileName }}</span>
          <div class="toolbar">
            <el-affix :offset="10">
              <el-button type="primary" @click="handleSave" :disabled="printLoading || saveLoading">Save</el-button>
            </el-affix>
            <el-backtop :right="20" :bottom="40" :visibility-height="10" target="#SnapshotDialog"/>
          </div>
        </el-col>
      </el-row>
      <el-row id="SnapshotDocument" v-loading="printLoading">
        <el-col v-for="(pageItem, index) in printList" style="margin: 10px 0;" :key="index" class="page-col">
          <div class="print-divider-title">
              <el-divider content-position="left">
                <h3> {{ pageItem.fullLabel }}</h3>
              </el-divider>
          </div>
          <my-main-page :propPage="pageItem.page" :mode="provideKeys.viewModePrintValue" />
        </el-col>
      </el-row>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, reactive, nextTick, defineProps, defineEmits } from 'vue'
import { useDasDevices } from "@/stores/das-devices";
import printHtml from "@/utils/print.js";
import { sleep } from "@/utils/index.js";
import provideKeys from '@/utils/provideKeys.js'
import { useAuthStore } from '@/stores/auth';

const props = defineProps(['snapshotDialogVisible']);

const emit = defineEmits(['update:snapshotDialogVisible', 'closePrintView']);

const dialogVisible = computed({
  get() {
    return props.snapshotDialogVisible
  },
  set(value) {
    emit('update:snapshotDialogVisible', value)
  }
})

let defaultProps = reactive({
  children: 'children',
  label: 'label',
});

let selectTreeRef = ref();
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
  printList.value = [];
  queryProgress.value = 0;
  nextTick(() => {
    selectTreeRef.value.filter();
  })

}

let pageTreeData = computed(() => dev.layout.pageTreeData());

function handleCheck(v) {
  selectNodes.value = selectTreeRef.value && selectTreeRef.value.getCheckedNodes(true);
}

function handleSelectNext() {
  selectNodes.value = JSON.parse(JSON.stringify(selectNodes.value)).filter(item => filterNode(null, item));
  handleSelectRetry();
}

let scrollbarRef = ref();
let updateParameterIsFail = ref(false);
async function handleSelectRetry() {
  // console.log('selectNodes', selectNodes.value);
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
    // console.log('nodeItem', nodeItem)
    let res = true;
    if (nodeItem.page && nodeItem.page.oids && Object.keys(nodeItem.page.oids).length > 0) {
      let resultValues = await dev.params.getParameterValues({
        oids: nodeItem.page.oids,
      });
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
  nextToResultStep()
}

function nextToResultStep() {
  if (closeViewFlag.value) {
    return
  }
  activeStep.value = 3
}

let printList = ref([])
let printLoading = ref(false);
function handleView() {
  nextTick(() => {
    activeStep.value = 100;
    printLoading.value = true
  })
  setTimeout(async () => {
    await nextTick(() => {
      printList.value = queryList.value;
    })
    setTimeout(() => {
      printLoading.value = false
    }, 1000);
  }, 1000);

}

let fileName = ref('');
function getFileName() {
  const dasDevices = useDasDevices();
  fileName.value = dasDevices.currentDeviceFileName('Snapshot');
}

let saveLoading = ref(false);
async function handleSave() {
  await nextTick(()=>{
    saveLoading.value = true;
  })

  //iDAS_A402_Screenshot_123pa3_2024_1_22
  const doc = document.getElementById("SnapshotDocument");
  printHtml(doc, fileName.value, true);
  setTimeout(() => {
    nextTick(()=>{
      saveLoading.value = false;
    })
  }, 500);
}

function closeView() {
  emit('closePrintView')
}

function handleCloseDialog(close) {
  close()
}

function filterNode(value, data) {
  const auth = useAuthStore();
  if (auth.superModeDisabled && data.superModeOnly) {
    return false;
  }
  if (data.key == '/overview/das_topo') return true;
  if (data.key == '/maintenance/firmware_information') return true;
  if (data.key == '/system_settings/stats') return true;
  if (data.page && data.page.rOids.length == 0) return false;
  return true;
}

</script>

<style scoped lang="scss"></style>

