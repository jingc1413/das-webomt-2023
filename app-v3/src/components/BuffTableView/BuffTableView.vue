<template>
  <el-row style="width: 100%;">
    <el-col>
      <el-form :inline="true" :model="{}" class="demo-form-inline">
        <el-form-item label="Offset">
          <el-input v-model.trim="tableBufferOffset" placeholder="Offset" maxlength="4" style="width: 80px;"
            :formatter="(value) => value.toLocaleUpperCase()" :parser="(value) => value.replace(/[^0-9A-Fa-f]/g, '')"
            @blur="updateBufferOffset()" />
        </el-form-item>
        <el-form-item label="Length">
          <el-select v-model="tableBufferLength" placeholder="Length" style="width:60px">
            <el-option v-for="lengthItem in lengthOption" :key="lengthItem" :label="lengthItem" :value="lengthItem" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-col>
    <el-col style="margin-top: 12px; maxHeight: calc(100vh - 280px);">
      <el-table :data="tableData" style="width: 100%" height="100%"
        :header-cell-style="{backgroundColor: '#f4f4f5'}"
        :cell-style="tableCellStyle">
        <el-table-column prop="offset" label="Offset" width="60" fixed align="center">
          <template #default="scope">
            {{ getHexString(bufferOffset + scope.$index * rowSize, 4) }}
          </template>
        </el-table-column>
        <el-table-column v-for="(_, lenIndex) in rowSize" :key="lenIndex" :prop="lenIndex.toString()"
          :label="getHexString(lenIndex)" width="50" align="center">
          <template #default="scope">
            <span :class="lenIndex % 2 ? 'buff_string1' : 'buff_string2'"
              @dblclick="switchInput(scope.$index * rowSize, lenIndex, true)"
              v-show="!inputState[scope.$index * rowSize + lenIndex].value && ((scope.$index * rowSize + lenIndex) < bufferLength)">
              <!-- {{ getHexString(bufferOffset + scope.$index * tableColSize+ lenIndex) }} -->
              {{ tableData[scope.$index][lenIndex]?.value ?? '--' }}
            </span>
            <el-input v-if="inputState[scope.$index * rowSize + lenIndex].value" ref="buffInput"
              v-model.trim="tableData[scope.$index][lenIndex].value" maxlength="2" style="width: 30px;"
              :formatter="(value) => value.toLocaleUpperCase()"
              :parser="(value) => value.replace(/[^0-9A-Fa-f]/g, '').toLocaleUpperCase()"
              @blur="switchInput(scope.$index * rowSize, lenIndex, false)" />
          </template>
        </el-table-column>
      </el-table>
    </el-col>
  </el-row>
</template>

<script>
export default {
  props: {
    bufferData: {
      type: Array,
      default: () => {
        return []
      }
    },
    bufferOffset: {
      type: Number,
      default: 0,
    },
    bufferLength: {
      type: Number,
      default: 4,
    },
    colSize: {
      type: Number,
      default: 4,
    },
    rowSize: {
      type: Number,
      default: 8,
    },
  },
  emits: ['update:bufferData', 'update:bufferOffset', 'update:bufferLength'],
  setup() {
  },
  data() {
    let filterBufferData = this.bufferData.map(item => {
      let offset = item.offset;
      let value = item.value;
      if (value !== null && value !== undefined) {
        value = this.getHexString(value);
      }
      return {
        offset,
        value
      };
    })
    return {
      lengthOption: [4, 8, 16, 32],
      tableBufferOffset: this.getHexString(this.bufferOffset || 0, 4),
      inputState: this.chunkArray([], this.rowSize, false),
      tableData: this.chunkArray(filterBufferData, this.rowSize),
    }
  },
  computed: {
    tableBufferLength: {
      get() {
        return this.bufferLength;
      },
      set(val) {
        console.log('tableBufferLength:', val);
        this.$emit('update:bufferLength', val)
      }
    }
  },
  watch: {
    bufferOffset() {
      if (this.getHexString(this.bufferOffset || 0, 4) != this.tableBufferOffset) {
        this.tableBufferOffset = this.getHexString(this.bufferOffset || 0, 4);
      }
    },
    bufferData() {
      this.reInitTableData();
    },
    bufferLength() {
      this.reInitTableData();
    }
  },
  created() {
  },
  methods: {
    reInitTableData() {
      let filterBufferData = this.bufferData.map(item => {
        let offset = item.offset;
        let value = item.value;
        if (value !== null && value !== undefined) {
          value = this.getHexString(value)
        }
        return {
          offset,
          value
        };
      })
      console.log('reInitTableData', { filterBufferData });
      this.inputState = this.chunkArray(this.inputState, this.rowSize, false),
        this.tableData = this.chunkArray(filterBufferData, this.rowSize)
    },
    getHexString(num, stringLength = 2) {
      return Number(num).toString(16).padStart(stringLength, '0').toLocaleUpperCase()
    },
    updateBufferOffset() {
      this.$emit('update:bufferOffset', parseInt(this.tableBufferOffset || 0, 16))
    },
    chunkArray(arrayBody = [], chunkSize, isChunk = true) {
      let length = this.colSize * this.rowSize;
      let valueIndex = this.bufferOffset;
      let offset = 0;
      let endIndex = this.bufferOffset + length;
      while ((valueIndex + offset) < endIndex) {
        let tempItem = arrayBody[offset];
        let trueIndex = valueIndex + offset;
        if (!tempItem) {
          arrayBody.push({
            offset: trueIndex,
            value: null
          });
        }
        tempItem = arrayBody[offset];
        if (tempItem.offset > trueIndex) {
          arrayBody.splice(offset, 0, {
            offset: trueIndex,
            value: null
          });
        }
        offset++;
      }
      while (arrayBody.length > length) {
        arrayBody.pop();
      }
      if (!isChunk) {
        return arrayBody
      }
      return arrayBody.reduce((chunks, item, index) => {
        const chunkIndex = Math.floor(index / chunkSize);
        if (!chunks[chunkIndex]) {
          chunks[chunkIndex] = [];
        }
        chunks[chunkIndex].push(item);
        return chunks;
      }, []);
    },
    unChunkArray(arrayBody) {
      return arrayBody.reduce((chunks, item, index) => {
        if (Array.isArray(item)) {
          chunks.push(...item);
        } else {
          chunks.push(item);
        }
        return chunks;
      }, []);
    },
    switchInput(startOffset, index, flag) {
      if (startOffset + index >= this.bufferLength) {
        return
      }
      this.inputState[startOffset + index].value = flag
      if (flag == false) {
        let temp = this.unChunkArray(this.tableData).map(item => {
          let offset = item.offset;
          let value = item.value;
          if (value !== null && value !== undefined) {
            value = parseInt(value, 16);
          }
          return {offset, value};
        });
        this.$emit('update:bufferData', temp.slice(0, this.bufferLength))
      }
      this.$nextTick(() => {
        if (flag == true) {
          this.$refs['buffInput'][0]?.focus();
        }
      })
    },
    tableCellStyle({row, column, rowIndex, columnIndex}) {
      if (columnIndex == 0) return {backgroundColor: '#f4f4f5'};
      return {};
    }
  }
}
</script>

<style scoped>
.buff_string1 {
  color: #606266;
}

.buff_string2 {
  color: #909399;
}

.buff_string1,
.buff_string2 {
  font-size: 12px;
  font-weight: bold;
}
</style>