<#
update-deps.ps1 —— 读取 go.mod 的 require 列表，逐条执行 `go get -u <module>@latest`。

适用：Windows 原生 PowerShell（Windows PowerShell 5.1 与 PowerShell 7）。
Bash 版本见同目录 update-deps.sh，二者默认行为与输出保持一致。

环境约定（见仓库 AGENTS.md）：
  GOCACHE / GOTMPDIR / GOMODCACHE 统一落在 <模块根>\runtime\.go\ 下，
  GOTOOLCHAIN=local 禁止 Go 自动联网下载工具链。

常用示例：
  powershell -File update-deps.ps1 -DryRun          # 只打印将要执行的命令
  powershell -File update-deps.ps1                 # 全量（直接 + 间接，当前 168 条）
  powershell -File update-deps.ps1 -DirectOnly     # 只更新 49 条直接依赖
  powershell -File update-deps.ps1 -Tidy -Verify   # 更新后收敛并编译校验
  powershell -File update-deps.ps1 -Proxy https://goproxy.cn,direct
#>

param(
    [Alias('C')]
    [string]$Dir = '',

    [Alias('direct-only')]
    [switch]$DirectOnly,

    [Alias('dry-run')]
    [switch]$DryRun,

    [Alias('fail-fast')]
    [switch]$FailFast,

    [switch]$Tidy,

    [switch]$Verify,

    [Alias('no-backup')]
    [switch]$NoBackup,

    [string]$Proxy = '',

    [switch]$Help,

    # 任意未声明的参数（如误写 -Bogus）都会被收集到这里，而不是被静默忽略后执行真实更新
    [Parameter(ValueFromRemainingArguments = $true)]
    [object[]]$RemainingArgs
)

$ErrorActionPreference = 'Continue'

if ($RemainingArgs -and $RemainingArgs.Count -gt 0) {
    [Console]::Error.WriteLine("错误: 未知参数: $($RemainingArgs -join ' ')（使用 -Help 查看用法）")
    exit 2
}

function Show-Usage {
    Write-Host @'
用法: powershell -File update-deps.ps1 [选项]

读取 go.mod 中 require 的每个模块，按出现顺序执行:
    go get -u <module>@latest

选项:
  -Dir <目录>       指定 Go 模块目录（默认: 脚本所在目录，即仓库根目录）
  -DirectOnly       只更新直接依赖，跳过带 // indirect 的条目
  -DryRun           只打印将要执行的命令，不创建目录、不备份、不写日志
  -FailFast         某条失败后立即停止（默认: 记录失败并继续）
  -Tidy             全部成功后执行 go mod tidy
  -Verify           全部成功后执行 go build ./...
  -NoBackup         不备份 go.mod / go.sum（不推荐）
  -Proxy <URL>      本次运行临时设置 GOPROXY
  -Help             显示本帮助

兼容写法: --dir / --direct-only / --dry-run / --fail-fast / --no-backup / --help

环境变量:
  GOCACHE / GOTMPDIR / GOMODCACHE 会被强制指向 <模块根>\runtime\.go\，
  GOTOOLCHAIN 固定为 local；GOPROXY 默认沿用当前 go env 的配置。

退出码:
  0  全部成功
  1  至少一条 go get（或 go mod tidy）失败
  2  用法或环境错误（找不到 go.mod、目录不存在、未知参数等）
'@
}

function Stop-WithError {
    param([string]$Message)
    [Console]::Error.WriteLine("错误: $Message")
    exit 2
}

# 输出 include 的 require 条目：Path / Version / Indirect
function Get-RequireEntries {
    param([string]$Path)

    $inBlock = $false
    foreach ($raw in Get-Content -LiteralPath $Path -Encoding UTF8) {
        $line = $raw.Trim()
        if ($line -eq '' -or $line.StartsWith('//')) { continue }

        if ($line -match '^require\s*\(\s*$') {
            $inBlock = $true
            continue
        }
        if ($inBlock -and $line -match '^\)') {
            $inBlock = $false
            continue
        }
        if ($line -match '^require\s+(.+)$') {
            $line = $Matches[1]
        }
        elseif (-not $inBlock) {
            continue
        }

        $tail = ''
        $idx = $line.IndexOf('//')
        if ($idx -ge 0) {
            $tail = $line.Substring($idx + 2)
            $line = $line.Substring(0, $idx)
        }

        $fields = $line.Trim() -split '\s+'
        if ($fields.Count -lt 2 -or [string]::IsNullOrWhiteSpace($fields[0])) { continue }

        [pscustomobject]@{
            Path     = $fields[0]
            Version  = $fields[1]
            Indirect = ($tail -match 'indirect')
        }
    }
}

