<template>
  <el-tooltip effect="dark" placement="bottom" :show-after="500">
    <template #content>
      <template v-if="sumStatsDataValue">
        <span v-for="(item, index) in statsData" :key="item.label.Key">
          <template v-if="item.label.Value || item.value.Value">
            <span>{{ `${item.label.Value ?? ''}: ${item.value.Value ?? ''} (${statsDataProgress[index]}%)` }}</span><br>
          </template>
        </span>
      </template>
      <template v-else>
        No Data
      </template>
    </template>
    <div class="statistic_group">
      <span v-for="(item, index) in statsData" :key="item.label.Key" :style="getSpanStyle(item, index)"
        class="statistic_item">
        <el-text v-if="item.value.Value" truncated>{{ `${item.label.Value ?? ''}: ${statsDataProgress[index]}%`
          }}</el-text>
      </span>
    </div>
  </el-tooltip>
</template>

<script>
import { useDasDevices } from '@/stores/das-devices';

export default {
  name: 'MyStatisticGroup',
  components: {},
  props: {
    owner: Object,
    data: Object,
  },
  setup() {
    const dev = useDasDevices().currentDevice;
    return {
      dev,
    };
  },
  computed: {
    colors() {
      return [
        "#3fb1e3cc", "#6be6c1cc", "#626c91cc", "#a0a7e6cc", "#c4ebadcc", "#96dee8cc"
      ]
    },
    statsData() {
      const out = [];
      this.data.Items.forEach(item => {
        if (item.Type === 'Component:Statistic') {
          const label = this.dev.params.getParam(item.Items[0].OID);
          const value = this.dev.params.getParam(item.Items[1].OID);
          if (label && value) {
            out.push({
              label: label,
              value: value
            });
          }
        }
      })
      return out;
    },
    sumStatsDataValue() {
      return this.statsData.reduce((accumulator, currentItem) => {
        let v = currentItem.value.Value ?? 0;
        return accumulator + v
      }, 0);
    },
    statsDataProgress() {
      let out = [];
      let sum = 0;
      let max = 0;
      let maxIndex = 0;
      this.statsData.forEach((currentItem, index) => {
        let v = this.getValueProgress(currentItem.value.Value);
        if (max <= v) {
          max = v;
          maxIndex = index;
        }
        sum += v;
        out.push(v);
      });
      if (sum != 100) {
        out[maxIndex] = 100 - (sum - out[maxIndex])
      }
      return out;
    },
  },
  methods: {
    getSpanStyle(item, index) {
      let colorIndex = index % 6;
      let color = this.colors[colorIndex];
      let height = item.value?.Height ?? '20px';
      let width = '0%';
      if (item.value?.Value) {
        width = (item.value?.Value / this.sumStatsDataValue) * 100 + '%';
      }
      if (this.sumStatsDataValue == 0) {
        width = (1 / this.statsData.length) * 100 + '%';
        color = '#c8c9cc';
      }
      return {
        height,
        background: color,
        width: width,
        'line-height': height
      }
    },
    getValueProgress(v) {
      if (!v) return 0;
      let p = Math.round((v / this.sumStatsDataValue) * 100);
      return p ? p : 1;
    }
  }
}
</script>
<style lang="scss" scoped>
.statistic_group {
  width: 640px;
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
}

.statistic_item {
  text-align: center;

  span {
    font-size: 10px;
    font-weight: bold;
    color: #fff;
  }

  &:first-child {
    border-bottom-left-radius: 4px;
    border-top-left-radius: 4px;
  }

  &:last-child {
    border-bottom-right-radius: 4px;
    border-top-right-radius: 4px;
  }
}
</style>
