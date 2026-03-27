export function newDeviceImportor(sub, params) {
  const m = {
    sub,
    params,
  };

  m.getFixedTableColumnes = function (tableCols) {
    const out = [];
    tableCols.forEach((col) => {
      if (col.Style?.fixed) {
        out.push(col.Key);
      }
    });
    return out;
  };
  m.importTableDataFromRaw = function (table, raw) {
    if (!table || !raw || raw.length <= 1) {
      return undefined;
    }
    const inputHeader = raw[0];
    const inputRows = raw.slice(1);
    const inputValues = [];
    for (const i in inputRows) {
      const inputRow = inputRows[i];
      if (inputRow) {
        const values = {};
        for (const j in inputHeader) {
          const key = inputHeader[j];
          const col = table.Items.find((v) => v.Key === key);
          if (col && inputRow[j] !== undefined) {
            values[col.Key] = inputRow[j].toString();
          }
        }
        inputValues.push(values);
      }
    }
    const tableData = this.makeTableDataForValues(table, inputValues);
    const tableColumneKeys = inputHeader.filter((key) => {
      let count = 0;
      tableData.forEach((rowData) => {
        if (rowData[key] != undefined) {
          count += 1;
        }
      });
      return count > 0;
    });
    return { columneKeys: tableColumneKeys, data: tableData };
  };
  m.makeTableDataForValues = function (table, inputValues) {
    const uniqueKey = table.Unique;
    const outputData = [];
    if (!uniqueKey) {
      return outputData;
    }
    inputValues.forEach((rowValues) => {
      const tableDataRow = table.Data?.find((rowScope) => {
        const tmp = rowScope[uniqueKey]?.Value;
        return tmp != undefined && tmp === rowValues[uniqueKey];
      });
      if (tableDataRow) {
        const outputDataRow = {};
        const uniqueItem = tableDataRow[uniqueKey];
        outputDataRow[uniqueKey] = { Type: "Text", Value: uniqueItem.Value, Key: uniqueKey };
        Object.keys(tableDataRow).forEach((key) => {
          const item = JSON.parse(JSON.stringify(tableDataRow[key]));
          if (item.Style?.fiexed) {
            outputDataRow[key] = item;
          } else {
            if (item.Type === "Param" && (item.Access === "rw" || item.Access === "wo")) {
              const param = this.params.getParam(item.OID);
              if (param) {
                item.Value = rowValues[key];
                item.isChanged = String(item.Value) !== String(param.Value);
                item.validate = param.validateInputValue(item.Value);
                outputDataRow[key] = item;
              }
            }
            if (item.Type?.startsWith("Component:ParamGroup") && (item.Access === "rw" || item.Access === "wo")) {
              item.Value = rowValues[key];
              if (item.Value) {
                const _values = item.Value.split(",");
                for (const i in item.Items) {
                  const _value = _values[i];
                  if (_value !== undefined) {
                    const param = this.params.getParam(item.Items[i].OID);
                    if (param) {
                      item.Items[i].Value = _value;
                      item.Items[i].isChanged = String(_value) !== String(param.Value);
                      item.Items[i].validate = param.validateInputValue(_value);
                      if (!item.isChanged) {
                        item.isChanged = item.Items[i].isChanged
                      }
                      if (!item.validate) {
                        item.validate = item.Items[i].validate
                      }
                    }
                  }
                }
              }
              outputDataRow[key] = item;
            }
          }
        });
        outputData.push(outputDataRow);
      }
    });
    return outputData;
  };
  return m;
}
