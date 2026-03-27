<template>
  <div>
    <el-input v-model="filterText" style="width: 200px" clearable />
    <el-scrollbar style="height: 240px;  margin-top: 12px;">
      <el-tree :data="dasTopo.treeTopoData" :props="{ class: customNodeClass }" @node-click="handleNodeClick"
        highlight-current :expand-on-click-node="false" :default-expand-all="true" :filter-node-method="filterNode"
        ref="treeRef" class="device_tree_selected" node-key="id">
        <template #default="{ data }">
          <el-icon>
            <CircleCloseFilled v-if="data.info.ConnectState >= 6" style="color: rgb(149, 153, 152);" />
            <WarningFilled v-else-if="data.info.AlarmState === 1" style="color:red;" />
            <SuccessFilled v-else style="color: rgb(20, 185, 163);" />
          </el-icon>
          <el-tag v-if="data.id != 0" type="info" style="margin-left:4px">
            {{ data.info.OpticalPort + `${data.info?.OpticalInputPort && ('>>'+data.info.OpticalInputPort)}` }}{{ " | L" }}{{ data.info.CascadingLevel }}
          </el-tag>
          <span style="margin-left:4px; margin-right:2px">{{ data.info.DeviceTypeName + ": " }}</span>
          <span v-if="treeDeviceName == 'Device Name'">{{ data.info.DeviceName }}</span>
          <span v-else-if="treeDeviceName == 'Location'">{{ data.info.InstalledLocation }}</span>
          <span v-else>{{ data.info.DeviceName }}</span>
        </template>
      </el-tree>
    </el-scrollbar>
  </div>
</template>
<script>
import { useDasDevices } from '@/stores/das-devices'
import { useDasTopo } from "@/stores/topo"
import { ElMessage } from 'element-plus'

export default {
  name: 'DeviceTreeSelected',
  props: {
    showDefaultSelected: {
      type: Boolean,
      default: true
    }
  },
  setup() {
    const dasDevices = useDasDevices() // 实例化userstore
    const dasTopo = useDasTopo();
    return {
      dasDevices,
      dasTopo,
    }
  },
  data() {
    return {
      defaultProps: {
        disabled(data) {
          return data.info.ConnectState < 6
        }
      },
      filterText: '',
      filterOptions: {
        online: true,
      }
    }
  },
  computed: {
    treeDeviceName() {
      return this.dasTopo.treeDeviceName
    }
  },
  watch: {
    filterText(val) {
      this.$refs.treeRef.filter(val)
    },
    'dasTopo.treeTopoData'() {
      this.$nextTick(()=>{
        this.$refs.treeRef.filter(this.filterText)
      })
    }
  },
  mounted() {
    this.$nextTick(()=>{
      this.$refs.treeRef.filter(this.filterText)
    })
  },
  methods: {
    customNodeClass(data, node) {
      if (data.info.ConnectState >= 6) {
        return 'tree_select_disable_click'
      }
      if (this.dasDevices?.isCurrentDevice(data.id) && this.showDefaultSelected) {
        return 'tree_select_is_penultimate'
      }
      return null
    },
    handleNodeClick(data) {
      this.$emit('selectDevice', data.info)
    },
    filterNode(value, data) {
      if (this.filterOptions.online !== undefined) {
        if (this.filterOptions.online) {
          if (data.info.ConnectState >= 6 || data.info.ConnectState < 1) {
            return false
          }
        } else {
          if (data.info.ConnectState < 6 && data.info.ConnectState > 0) {
            return false
          }
        }
      }
      
      if (!value) return true
      return data.info.DeviceName.includes(value)
    }
  },
}
</script>

<style lang="scss" scoped>
.device_tree_selected {
  min-height: 240px;

  & :deep(.el-tree-node) {
    .el-tree-node__children {
      display: table;
      padding-right: 8px;
    }
  }
}
</style>
<style lang="scss">
.tree_select_is_penultimate>.el-tree-node__content {
  background-color: #ebe2f696;
}

.tree_select_disable_click>.el-tree-node__content {
  cursor: not-allowed;
}
</style>