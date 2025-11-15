# 锁定策略管理器使用文档

## 概述

锁定策略管理器是一个用于管理多层次锁定策略的系统，支持独立锁定和升级链锁定两种模式。主要用于安全防护场景，如登录失败锁定、操作异常锁定等。

## 核心概念

### 锁定目标 (LockTarget)
- `session` - 会话层锁定
- `ip` - IP层锁定
- `user` - 用户层锁定

### 锁定类型 (LockType)
- `direct` - 直接触发锁定
- `escalation` - 升级触发锁定
- `memory` - 记忆效应锁定

## 快速开始

### 1. 初始化管理器

```go
import "github.com/your-project/lockpolicy"

// 定义锁定策略
policies := lockpolicy.Policies{
    {
        Target:      lockpolicy.LockTargetSession,
        Trigger:     3,                    // 3次失败触发锁定
        WindowTime:  5 * time.Minute,      // 5分钟窗口
        LockoutTime: 10 * time.Minute,     // 锁定10分钟
        Priority:    10,
    },
    {
        Target:      lockpolicy.LockTargetIP,
        Trigger:     5,
        WindowTime:  10 * time.Minute,
        LockoutTime: 30 * time.Minute,
        Priority:    20,
        Escalation: &lockpolicy.EscalationRule{
            UpgradeTo:   lockpolicy.LockTargetUser,
            MemoryEffect: true,
        },
    },
}

// 创建管理器
manager := lockpolicy.NewManager(policies)

// 程序重启后恢复锁定状态
func restoreLocks(manager *lockpolicy.Manager) {
    // 从数据库加载活跃锁定
    activeLocks := loadActiveLocksFromDB()
    for _, lock := range activeLocks {
        manager.RestoreLock(lock.Target, lock.Identifier, lock.Expire)
    }
}
```

### 2. 记录失败事件

```go
// 记录单目标失败
targets := map[lockpolicy.LockTarget]string{
    lockpolicy.LockTargetSession: "session123",
    lockpolicy.LockTargetIP:      "192.168.1.100",
    lockpolicy.LockTargetUser:    "user456",
}

locked, event := manager.RecordFailures(targets, func(event lockpolicy.LockEvent) {
    // 锁定回调 - 可以记录日志或发送通知
    log.Printf("目标 %s 被锁定, 标识: %s, 时长: %v", 
        event.Target, event.Identifier, event.LockoutTime)
})

if locked {
    log.Printf("检测到锁定: %s - %s", event.Target, event.Identifier)
}
```

### 3. 检查锁定状态

```go
// 检查多个目标是否被锁定
targets := map[lockpolicy.LockTarget]string{
    lockpolicy.LockTargetSession: "session123",
    lockpolicy.LockTargetIP:      "192.168.1.100",
}

locked, event := manager.IsLocked(targets)
if locked {
    log.Printf("目标已被锁定: %s, 剩余时间: %v", event.Target, event.RemainingTime)
    return errors.New("账户已被锁定")
}
```

### 4. 清理锁定状态

```go
// 操作成功后清理锁定计数
targets := map[lockpolicy.LockTarget]string{
    lockpolicy.LockTargetSession: "session123",
    lockpolicy.LockTargetIP:      "192.168.1.100",
}
manager.Clear(targets)
```

## 策略配置详解

### 独立策略配置

```yaml
- target: session
  trigger: 3           # 3次失败触发锁定
  window_time: 5m      # 5分钟统计窗口
  lockout_time: 10m    # 锁定10分钟
  priority: 10         # 优先级
  # 无escalation配置即为独立策略
```

### 升级链策略配置

```yaml
- target: session
  trigger: 3
  window_time: 5m
  lockout_time: 10m
  priority: 10
  escalation:
    upgrade_to: ip     # 升级到IP锁定
    memory_effect: false

- target: ip
  trigger: 5  
  window_time: 15m
  lockout_time: 30m
  priority: 20
  escalation:
    upgrade_to: user   # 升级到用户锁定
    memory_effect: true

- target: user
  trigger: 8
  window_time: 30m
  lockout_time: 24h    # 锁定24小时
  priority: 30
  # 无升级目标，升级链终点
```

## 使用场景示例

### 登录失败锁定

```go
func HandleLogin(username, password, ip, sessionID string) error {
    targets := map[lockpolicy.LockTarget]string{
        lockpolicy.LockTargetSession: sessionID,
        lockpolicy.LockTargetIP:      ip,
        lockpolicy.LockTargetUser:    username,
    }
    
    // 先检查是否已被锁定
    if locked, _ := manager.IsLocked(targets); locked {
        return errors.New("账户已被锁定，请稍后重试")
    }
    
    // 验证登录
    if !validateLogin(username, password) {
        // 记录失败
        locked, event := manager.RecordFailures(targets, func(event lockpolicy.LockEvent) {
            notifyAdmin(event) // 通知管理员
        })
        
        if locked {
            return fmt.Errorf("登录失败次数过多，账户已被锁定 %v", event.LockoutTime)
        }
        return errors.New("用户名或密码错误")
    }
    
    // 登录成功，清理失败计数
    manager.Clear(targets)
    return nil
}
```

### 操作异常锁定

```go
func HandleSensitiveOperation(userID, operation, ip string) error {
    targets := map[lockpolicy.LockTarget]string{
        lockpolicy.LockTargetUser: userID,
        lockpolicy.LockTargetIP:   ip,
    }
    
    if err := performSensitiveOperation(); err != nil {
        locked, _ := manager.RecordFailures(targets)
        if locked {
            // 触发锁定后的处理
            securityAlert(userID, operation)
        }
        return err
    }
    
    manager.Clear(targets)
    return nil
}
```

## 最佳实践

### 1. 策略设计建议
- 窗口时间应大于锁定时长，避免连续锁定
- 升级链策略应设置合理的记忆效应
- 优先级高的策略应设置更大的触发次数

### 2. 性能考虑
- 使用合适的窗口时间，避免内存占用过大
- 定期清理过期的锁定记录
- 在分布式环境中确保锁定状态同步

### 3. 安全建议
- 程序重启后务必恢复锁定状态
- 记录锁定事件用于审计
- 设置合理的最大锁定时长

## API参考

### Manager 主要方法

| 方法 | 说明 |
|------|------|
| `NewManager(policies)` | 创建策略管理器 |
| `UpdatePolices(policies)` | 更新策略配置 |
| `RecordFailures(targets, callbacks)` | 记录失败事件 |
| `IsLocked(targets)` | 检查锁定状态 |
| `Clear(targets)` | 清理锁定状态 |
| `RestoreLock(target, identifier, expire)` | 恢复锁定状态 |

### 数据结构

| 类型 | 说明 |
|------|------|
| `LockEvent` | 锁定事件信息 |
| `Policy` | 单条锁定策略 |
| `Targets` | 目标标识映射 |

## 注意事项

1. **并发安全**: 所有操作都是线程安全的
2. **状态持久化**: 程序重启后锁定状态会丢失，需要手动恢复
3. **内存管理**: 长时间运行需监控内存使用情况
4. **策略更新**: 更新策略会重建所有内部状态