<template>
  <div v-loading.fullscreen.lock="fullscreenLoading"></div>
</template>

<script setup>

import { ref } from "vue";
import { onMounted } from "vue";

const fullscreenLoading = ref(true);

onMounted(async () => {
  // await useAppStore().gotoDefaultRoute({ checkAvailable: true });
})

</script>
