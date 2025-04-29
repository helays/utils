package electMaster

import (
	"github.com/helays/utils/logger/ulogs"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sort"
	"time"
)

type nacosElect struct {
	client        naming_client.INamingClient
	electionKey   string // Nacos服务名
	candidateInfo string // 当前节点唯一编号
	leaderChs     []chan<- bool
	stopCh        chan struct{}
}

// ElectMaster Nacos选举实现
func ElectMaster(client naming_client.INamingClient, electionKey, candidateInfo string, leaderChs ...chan<- bool) {
	e := &nacosElect{
		client:        client,
		electionKey:   electionKey,
		candidateInfo: candidateInfo,
		leaderChs:     leaderChs,
		stopCh:        make(chan struct{}),
	}
	go e.process()
}

func (e *nacosElect) process() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// 先注册自己
	if err := e.registerSelf(); err != nil {
		e.error(err, "注册实例失败")
		return
	}

	for {
		select {
		case <-e.stopCh:
			e.log("停止Nacos选举")
			return
		case <-ticker.C:
			e.checkLeadership()
		}
	}
}

func (e *nacosElect) checkLeadership() {
	// 获取所有实例
	param := vo.SelectInstancesParam{
		ServiceName: e.electionKey,
		HealthyOnly: true,
	}
	instances, err := e.client.SelectInstances(param)
	if err != nil {
		e.error(err, "获取服务实例失败")
		return
	}

	if len(instances) == 0 {
		e.log("没有找到服务实例，重新注册自己")
		if err := e.registerSelf(); err != nil {
			e.error(err, "重新注册实例失败")
		}
		return
	}

	// 按元数据中的注册时间排序
	sort.Slice(instances, func(i, j int) bool {
		// 使用实例注册时间戳排序
		return instances[i].Metadata["registerTime"] < instances[j].Metadata["registerTime"]
	})

	// 第一个实例是leader
	isLeader := instances[0].Metadata["candidateInfo"] == e.candidateInfo
	e.notifyLeaderChange(isLeader)
}

func (e *nacosElect) registerSelf() error {
	// 使用当前时间作为注册时间
	registerTime := time.Now().Format(time.RFC3339Nano)

	// 注册当前实例
	_, err := e.client.RegisterInstance(vo.RegisterInstanceParam{
		ServiceName: e.electionKey,
		Ip:          "127.0.0.1", // 使用本地IP，实际不重要
		Port:        0,           // 端口设为0，因为我们不真正提供服务
		Metadata: map[string]string{
			"candidateInfo": e.candidateInfo,
			"registerTime":  registerTime,
		},
		Ephemeral: true, // 临时实例
	})
	return err
}

func (e *nacosElect) notifyLeaderChange(isLeader bool) {
	for _, ch := range e.leaderChs {
		go func(ch chan<- bool) {
			ch <- isLeader
		}(ch)
	}
}

func (e *nacosElect) error(err error, msg ...any) {
	ulogs.Error(append([]any{"【Nacos选leader】", err.Error()}, msg...)...)
}

func (e *nacosElect) log(args ...any) {
	ulogs.Log(append([]any{"【Nacos选leader】"}, args...)...)
}