# 输出 go.mod 中 replace 指令左侧的模块路径（用于跳过被替换的依赖）
function Get-ReplacePaths {
    param([string]$Path)

    $inBlock = $false
    foreach ($raw in Get-Content -LiteralPath $Path -Encoding UTF8) {
        $line = $raw.Trim()
        if ($line -eq '' -or $line.StartsWith('//')) { continue }

        if ($line -match '^replace\s*\(\s*$') {
            $inBlock = $true
            continue
        }
        if ($inBlock -and $line -match '^\)') {
            $inBlock = $false
            continue
        }
        if ($line -match '^replace\s+(.+)$') {
            $line = $Matches[1]
        }
        elseif (-not $inBlock) {
            continue
        }

        $arrow = $line.IndexOf('=>')
        $left = if ($arrow -ge 0) { $line.Substring(0, $arrow) } else { $line }
        $fields = $left.Trim() -split '\s+'
        if ($fields.Count -ge 1 -and -not [string]::IsNullOrWhiteSpace($fields[0])) { $fields[0] }
    }
}

if ($Help) {
    Show-Usage
    exit 0
}

try {
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
}
catch {
    # 控制台不支持时忽略
}

# ---------- 定位模块 ----------
$moduleDir = if ([string]::IsNullOrWhiteSpace($Dir)) {
    $PSScriptRoot
}
else {
    try {
        (Resolve-Path -LiteralPath $Dir).Path
    }
    catch {
        Stop-WithError "模块目录不存在: $Dir"
    }
}
if (-not (Test-Path -LiteralPath $moduleDir -PathType Container)) {
    Stop-WithError "模块目录不存在: $moduleDir"
}
$gomod = Join-Path $moduleDir 'go.mod'
if (-not (Test-Path -LiteralPath $gomod -PathType Leaf)) {
    Stop-WithError "找不到 go.mod: $gomod"
}

$replacedPaths = @(Get-ReplacePaths -Path $gomod)

$selected = @(foreach ($entry in @(Get-RequireEntries -Path $gomod)) {
        if ($DirectOnly -and $entry.Indirect) { continue }
        if ($replacedPaths -contains $entry.Path) {
            Write-Warning "跳过（go.mod 使用了 replace）: $($entry.Path)"
            continue
        }
        $entry
    })
$total = $selected.Count
if ($total -eq 0) {
    Stop-WithError '没有可更新的 require 条目'
}

Write-Host "模块: $moduleDir"
if ($DirectOnly) {
    Write-Host "范围: 仅直接依赖，共 $total 条"
}
else {
    Write-Host "范围: 直接 + 间接依赖，共 $total 条"
}

# ---------- dry-run：不产生任何副作用 ----------
if ($DryRun) {
    Write-Host '模式: dry-run（只打印命令，不修改文件）'
    for ($i = 0; $i -lt $total; $i++) {
        $entry = $selected[$i]
        $tag = if ($entry.Indirect) { ' [indirect]' } else { '' }
        Write-Host ('[{0}/{1}] go get -u {2}@latest{3}  # 当前 {4}' -f ($i + 1), $total, $entry.Path, $tag, $entry.Version)
    }
    exit 0
}

# ---------- 运行环境 ----------
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Stop-WithError '未找到 go 命令，请先安装 Go 并加入 PATH'
}

$runtimeRoot = Join-Path $moduleDir 'runtime\.go'
foreach ($sub in @('cache', 'tmp', 'mod')) {
    New-Item -ItemType Directory -Force -Path (Join-Path $runtimeRoot $sub) | Out-Null
}
$env:GOCACHE = Join-Path $runtimeRoot 'cache'
$env:GOTMPDIR = Join-Path $runtimeRoot 'tmp'
$env:GOMODCACHE = Join-Path $runtimeRoot 'mod'
$env:GOTOOLCHAIN = 'local'
if (-not [string]::IsNullOrWhiteSpace($Proxy)) {
    $env:GOPROXY = $Proxy
}

$ts = [DateTime]::UtcNow.ToString("yyyyMMdd'T'HHmmss'Z'")
$logDir = Join-Path $runtimeRoot "deps-update\$ts"
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
$logPath = Join-Path $logDir 'go-get.log'

$header = @(
    "模块: $moduleDir",
    "时间(UTC): $ts",
    "模块数: $total",
    "GOCACHE=$($env:GOCACHE)",
    "GOTMPDIR=$($env:GOTMPDIR)",
    "GOMODCACHE=$($env:GOMODCACHE)",
    "GOTOOLCHAIN=$($env:GOTOOLCHAIN)",
    "GOPROXY=$(if ($env:GOPROXY) { $env:GOPROXY } else { '未设置(沿用 go env)' })"
)
Add-Content -LiteralPath $logPath -Value $header -Encoding UTF8

