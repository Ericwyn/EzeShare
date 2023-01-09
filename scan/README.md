# scan

策略模式实现不同扫描方法

入口在 scan.go 里面的 StartScan 方法

不同的实现有
- udpscan.go `局域网扫描` 
- blescan.go `BLE 扫描`