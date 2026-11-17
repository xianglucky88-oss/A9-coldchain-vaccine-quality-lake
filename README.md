# A9 疫苗冷链运输温控数据湖

面向疫苗冷链运输的温控数据质量基线，接入温度记录仪、车门事件和运输交接数据，输出温度偏移窗口、开门归因、平均动力学温度与数据质量摘要。

## 运行

```bash
go test ./...
go run ./cmd/coldchain-lake demo --data-dir .demo-data
```

真实接入：

```bash
go run ./cmd/coldchain-lake ingest --samples testdata/samples.jsonl --events testdata/events.jsonl --data-dir data
go run ./cmd/coldchain-lake compute --data-dir data --min-temp 2 --max-temp 8
```

## 目录与职责

```text
internal/model/       温控样本与车门事件模型
internal/engine/      时钟对齐、去重、偏移窗口、开门归因和 MKT
internal/store/       原始记录与派生结果持久化
internal/service/     接入与计算用例编排
internal/demo/        Demo 数据
cmd/coldchain-lake/   参数解析与输出
```

## 单一职责说明

- 模型只负责数据结构与校验。
- 归一化、偏移检测、事件归因和 MKT 计算相互独立。
- 原始记录与派生结果由不同存储对象管理。
- CLI 不包含温度规则。
