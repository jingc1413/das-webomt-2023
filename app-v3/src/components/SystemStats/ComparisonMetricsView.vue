<template>
  <el-col>
    <div class="chart_key_select">
      <el-select v-model="leftKey" style="width: 240px;" @change="changeChartKey()">
        <el-option v-for="keyItem in chartKeys" :key="keyItem" :label="chartMetrics[keyItem]?.name" :value="keyItem"
          :disabled="keyItem == rightKey" />
      </el-select>

      <el-select v-model="rightKey" style="width: 240px;" @change="changeChartKey()">
        <el-option v-for="keyItem in chartKeys" :key="keyItem" :label="chartMetrics[keyItem]?.name" :value="keyItem"
          :disabled="keyItem == leftKey" />
      </el-select>
    </div>
    <div id="comparisonChart" class="chart_size" style="width: 100%;">
    </div>
  </el-col>
</template>

<script>
import * as echarts from 'echarts';
import { dayjs } from 'element-plus';

import { useWindowSize } from '@vueuse/core'
import { useDasDevices } from '@/stores/das-devices';

let comparisonChartDom = null;

export default {
  name: 'ComparisonMetricsView',
  props: {
    chartKeys: {
      default: () => { return [] },
      type: Array
    },
    beginTime: {
      default: null,
      type: Number
    },
    endTime: {
      default: null,
      type: Number
    },
    activeTabName: {
      default: 'system',
      type: String
    },
    resizeTabName: {
      default: 'comparison',
      type: String
    }
  },
  setup() {
    let dev = useDasDevices().currentDevice;
    const { width: windowWidth, height: windowHeight } = useWindowSize();
    return { dev, windowWidth, windowHeight }
  },
  data() {
    return {
      sub: '',
      chartMetrics: {},
      leftKey: '',
      rightKey: '',
    }
  },
  watch: {
    windowWidth() {
      this.watchResizeChartDom()
    },
    windowHeight() {
      this.watchResizeChartDom()
    },
    activeTabName() {
      this.watchResizeChartDom()
    }
  },
  mounted() {
    this.sub = this.$route.params.sub;
    this.getChartKeys();
    this.initChartDom();
  },
  beforeUnmount() {
    comparisonChartDom && comparisonChartDom.dispose();
    comparisonChartDom = null;
  },
  methods: {
    watchResizeChartDom() {
      if (this.activeTabName != this.resizeTabName) {
        return
      }
      this.$nextTick(() => {
        comparisonChartDom && comparisonChartDom.resize();
      })
    },
    getChartKeys() {
      return new Promise((resolve, reject) => {
        this.chartMetrics = this.dev.stats.getMetrics(this.sub);
        resolve(null);
      })
    },
    initChartDom() {
      let dom = echarts.init(document.getElementById('comparisonChart'), null, { locale: 'EN' });
      let option = this.getEchartOptionWithXAxisTime();
      comparisonChartDom = dom;
      dom.setOption(option);
    },
    getEchartOptionWithXAxisTime() {
      let option = {
        tooltip: {
          trigger: 'axis',
          confine: true,
        },
        grid: {
          left: 64,
          right: 64,
          top: 48,
        },
        legend: {
          show: true,
          left: 64,
          right: 64,
          bottom: 40
        },
        xAxis: {
          type: 'time',
          boundaryGap: false,
          axisLabel: {
            formatter: {
              year: '{yyyy}',
              month: '{M}',
              day: '{M}/{dd}',
              hour: '{HH}:{mm}',
              minute: '{HH}:{mm}',
              second: '{HH}:{mm}:{ss}',
            }
          },
        },
        yAxis: [{
          id: 'left',
          name: 'left',
          type: 'value',
          scale: true,
          yAxisIndex: 0,
          nameTextStyle: {
            align: "left"
          }
        },
        {
          id: 'right',
          name: 'right',
          type: 'value',
          scale: true,
          yAxisIndex: 1,
          nameTextStyle: {
            align: "right"
          }
        }],
        series: []
      };
      return option;
    },
    filterChartDate(chartDates, metricsItem) {
      let { name, key, items } = metricsItem;
      let cpItems = items;
      if (!cpItems) {
        cpItems = [{ key, name }]
      }
      let seriesData = [];
      let firstTime = dayjs(this.beginTime, 'X').format('YYYY-MM-DD HH:mm:ss');
      let lastTime = dayjs(this.endTime, 'X').format('YYYY-MM-DD HH:mm:ss');
      let addHeaderTime = dayjs(this.beginTime, 'X').subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss');
      let addEndTime = dayjs(this.endTime, 'X').add(1, 'hour').format('YYYY-MM-DD HH:mm:ss');
      cpItems.forEach(({ key: itemKey, name: itemName }) => {
        let chartData = chartDates[itemKey] ?? [];
        if (chartData.length == 0) {
          chartData.push([firstTime, null])
          chartData.push([lastTime, null])
          seriesData.push({ key: itemKey, chartData, name: itemName })
          return
        }
        let first = chartData[0][0] ?? firstTime;
        if (first.localeCompare(firstTime) > 0) {
          chartData = [[firstTime, null]].concat(chartData);
        } else {
          chartData = [[addHeaderTime, null]].concat(chartData);
        }
        let last = chartData[chartData.length - 1][0] ?? lastTime;
        if (last.localeCompare(lastTime) < 0) {
          chartData = chartData.concat([[lastTime, null]])
        } else {
          chartData = chartData.concat([[addEndTime, null]])
        }
        seriesData.push({ key: itemKey, chartData, name: itemName })
      })
      return seriesData;
    },
    getSeriesDate(seriesData, yAxisIndex) {
      let series = [];
      seriesData.forEach(({ key, chartData, name }) => {
        series.push({
          yAxisIndex,
          id: key,
          name: name,
          type: 'line',
          data: chartData,
          smooth: true,
          symbolSize: 2,
          showSymbol: true,
          showAllSymbol: 'auto',
          tooltip: {
            valueFormatter: (value) => {
              if (typeof value == 'string' && value.length > 16) {
                return value.substring(13) + '...'
              } else {
                return value
              }
            }
          },
          markPoint: {
            data: [
              { type: 'max', name: 'Max', label: { formatter: 'Max' } },
              { type: 'min', name: 'Min', label: { formatter: 'Min' } }
            ]
          }
        })
      })
      return series;
    },
    changeChartKey() {
      if (!this.beginTime || !this.endTime) {
        return;
      }
      let series = [];
      let yAxis = [];
      if (this.leftKey) {
        let { seriesData, yAxisItem } = this.getUpdateOptions(this.leftKey, 0);
        series.push(...seriesData);
        yAxis.push(yAxisItem);
      }

      if (this.rightKey) {
        let { seriesData, yAxisItem } = this.getUpdateOptions(this.rightKey, 1);
        series.push(...seriesData);
        yAxis.push(yAxisItem);
      }
      this.updateChartDom(series, yAxis)

    },
    getUpdateOptions(chartKey, yAxisIndex) {
      let seriesData = [];
      let metricsItem = this.chartMetrics[chartKey];
      seriesData = this.filterChartDate(this.dev.stats.getChartData(), metricsItem);
      seriesData = this.getSeriesDate(seriesData);
      let { unit, yInterval, max, min } = metricsItem
      let yAxisItem = {
        id: yAxisIndex ? 'right' : 'left',
        type: 'value',
        name: metricsItem.name + ' ' + unit,
        scale: true,
        max,
        min,
        interval: yInterval,
        yAxisIndex: yAxisIndex,
      }

      return { seriesData, yAxisItem };
    },
    updateChartDom(series, yAxis) {
      let dom = comparisonChartDom;
      let bottom = 64;
      bottom += Math.ceil(series.length / 8) * 30;
      dom.setOption({
        series,
        yAxis,
        grid: {
          bottom
        }
      }, {
        replaceMerge: ['series']
      });
    },
  }
}
</script>

<style scoped lang="scss">
.chart_key_select {
  display: flex;
  justify-content: space-between;
  margin: 4px 4% 8px 4%;
}

.chart_size {
  height: 480px;
  width: 100%;
}
</style>