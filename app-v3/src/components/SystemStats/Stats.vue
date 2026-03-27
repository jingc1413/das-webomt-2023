<template>
  <div v-loading="updateLoading">
    <el-row v-if="!isPrintMode">
      <el-col :span="23" style="text-align: right;margin-top: 12px;">
        <div class="toolbar">
          <el-date-picker v-model="eventTime" type="datetimerange" range-separator="To" start-placeholder="Start date"
            end-placeholder="End date" :shortcuts="shortcuts" :disabled-date="disabledEventTime" value-format="X"
            time-format="HH:mm" :clearable="false" @change="getChartData([], true)" />
          <el-button style="margin-left: 8px;" @click="getChartData([], true)">Refresh</el-button>
          <el-button @click="exportChartData()" type="primary" plain>Export</el-button>
        </div>
      </el-col>
    </el-row>
    <template v-if="!isPrintMode">
      <el-tabs v-model="activeTabName" class="chart_tab" @tab-change="resizeChartDom()">
        <el-tab-pane v-for="tabItem in chartTabs" :key="tabItem.key" :label="tabItem.name" :name="tabItem.key">
          <el-row class="chart_tab_body">
            <template v-if="tabItem.key != 'comparison'">
              <el-col v-for="chartKey in tabItem.chartKeys" :span="chartMetrics[chartKey]?.span ?? 12" :key="chartKey">
                <div :id="chartDomIdHeader + chartKey" class="chart_size"
                  :style="{ height: chartMetrics[chartKey]?.height ?? '320px', width: '100%' }">
                </div>
              </el-col>
            </template>
            <template v-else>
              <el-col>
                <comparison-metrics-view :beginTime="beginTime" :endTime="beginTime" :chartKeys="tabItem.chartKeys"
                  :activeTabName="activeTabName">
                </comparison-metrics-view>
              </el-col>
            </template>
          </el-row>
        </el-tab-pane>
      </el-tabs>
    </template>

    <template v-else>
      <template v-for="tabItem in chartTabs">
        <el-row v-if="tabItem.key != 'comparison'" :key="tabItem.key" style="width: calc(100% - 30px);">
          <el-col>
            <el-divider content-position="left">
              <h3>{{ tabItem.name }}</h3>
            </el-divider>
          </el-col>
          <el-col v-for="chartKey in tabItem.chartKeys" :span="chartMetrics[chartKey]?.span ?? 12" :key="chartKey">
            <div :id="chartDomIdHeader + chartKey" class="chart_size"
              :style="{ height: chartMetrics[chartKey]?.height ?? '320px' }" />
            <img v-if="isPrintMode && chartImgList[chartKey]" :src="chartImgList[chartKey]" style="width: 100%;">
          </el-col>
        </el-row>
      </template>
    </template>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import provideKeys from '@/utils/provideKeys.js'
import { dayjs } from 'element-plus';
import { useDasDevices } from "@/stores/das-devices";
import { useWindowSize } from '@vueuse/core'
import { utils, writeFile } from 'xlsx'
import model from '@/stores/model'

import ComparisonMetricsView from "./ComparisonMetricsView.vue";


var utc = require('dayjs/plugin/utc');
dayjs.extend(utc);

let chartDomList = {};

let timestampWith14DaysAgo = dayjs().subtract(14, 'd').valueOf();

