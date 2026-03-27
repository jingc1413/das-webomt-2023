<template>
  <span>
    {{ formatTime }}
  </span>
</template>

<script setup>
import { dayjs } from 'element-plus';
import { computed, onMounted, ref } from 'vue';
import { useRafFn } from '@vueuse/core'

const beginTime = ref(null);
const nowTime = ref(Date.now());

const { pause, resume } = useRafFn(() => {
  nowTime.value = Date.now();
})

let formatTime = computed(()=>{
  return filterFormatTime(dayjs(nowTime.value).diff(beginTime.value), 's')
})

onMounted(()=>{
  beginTime.value = Date.now();
})

function filterFormatTime(timeNumber) {
  let temp = timeNumber/1000;
  let m = parseInt((temp / 60).toFixed(1)).toString().padStart(2, '0');
  temp = temp % 60;
  let s = parseInt(temp).toString().padStart(2, '0')

  return `${m ?? '00'}:${s ?? '00'}`
}

</script>

<style scoped>
</style>