<template>
  <div id="topoMain">
    <el-tree :data="dasTopo.treeTopoData" :props="{ class: customNodeClass }" @node-click="handleNodeClick"
      highlight-current @node-contextmenu="handleContextmenu" :expand-on-click-node="false" :default-expand-all="true">
      <template #default="{ data, node }">
        <el-icon>
          <CircleCloseFilled v-if="data.info.ConnectState >= 6" style="color: rgb(149, 153, 152);" />
          <WarningFilled v-else-if="data.info.AlarmState === 1" style="color:red;" />
          <SuccessFilled v-else style="color: rgb(20, 185, 163);" />
        </el-icon>
        <el-tag v-if="data.id != 0" type="info" style="margin-left:4px">
          {{ data.info.OpticalPort + `${data.info?.OpticalInputPort && ('>>'+data.info.OpticalInputPort)}` }}{{ " | L" }}{{ data.info.CascadingLevel }}
        </el-tag>
        <span style="margin-left:4px; margin-right:2px">{{ data.info.DeviceTypeName + ": "}}</span>
        <span v-if="secondContent == 'Device Name'">{{ data.info.DeviceName }}</span>
        <span v-else-if="secondContent == 'Location'">{{ data.info.InstalledLocation }}</span>
        <span v-else>{{ data.info.DeviceName }}</span>
      </template>
    </el-tree>
  </div>
</template>
<script>
import { useDasDevices } from '@/stores/das-devices'
import { useDasTopo } from "@/stores/topo"

export default {
  name: 'treeTopo',
  props: {
    page: Object,
    secondContent: String,
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
        children: 'children',
        label: 'id',
        class: 'customNodeClass',
      },
    }
  },
  methods: {
    customNodeClass(data, node) {
      var self = this;
      if (self.dasDevices?.isCurrentDevice(data.id)) {
        return 'is-penultimate'
      }
      return null
    },
    updateClickedDeviceInfo(model) {
      this.$emit('selectDevice', model.SubID)
    },
    handleNodeClick(data) {
      this.$emit('selectDevice', data.info.SubID)
    },
    handleContextmenu(data, node) {
      if (!this.dasDevices.isCurrentDevice(0)) {
        return
      }
      if (node.id == 0) {
        return
      }
      this.deleteitem(node.id)
    },
    async deleteitem(id) {
      this.$emit('deleteDevice', id)
    }
  },
}
</script>
<style>
.is-penultimate>.el-tree-node__content {
  background-color: #ebe2f696;
}
</style>