$backupDir = $null
$backedUpSum = $false
if (-not $NoBackup) {
    $backupDir = Join-Path $runtimeRoot "deps-backup\$ts"
    New-Item -ItemType Directory -Force -Path $backupDir | Out-Null
    Copy-Item -LiteralPath $gomod -Destination (Join-Path $backupDir 'go.mod') -Force
    $goSum = Join-Path $moduleDir 'go.sum'
    if (Test-Path -LiteralPath $goSum -PathType Leaf) {
        Copy-Item -LiteralPath $goSum -Destination (Join-Path $backupDir 'go.sum') -Force
        $backedUpSum = $true
    }
    Write-Host "备份: $backupDir"
}
Write-Host "日志: $logPath"
Write-Host ''

# ---------- 逐条 go get ----------
Push-Location -LiteralPath $moduleDir
$ok = 0
$fail = 0
$failedModules = @()
$tidyFailed = $false

try {
    for ($i = 0; $i -lt $total; $i++) {
        $entry = $selected[$i]
        $current = $i + 1
        $line = '[{0}/{1}] go get -u {2}@latest（当前 {3}）' -f $current, $total, $entry.Path, $entry.Version
        Write-Host ''
        Write-Host $line
        Add-Content -LiteralPath $logPath -Value ("[{0}/{1}] go get -u {2}@latest (current {3})" -f $current, $total, $entry.Path, $entry.Version) -Encoding UTF8

        $output = & go get -u "$($entry.Path)@latest" 2>&1
        $exitCode = $LASTEXITCODE
        if ($null -ne $output) {
            $text = @($output | ForEach-Object { $_.ToString() })
            foreach ($chunk in $text) { Write-Host $chunk }
            Add-Content -LiteralPath $logPath -Value $text -Encoding UTF8
        }

        if ($exitCode -eq 0) {
            $ok++
            Write-Host '  -> 成功'
        }
        else {
            $fail++
            $failedModules += $entry.Path
            Write-Host '  -> 失败（已记录，继续执行）' -ForegroundColor Yellow
            if ($FailFast) {
                Write-Host '已按 -FailFast 停止。' -ForegroundColor Yellow
                break
            }
        }
    }

    if ($Tidy) {
        if ($fail -eq 0) {
            Write-Host ''
            Write-Host '[go mod tidy]'
            $output = & go mod tidy 2>&1
            $exitCode = $LASTEXITCODE
            if ($null -ne $output) {
                $text = @($output | ForEach-Object { $_.ToString() })
                foreach ($chunk in $text) { Write-Host $chunk }
                Add-Content -LiteralPath $logPath -Value $text -Encoding UTF8
            }
            if ($exitCode -eq 0) {
                Write-Host '  -> go mod tidy 完成'
            }
            else {
                Write-Host '  -> go mod tidy 失败' -ForegroundColor Red
                $tidyFailed = $true
            }
        }
        else {
            Write-Host ''
            Write-Host '存在失败项，跳过 go mod tidy。' -ForegroundColor Yellow
        }
    }

    if ($Verify) {
        if ($fail -eq 0) {
            Write-Host ''
            Write-Host '[go build ./...]'
            $output = & go build ./... 2>&1
            $exitCode = $LASTEXITCODE
            if ($null -ne $output) {
                $text = @($output | ForEach-Object { $_.ToString() })
                foreach ($chunk in $text) { Write-Host $chunk }
                Add-Content -LiteralPath $logPath -Value $text -Encoding UTF8
            }
            if ($exitCode -eq 0) {
                Write-Host '  -> go build 通过'
            }
            else {
                Write-Host '  -> go build 失败（可能缺 CGO/构建标签，请人工确认）' -ForegroundColor Yellow
            }
        }
        else {
            Write-Host ''
            Write-Host '存在失败项，跳过 go build。' -ForegroundColor Yellow
        }
    }
}
finally {
    Pop-Location
}

# ---------- 汇总 ----------
Write-Host ''
Write-Host '================ 汇总 ================'
Write-Host ('预选: {0}  执行: {1}  成功: {2}  失败: {3}' -f $total, ($ok + $fail), $ok, $fail)
Write-Host "日志: $logPath"
if ($backupDir) {
    Write-Host "备份: $backupDir"
    $restore = if ($backedUpSum) {
        'Copy-Item -LiteralPath "{0}\go.mod","{0}\go.sum" -Destination "{1}" -Force' -f $backupDir, $moduleDir
    }
    else {
        'Copy-Item -LiteralPath "{0}\go.mod" -Destination "{1}" -Force' -f $backupDir, $moduleDir
    }
    Write-Host "恢复: $restore"
}
if ($fail -gt 0) {
    Write-Host ''
    Write-Host '失败模块（可单独重试）:'
    foreach ($module in $failedModules) {
        Write-Host "  go get -u $module@latest"
    }
}
if (Get-Command git -ErrorAction SilentlyContinue) {
    & git -C $moduleDir rev-parse --git-dir *> $null
    if ($LASTEXITCODE -eq 0) {
        Write-Host ''
        Write-Host 'go.mod 变更统计:'
        & git -C $moduleDir --no-pager diff --stat -- go.mod
    }
}

if ($fail -gt 0 -or $tidyFailed) {
    exit 1
}
exit 0
