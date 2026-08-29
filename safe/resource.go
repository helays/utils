package safe

import (
	"sync/atomic"
)

// Resource 基于 atomic.Pointer 的原子、无锁（lock-free）资源。
//
// # 线程安全说明
//
// 只要通过本类型的方法（Read/Write/Get/Set/Reset/Update/Swap 等）访问资源，
// 操作就是原子且线程安全的，且完全无锁（多个 goroutine 读可完美并行）。
//
// 需要区分两种"不安全"：
//   - 如果调用方拿到 Read() 返回的 *T 指针后，在方法外部并发修改指针所指向的内容，
//     这一段改写不是原子的，会引发数据竞争。此时局部指针只应被单 goroutine 独占使用。
//   - "读出来 → 函数外比较/修改 → 再写回"这类 read-modify-write 拆分到多个方法调用，
//     不是原子的，多 goroutine 并发会丢更新。需要原子 RMW 时请用 Update。
//
// 因此，纯读推荐 Get（值拷贝）或 Read（指针，需自行保证不使用），写路径推荐 Update。
//
// 注意：本类型是通用原子资源，不做编译期数值约束，也未内置数值累加方法。
// 若需要原子的数值累加（计数器等），请在业务端用标准库 sync/atomic 的整型专用
// 方法（如 atomic.Int64.Add），或用本类型的 Update 实现 read-modify-write。
type Resource[T any] struct {
	_        noCopy
	resource atomic.Pointer[T]
}

// NewResource 创建一个原子操作的资源。
// resource 允许传入 nil，此时内部会存储 T 的零值。
func NewResource[T any](resource *T) *Resource[T] {
	s := &Resource[T]{}
	if resource == nil {
		var zero T
		s.resource.Store(&zero)
	} else {
		s.resource.Store(resource)
	}
	return s
}

// NewResourceOf 按值创建一个原子操作的资源。
// 与 NewResource 的区别在于直接传入值，风格同 NewResourceMutex/NewResourceRWMutex。
func NewResourceOf[T any](resource T) *Resource[T] {
	return NewResource(&resource)
}

// Read 返回当前资源的指针。
// 返回值是原始指针，多 goroutine 并发读取指针本身是安全的；
// 但在外部对该指针指向的内容做并发修改并不安全，请谨慎使用或改用 Get。
func (s *Resource[T]) Read() *T {
	return s.resource.Load()
}

// Ptr 是 Read 的别名，返回当前资源的原始指针。
func (s *Resource[T]) Ptr() *T {
	return s.Read()
}

// Get 返回当前资源的拷贝值。
// 比 Read 更安全：即使内部为空指针，也返回 T 的零值而不会 panic。
func (s *Resource[T]) Get() T {
	p := s.resource.Load()
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// Write 原子地写入新资源（指针）。
// 允许写入 nil，此时后续 Get 会返回 T 的零值。
func (s *Resource[T]) Write(newResource *T) {
	s.resource.Store(newResource)
}

// Set 原子地按值写入新资源。
// 等价于调用 Write(&v)，但调用方无需手动取地址。
func (s *Resource[T]) Set(v T) {
	s.Write(&v)
}

// Reset 原子地将资源重置为 T 的零值。
func (s *Resource[T]) Reset() {
	var zero T
	s.Set(zero)
}

// Update 原子地执行修改（read-modify-write）。
// 通过对当前指针快照调用 fn 得到新值，再用 CompareAndSwap 写回；
// 若期间被其他 goroutine 修改则重试，直至成功。
//
// 注意：每次重试都会重新读取最新值并重新调用 fn，因此 fn 应尽量是
// 纯函数（只依赖入参、不修改外部状态），否则可能导致 CAS 反复失败甚至死循环。
// 适用于低竞争热路径，不承诺 wait-free。
func (s *Resource[T]) Update(fn func(t T) T) {
	for {
		old := s.Read()
		var oldVal T
		if old != nil {
			oldVal = *old
		}
		newVal := fn(oldVal)
		if s.resource.CompareAndSwap(old, &newVal) {
			return
		}
	}
}

// Swap 原子地写入新指针并返回旧指针。
func (s *Resource[T]) Swap(newValue *T) *T {
	return s.resource.Swap(newValue)
}

// SwapVal 原子地按值写入并返回旧值。
// 若旧指针为空，返回 T 的零值。
func (s *Resource[T]) SwapVal(v T) T {
	old := s.resource.Swap(&v)
	if old == nil {
		var zero T
		return zero
	}
	return *old
}

// CompareAndSwap 原子地比较并交换指针。
// 仅当当前指针与 old 相等时，才替换为 new 并返回 true。
func (s *Resource[T]) CompareAndSwap(old, new *T) bool {
	return s.resource.CompareAndSwap(old, new)
}