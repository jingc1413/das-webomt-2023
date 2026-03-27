<template>
  <el-dialog v-model="dialogVisible" :title="props.title" @open="handleOpen" :close-on-click-modal="false"
    :close-on-press-escape="false" :before-close="handleCloseDialog">
    <el-row justify="center">
      <el-col :span="16">
        <el-progress :percentage="queryProgress" :status="queryProgress === 100 ? 'success' : ''" />
        <el-scrollbar style="height:calc(100vh - 15vh - 15vh - 54px - 66px - 60px - 16px);">
          <template v-for="(item, index) in queryList">
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
            <p v-else :key="index + 'wait'">
              <el-icon style="margin-right: 32px;">
                <Loading />
              </el-icon>
              {{ item.name }}
            </p>
          </template>
        </el-scrollbar>

      </el-col>
    </el-row>
    <el-row justify="center">
      <el-button type="primary" plain v-if="actionIsFail" @click="handleActionRetry">{{ $t("button.Retry") }}</el-button>
    </el-row>
  </el-dialog>
</template>

<script setup>
import { ElMessage } from "element-plus";
import { ref, onMounted, onBeforeUnmount, computed, reactive, nextTick, defineProps, defineEmits } from 'vue'
import { sleep } from "@/utils/index.js";

import { useI18n } from 'vue-i18n'
const { t } = useI18n();

const props = defineProps({
  batchDialogVisible: Boolean,
  batchActionCallBack: Function, //async  return Promise<Boolean>
  batchActionArguments: Array,
  batchActionNames: Array,
  title: {
    default: 'Batch Action'
  }
});

const emit = defineEmits(['update:batchDialogVisible', 'closeBatchView']);

const dialogVisible = computed({
  get() {
    return props.batchDialogVisible
  },
  set(value) {
    emit('update:batchDialogVisible', value)
  }
})

let queryProgress = ref(0);
let queryList = ref([]);
let completionList = ref([]);



let closeViewFlag = ref(false);

onMounted(() => {
  handleOpen()
})
onBeforeUnmount(() => {
  closeViewFlag.value = true;
})



function handleOpen() {
  queryList.value = props.batchActionNames.map((item) => {
    return {
      ...item,
      status: 'waiting'
    }
  });
  queryProgress.value = 0;
  nextTick(() => {
    handleActionRetry();
  })
}


let actionIsFail = ref(false);
async function handleActionRetry() {
  actionIsFail.value = false;
  queryList.value.forEach(v => {
    if (v.status !== "success") {
      v.status = "waiting";
    }
  });
  let count = completionList.value.length;
  let len = queryList.value.length;
  queryProgress.value = Number((count * 100 / len).toFixed(0)) || 0;
  for (let i = completionList.value.length; i < len; i++) {
    let nodeItem = queryList.value[i];
    if (closeViewFlag.value) {
      return;
    }
    if (nodeItem.status == 'success') {
      continue;
    }
    let res = await props.batchActionCallBack(...(props.batchActionArguments[i]));
    if ((process.env.NODE_ENV === 'development') || res) {
      nodeItem.status = "success";
      completionList.value.push(nodeItem);
      count += 1;
      queryProgress.value = Number((count * 100 / len).toFixed(0)) || 0;
    } else {
      nodeItem.status = "error";
      actionIsFail.value = true;
      return
    }
  }
  await sleep(1000);
  nextToCloseView()
}

function nextToCloseView() {
  if (closeViewFlag.value) {
    return
  }
  ElMessage.success(t("tip.successfully"))
  dialogVisible.value = false;
}


function handleCloseDialog(close) {
  close()
}

</script>

<style scoped lang="scss"></style>

