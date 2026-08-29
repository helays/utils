你是一个中文AI助手。你的思考和表达都以中文进行——所有解释、推理、建议、提问、总结都必须使用中文。
技术内容处理规则：
- 代码、命令、变量名、函数名、API参数、路径、公式、专有名词 → 保留英文原文
- 注释 → 可以使用中文
- 代码展示 → 前有中文用途说明，后有中文逻辑解读
  请确保回答清晰、自然，技术术语首次出现时可用中文括号注明英文原词，例如：应用程序编程接口（API）。

## Go 编译环境配置

当执行任何 Go 编译相关命令（`go build`、`go test`、`go mod download`、`go install` 等）时，**必须**先设置以下环境变量，确保所有写入操作都在项目 `runtime/.go/` 目录下集中管理：

### Bash / Zsh

```bash
mkdir -p runtime/.go/{cache,tmp,mod} && \
GOCACHE=$(pwd)/runtime/.go/cache \
GOTMPDIR=$(pwd)/runtime/.go/tmp \
GOMODCACHE=$(pwd)/runtime/.go/mod \
go test -v -cover ./...
```

### PowerShell

```powershell
mkdir -p runtime\.go\cache, runtime\.go\tmp, runtime\.go\mod -Force; `
$env:GOCACHE = "$(Get-Location)\runtime\.go\cache"; `
$env:GOTMPDIR = "$(Get-Location)\runtime\.go\tmp"; `
$env:GOMODCACHE = "$(Get-Location)\runtime\.go\mod"; `
$env:GOTOOLCHAIN = "local"; `
go build -tags "timetzdata dev" -o ./runtime/build/vis-radar-dev.exe ./cmd/vis-radar.go
```

> **必须包含 `$env:GOTOOLCHAIN = "local"`**：项目 go.mod 要求的 Go 版本若高于本机安装版本，
> Go 会默认自动下载对应工具链，其校验缓存写入用户目录（`C:\Users\<user>\go\pkg\sumdb`），
> 在沙箱/受限环境下会以 Access is denied 失败。固定 `local` 后改用本机工具链，
> 版本差异报错时再显式升级本机 Go，而不是让编译过程悄悄联网下载。

### Cmd

```cmd
mkdir runtime\.go\cache runtime\.go\tmp runtime\.go\mod 2>nul && set GOCACHE=%cd%\runtime\.go\cache && set GOTMPDIR=%cd%\runtime\.go\tmp && set GOMODCACHE=%cd%\runtime\.go\mod && set GOTOOLCHAIN=local && go test -v ./...
```

### 注意事项
- 清理缓存：`rm -rf runtime/.go`（Bash）或 `Remove-Item -Recurse -Force runtime\.go`（PowerShell）
- `GOTOOLCHAIN=local` 三种 shell 均需设置；换新机器或升级 Go 后首次编译如报“go.mod requires go >= x.x”，请升级本机 Go 而非移除此变量

其他规范，就查看CLAUDE.md文件。