package electMaster

import (
	"fmt"
	"github.com/helays/utils/logger/ulogs"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
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
	param := vo.SelectInstancesParam{
		ServiceName: e.ElectionKey,
		HealthyOnly: true,
		GroupName:   "DEFAULT_GROUP",
	}
	instances, err := e.Client.SelectInstances(param)
	if err != nil {
		e.error(err, fmt.Sprintf("获取服务实例失败 %s", param.ServiceName))
		return
	}

	if len(instances) == 0 {
		e.log("没有找到服务实例，重新注册自己")
		if err = e.registerSelf(); err != nil {
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
	isLeader := instances[0].Metadata["candidateInfo"] == e.CandidateInfo
	e.notifyLeaderChange(isLeader)
}

func (e *NacosElect) registerSelf() error {
	registerTime := time.Now().Format(time.RFC3339Nano)

	param := vo.RegisterInstanceParam{
		ServiceName: e.ElectionKey,
		Ip:          e.Ip,
		Port:        e.Port,
		Metadata: map[string]string{
			"candidateInfo": e.CandidateInfo,
			"registerTime":  registerTime,
		},
		Ephemeral: true,
		Healthy:   true,
		Enable:    true,
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
