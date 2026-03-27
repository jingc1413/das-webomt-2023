<template>
  <el-row>
    <el-col>
      <el-form :inline="true" :model="bandwidthForm">
        <el-form-item label="Current FC">
          <el-input-number v-model.lazy="bandwidthForm.currentFrequencyCenter" size="small" :controls="false"
            style="width: 64px;" :max="currentFrequencyMax" :min="currentFrequencyMin"
            @change="updateCurrentFrequencyCenter" />
        </el-form-item>
        <el-form-item label="Bandwidth">
          <el-input-number v-model.lazy="bandwidthForm.currentBandwidth" size="small" :controls="false"
            style="width: 64px;" @change="updateCurrentFrequencyCenter" />
        </el-form-item>
        <!-- <el-form-item>
          <el-button type="primary" @click="onSubmit">Set</el-button>
        </el-form-item> -->
      </el-form>
    </el-col>
    <el-col>
      <div ref="bandwidthSelectorChartRef" class="chart_size">
      </div>
    </el-col>
  </el-row>
</template>

<script>
import * as echarts from 'echarts';
import { ref } from 'vue'
import { useElementSize } from '@vueuse/core'

export default {
  name: "BandwidthSelectorView",
  props: {
    bandwidth: {
      type: Number,
      required: true
    },
    frequencyCenter: {
      type: Number,
      required: true
    },
    frequencyMax: {
      type: Number,
      required: true
    },
    frequencyMin: {
      type: Number,
      required: true
    },
    otherFrequencies: {
      type: Array,//[{BW,FC}...]
      default: () => {
        return []
      }
    }
  },
  setup() {
    let bandwidthSelectorChartRef = ref(null);
    const { width: bandwidthSelectorChartWidth } = useElementSize(bandwidthSelectorChartRef);
    return { bandwidthSelectorChartRef, bandwidthSelectorChartWidth }
  },
  data() {
    return {
      bandwidthForm: {
        currentFrequencyCenter: -1,
        currentBandwidth: -1,
      },
    }
  },
  computed: {
    currentHalf_BW() {
      return this.bandwidthForm.currentBandwidth / 2;
    },
    currentFrequencyMax() {
      return this.frequencyMax - this.currentHalf_BW;
    },
    currentFrequencyMin() {
      return this.frequencyMin + this.currentHalf_BW;
    },
    currentFrequencyCenterNotFault() {
      let currentMax = this.bandwidthForm.currentFrequencyCenter + this.currentHalf_BW;
      let currentMin = this.bandwidthForm.currentFrequencyCenter - this.currentHalf_BW;
      return this.otherFrequencies.every(item => {
        let half_BW = item.BW / 2;
        let max = item.FC + half_BW;
        let min = item.FC - half_BW;
        return (currentMin > max && currentMin > min) || (currentMax < min && currentMax < max)
      })
    },
    chartColor_Other() {
      return '#eeeeee'
    },
    chartColor_Fault() {
      return '#fcd3d3'
    },
    chartColor_Properly() {
      return '#c6e2ff'
    },
    all_FC_list() {
      return [this.bandwidthForm.currentFrequencyCenter].concat(this.otherFrequencies.map(item => item.FC));
    }
  },
  watch: {
    bandwidthSelectorChartWidth() {
      this.watchResizeChartDom();
    },
  },
  mounted() {
    this.bandwidthForm.currentFrequencyCenter = this.frequencyCenter;
    this.bandwidthForm.currentBandwidth = this.bandwidth;
    this.initChartDom();
  },
  beforeUnmount() {
    let chartDom = echarts.getInstanceByDom(this.$refs['bandwidthSelectorChartRef']);
    if (chartDom) {
      chartDom.dispose();
    }
  },
  methods: {
    watchResizeChartDom() {
      let chartDom = echarts.getInstanceByDom(this.$refs['bandwidthSelectorChartRef']);
      this.$nextTick(() => {
        chartDom && chartDom.resize({ width: this.bandwidthSelectorChartWidth });
      })
    },
    initChartDom() {
      let chartDom = echarts.init(this.$refs['bandwidthSelectorChartRef'], null, {
        locale: 'EN',
      });
      let option = this.getChartOption();
      chartDom.setOption(option);
      setTimeout(() => {
        chartDom && chartDom.resize();
      }, 100);
    },
    getChartOption() {
      let series = [];
      let self = this;
      let seriesOption = {
        type: 'line',
        symbolSize: 0,
        label: {
          show: true,
          color: '#000',
          formatter(params) {
            let { value } = params;
            console.log({params});
            if (self.all_FC_list.includes(value[0])) {
              return value[0]
            }
            return ''
          }
        }
      }
      series.push({
        name: `current`,
        data: this.getChartSeriesData(this.bandwidthForm.currentFrequencyCenter, this.bandwidthForm.currentBandwidth),
        ...seriesOption,
        lineStyle: {
          color: this.chartColor_Properly,
        },
        areaStyle: {
          color: this.chartColor_Properly,
        }
      })
      this.otherFrequencies.forEach((item, index) => {
        series.push({
          name: `other-${index}`,
          data: this.getChartSeriesData(item.FC, item.BW),
          ...seriesOption,
          lineStyle: {
            color: this.chartColor_Other,
          },
          areaStyle: {
            color: this.chartColor_Other,
          }
        })
      })
      let option = {
        tooltip: {
          show: false
        },
        grid: {
          left: 24,
          right: 36,
          top: 24,
          bottom: 24,
        },
        xAxis: {
          type: 'value',
          name: 'MHz',
          show: true,
          min: this.frequencyMin,
          max: this.frequencyMax,
        },
        yAxis: {
          type: 'value',
          show: false,
          min: -0.4,
          max: 100,
        },
        series: series
      }
      return option;
    },
    getChartSeriesData(FC, BW) {
      let data = [];
      let half_BW = BW / 2;
      data.push([FC - half_BW - 0.01, 0]);
      data.push([FC - half_BW, 60]);
      data.push([FC, 60]);
      data.push([FC + half_BW, 60]);
      data.push([FC + half_BW + 0.01, 0]);
      return data;
    },
    updateCurrentFrequencyCenter() {
      let data = this.getChartSeriesData(this.bandwidthForm.currentFrequencyCenter, this.bandwidthForm.currentBandwidth);
      let color = this.chartColor_Properly;
      if (this.currentFrequencyCenterNotFault == false) {
        color = this.chartColor_Fault;
      }
      let chartDom = echarts.getInstanceByDom(this.$refs['bandwidthSelectorChartRef']);
      chartDom.setOption({
        series: [
          {
            name: `current`,
            data,
            lineStyle: {
              color,
            },
            areaStyle: {
              color,
            }
          }
        ]
      })
    },
    onSubmit() {
      this.$emit('submitFrequencyCenter', this.bandwidthForm.currentFrequencyCenter);
    }
  }
}
</script>

<style lang="scss" scoped>
.chart_size {
  height: 200px;
  width: 100%;
}
</style>