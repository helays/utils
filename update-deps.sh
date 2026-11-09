#!/usr/bin/env bash
# update-deps.sh —— 读取 go.mod 的 require 列表，逐条执行 `go get -u <module>@latest`。
#
# 适用：Git Bash / WSL / Linux / macOS。Windows 原生请用同目录的 update-deps.ps1。
# 依赖：bash 4+、go。脚本不会自动提交 git，更新结果请用 `git diff go.mod` 审核。
#
# 环境约定（见仓库 AGENTS.md）：
#   GOCACHE / GOTMPDIR / GOMODCACHE 统一落在 <模块根>/runtime/.go/ 下，
#   GOTOOLCHAIN=local 禁止 Go 自动联网下载工具链。
#
# 常用示例：
#   bash update-deps.sh --dry-run          # 只打印将要执行的命令
#   bash update-deps.sh                    # 全量（直接 + 间接，当前 168 条）
#   bash update-deps.sh --direct-only      # 只更新 49 条直接依赖
#   bash update-deps.sh --tidy --verify    # 更新后收敛并编译校验
#   bash update-deps.sh --proxy https://goproxy.cn,direct

set -uo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
MODULE_DIR="$SCRIPT_DIR"

DIRECT_ONLY=0
DRY_RUN=0
FAIL_FAST=0
RUN_TIDY=0
RUN_VERIFY=0
DO_BACKUP=1
GOPROXY_OVERRIDE=""
GOMOD=""

usage() {
    cat <<'EOF'
用法: bash update-deps.sh [选项]

读取 go.mod 中 require 的每个模块，按出现顺序执行:
    go get -u <module>@latest

选项:
  -C, --dir <目录>      指定 Go 模块目录（默认: 脚本所在目录，即仓库根目录）
      --direct-only     只更新直接依赖，跳过带 // indirect 的条目
      --dry-run         只打印将要执行的命令，不创建目录、不备份、不写日志
      --fail-fast       某条失败后立即停止（默认: 记录失败并继续）
      --tidy            全部成功后执行 go mod tidy
      --verify          全部成功后执行 go build ./...
      --no-backup       不备份 go.mod / go.sum（不推荐）
      --proxy <URL>     本次运行临时设置 GOPROXY
  -h, --help            显示本帮助

环境变量:
  GOCACHE / GOTMPDIR / GOMODCACHE 会被强制指向 <模块根>/runtime/.go/，
  GOTOOLCHAIN 固定为 local；GOPROXY 默认沿用当前 go env 的配置。

退出码:
  0  全部成功
  1  至少一条 go get（或 go mod tidy）失败
  2  用法或环境错误（找不到 go.mod、目录不存在、未知参数等）
EOF
}

die() {
    printf '错误: %s\n' "$*" >&2
    exit 2
}

# Git Bash 下的 go.exe 需要 Windows 原生路径；Linux/macOS/WSL 保持 POSIX 路径。
native_path() {
    if command -v cygpath >/dev/null 2>&1; then
        cygpath -w -- "$1"
    else
        printf '%s\n' "$1"
    fi
}

