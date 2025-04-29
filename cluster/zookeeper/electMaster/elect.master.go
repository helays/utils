package electMaster

import (
	"errors"
	"github.com/go-zookeeper/zk"
	"github.com/helays/utils/db/zookeeper"
	"github.com/helays/utils/logger/ulogs"
	"path"
	"sort"
	"strings"
	"time"
)

type zkElect struct {
	conn          *zookeeper.Client
	electionKey   string // 选举路径
	candidateInfo string // 当前节点信息
	leaderChs     []chan<- bool
	nodePath      string // 当前节点路径
	stopCh        chan struct{}
}

// ElectMaster Zookeeper选举实现
func ElectMaster(conn *zookeeper.Client, electionKey, candidateInfo string, leaderChs ...chan<- bool) {
	e := &zkElect{
		conn:          conn,
		electionKey:   electionKey,
		candidateInfo: candidateInfo,
		leaderChs:     leaderChs,
		stopCh:        make(chan struct{}),
	}
	go e.process()
}

func (e *zkElect) process() {
	// 确保选举路径存在
	if err := e.ensurePathExists(e.electionKey); err != nil {
		e.error(err, "创建选举路径失败")
		return
	}

	// 创建临时顺序节点
	nodePath, err := e.conn.GetConn().Create(
		path.Join(e.electionKey, "node-"),
		[]byte(e.candidateInfo),
		zk.FlagEphemeral|zk.FlagSequence,
		zk.WorldACL(zk.PermAll),
	)
	if err != nil {
		e.error(err, "创建临时节点失败")
		return
	}
	e.nodePath = nodePath

	// 开始选举循环
	e.electLoop()
}

func (e *zkElect) electLoop() {
	for {
		select {
		case <-e.stopCh:
			e.log("停止Zookeeper选举")
			return
		default:
			isLeader, watchNode, err := e.checkLeadership()
			if err != nil {
				e.error(err, "检查领导权失败")
				time.Sleep(1 * time.Second)
				continue
			}

			e.notifyLeaderChange(isLeader)

			// 如果不是leader，监听前一个节点
			if !isLeader && watchNode != "" {
				_, _, ch, err := e.conn.GetConn().GetW(watchNode)
				if err != nil {
					e.error(err, "监听节点失败", watchNode)
					time.Sleep(1 * time.Second)
					continue
				}

				select {
				case <-ch: // 前一个节点变化，重新检查
				case <-e.stopCh:
					return
				}
			} else {
				// 是leader，等待停止信号
				select {
				case <-e.stopCh:
					return
				}
			}
		}
	}
}

func (e *zkElect) checkLeadership() (bool, string, error) {
	children, _, err := e.conn.GetConn().Children(e.electionKey)
	if err != nil {
		return false, "", err
	}

	if len(children) == 0 {
		return false, "", nil
	}

	// 获取所有节点并排序
	sort.Strings(children)
	currentNode := path.Base(e.nodePath)

	// 找到当前节点的位置
	var currentIndex int = -1
	for i, child := range children {
		if child == currentNode {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		return false, "", nil
	}

	// 第一个节点是leader
	if currentIndex == 0 {
		return true, "", nil
	}

	// 不是leader，返回前一个节点路径用于监听
	prevNode := children[currentIndex-1]
	return false, path.Join(e.electionKey, prevNode), nil
}

func (e *zkElect) ensurePathExists(path string) error {
	exists, _, err := e.conn.GetConn().Exists(path)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// 递归创建路径
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentPath := ""
	for _, part := range parts {
		currentPath += "/" + part
		exists, _, err = e.conn.GetConn().Exists(currentPath)
		if err != nil {
			return err
		}
		if !exists {
			_, err = e.conn.GetConn().Create(currentPath, nil, 0, zk.WorldACL(zk.PermAll))
			if err != nil && !errors.Is(err, zk.ErrNodeExists) {
				return err
			}
		}
	}
	return nil
}

func (e *zkElect) notifyLeaderChange(isLeader bool) {
	for _, ch := range e.leaderChs {
		go func(ch chan<- bool) {
			ch <- isLeader
		}(ch)
	}
}

func (e *zkElect) error(err error, msg ...any) {
	ulogs.Error(append([]any{"【Zookeeper选leader】", err.Error()}, msg...)...)
}

func (e *zkElect) log(args ...any) {
	ulogs.Log(append([]any{"【Zookeeper选leader】"}, args...)...)
}
