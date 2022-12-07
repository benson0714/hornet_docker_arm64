---
title: 'resberrypi (arm64)下的hornet docker 建置'
disqus: Benson
---

resberrypi (arm64)下的hornet docker 建置
===

## 1. docker & ubuntu版本
### 1.1 docker setup
- 下載docker
```shell=
sudo apt-get update
sudo apt-get install docker
sudo apt-get install docker-compose
```
- 將自己的使用者帳號加入至 docker 群組：
```shell=
sudo usermod -aG docker your_account
```
### 1.2 docker版本截圖
```shell=
docker version
```
![](https://i.imgur.com/0H2miGD.png)


### 1.3 ubuntu版本截圖
```shell=
lsb_release -a
```
![](https://i.imgur.com/6eYgBAW.png)


## 2. gohornet 
### 2.1 clone:
- git clone hornet，取得需要的檔案
```shell=
git clone https://github.com/gohornet/hornet && cd hornet && git checkout production
```
## 3. 修改檔案
可參考 : https://github.com/benson0714/hornet_docker_windows

### 3.1 確認必要檔案
- 需要新增一個docker-compose.yml檔 (3.2)
- 需要新增config.json檔，可參考config_defaults.json去設定 (3.3)
- 官方網站會寫說需要手動新增**mainnetdb、snapshot、p2pstore**資料夾，略過此步驟，否則會衝突!!!

```shell=
.
├── config.json            <NEWLY ADDED FILE>
├── peering.json
├── profiles.json
├── docker-compose.yml      <NEWLY ADDED FILE>
```
### 3.2 修改docker-compose.yml內容
```yaml=
version: '3'
services:
  hornet:
    container_name: hornet
    image: gohornet/hornet:latest
    restart: always
    ulimits:
      nofile:
        soft: 8192
        hard: 8192
    stop_grace_period: 5m
    cap_drop:
      - ALL
    ports:
      - "14265:14265/tcp"
      - "8081:8081/tcp"
      - "8091:8091/tcp"
      - "9029:9029/tcp"
    volumes:
      - ./config.json:/app/config.json:ro
      - ./peering.json:/app/peering.json
      - ./profiles.json:/app/profiles.json:ro
      - ./mainnetdb:/app/mainnetdb
      - ./p2pstore:/app/p2pstore
      - ./snapshots/mainnet:/app/snapshots/mainnet
```
### 3.3 修改config.json內容
- dashboard會需要先創密碼才可以進行登入，可參考 :
https://wiki.iota.org/hornet/how_tos/using_docker#create-username-and-password-for-the-hornet-dashboard
```json
{
  "restAPI": {
    "bindAddress": "0.0.0.0:14265",
    "jwtAuth": {
      "salt": "HORNET"
    },
    "publicRoutes": [
      "/health",
      "/mqtt",
      "/api/v1/info",
      "/api/v1/tips",
      "/api/v1/messages*",
      "/api/v1/transactions*",
      "/api/v1/milestones*",
      "/api/v1/outputs*",
      "/api/v1/addresses*",
      "/api/v1/treasury",
      "/api/v1/receipts*",
      "/api/plugins/participation/events*",
      "/api/plugins/participation/outputs*",
      "/api/plugins/participation/addresses*"
    ],
    "protectedRoutes": [
      "/api/v1/*",
      "/api/plugins/*"
    ],
    "powEnabled": true,
    "powWorkerCount": 1,
    "limits": {
      "bodyLength": "1M",
      "maxResults": 1000
    }
  },
  "dashboard": {
    "bindAddress": "0.0.0.0:8081",
    "dev": false,
    "auth": {
      "sessionTimeout": "72h",
      "username": "admin",
      "passwordHash": "749907e87de741c73ff12f488c4a86da463bf8070fd62408b9eb7e5280f67312",
      "passwordSalt": "960d8c26555fd43cfd274e8403852bed52a1573f0f648552aacb958ed9fc0fc9"
    }
  },
  "db": {
    "engine": "rocksdb",
    "path": "mainnetdb",
    "autoRevalidation": false
  },
  "snapshots": {
    "depth": 50,
    "interval": 200,
    "fullPath": "snapshots/mainnet/full_snapshot.bin",
    "deltaPath": "snapshots/mainnet/delta_snapshot.bin",
    "deltaSizeThresholdPercentage": 50.0,
    "downloadURLs": [
      {
        "full": "https://chrysalis-dbfiles.iota.org/snapshots/hornet/latest-full_snapshot.bin",
        "delta": "https://chrysalis-dbfiles.iota.org/snapshots/hornet/latest-delta_snapshot.bin"
      },
      {
        "full": "https://cdn.tanglebay.com/snapshots/mainnet/full_snapshot.bin",
        "delta": "https://cdn.tanglebay.com/snapshots/mainnet/delta_snapshot.bin"
      }
    ]
  },
  "pruning": {
    "milestones": {
      "enabled": false,
      "maxMilestonesToKeep": 60480
    },
    "size": {
      "enabled": true,
      "targetSize": "30GB",
      "thresholdPercentage": 10.0,
      "cooldownTime": "5m"
    },
    "pruneReceipts": false
  },
  "protocol": {
    "networkID": "chrysalis-mainnet",
    "bech32HRP": "iota",
    "minPoWScore": 4000.0,
    "milestonePublicKeyCount": 2,
    "publicKeyRanges": [
      {
        "key": "a9b46fe743df783dedd00c954612428b34241f5913cf249d75bed3aafd65e4cd",
        "start": 0,
        "end": 777600
      },
      {
        "key": "365fb85e7568b9b32f7359d6cbafa9814472ad0ecbad32d77beaf5dd9e84c6ba",
        "start": 0,
        "end": 1555200
      },
      {
        "key": "ba6d07d1a1aea969e7e435f9f7d1b736ea9e0fcb8de400bf855dba7f2a57e947",
        "start": 552960,
        "end": 2108160
      },
      {
        "key": "760d88e112c0fd210cf16a3dce3443ecf7e18c456c2fb9646cabb2e13e367569",
        "start": 1333460,
        "end": 2888660
      },
      {
        "key": "7bac2209b576ea2235539358c7df8ca4d2f2fc35a663c760449e65eba9f8a6e7",
        "start": 2108160,
        "end": 3359999
      },
      {
        "key": "edd9c639a719325e465346b84133bf94740b7d476dd87fc949c0e8df516f9954",
        "start": 2888660,
        "end": 3359999
      },
      {
        "key": "47a5098c696e0fb53e6339edac574be4172cb4701a8210c2ae7469b536fd2c59",
        "start": 3360000,
        "end": 0
      },
      {
        "key": "ae4e03072b4869e87dd4cd59315291a034493a8c202b43b257f9c07bc86a2f3e",
        "start": 3360000,
        "end": 0
      }
    ]
  },
  "pow": {
    "refreshTipsInterval": "5s"
  },
  "requests": {
    "discardOlderThan": "15s",
    "pendingReEnqueueInterval": "5s"
  },
  "receipts": {
    "backup": {
      "enabled": false,
      "path": "receipts"
    },
    "validator": {
      "validate": false,
      "ignoreSoftErrors": false,
      "api": {
        "address": "http://localhost:14266",
        "timeout": "5s"
      },
      "coordinator": {
        "address": "UDYXTZBE9GZGPM9SSQV9LTZNDLJIZMPUVVXYXFYVBLIEUHLSEWFTKZZLXYRHHWVQV9MNNX9KZC9D9UZWZ",
        "merkleTreeDepth": 24
      }
    }
  },
  "tangle": {
    "milestoneTimeout": "30s"
  },
  "tipsel": {
    "maxDeltaMsgYoungestConeRootIndexToCMI": 8,
    "maxDeltaMsgOldestConeRootIndexToCMI": 13,
    "belowMaxDepth": 15,
    "nonLazy": {
      "retentionRulesTipsLimit": 100,
      "maxReferencedTipAge": "3s",
      "maxChildren": 30,
      "spammerTipsThreshold": 0
    },
    "semiLazy": {
      "retentionRulesTipsLimit": 20,
      "maxReferencedTipAge": "3s",
      "maxChildren": 2,
      "spammerTipsThreshold": 30
    }
  },
  "node": {
    "alias": "HORNET mainnet node",
    "profile": "auto",
    "disablePlugins": [],
    "enablePlugins": [
      "Spammer"
    ]
  },
  "p2p": {
    "bindMultiAddresses": [
      "/ip4/0.0.0.0/tcp/15600",
      "/ip6/::/tcp/15600"
    ],
    "connectionManager": {
      "highWatermark": 10,
      "lowWatermark": 5
    },
    "gossip": {
      "unknownPeersLimit": 4,
      "streamReadTimeout": "1m0s",
      "streamWriteTimeout": "10s"
    },
    "db": {
      "path": "p2pstore"
    },
    "reconnectInterval": "30s",
    "autopeering": {
      "bindAddress": "0.0.0.0:14626",
      "entryNodes": [
        "/dns/lucamoser.ch/udp/14826/autopeering/4H6WV54tB29u8xCcEaMGQMn37LFvM1ynNpp27TTXaqNM",
        "/dns/entry-hornet-0.h.chrysalis-mainnet.iotaledger.net/udp/14626/autopeering/iotaPHdAn7eueBnXtikZMwhfPXaeGJGXDt4RBuLuGgb",
        "/dns/entry-hornet-1.h.chrysalis-mainnet.iotaledger.net/udp/14626/autopeering/iotaJJqMd5CQvv1A61coSQCYW9PNT1QKPs7xh2Qg5K2",
        "/dns/entry-0.mainnet.tanglebay.com/udp/14626/autopeering/iot4By1FD4pFLrGJ6AAe7YEeSu9RbW9xnPUmxMdQenC",
        "/dns/entry-1.mainnet.tanglebay.com/udp/14636/autopeering/CATsx21mFVvQQPXeDineGs9DDeKvoBBQdzcmR6ffCkVA"
      ],
      "entryNodesPreferIPv6": false,
      "runAsEntryNode": false
    }
  },
  "logger": {
    "level": "info",
    "disableCaller": true,
    "encoding": "console",
    "outputPaths": [
      "stdout"
    ]
  },
  "warpsync": {
    "advancementRange": 150
  },
  "spammer": {
    "message": "IOTA - A new dawn",
    "index": "HORNET Spammer",
    "indexSemiLazy": "HORNET Spammer Semi-Lazy",
    "cpuMaxUsage": 0.8,
    "mpsRateLimit": 0.0,
    "workers": 0,
    "autostart": false
  },
  "faucet": {
    "amount": 10000000,
    "smallAmount": 1000000,
    "maxAddressBalance": 20000000,
    "maxOutputCount": 127,
    "indexationMessage": "HORNET FAUCET",
    "batchTimeout": "2s",
    "powWorkerCount": 0,
    "website": {
      "bindAddress": "localhost:8091",
      "enabled": true
    }
  },
  "mqtt": {
    "bindAddress": "localhost:1883",
    "wsPort": 1888,
    "workerCount": 100
  },
  "profiling": {
    "bindAddress": "localhost:6060"
  },
  "prometheus": {
    "bindAddress": "localhost:9311",
    "fileServiceDiscovery": {
      "enabled": false,
      "path": "target.json",
      "target": "localhost:9311"
    },
    "databaseMetrics": true,
    "nodeMetrics": true,
    "gossipMetrics": true,
    "cachesMetrics": true,
    "restAPIMetrics": true,
    "migrationMetrics": true,
    "coordinatorMetrics": true,
    "mqttBrokerMetrics": true,
    "debugMetrics": false,
    "goMetrics": false,
    "processMetrics": false,
    "promhttpMetrics": false
  },
  "debug": {
    "whiteFlagParentsSolidTimeout": "2s"
  }
}
```
### 3.4 預先給定docker存取資料夾權限
```shell=
sudo chown 65532:65532 mainnetdb
sudo chown 65532:65532 p2pstore
sudo chown -R 65532:65532 snapshots
```
## 4. docker-compose 安裝
```shell=
docker-compose up -d
```

## 5. 參考
- hornet docker建置教學:
https://wiki.iota.org/hornet/how_tos/using_docker
- hornet參考github
https://github.com/gohornet/hornet
