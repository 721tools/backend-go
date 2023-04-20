## 代码结构介绍

```
├── README.md
├── cmd
│   └── indexer
│   └── xxxx
├── indexer
│   ├── implement
│   └── pkg
├── go.mod
└── go.sum

1. cmd 下的每个文件夹可单独打包成一个可执行文件，用于执行某一具体任务
2. cmd 下的每个具体任务在一级目录下同名文件夹内具体实现
3. config 文件夹中存放代码运行所需的配置文件（or 可从环境变量中读取）

```

## 编译

### indexer

```
go build -o ./bin/indexer cmd/indexer/main.go
```

## 部署
