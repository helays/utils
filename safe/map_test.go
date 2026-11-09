package safe

import (
	"context"
	"sync"
	"testing"
	"time"
)

// newUpdateTestMap 构造测试用 Map，并在用例结束时停止清理协程。
func newUpdateTestMap(t *testing.T, cfg ...CacheConfig) *Map[string, int] {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return NewMap[string, int](ctx, StringHasher{}, cfg...)
}

// TestMapUpdateBasic 验证存在键时写回新值并返回 true。
func TestMapUpdateBasic(t *testing.T) {
	m := newUpdateTestMap(t)
	m.Store("k", 1)

	got, ok := m.Update("k", func(old int) (int, bool) { return old + 1, false })
	if !ok || got != 2 {
		t.Fatalf("Update() = (%v, %v), want (2, true)", got, ok)
	}
	if v, loaded := m.Load("k"); !loaded || v != 2 {
		t.Fatalf("Load() = (%v, %v), want (2, true)", v, loaded)
	}
}

// TestMapUpdateMissing 验证键不存在时不调用回调、不创建新键。
func TestMapUpdateMissing(t *testing.T) {
	m := newUpdateTestMap(t)

	called := false
	got, ok := m.Update("missing", func(old int) (int, bool) {
		called = true
		return old + 1, false
	})
	if ok || got != 0 {
		t.Fatalf("Update() = (%v, %v), want (0, false)", got, ok)
	}
	if called {
		t.Fatalf("键不存在时不应调用 fn")
	}
	if _, loaded := m.Load("missing"); loaded {
		t.Fatalf("Update 不应创建新键")
	}
}

// TestMapUpdateExpired 验证键已过期时返回 false、不调用回调并删除过期条目。
func TestMapUpdateExpired(t *testing.T) {
	m := newUpdateTestMap(t, CacheConfig{
		EnableCleanup: true,
		TTL:           20 * time.Millisecond,
		ClearInterval: time.Hour, // 避免清理协程在用例内提前删除，交由 Update 自己处理
	})
	m.Store("k", 1)
	time.Sleep(40 * time.Millisecond)

	called := false
	got, ok := m.Update("k", func(old int) (int, bool) {
		called = true
		return old + 1, false
	})
	if ok || got != 0 {
		t.Fatalf("过期键 Update() = (%v, %v), want (0, false)", got, ok)
	}
	if called {
		t.Fatalf("键过期时不应调用 fn")
	}
	if _, loaded := m.Load("k"); loaded {
		t.Fatalf("过期条目应被删除")
	}
}

// TestMapUpdateRefresh 验证 refresh 能延长 TTL，refresh=false 不延长。
func TestMapUpdateRefresh(t *testing.T) {
	m := newUpdateTestMap(t, CacheConfig{
		EnableCleanup: true,
		TTL:           500 * time.Millisecond,
		ClearInterval: time.Hour,
	})
	m.Store("k", 1)

	before, ok := m.GetTTL("k")
	if !ok || before <= 0 {
		t.Fatalf("GetTTL() = (%v, %v), want 正值", before, ok)
	}
	time.Sleep(20 * time.Millisecond)

	// refresh=false：只改值，不应延长 TTL。
	if _, updated := m.Update("k", func(old int) (int, bool) { return old + 1, false }); !updated {
		t.Fatalf("Update 未命中已存在的键")
	}
	after, ok := m.GetTTL("k")
	if !ok {
		t.Fatalf("refresh=false 后键不应失效")
	}
	if after > before {
		t.Fatalf("refresh=false 不应延长 TTL：before=%v after=%v", before, after)
	}

	// refresh=true + duration：改用新的 TTL。
	m.Store("k2", 1, 30*time.Millisecond)
	if _, updated := m.Update("k2", func(old int) (int, bool) { return old + 1, true }, 200*time.Millisecond); !updated {
		t.Fatalf("Update 未命中已存在的键")
	}
	ttl, ok := m.GetTTL("k2")
	if !ok || ttl < 150*time.Millisecond {
		t.Fatalf("续期后 GetTTL() = (%v, %v), want >= 150ms", ttl, ok)
	}
	time.Sleep(50 * time.Millisecond)
	if _, loaded := m.Load("k2"); !loaded {
		t.Fatalf("续期后的键不应过期")
	}
}

// TestMapUpdateCleanupDisabled 验证关闭自动清理时 refresh 不写入过期时间。
// store 在关闭清理时忽略 duration，Update 也必须保持同样的语义，
// 否则会留下“已过期但永不被清理”的条目。
func TestMapUpdateCleanupDisabled(t *testing.T) {
	m := newUpdateTestMap(t) // 默认关闭 EnableCleanup
	m.Store("k", 1, 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	got, ok := m.Update("k", func(old int) (int, bool) { return old + 1, true }, 10*time.Millisecond)
	if !ok || got != 2 {
		t.Fatalf("Update() = (%v, %v), want (2, true)", got, ok)
	}
	if _, loaded := m.Load("k"); !loaded {
		t.Fatalf("关闭自动清理时 Update 不应使键过期")
	}
	if ttl, ttlOK := m.GetTTL("k"); !ttlOK || ttl != 0 {
		t.Fatalf("关闭自动清理时 GetTTL() = (%v, %v), want (0, true)", ttl, ttlOK)
	}
}

// TestMapUpdateNilFunc 验证 fn 为 nil 时是空操作。
func TestMapUpdateNilFunc(t *testing.T) {
	m := newUpdateTestMap(t)
	m.Store("k", 7)

	got, ok := m.Update("k", nil)
	if !ok || got != 7 {
		t.Fatalf("Update(nil) = (%v, %v), want (7, true)", got, ok)
	}
	if v, loaded := m.Load("k"); !loaded || v != 7 {
		t.Fatalf("Load() = (%v, %v), want (7, true)", v, loaded)
	}
}

// TestMapUpdateConcurrent 并发验证读改写的原子性。
func TestMapUpdateConcurrent(t *testing.T) {
	m := newUpdateTestMap(t)
	m.Store("k", 0)

	const goroutines = 16
	const iters = 1000
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			for range iters {
				m.Update("k", func(v int) (int, bool) { return v + 1, false })
			}
		}()
	}
	wg.Wait()

	if v, loaded := m.Load("k"); !loaded || v != goroutines*iters {
		t.Fatalf("并发 Update 最终值 = (%v, %v), want (%v, true)", v, loaded, goroutines*iters)
	}
}
