package lockpolicy

import (
	"time"

	"helay.net/go/utils/v3/tools"
)

// Valid 锁定策略验证
// 对于升级类型策略，并且有上一级的情况，窗口时间算法建议如下：
// policy.window_time > (2 * pre.lockout_time * policy.lockout_count)
// 大于上级策略中的，2*锁定时长*当前策略的触发次数；
func (p *Policy) Valid() {
	if p.Trigger < 1 {
		return // 未启用的策略
	}
	p.LockoutTime = tools.AutoTimeDuration(p.LockoutTime, time.Second, 5*time.Minute) // 策略锁定时长
	p.WindowTime = tools.AutoTimeDuration(p.WindowTime, time.Second, 10*time.Minute)  // 异常记录窗口时间窗口时间

	return
}