# 输出: path|indirect|version。仅解析 require 块与单行 require，忽略注释行与行尾注释。
parse_requires() {
    local line rest path ver in_block=0
    while IFS= read -r line || [[ -n $line ]]; do
        line=${line%$'\r'}
        line=${line#"${line%%[![:space:]]*}"}
        [[ -z $line || $line == //* ]] && continue

        if [[ $line =~ ^require[[:space:]]*\($ ]]; then
            in_block=1
            continue
        fi
        if [[ $in_block -eq 1 && $line =~ ^\) ]]; then
            in_block=0
            continue
        fi
        if [[ $line =~ ^require[[:space:]]+(.+)$ ]]; then
            line=${BASH_REMATCH[1]}
        elif [[ $in_block -eq 0 ]]; then
            continue
        fi

        rest=""
        if [[ $line == *"//"* ]]; then
            rest=${line#*//}
            line=${line%%//*}
        fi
        read -r path ver _ <<< "$line"
        [[ -z $path ]] && continue

        if [[ $rest == *indirect* ]]; then
            printf '%s|1|%s\n' "$path" "${ver:-unknown}"
        else
            printf '%s|0|%s\n' "$path" "${ver:-unknown}"
        fi
    done < "$GOMOD"
}

# 输出 go.mod 中 replace 指令左侧的模块路径（用于跳过被替换的依赖）。
parse_replaces() {
    local line left path in_block=0
    while IFS= read -r line || [[ -n $line ]]; do
        line=${line%$'\r'}
        line=${line#"${line%%[![:space:]]*}"}
        [[ -z $line || $line == //* ]] && continue

        if [[ $line =~ ^replace[[:space:]]*\($ ]]; then
            in_block=1
            continue
        fi
        if [[ $in_block -eq 1 && $line =~ ^\) ]]; then
            in_block=0
            continue
        fi
        if [[ $line =~ ^replace[[:space:]]+(.+)$ ]]; then
            line=${BASH_REMATCH[1]}
        elif [[ $in_block -eq 0 ]]; then
            continue
        fi

        left=${line%%=>*}
        read -r path _ <<< "$left"
        [[ -n ${path:-} ]] && printf '%s\n' "$path"
    done < "$GOMOD"
}

# 返回 0 表示该模块需要跳过。
should_skip() {
    local path=$1 indirect=$2 warn=$3
    if [[ $DIRECT_ONLY -eq 1 && $indirect -eq 1 ]]; then
        return 0
    fi
    if [[ $REPLACED == *$'\n'"$path"$'\n'* ]]; then
        if [[ -n $warn ]]; then
            printf '  跳过（go.mod 使用了 replace）: %s\n' "$path" >&2
        fi
        return 0
    fi
    return 1
}

# ---------- 参数解析 ----------
while [[ $# -gt 0 ]]; do
    case $1 in
        -C|--dir)
            [[ $# -ge 2 ]] || die "$1 需要一个目录参数"
            MODULE_DIR=$2
            shift 2
            ;;
        --direct-only) DIRECT_ONLY=1; shift ;;
        --dry-run) DRY_RUN=1; shift ;;
        --fail-fast) FAIL_FAST=1; shift ;;
        --tidy) RUN_TIDY=1; shift ;;
        --verify) RUN_VERIFY=1; shift ;;
        --no-backup) DO_BACKUP=0; shift ;;
        --proxy)
            [[ $# -ge 2 ]] || die "$1 需要一个 URL 参数"
            GOPROXY_OVERRIDE=$2
            shift 2
            ;;
        -h|--help) usage; exit 0 ;;
        *) die "未知参数: $1（使用 --help 查看用法）" ;;
    esac
done

[[ -d $MODULE_DIR ]] || die "模块目录不存在: $MODULE_DIR"
MODULE_DIR=$(cd -- "$MODULE_DIR" && pwd)
GOMOD="$MODULE_DIR/go.mod"
[[ -f $GOMOD ]] || die "找不到 go.mod: $GOMOD"

# 收集被 replace 的模块，组成 "\npath\n" 形式的集合
REPLACED=$'\n'
while IFS= read -r replaced_path; do
    [[ -n $replaced_path ]] && REPLACED+="$replaced_path"$'\n'
done < <(parse_replaces)

TOTAL=0
while IFS='|' read -r path indirect _version; do
    [[ -z $path ]] && continue
    should_skip "$path" "$indirect" "" && continue
    TOTAL=$((TOTAL + 1))
done < <(parse_requires)
((TOTAL > 0)) || die "go.mod 中没有可更新的 require 条目"

if ((DIRECT_ONLY == 1)); then
    printf '模块: %s\n范围: 仅直接依赖，共 %d 条\n' "$MODULE_DIR" "$TOTAL"
else
    printf '模块: %s\n范围: 直接 + 间接依赖，共 %d 条\n' "$MODULE_DIR" "$TOTAL"
fi

# ---------- dry-run：不产生任何副作用 ----------
if ((DRY_RUN == 1)); then
    printf '模式: dry-run（只打印命令，不修改文件）\n'
    idx=0
    while IFS='|' read -r path indirect version; do
        [[ -z $path ]] && continue
        should_skip "$path" "$indirect" "1" && continue
        idx=$((idx + 1))
        tag=""
        [[ $indirect -eq 1 ]] && tag=" [indirect]"
        printf '[%d/%d] go get -u %s@latest%s  # 当前 %s\n' "$idx" "$TOTAL" "$path" "$tag" "$version"
    done < <(parse_requires)
    exit 0
fi

# ---------- 运行环境 ----------
command -v go >/dev/null 2>&1 || die "未找到 go 命令，请先安装 Go 并加入 PATH"

RUNTIME_DIR="$MODULE_DIR/runtime/.go"
mkdir -p -- "$RUNTIME_DIR/cache" "$RUNTIME_DIR/tmp" "$RUNTIME_DIR/mod" \
    || die "无法创建运行时目录: $RUNTIME_DIR"

GOCACHE=$(native_path "$RUNTIME_DIR/cache")
GOTMPDIR=$(native_path "$RUNTIME_DIR/tmp")
GOMODCACHE=$(native_path "$RUNTIME_DIR/mod")
export GOCACHE GOTMPDIR GOMODCACHE
export GOTOOLCHAIN=local
if [[ -n $GOPROXY_OVERRIDE ]]; then
    export GOPROXY="$GOPROXY_OVERRIDE"
fi

TS=$(date -u +%Y%m%dT%H%M%SZ)
LOG_DIR="$RUNTIME_DIR/deps-update/$TS"
mkdir -p -- "$LOG_DIR" || die "无法创建日志目录: $LOG_DIR"
LOG="$LOG_DIR/go-get.log"
{
    printf '模块: %s\n' "$MODULE_DIR"
    printf '时间(UTC): %s\n' "$TS"
    printf '模块数: %d\n' "$TOTAL"
    printf 'GOCACHE=%s\nGOTMPDIR=%s\nGOMODCACHE=%s\nGOTOOLCHAIN=%s\n' \
        "$GOCACHE" "$GOTMPDIR" "$GOMODCACHE" "$GOTOOLCHAIN"
    printf 'GOPROXY=%s\n' "${GOPROXY:-未设置(沿用 go env)}"
} > "$LOG"

BACKUP_DIR=""
BACKED_UP_SUM=0
if ((DO_BACKUP == 1)); then
    BACKUP_DIR="$RUNTIME_DIR/deps-backup/$TS"
    mkdir -p -- "$BACKUP_DIR" || die "无法创建备份目录: $BACKUP_DIR"
    cp -p -- "$GOMOD" "$BACKUP_DIR/go.mod" || die "备份 go.mod 失败"
    if [[ -f "$MODULE_DIR/go.sum" ]]; then
        cp -p -- "$MODULE_DIR/go.sum" "$BACKUP_DIR/go.sum" || die "备份 go.sum 失败"
        BACKED_UP_SUM=1
    fi
    printf '备份: %s\n' "$BACKUP_DIR"
fi
printf '日志: %s\n\n' "$LOG"

# ---------- 逐条 go get ----------
cd -- "$MODULE_DIR" || die "无法进入模块目录: $MODULE_DIR"

idx=0
ok=0
fail=0
skipped=0
failures=()
TIDY_FAILED=0

while IFS='|' read -r path indirect version; do
    [[ -z $path ]] && continue
    if should_skip "$path" "$indirect" "1"; then
        skipped=$((skipped + 1))
        continue
    fi

    idx=$((idx + 1))
    printf '\n[%d/%d] go get -u %s@latest（当前 %s）\n' "$idx" "$TOTAL" "$path" "$version" | tee -a "$LOG"

    if go get -u "${path}@latest" 2>&1 | tee -a "$LOG"; then
        ok=$((ok + 1))
        printf '  -> 成功\n'
    else
        fail=$((fail + 1))
        failures+=("$path")
        printf '  -> 失败（已记录，继续执行）\n' >&2
        if ((FAIL_FAST == 1)); then
            printf '已按 --fail-fast 停止。\n' >&2
            break
        fi
    fi
done < <(parse_requires)

# ---------- 收尾 ----------
if ((RUN_TIDY == 1)); then
    if ((fail == 0)); then
        printf '\n[go mod tidy]\n'
        if go mod tidy 2>&1 | tee -a "$LOG"; then
            printf '  -> go mod tidy 完成\n'
        else
            printf '  -> go mod tidy 失败\n' >&2
            TIDY_FAILED=1
        fi
    else
        printf '\n存在失败项，跳过 go mod tidy。\n' >&2
    fi
fi

if ((RUN_VERIFY == 1)); then
    if ((fail == 0)); then
        printf '\n[go build ./...]\n'
        if go build ./... 2>&1 | tee -a "$LOG"; then
            printf '  -> go build 通过\n'
        else
            printf '  -> go build 失败（可能缺 CGO/构建标签，请人工确认）\n' >&2
        fi
    else
        printf '\n存在失败项，跳过 go build。\n' >&2
    fi
fi

printf '\n================ 汇总 ================\n'
printf '预选: %d  执行: %d  成功: %d  失败: %d  跳过: %d\n' \
    "$TOTAL" "$idx" "$ok" "$fail" "$skipped"
printf '日志: %s\n' "$LOG"
if [[ -n $BACKUP_DIR ]]; then
    printf '备份: %s\n' "$BACKUP_DIR"
    if ((BACKED_UP_SUM == 1)); then
        printf '恢复: cp -p "%s/go.mod" "%s/go.sum" "%s/"\n' \
            "$BACKUP_DIR" "$BACKUP_DIR" "$MODULE_DIR"
    else
        printf '恢复: cp -p "%s/go.mod" "%s/"\n' "$BACKUP_DIR" "$MODULE_DIR"
    fi
fi
if ((fail > 0)); then
    printf '\n失败模块（可单独重试）:\n'
    for path in "${failures[@]}"; do
        printf '  go get -u %s@latest\n' "$path"
    done
fi
if command -v git >/dev/null 2>&1 && git -C "$MODULE_DIR" rev-parse --git-dir >/dev/null 2>&1; then
    printf '\ngo.mod 变更统计:\n'
    git -C "$MODULE_DIR" --no-pager diff --stat -- go.mod 2>/dev/null || true
fi

if ((fail > 0 || TIDY_FAILED == 1)); then
    exit 1
fi
exit 0
