package electMaster

import (
	"fmt"
	"github.com/helays/utils/logger/ulogs"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sort"
	"time"
)

type NacosElect struct {
	Client        naming_client.INamingClient
	ElectionKey   string // Nacos服务名
	CandidateInfo string // 当前节点唯一编号
	LeaderChs     []chan<- bool
	stopCh        chan struct{}
	Ip            string
	Port          uint64
	isLeader      bool // 标记当前是否为 Leader
}

// ElectMaster Nacos选举实现
func ElectMaster(e *NacosElect) {
	e.stopCh = make(chan struct{})
	go e.process()
}

func (e *NacosElect) process() {
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

func (e *NacosElect) checkLeadership() {
	// 获取所有实例
	instances, err := e.Client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: e.ElectionKey,
		HealthyOnly: true,
		GroupName:   "DEFAULT_GROUP",
	})
	if err != nil {
		e.error(err, fmt.Sprintf("获取服务实例失败 %s", e.ElectionKey))
		return
	}
	if len(instances) == 0 {
		e.log("没有找到服务实例，重新注册自己")
		if err = e.registerSelf(); err != nil {
			e.error(err, "重新注册实例失败")
		}
		return
	}
	// 复制实例数据
	copiedInstances := make([]model.Instance, len(instances))
	copy(copiedInstances, instances)

	// 按元数据中的注册时间排序
	sort.Slice(copiedInstances, func(i, j int) bool {
		// 使用实例注册时间戳排序
		return copiedInstances[i].Metadata["registerTime"] < copiedInstances[j].Metadata["registerTime"]
	})

	// 第一个实例是leader
	newLeader := copiedInstances[0].Metadata["candidateInfo"] == e.CandidateInfo
	if newLeader != e.isLeader {
		e.isLeader = newLeader
		if e.isLeader {
			e.log("当前节点成为leader")
		}
		e.notifyLeaderChange(newLeader)
	}
}

func (e *NacosElect) registerSelf() error {
	param := vo.RegisterInstanceParam{
		ServiceName: e.ElectionKey,
		Ip:          e.Ip,
		Port:        e.Port,
		Metadata: map[string]string{
			"candidateInfo": e.CandidateInfo,
			"registerTime":  time.Now().Format(time.RFC3339Nano),
		},
		Ephemeral: true,
		Healthy:   true,
		Enable:    true,
		Weight:    0.1, // 必须 >0
	}

	success, err := e.Client.RegisterInstance(param)
	if err != nil {
		return fmt.Errorf("注册实例失败: %v", err)
	}

	if !success {
		return fmt.Errorf("注册实例返回失败")
	}

	e.log("成功注册实例:", e.CandidateInfo)
	return nil
}

func (e *NacosElect) notifyLeaderChange(isLeader bool) {
	for _, ch := range e.LeaderChs {
		go func(ch chan<- bool) {
			ch <- isLeader
		}(ch)
	}
}

func (e *NacosElect) error(err error, msg ...any) {
	ulogs.Error(append([]any{"【Nacos选leader】", err.Error()}, msg...)...)
}

func (e *NacosElect) log(args ...any) {
	ulogs.Log(append([]any{"【Nacos选leader】"}, args...)...)
}