export default {
  name: 'MyStatisticsPage',
  components: ['ComparisonMetricsView'],
  inject: ['viewMode'],
  setup() {
    let dasDevices = useDasDevices();
    const dev = dasDevices.currentDevice;
    const { width: windowWidth, height: windowHeight } = useWindowSize();
    return { dev, dasDevices, windowWidth, windowHeight }
  },
  data() {
    return {
      chartDomIdHeader: dayjs().format('X'),
      beginTime: null,
      endTime: null,
      sub: '',
      chartImgList: {},
      chartMetrics: {},
      chartTabs: {},
      updateLoading: false,
      resizeTimer: null,
      activeTabName: 'system'
    }
  },
  computed: {
    eventTime: {
      set(value) {
        if (value) {
          this.beginTime = value[0];
          this.endTime = value[1];
        } else {
          this.beginTime = null;
          this.endTime = null;
        }
      },
      get() {
        return [this.beginTime, this.endTime]
      }
    },
    isPrintMode() {
      return this.viewMode != provideKeys.viewModeDefaultValue
    },
    shortcuts() {
      return [
        {
          text: '24 Hours',
          value: () => {
            let nowTime = dayjs().set('m', 0).set('s', 0).set('ms', 0).add(1, 'h')
            const end = nowTime.valueOf();
            const start = nowTime.subtract(1, 'd').valueOf();
            return [start, end]
          },
        },
        {
          text: '48 Hours',
          value: () => {
            let nowTime = dayjs().set('m', 0).set('s', 0).set('ms', 0).add(1, 'h')
            const end = nowTime.valueOf();
            const start = nowTime.subtract(2, 'd').valueOf();
            return [start, end]
          },
        },
        {
          text: '3 Days',
          value: () => {
            let nowTime = dayjs().set('h', 0).set('m', 0).set('s', 0).set('ms', 0).add(1, 'd')
            const end = nowTime.valueOf();
            const start = nowTime.subtract(3, 'd').valueOf();
            return [start, end]
          },
        },
        {
          text: '7 Days',
          value: () => {
            let nowTime = dayjs().set('h', 0).set('m', 0).set('s', 0).set('ms', 0).add(1, 'd')
            const end = nowTime.valueOf();
            const start = nowTime.subtract(7, 'd').valueOf();
            return [start, end]
          },
        },
        {
          text: '14 Days',
          value: () => {
            let nowTime = dayjs().set('h', 0).set('m', 0).set('s', 0).set('ms', 0).add(1, 'd')
            const end = nowTime.valueOf();
            const start = nowTime.subtract(14, 'd').valueOf();
            return [start, end]
          },
        },
      ]
    }
  },
  watch: {
    windowWidth() {
      if (this.resizeTimer) {
        clearTimeout(this.resizeTimer);
      }
      this.resizeTimer = setTimeout(() => {
        this.resizeChartDom()
      }, 200);
    },
    windowHeight() {
      if (this.resizeTimer) {
        clearTimeout(this.resizeTimer);
      }
      this.resizeTimer = setTimeout(() => {
        this.resizeChartDom()
      }, 200);
    }
  },
  created() {
    this.changeDataTime(1, false);
  },
  async mounted() {
    this.sub = this.$route.params.sub;
    await this.getChartKeys();
    this.batchActionByChartKey(({ key, name, unit, max, min, yInterval, bottom }) => {
      this.initChartDom({ domId: key, chartName: name, unit, max, min, yInterval, bottom });
    })
    this.getChartData([], false);
  },
  beforeUnmount() {
    this.batchActionByChartKey(({ key }) => {
      chartDomList[this.chartDomIdHeader + key] && chartDomList[this.chartDomIdHeader + key].dispose();
      delete chartDomList[this.chartDomIdHeader + key];
    })
  },
  methods: {
    disabledEventTime(time) {
      if (time.getTime() > Date.now()) {
        return true;
      }
      if (time.getTime() < timestampWith14DaysAgo) {
        return true
      }
      return false;
    },
    changeDataTime(day, isRefresh = true) {
      let nowTime = dayjs().set('m', 0).set('s', 0).set('ms', 0).add(1, 'h')
      this.endTime = nowTime.unix();
      this.beginTime = nowTime.subtract(day, 'd').unix();
      isRefresh && this.getChartData([], true)
    },
    getChartKeys() {
      return new Promise((resolve, reject) => {
        this.chartTabs = this.dev.stats.getTabs(this.sub);
        this.chartMetrics = this.dev.stats.getMetrics(this.sub);
        resolve(null);
      })
    },
    getEchartOptionWithXAxisTime({ id, title, xAxisName = '', yAxisName = '', unit = '', max, min, yInterval, bottom = 64 }) {
      let option = {
        tooltip: {
          trigger: 'axis',
          confine: true,
        },
        title: {
          left: 'center',
          text: title
        },
        grid: {
          left: 64,
          right: 64,
          top: 32,
          bottom: bottom,
        },
        legend: {
          show: false,
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
          name: xAxisName
        },
        yAxis: {
          type: 'value',
          name: yAxisName + ' ' + unit,
          scale: true,
          max,
          min,
          interval: yInterval,
        },
        series: [{
          id: id,
          name: title,
          type: 'line',
          data: [],
          smooth: true,
          showSymbol: false,
        }]
      };
      return option;
    },
    initChartDom({ domId, chartName = '', unit, max, min, yInterval, bottom }) {
      let dom = echarts.init(document.getElementById(this.chartDomIdHeader + domId), null, { locale: 'EN' });
      let option = this.getEchartOptionWithXAxisTime({ id: domId, title: chartName, unit: unit, max, min, yInterval, bottom });
      chartDomList[this.chartDomIdHeader + domId] = dom;
      dom.setOption(option);
    },
    resizeChartDom() {
      this.batchActionByChartKey(({ key }) => {
        this.$nextTick(() => {
          chartDomList[this.chartDomIdHeader + key] && chartDomList[this.chartDomIdHeader + key].resize();
        })
      })
    },
    getChartData(keys = [], showMessage) {
      this.updateLoading = true
      if (keys.length == 0) {
        this.batchActionByChartKey(({ key }) => {
          keys.push(key)
        })
      }
      let query = {
        beginTime: this.beginTime,
        endTime: this.endTime,
        keys: keys,
        sub: this.sub,
        showMessage: showMessage,
        isPrintMode: this.isPrintMode
      }
      this.dev.stats.queryStats(query).then(res => {
        // console.log('queryDeviceStats', {query}, res)
        if (res) {
          this.filterChartDate(res);
        }
        setTimeout(() => {
          this.updateLoading = false
        }, 1000);
      }).catch((e) => {
        console.error(e);
        setTimeout(() => {
          this.updateLoading = false
        }, 1000);
      })
    },
    filterChartDate(chartDates) {
      this.batchActionByChartKey(({ key, items, name }) => {
        let cpItems = items;
        if (!cpItems) {
          cpItems = [{ key, name }]
        }
        let itemsData = [];
        let firstTime = dayjs(this.beginTime, 'X').format('YYYY-MM-DD HH:mm:ss');
        let lastTime = dayjs(this.endTime, 'X').format('YYYY-MM-DD HH:mm:ss');

        let addHeaderTime = dayjs(this.beginTime, 'X').subtract(1, 'hour').format('YYYY-MM-DD HH:mm:ss');
        let addEndTime = dayjs(this.endTime, 'X').add(1, 'hour').format('YYYY-MM-DD HH:mm:ss');
        cpItems.forEach(({ key: itemKey, name: itemName }) => {
          let chartData = chartDates[itemKey] ?? [];
          if (chartData.length == 0) {
            chartData.push([firstTime, null])
            chartData.push([lastTime, null])
            itemsData.push({ key: itemKey, chartData, name: itemName })
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
          itemsData.push({ key: itemKey, chartData, name: itemName })
        })
        this.updateChartDate(key, itemsData)
      })
    },
    updateChartDate(domId, seriesData) {
      let series = [];
      seriesData.forEach(({ key, chartData, name }) => {
        series.push({
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
          }
        })
      })
      let legendShow = series.length > 1;
      let dom = chartDomList[this.chartDomIdHeader + domId];
      dom.setOption({
        series,
        legend: {
          show: legendShow
        },
      }, {
        replaceMerge: ['series']
      });
      if (this.isPrintMode) {
        let self = this;
        dom.on('finished', function () {
          self.chartImgList[domId] = dom.getDataURL({ pixelRatio: 2, backgroundColor: '#fff' });
          document.getElementById(self.chartDomIdHeader + domId).style.display = "none";
        });
      }
    },
    batchActionByChartKey(cb) {
      for (const chartTabKey in this.chartTabs) {
        this.chartTabs[chartTabKey].chartKeys.forEach(keysItem => {
          if (cb && this.chartMetrics[keysItem]) {
            cb(this.chartMetrics[keysItem])
          }
        })
      }
    },
    exportChartData() {
      let data = this.dev.stats.getChartData();
      const workbook = utils.book_new();
      let keyMap = {};
      this.batchActionByChartKey(({ key, name, items }) => {
        let keys = [{ key, name }];
        if (items) {
          keys = items.map(item => { return { key: item.key, name: item.name } });
        }
        keys.forEach(({ key: itemKey, name: itemName }) => {
          if (itemName != name) {
            itemName = itemName + ' ' + name
          }
          keyMap[itemKey] = itemName
        })
      })
      let keys = Object.keys(keyMap);
      // console.log({keys}, this.deviceStats.sourceData)
      let dataJson = this.dev.stats.sourceData.map(item => {
        let filterItem = {}
        if (item['t']) {
          filterItem["Time"] = model.unixTimestampWithoutTimezones2nowTime(item['t'], false);
        }
        keys.forEach(key => {
          if (item[key] != undefined) {
            filterItem[keyMap[key]] = item[key]
          }
        })
        return filterItem
      })
      let worksheet = utils.json_to_sheet(dataJson)
      utils.book_append_sheet(workbook, worksheet);
      writeFile(workbook, this.dasDevices.currentDeviceFileName('Stats') + '.xlsx', { compression: true });
    }
  }
}

</script>

<style lang="scss" scoped>
.chart_size {
  height: 320px;
  width: 100%;
}

.chart_tab_body {
  width: calc(100% - 42px);
  margin-top: 24px;
}

.chart_tab {
  margin-left: 10px;

  :deep(.el-tabs__content) {
    height: calc(100vh - 140px);
    overflow-x: hidden;
    overflow-y: auto;
  }
}
</style>
