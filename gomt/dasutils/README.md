# dasutil工具使用说明

---
## 启动方式

默认提供windows/linux平台可执行文件 dasutil.exe dasutil
```shell
dasutil.exe -device-addr=10.7.3.164 -schema=corning
dasutil -device-addr=10.7.3.164 -schema=corning

可选参数
-csv-path="/home/csvdir/" 
```

## 输出样式
```golang
GET: DeviceTypeName=Primary A3, SubID=0, Combiner=1
GET: SerialNumber=********************(00), DownlinkFrequencyStart=862000(00), DownlinkFrequencyEnd=894000(00)
Query ModuleTypeId from combined data csv file empty, TB2-COMBINER1.P00E4=********************(00)

GET: DeviceTypeName=Primary A3, SubID=0, Combiner=2
GET: SerialNumber=poi-3(00), DownlinkFrequencyStart=859000(00), DownlinkFrequencyEnd=894000(00)
Query ModuleTypeId from combined data csv file empty, TB2-COMBINER2.P04E4=poi-3(00)

GET: DeviceTypeName=Primary A3, SubID=0, Combiner=3
GET: SerialNumber=(00), DownlinkFrequencyStart=0(00), DownlinkFrequencyEnd=0(00)
Query ModuleTypeId from combined data csv file empty, TB2-COMBINER3.P08E4=(00)

GET: DeviceTypeName=Primary A3, SubID=0, Combiner=4
GET: SerialNumber=(00), DownlinkFrequencyStart=0(00), DownlinkFrequencyEnd=0(00)
Query ModuleTypeId from combined data csv file empty, TB2-COMBINER4.P0CE4=(00)
```
## 云文档下载地址
### [dasutil下载](http://doccloud.sunwave.com.cn/index.html#doc/enterprise/4792?key=17210373851589726)

## 说明
> 程序使用默认的内嵌csv文件, combined_data.csv dup_data.csv dup_id.csv支持自定义路径<br>
> 修改后支持外部替换csv文件