# helays/utils

通用 Go 工具包集合 —— 一组实用、模块化的基础库，用于加速后端服务和工具类开发。模块以 `helay.net/go/utils/v3` 作为 module 前缀导入，包含文件/IO 工具、模板助手、HTTP 层辅助、文件存储驱动、数据库辅助、规则验证引擎、锁定策略、安全工具与日志封装等。

---

## 快速开始

1. **安装（模块导入）：**

```bash
go get helay.net/go/utils/v3
```

2. **典型导入示例：**

```go
import (
    "helay.net/go/utils/v3/tools"
    "helay.net/go/utils/v3/template/template_engine"
    "helay.net/go/utils/v3/rule-engine/validator"
)
```

3. **常见使用示例：**

```go
// 文件读写
err := tools.FilePutContents("/tmp/hello.txt", "hello world")
b, _ := tools.FileGetContents("/tmp/hello.txt")

// 模板函数注册
t := template.New("page").Funcs(template_engine.BuiltinFuncMap())

// 规则验证示例
rule := &validator.Rule{
    Field: "age",
    FieldDataType: "int",
    Category: types.CategoryContent,
    Operator: types.GreaterEqual,
    Value: []any{18},
}
msg, ok := rule.Validate(map[string]any{"age": 20})
```

---

## 功能与模块

下面按模块列出主要职责、关键文件与典型用法，便于快速定位与引用。

1. **tools — 通用辅助工具**  

   - 包路径：`helay.net/go/utils/v3/tools`  

   - 功能：字符串与命名转换、文件读写、IO 辅助、map/切片工具、环形缓冲索引、深拷贝等。  

   - 关键文件：`tools/file.go`、`tools/io.go`、`tools/map.go`、`tools/func.go`、`tools/ringbuffer/...`。  

   - 典型用法：
     - `tools.FileGetContents(path)`、`tools.FilePutContents(path, content)`、`tools.FileAppendContents(path, content)`


2. **net/http — HTTP 层工具与中间件**  

   - 包路径：`helay.net/go/utils/v3/net/http`（及其子包）  

   - 功能：统一响应（response）、路由处理、http server 支持、session、mime、route 中间件（日志/metrics）等。  

   - 关键文件：`net/http/response/*`、`net/http/httpServer/router/*`、`net/http/route/middleware/*`。  

   - 典型用法：在 handler 中使用统一响应 `response.SetReturnData(w, code, data)`。


3. **template/template_engine — 模板内置函数**  

   - 包路径：`helay.net/go/utils/v3/template/template_engine`  

   - 功能：为 `html/template` 提供内置函数（时间、格式化、链接生成、dict、循环等）。  

   - 关键文件：`enginetools.go`、`utils.go`。  

   - 典型用法：`t := template.New("tpl").Funcs(template_engine.BuiltinFuncMap())`


4. **file/filesaver & db/fileSaver — 文件存储抽象与驱动**  

   - 包路径：
     - 抽象：`helay.net/go/utils/v3/file/filesaver`
     - 本地驱动：`helay.net/go/utils/v3/file/filesaver/localfile`
     - MinIO 驱动：`helay.net/go/utils/v3/db/fileSaver/minio`  

   - 功能：统一文件写入/读取/列举/删除接口，支持本地、FTP、SFTP、HDFS、MinIO 等后端。  

   - 关键方法（MinIO）：`Write`, `Read`, `ListFiles`, `Delete`。


5. **db — 数据库辅助与查询构造器**  

   - 包路径：`helay.net/go/utils/v3/db`，查询构造器在 `db/query`  

   - 功能：DSN 构建（MySQL/Postgres/SQLite）、查询构造器到 GORM clause 的转换、字段解析工具。  

   - 关键文件：`db/dsn.go`、`db/query/gorm.go`、`db/query/helper.go`。


