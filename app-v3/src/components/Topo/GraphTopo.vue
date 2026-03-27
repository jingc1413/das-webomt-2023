<template>
  <div id="topoMain">
    <div id="mountNode" ref="mountNodeRef" @contextmenu.prevent></div>
    <img :src="g6toImgUrl" v-if="viewMode != provideKeys.viewModeDefaultValue" />
  </div>
</template>
<script>
import G6 from '@antv/g6'
import { useDasDevices } from '@/stores/das-devices'
import { useDasTopo } from "@/stores/topo"
import provideKeys from '@/utils/provideKeys.js'
import grayFull from '@/assets/gray.png'
import redFull from '@/assets/red.png'
import greenFull from '@/assets/green.png'
import { ref } from "vue";
import { useElementSize, useMutationObserver } from '@vueuse/core'

let topograph;
let FONTSIZE = 12;
let initCharNum = 7;
export default {
  name: 'graphTopo',
  inject: ['viewMode'],
  props: {
    page: Object,
    secondContent: String,
    isGraphTopo: Boolean
  },
  setup() {
    const dasDevices = useDasDevices() // 实例化userstore
    const dasTopo = useDasTopo();
    let mountNodeRef = ref(null)
    const { width: mountNodeWidth, height: mountNodeHeight } = useElementSize(mountNodeRef);
    return {
      mountNodeRef,
      dasDevices,
      dasTopo,
      provideKeys,
      mountNodeWidth,
      mountNodeHeight
    }
  },
  data() {
    return {
      g6toImgUrl: '',
      mountNodeSizeChange: null,
      nextShowRefresh: false,
    }
  },
  computed: {
    graphTopoData() {
      return this.dasTopo.graphTopoData;
    },
  },
  watch: {
    secondContent(newsecondContent, oldsecondContent) {
      this.changeswitch()
    },
    graphTopoData() {
      console.log("update graph", this.isGraphTopo, this.nextShowRefresh)
      if (this.isGraphTopo == false) {
        this.nextShowRefresh = true;
      } else {
        this.updateGraph();
      }
    },
    mountNodeWidth() {
      this.changeGraphSize()
    },
    mountNodeHeight() {
      this.changeGraphSize()
    },
    isGraphTopo() {
      if (!topograph && this.isGraphTopo == true) {
        this.$nextTick(()=>{
          this.initG6();
          this.updateGraph();
        })
      }
      if (this.nextShowRefresh && this.isGraphTopo == true) {
        this.nextShowRefresh = false;
        setTimeout(() => {
          this.updateGraph();
        }, 250);
      }
    },
  },
  mounted() {
    if (this.isGraphTopo == true) {
      this.initG6();
      this.updateGraph();
    }
  },
  methods: {
    changeswitch() {
      if (!topograph) return;
      const nodes = topograph.getNodes();
      nodes.forEach((node) => {
        const newNodeData = node.getModel()
        node.update(newNodeData);
      });
    },
    fittingStringLength(str, maxWidth = 70) {
      const ellipsis = '...'
      const ellipsisLength = G6.Util.getTextSize(ellipsis, FONTSIZE)[0]
      let currentWidth = 0
      let res = str
      const HanPattern = /[\u4e00-\u9fff]/g;
      const EnPattern = /[a-zA-Z_-]/g;
      str.split('').forEach((letter, i) => {
        if (currentWidth > maxWidth - ellipsisLength) return
        if (HanPattern.test(letter)) {
          currentWidth += FONTSIZE
        } else if (EnPattern.test(letter)) {
          currentWidth += (G6.Util.getLetterWidth(letter, FONTSIZE) / 2)
        } else {
          currentWidth += G6.Util.getLetterWidth(letter, FONTSIZE)
        }
        if (currentWidth > maxWidth - ellipsisLength) {
          res = `${str.substr(0, i)}${ellipsis}`
        }
      })
      return res
    },
    installContent2ByOption(isFull, cfg) {
      var content2;
      var dispalystr;
      let w = cfg.size[0];
      switch (this.secondContent) {
        default:
        case 'Device Name':
          if (isFull) {
            return cfg.info.DeviceName;
          }
          content2 = this.fittingStringLength(cfg.info.DeviceName, w - 70);
          break;
        case 'Location':
          if (isFull) {
            return cfg.info.InstalledLocation;
          }
          content2 = this.fittingStringLength(cfg.info.InstalledLocation, w - 70);
          break;
      }
      return content2;
    },
    handleSelectDevice(subID) {
      this.$emit('selectDevice', subID)
    },
    updateCurrentDevice: async function () {
      if (this.dasDevices.currentDeviceInfo?.SubID !== undefined) {
        this.handleSelectDevice(this.dasDevices.currentDeviceInfo.SubID);
      }
    },
    initG6() {
      const self = this;
      const defaultStateStyles = {
        hover: {
          stroke: '#1890ff',
          lineWidth: 2,
        },
      }
      const defaultNodeStyle = {
        fill: '#91d5ff',
        stroke: '#40a9ff',
        radius: 5,
      }

      const defaultNodeLostStyle = {
        fill: '#000000',
        stroke: '#FFFFFF',
        radius: 5,
      }

      const defaultEdgeStyle = {
        stroke: '#91d5ff',
        endArrow: {
          path: 'M 0,0 L 12, 6 L 9,0 L 12, -6 Z',
          fill: '#91d5ff',
          d: -20,
        },
      }
      const defaultLayout = {
        type: 'compactBox',
        direction: 'TB',
        getId: function getId(d) {
          return d.id
        },
        getHeight: function getHeight() {
          return 48
        },
        getWidth: function getWidth() {
          return 16
        },
        getVGap: function getVGap() {
          return 40
        },
        getHGap: function getHGap() {
          return 70
        },
      }
      const defaultLabelCfg = {
        style: {
          fill: '#303133',
          fontSize: 12,
        },
      }
      G6.registerNode(
        'icon-node',
        {
          options: {
            size: [60, 20],
            stroke: '#91d5ff',
            fill: '#91d5ff',
          },
          draw(cfg, group) {
            const styles = this.getShapeStyle(cfg)
            const { labelCfg = {} } = cfg
            const w = cfg.size[0]
            const h = cfg.size[1]
            if (self.dasDevices?.isCurrentDevice(cfg.id)) {
              styles.stroke = '#7426d9'
            }
            var imgurl = ''
            if (cfg.info.AlarmState == 0) {
              imgurl = greenFull
            } else {
              imgurl = redFull
            }
            if (cfg.info.ConnectState >= 6) {
              imgurl = grayFull
            }
            const keyShape = group.addShape('rect', {
              attrs: {
                ...styles,
                lineWidth: 1,
                x: -w / 2,
                y: -h / 2,
              },
              name: 'rect1'
            })
            // state logo


            const style = {
              fill: '#F0FFFF',
              stroke: '#F0FFFF',
            }
            group.addShape('rect', {
              attrs: {
                x: 1 - w / 2,
                y: 1 - h / 2,
                width: 38,
                height: styles.height - 2,
                ...style,
              },
              name: 'rect2'
            })

            group.addShape('image', {
              attrs: {
                x: 8 - w / 2,
                y: 8 - h / 2,
                width: 24,
                height: 24,
                img: imgurl,
              },
              name: 'image-shape',
            })
            var content;
            if (cfg.info.DeviceTypeName) {
              content = cfg.info.CascadingLevel ? `L${cfg.info.CascadingLevel}: ${cfg.info.DeviceTypeName}` : `${cfg.info.DeviceTypeName}`
            } else {
              content = "undefined"
            }
            group.addShape('text', {
              attrs: {
                ...labelCfg.style,
                text: content,
                x: 45 - w / 2,
                y: 18 - h / 2,
              },
              name: 'text1'
            })
            var content2 = self.installContent2ByOption(false, cfg)
            group.addShape('text', {
              attrs: {
                ...labelCfg.style,
                text: content2,
                x: 45 - w / 2,
                y: 34 - h / 2,
                fontSize: FONTSIZE,
              },
              name: 'text2'
            })
            return keyShape
          },
          update(cfg, node) {
            const group = node.getContainer(); // 获取容器
            const label = group.get('children')[4]; // 按照添加的顺序
            const styles = this.getShapeStyle(cfg);
            const { labelCfg = {} } = cfg;
            const w = cfg.size[0];
            const h = cfg.size[1];

            var content2 = self.installContent2ByOption(false, cfg)
            const style = {
              ...labelCfg.style,
              text: content2,
              x: 45 - w / 2,
              y: 34 - h / 2,
              fontSize: FONTSIZE
            };
            label.attr(style);

          },
        },
        'rect'
      )
      G6.registerEdge('flow-line', {
        draw(cfg, group) {
          const startPoint = cfg.startPoint
          const endPoint = cfg.endPoint

          const { style } = cfg
          const shape = group.addShape('path', {
            attrs: {
              stroke: style.stroke,
              lineWidth: 2,
              endArrow: style.endArrow,
              path: [
                ['M', startPoint.x, startPoint.y],
                ['L', startPoint.x, (startPoint.y + endPoint.y)/2 - 16],
                ['L', endPoint.x, (startPoint.y + endPoint.y)/2 - 16],
                ['L', endPoint.x, endPoint.y],
              ],
            },
          })
          var OpticalPort = cfg.targetNode._cfg.model.info.OpticalPort;
          if (cfg.targetNode._cfg.model.info?.OpticalInputPort) {
            OpticalPort = `${cfg.targetNode._cfg.model.info.OpticalPort}

${cfg.targetNode._cfg.model.info.OpticalInputPort}`
          }
          group.addShape('text', {
            attrs: {
              text: OpticalPort,
              fill: '#304156',
              textAlign: 'bottom',
              //          textBaseline: "middle",
              x: endPoint.x + 5,
              y: endPoint.y - 32,
            },
            name: 'left-text-shape',
          })
          return shape
        },
      })
      const width = document.getElementById('topoMain').offsetWidth - 20
      const height = window.innerHeight - 148

      const tooltip = new G6.Tooltip({
        className: 'component-tooltip',
        fixToNode: [1, 0.5],
        itemTypes: ['node'],
        getContent(e) {
          return self.installContent2ByOption(true, e.item._cfg.model);
        },
        shouldBegin(e) {
          if (self.installContent2ByOption(true, e.item._cfg.model) == '') {
            return false;
          } else {
            return true;
          }
        }
      });
      let topoMode = this.viewMode == provideKeys.viewModeDefaultValue ? ['drag-canvas', 'zoom-canvas'] : [];
      topograph = new G6.TreeGraph({
        container: 'mountNode',
        width,
        height,
        linkCenter: true,
        disableContextMenu: true,
        plugins: [tooltip],
        modes: {
          default: [
            ...topoMode,
          ],
        },
        defaultNode: {
          type: 'icon-node',
          size: [140, 40],
          style: defaultNodeStyle,
          labelCfg: defaultLabelCfg,
        },
        defaultEdge: {
          type: 'flow-line',
          style: defaultEdgeStyle,
        },
        nodeStateStyles: {
          hover: {
            //fill: '#d3adf7',
            'rect1': {
              lineWidth: 2,
              // stroke: '#d3adf7',
            }
          },
          selected: {
            'rect1': {
              lineWidth: 2.5,
              // stroke: '#d3adf7',
            },
            'rect2': {
              height: 40 - 3,
            },
            //fillOpacity: 0.8,
          }
        },
        edgeStateStyles: defaultStateStyles,
        layout: defaultLayout,
        maxZoom: 1,
      })

      topograph.on('node:mouseenter', (evt) => {
        const { item } = evt
        topograph.setItemState(item, 'hover', true)
      })

      topograph.on('node:mouseleave', (evt) => {
        const { item } = evt
        const model = item.getModel()
        topograph.setItemState(item, 'hover', false)
      })

      topograph.on('node:click', (evt) => {
        const node = evt.item
        const model = node.getModel()
        const nodes = topograph.getNodes();
        nodes.forEach(ele => { topograph.clearItemStates(ele) })

        topograph.setItemState(node, 'selected', true);
        this.handleSelectDevice(model.info.SubID)
      })
      topograph.on('node:contextmenu', (evt) => {
        const { item } = evt
        const model = item.getModel()
        if (!this.dasDevices.isCurrentDevice(0)) {
          return
        }
        if (model.id == 0) {
          return
        }
        this.deleteitem(model.id)
      })

      topograph.on('afterrender', ev => {
        if (!this.g6toImgUrl) {
          this.graph2img();
        }
      })

    },
    updateGraph() {
      console.log("update graph", this.graphTopoData)
      if (this.graphTopoData !== undefined) {
        topograph.data(this.graphTopoData);
        topograph.render();
        topograph.fitView();
        this.updateCurrentDevice();
      }
    },
    graph2img() {
      if (this.viewMode !== this.provideKeys.viewModeDefaultValue) {
        setTimeout(() => {
          this.g6toImgUrl = topograph.toDataURL({ type: '' });
          this.mountNodeRef.style.display = "none";
        }, 2000);
      }
    },
    async deleteitem(id) {
      this.$emit('DeleteDevice', id)
    },
    changeGraphSize() {
      if (!topograph) {
        return
      }
      if (this.mountNodeSizeChange) {
        clearTimeout(this.mountNodeSizeChange)
      }
      this.mountNodeSizeChange = setTimeout(() => {
        topograph.changeSize(this.mountNodeWidth, this.mountNodeHeight);
      }, 100);
    }
  }
}
</script>
<style lang="scss">
.component-tooltip {
  border: 1px solid #e2e2e2;
  border-radius: 4px;
  font-size: 12px;
  color: #545454;
  background-color: white;
  padding: 6px 5px;
  box-shadow: rgb(174 174 174) 0px 0px 0px;
}

#mountNode {
  width: 100%;
  height: calc(100vh - 150px);
}
</style>