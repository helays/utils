package electMaster

import (
	"context"
	"github.com/helays/utils/close/vclose"
	"github.com/helays/utils/logger/ulogs"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"time"
)

type elect struct {
	cli           *clientv3.Client
	electionKey   string
	candidateInfo string
	leaderChs     []chan<- bool
}

// ElectMaster 选leader
func ElectMaster(cli *clientv3.Client, electionKey, candidateInfo string, leaderChs ...chan<- bool) {
	e := elect{
		cli:           cli,
		electionKey:   electionKey,
		candidateInfo: candidateInfo,
		leaderChs:     leaderChs,
	}
	go e.process()
}

func (this *elect) process() {
	for {
		this.run()
		// 添加一个延迟，避免在短时间内频繁重试
		time.Sleep(1 * time.Second)
	}
}

func (this *elect) run() {
	session, err := concurrency.NewSession(this.cli, concurrency.WithTTL(10))
	if err != nil {
		this.error(err, "创建会话失败")
		return
	}
	defer vclose.Close(session)
	elector := concurrency.NewElection(session, this.electionKey)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保在函数退出时取消context，避免goroutine泄露
	go this.campaign(ctx, elector)

	// 等待session结束或竞选失败
	select {
	case <-session.Done():
		this.log("会话结束，可能是由于失去了Leader地位或竞选失败")
		return
	case <-ctx.Done():
		this.log("竞选失败或主动取消")
		return
	}
}

func (this *elect) campaign(ctx context.Context, elector *concurrency.Election) {
	err := elector.Campaign(ctx, this.candidateInfo)
	if err != nil {
		this.error(err, "竞选失败")
		return
	}
	for _, ch := range this.leaderChs {
		go func(ch chan<- bool) {
			ch <- true // 发送信号表示成为leader
		}(ch)
	}
	this.log("当前节点成为leader")
	observeCh := elector.Observe(ctx)
	for resp := range observeCh {
		this.log("新的leader是:", string(resp.Kvs[0].Value))
		// 当前节点不再是leader时，发送信号并退出循环
		if string(resp.Kvs[0].Value) != this.candidateInfo {
			for _, ch := range this.leaderChs {
				go func(ch chan<- bool) {
					ch <- false
				}(ch)
			}
			return
		}
	}
}

func (this *elect) error(err error, msg ...any) {
	ulogs.Error(append([]any{"【ETCD选leader】", err.Error()}, msg...)...)
}

func (this *elect) log(args ...any) {
	ulogs.Log(append([]any{"【ETCD选leader】"}, args...)...)
}
