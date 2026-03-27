<!-- eslint-disable vue/no-v-for-template-key -->
<template>
  <div class="table-search-div">

    <el-form :inline="true" :model="searchForm" class="table-search-form">
      <template v-for="(infoValue, infoKey, index) in searchItemsInfo" :key="infoKey">
        <el-form-item :label="infoValue.name" v-show="index < showSearchItemCount" style="margin-bottom: 8px;">
          <template v-if="infoValue.type == 'string'">
            <el-input v-model="searchForm[infoKey]" size="small" style="width: 90px;" clearable
              @blur="stringInputBlur(infoKey)" />
          </template>
          <template v-else-if="infoValue.type == 'number'">
            <el-input-number v-model.lazy="searchForm[infoKey]" size="small" :controls="false" style="width: 64px;" />
          </template>
          <template v-else-if="infoValue.type == 'select'">
            <el-select v-model="searchForm[infoKey]" style="width: 160px" size="small" clearable>
              <el-option v-for="optionItem in infoValue.options" :key="optionItem" :label="optionItem"
                :value="optionItem" />
            </el-select>
          </template>
        </el-form-item>
      </template>

      <el-form-item label="" style="margin-bottom: 8px;">
        <el-button circle @click="showAllItem = !showAllItem" v-if="showSearchItemCount < modelValue.length">
          <el-icon v-show="!showAllItem">
            <ArrowDownBold />
          </el-icon>
          <el-icon v-show="showAllItem">
            <ArrowUpBold />
          </el-icon>
        </el-button>
        <el-button type="primary" :disabled="loading" @click="startQuery">Query</el-button>
        <!-- <el-button @click="clearQuery">Clear</el-button> -->
      </el-form-item>
    </el-form>

    <el-form :inline="true" :model="searchForm" class="table-search-form">
      <TransitionGroup name="el-zoom-in-top">
        <template v-for="(infoValue, infoKey, index) in searchItemsInfo" :key="infoKey">

          <el-form-item :label="infoValue.name" v-show="index >= showSearchItemCount && showAllItem"
            style="margin-bottom: 8px;">
            <template v-if="infoValue.type == 'string'">
              <el-input v-model="searchForm[infoKey]" size="small" style="width: 90px;" clearable
                @blur="stringInputBlur(infoKey)" />
            </template>
            <template v-else-if="infoValue.type == 'number'">
              <el-input-number v-model.lazy="searchForm[infoKey]" size="small" :controls="false" style="width: 64px;" />
            </template>
            <template v-else-if="infoValue.type == 'select'">
              <el-select v-model="searchForm[infoKey]" style="width: 160px" size="small" clearable>
                <el-option v-for="optionItem in infoValue.options" :key="optionItem" :label="optionItem"
                  :value="optionItem" />
              </el-select>
            </template>
          </el-form-item>
        </template>
      </TransitionGroup>
    </el-form>

  </div>
</template>

<script>

export default {
  name: "BandwidthSelectorView",
  props: {
    showSearchItemCount: {
      type: Number,
      default: 3,
    },
    modelValue: {
      type: Array,//[{key,name,type, value}...]
      default: () => {
        return []
      }
    },
    loading:{
      type: Boolean,
      default: false
    }
  },
  emits: ['update:modelValue'],
  data() {
    return {
      searchForm: {
      },
      searchItemsInfo: {
      },
      showAllItem: false,
    }
  },
  watch: {
    modelValue: {
      handler() {
        this.modelValue.forEach(item => {
          let { key, name, value, ...options } = item;
          this.searchItemsInfo[key] = { name, ...options };
        })
      }
    }
  },
  mounted() {
    this.initSearchBar();
  },
  methods: {
    initSearchBar() {
      this.modelValue.forEach(item => {
        let { key, name, value, ...options } = item;
        this.searchForm[key] = value;
        this.searchItemsInfo[key] = { name, ...options };
      })
      console.log('searchForm', this.searchForm);
    },
    stringInputBlur(key) {
      if (this.searchForm[key]) {
        this.searchForm[key] = this.searchForm[key].trim();
      }
    },
    clearQuery() {
      this.modelValue.forEach(item => {
        let { key } = item;
        this.searchForm[key] = null;
      })
      this.startQuery();
    },
    startQuery() {
      let temp = [...this.modelValue];
      temp.forEach((item) => {
        let { key } = item;
        item.value = this.searchForm[key]
      })
      this.$emit('update:modelValue', temp);
    }
  }
}
</script>

<style lang="scss" scoped>
.table-search-div {
  border-bottom: 1px solid #eee;
  margin-bottom: 12px;
  // border-radius: 4px;
  // box-shadow: 0px 0px 6px rgba(0, 0, 0, .12);
}

.table-search-form {
  // background: #eee;
  // margin: 8px 0;
  padding: 8px 12px 0 12px;

  &:last-child {
    margin-bottom: 8px;
  }
}
</style>