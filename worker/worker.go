package worker

import "context"

type Job struct {
	// 工作函数
	Func func(interface{})
	// 工作参数
	Params interface{}
}

type worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	ctx        context.Context
	cancel     context.CancelFunc
}

func (w worker) start() {
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

// StartWorker 定义开始工作结构体
type StartWorker struct {
	// 最大运行数
	MaxSize    int `ini:"max_size"`
	WorkerPool chan chan Job
	workers    []*worker
}

// Init 初始化工作池
func (s *StartWorker) Init() {
	s.workers = make([]*worker, s.MaxSize)
	for i := 0; i < s.MaxSize; i++ {
		w := &worker{
			WorkerPool: s.WorkerPool,
			JobChannel: make(chan Job),
		}
		w.ctx, w.cancel = context.WithCancel(context.Background())
		w.start()
		s.workers[i] = w
	}
}

// Run 运行
func (s *StartWorker) Run(j *Job) {
	jobChannel := <-s.WorkerPool
	jobChannel <- *j
}

// Close 关闭所有worker
func (s *StartWorker) Close() error {
	for _, w := range s.workers {
		w.cancel()
	}
	return nil
}
