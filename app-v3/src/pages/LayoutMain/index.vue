<template>
  <div>
    <my-main-page />
  </div>
</template>

<script>
import { ref } from "vue";
import { useRoute } from "vue-router";

export default {
  name: "LayoutPage",
  setup: function () {
    const sub = ref(null);
    const moduleName = ref(null);
    const pageName = ref(null);
    const route = useRoute();
    sub.value = route.params.sub;
    moduleName.value = route.params.module;
    pageName.value = route.params.page;
    return {
      sub,
      moduleName,
      pageName,
    };
  },
};
</script>
<style lang="scss" scoped></style>
