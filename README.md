# EzeShare
EzeShare 一个跨平台文件传输项目，支持 Pc / Android

该 repo 是 PC 端实现，目前暂时只支持命令行操作

- Sender
  - 自动扫描并发现同个网络内其他节点
- Receiver
  - 支持鉴权 / 接收

## 使用教程
### 编译
```shell
go mod tidy
go build ./EzeShare.go
```
### 发送
```shell
./EzeShare -sender -f "/dev/Downloads/test-video.mp4"
```

shell 显示
```shell
root@mini-godzilla:/opt/EzeShare# ./EzeShare -sender -f "README.md"
  _____         ____  _
 | ____|_______/ ___|| |__   __ _ _ __ ___
 |  _| |_  / _ \___ \| '_ \ / _` | '__/ _ \
 | |___ / /  __/___) | | | | (_| | | |  __/
 |_____/___\___|____/|_| |_|\__,_|_|  \___|

          run in mode: Sender

当前设备 IP 为  192.168.199.100










当前 receiver 列表如下:
         address                 name                    deviceId                deviceType
[0]      192.168.199.10          Desktop-Godzilla        12f0ea82        windows
-----------------------
输入编号并回车, 选择具体 receiver
0
[01-15 22:54:39] [I] 上传 README.md, 进度: [==================== ]
root@mini-godzilla:/opt/EzeShare#
```


### 接收
```shell
./EzeShare -receiver
```

shell 显示
```shell
PS D:\Chaos\go\EzeShare> go run .\EzeShare.go -receiver
  _____         ____  _
 | ____|_______/ ___|| |__   __ _ _ __ ___
 |  _| |_  / _ \___ \| '_ \ / _` | '__/ _ \
 | |___ / /  __/___) | | | | (_| | | |  __/
 |_____/___\___|____/|_| |_|\__,_|_|  \___|

          run in mode: Receiver

[01-15 22:54:36] [I] 使用历史 IP 设置 : 192.168.199.10
[01-15 22:54:36] [I] start http server in Addr 192.168.199.10:23019
[01-15 22:54:39] [I] save file to new name: README.md(8)
[01-15 22:54:39] [I] save file success, filePath : C:\Users\Ericwyn\Downloads\EzeShareFiles/README.md(8)
```


## 基本发送逻辑
 - Receiver 开启一个 EzeShareServer (Http Port: 23019, Udp Port: 23010)
 - Receiver 广播自己的 ip 地址 (向 255.255.255.255:23010 发送 udp 广播)
 - Sender 扫描得到 Receiver 列表
 - Sender 选择一个 Receiver，访问 `premReq` 接口发送数据
 - Receiver 端手动确认，返回一个 token (经过 Sender 公钥加密) 
 - Sender 私钥解密出来 token，md5 计算得到 sign，将 sign 以及需要发送的文件吗，通过 `fileTransfer` 接口发送过去
 - Receiver 校验、接收数据

## 接收方接口
- `/api/premReq`

  - 请求
    - type
      - once 一次
      - always 永远
    - 文件名称
    - 文件大小
    - 发送方名称
    - 发送方公钥
  - 返回
    - 发送方公钥加密后的 token
      - token 由接收方使用 (自己密钥 + 接收方名称 + 接收方文件) md5 得到
    - type: once / always

- `/api/fileTransfer`

  - 往这个接口发送数据就可以了

  - 请求

    - key-type : always 还是 once
    - 接收方公钥加密后的 token
    - 文件


## 迭代计划

- 一期
  - 在同个内网环境下实现
  - PC 端支持命令行发送数据，或者是右键发送数据
  - webui 方案
    - 通过命令行调用浏览器展示某个 url 的形式，来展示特定 ui
- 二期
  - gtk ui 或者其他 ui