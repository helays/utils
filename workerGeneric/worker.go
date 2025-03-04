package workerGeneric

import "context"

// Job 结构体现在支持泛型
type Job[T any] struct {
	// 工作函数，接受泛型参数
	Func func(T)
	// 工作参数，类型为泛型 T
	Params T
}

type worker[T any] struct {
	WorkerPool chan chan Job[T]
	JobChannel chan Job[T]
	ctx        context.Context
	cancel     context.CancelFunc
}

func (w worker[T]) start() {
	go func() {
		for {
			select {
			// 将当前 worker 注册到 worker 队列中
			case w.WorkerPool <- w.JobChannel:
				job := <-w.JobChannel
				// 接收到工作请求
				job.Func(job.Params)
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

// StartWorker 定义开始工作结构体，支持泛型
type StartWorker[T any] struct {
	// 最大运行数
	MaxSize    int `ini:"max_size"`
	WorkerPool chan chan Job[T]
	workers    []*worker[T]
}

// Init 初始化工作池
func (s *StartWorker[T]) Init() {
	s.workers = make([]*worker[T], s.MaxSize)
	for i := 0; i < s.MaxSize; i++ {
		w := &worker[T]{
			WorkerPool: s.WorkerPool,
			JobChannel: make(chan Job[T]),
		}
		w.ctx, w.cancel = context.WithCancel(context.Background())
		w.start()
		s.workers[i] = w
	}
}

// Run 运行，支持泛型
func (s *StartWorker[T]) Run(j *Job[T]) {
	jobChannel := <-s.WorkerPool
	jobChannel <- *j
}

// Close 关闭所有worker
func (s *StartWorker[T]) Close() error {
	for _, w := range s.workers {
		w.cancel()
	}
	return nil
}