6. **db/localredis — 本地 Redis 适配（测试替代）**  

   - 包路径：`helay.net/go/utils/v3/db/localredis`  

   - 说明：提供本地实现以便在无 Redis 环境或测试中替代真实 redis 客户端。注意：部分方法为占位实现，生产使用前请确认完整性。


7. **rule-engine/validator — 规则验证引擎**  

   - 包路径：`helay.net/go/utils/v3/rule-engine/validator`  

   - 功能：支持数据类型、长度、格式、内容与高级表达式验证；支持 AND/OR 逻辑组合、通配符字段、CEL/Go/JSONLogic 等表达式类型；内置中文错误提示。  

   - 文档：`rule-engine/validator/readme.md`（含操作符表与示例）。


8. **security — 安全相关工具**  

   - 子模块：
     - lockpolicy（锁定策略管理，session/ip/user，多层锁定与升级链）
     - cors（CORS 策略与预检处理）  

   - 包路径：`helay.net/go/utils/v3/security/...`


9. **logger — 日志实现**  

   - 包路径：`helay.net/go/utils/v3/logger`（含 `ulogs`、`zaploger`）  

   - 功能：日志记录、文件滚动（结合 lumberjack）、HTTP 请求日志中间件等。


10. **crypto — 加密工具**  

    - 包路径：`helay.net/go/utils/v3/crypto/aes`  

    - 功能：AES CBC 简单加解密（包含 padding/unpadding，适用于简单场景）。


11. **config — 配置与正则常量**  

    - 包路径：`helay.net/go/utils/v3/config`  

    - 功能：常用正则（手机号、邮箱、页面解析）、YAML 加载器（支持 include 机制）。  

    - 关键文件：`config/regexp.go`、`config/loadYaml/load.go`。


12. **其它项目级文件**  

    - `var.go`：全局常量（Salt、Version、BuildTime）  

    - `init.go`：运行目录初始化（设置 `config.Appath`）  

    - `message/`：消息模块（当前 README 仅有标题，建议补充实现/使用说明）


---

## 注意事项

1. 许可证与依赖

   - 仓库采用 MIT 许可证（见仓库根 LICENSE）。  

   - 在打包/发布时，请保留第三方依赖的 LICENSE/NOTICE 文件（特别是 Apache-2.0 的 NOTICE）。

2. 性能与并发

   - 规则引擎：通配符匹配（如 `items.*.price`）在数据量大或规则复杂时会带来性能开销，建议缓存规则解析或避免在热路径频繁使用。  

   - 文件写入：并发写入到网络文件系统（NFS、MinIO、HDFS 等）时请考虑临时文件 + 重命名方案以保证原子性与一致性。

3. 安全与加密

   - 仓内提供的加密辅助（如 AES）为通用工具示例。选择真实生产加密方案时，请注意 IV、认证加密（AEAD）与密钥管理等安全需求。

4. Go 版本

   - 项目声明 `go 1.24.0`（见 go.mod），请在升级 Go 版本或依赖版本时做兼容性测试。


---

## 贡献

欢迎贡献，简要流程：

1. 提交 Issue
   - 先创建 issue 讨论需求或 bug，并给出复现步骤与期望行为。

2. Fork 并创建分支
   - Fork 仓库，基于 `master` 创建 feature 分支（例如 `feature/xxx` 或 `fix/xxx`）。

3. 编码与测试
   - 编写实现并添加测试，保持模块边界清晰。运行 `go test` 并确保通过。

4. 提交 PR
   - 提交 PR，描述变更内容、影响范围、测试步骤与兼容性说明。维护者会审阅并给出反馈。

5. 代码规范
   - 请使用 `gofmt`、`go vet`，保持清晰的注释与示例代码。新增的公共 API 请在对应子模块 README 中补充用法示例。

---

## License

© helays, 2024~time.Now

本项目采用 MIT 许可证，详见 [MIT License](https://github.com/helays/utils/blob/master/LICENSE)。

