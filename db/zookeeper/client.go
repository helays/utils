package zookeeper

import (
	"errors"
	"fmt"
	"github.com/helays/utils/v2/tools/backoff"
	"strings"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/helays/utils/v2/logger/ulogs"
	"github.com/helays/utils/v2/tools"
)

// Config Zookeeper客户端配置
type Config struct {
	Addrs           []string        `yaml:"addrs" json:"addrs" ini:"addrs"`                                  // ZK集群地址，格式：["host1:2181", "host2:2181"]
	SessionTimeout  time.Duration   `yaml:"session_timeout" json:"session_timeout" ini:"session_timeout"`    // 会话超时（建议15-30秒）
	BasePath        string          `yaml:"base_path" json:"base_path" ini:"base_path"`                      // 基础路径（如"/services/prod"）
	Auth            ACL             `yaml:"auth" json:"auth" ini:"auth"`                                     // 认证信息
	EnableEphemeral bool            `yaml:"enable_ephemeral" json:"enable_ephemeral" ini:"enable_ephemeral"` // 是否自动恢复临时节点
	Retries         backoff.Backoff `yaml:"retries" json:"retries" ini:"retries"`                            // 重试配置
}

// ACL 认证信息结构体
type ACL struct {
	Scheme   string `yaml:"scheme" json:"scheme" ini:"scheme"`       // 认证方案（digest/sasl）
	Password string `yaml:"password" json:"password" ini:"password"` // 认证密码
}

// Client 封装的Zookeeper客户端
type Client struct {
	config *Config
	conn   *zk.Conn      // 当前活跃的连接
	mu     sync.RWMutex  // 保护conn的读写锁
	stopCh chan struct{} // 停止信号通道

	// 临时节点管理
	ephemeralNodes   map[string][]byte // 需要重建的临时节点 path -> data
	ephemeralNodesMu sync.Mutex        // 保护临时节点map的锁
}

// NewClient 创建并返回一个新的Zookeeper客户端
func (c *Config) NewClient() (*Client, error) {
	// 参数校验
	if len(c.Addrs) == 0 {
		return nil, fmt.Errorf("zk服务器地址不能为空")
	}
	c.SessionTimeout = tools.AutoTimeDuration(c.SessionTimeout, time.Second, 15*time.Second)
	c.Retries.InitialSleep = tools.AutoTimeDuration(c.Retries.InitialSleep, time.Millisecond)
	c.Retries.MaxSleep = tools.AutoTimeDuration(c.Retries.MaxSleep, time.Second)
	c.Retries.Step = tools.AutoTimeDuration(c.Retries.Step, time.Millisecond)
	client := &Client{
		config:         c,
		stopCh:         make(chan struct{}),
		ephemeralNodes: make(map[string][]byte),
	}

	// 初始连接
	if err := client.connect(); err != nil {
		return nil, fmt.Errorf("初始连接失败: %v", err)
	}

	return client, nil
}

func (this *Config) getRetries() *backoff.Backoff {
	var more []any
	if this.Retries.Step > 0 {
		more = append(more, this.Retries.Step)
	}
	if this.Retries.Base > 0 {
		more = append(more, this.Retries.Base)
	}
	if this.Retries.StepThreshold > 0 {
		more = append(more, this.Retries.StepThreshold)
	}
	return backoff.NewBackoff(this.Retries.Type, this.Retries.InitialSleep, this.Retries.MaxSleep, more...)
}

// GetConn 获取当前活跃的连接（线程安全）
func (c *Client) GetConn() *zk.Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

// Close 关闭客户端
func (c *Client) Close() {
	close(c.stopCh)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.Close()
	}
}

// connect 内部连接/重连逻辑
func (c *Client) connect() error {
	// 创建新连接
	conn, eventChan, err := zk.Connect(c.config.Addrs, c.config.SessionTimeout)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}

	// 设置认证
	if c.config.Auth.Scheme != "" && c.config.Auth.Password != "" {
		if err = conn.AddAuth(c.config.Auth.Scheme, []byte(c.config.Auth.Password)); err != nil {
			conn.Close()
			return fmt.Errorf("认证失败: %v", err)
		}
	}

	// 启动事件监听
	go c.watchEvents(conn, eventChan)

	// 等待连接就绪
	if err = c.waitUntilConnected(conn, eventChan); err != nil {
		conn.Close()
		return err
	}

	// 初始化基础路径
	if c.config.BasePath != "" {
		if err = c.createPathRecursive(conn, c.config.BasePath); err != nil {
			conn.Close()
			return fmt.Errorf("创建基础路径失败: %v", err)
		}
	}

	// 更新当前连接
	c.mu.Lock()
	oldConn := c.conn
	c.conn = conn
	c.mu.Unlock()

	// 关闭旧连接
	if oldConn != nil {
		oldConn.Close()
	}

	return nil
}

// watchEvents 事件监听循环
func (c *Client) watchEvents(conn *zk.Conn, eventChan <-chan zk.Event) {
	for {
		select {
		case <-c.stopCh:
			return
		case event, ok := <-eventChan:
			if !ok {
				return
			}
			switch event.State {
			case zk.StateConnecting:
				ulogs.Info("[ZK] 正在连接服务器...")
			case zk.StateConnected:
				ulogs.Info("[ZK] 连接已建立")
				if c.config.EnableEphemeral {
					c.recreateEphemeralNodes(conn)
				}
			case zk.StateDisconnected:
				ulogs.Info("[ZK] 连接已断开，正在重建连接...")
			case zk.StateExpired:
				ulogs.Info("[ZK] 会话过期，需要重建连接")
				c.reconnect()
			case zk.StateAuthFailed:
				ulogs.Info("[ZK] 认证失败，请检查账号密码")
			}
		}
	}
}

// reconnect 带重试机制的连接重建
func (c *Client) reconnect() {
	reset := c.config.getRetries()
	defer reset.Reset()
	i := 1
	for {
		select {
		case <-c.stopCh:
			return
		default:
			if err := c.connect(); err == nil {
				ulogs.Info("[ZK] 重连成功")
				return
			}
			next := reset.Next()
			ulogs.Infof("[ZK] 第%d次重连失败，%s后重试", i+1, next.String())
			i++
			time.Sleep(next)
		}
	}
}

// waitUntilConnected 等待连接就绪
func (c *Client) waitUntilConnected(conn *zk.Conn, eventChan <-chan zk.Event) error {
	select {
	case event := <-eventChan:
		if event.State == zk.StateConnected {
			return nil
		}
		return fmt.Errorf("连接失败，状态: %v", event.State)
	case <-time.After(c.config.SessionTimeout):
		return fmt.Errorf("连接超时")
	case <-c.stopCh:
		return fmt.Errorf("客户端已关闭")
	}
}

// createPathRecursive 递归创建路径
func (c *Client) createPathRecursive(conn *zk.Conn, path string) error {
	current := ""
	for _, part := range strings.Split(strings.Trim(path, "/"), "/") {
		current += "/" + part
		exists, _, err := conn.Exists(current)
		if err != nil {
			return err
		}
		if !exists {
			if _, err := conn.Create(current, nil, 0, zk.WorldACL(zk.PermAll)); err != nil && err != zk.ErrNodeExists {
				return err
			}
		}
	}
	return nil
}

// 重建所有临时节点
func (c *Client) recreateEphemeralNodes(conn *zk.Conn) {
	c.ephemeralNodesMu.Lock()
	defer c.ephemeralNodesMu.Unlock()

	for path, data := range c.ephemeralNodes {
		go func(p string, d []byte) {
			reset := c.config.getRetries()
			defer reset.Reset()
			for {
				select {
				case <-c.stopCh:
					return
				default:
					_, err := conn.Create(p, d, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
					if err == nil || errors.Is(err, zk.ErrNodeExists) {
						return
					}
					time.Sleep(reset.Next())
				}
			}
		}(path, data)
	}
}

// CreateEphemeral 创建临时节点并记录以便重建
func (c *Client) CreateEphemeral(path string, data []byte) error {
	conn := c.GetConn()
	if conn == nil {
		return fmt.Errorf("连接未建立")
	}

	// 创建节点
	fullPath := c.config.BasePath + path
	_, err := conn.Create(fullPath, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	// 记录节点以便重建
	c.ephemeralNodesMu.Lock()
	c.ephemeralNodes[fullPath] = data
	c.ephemeralNodesMu.Unlock()

	return nil
}

// DeleteEphemeral 删除临时节点并从重建列表中移除
func (c *Client) DeleteEphemeral(path string) error {
	conn := c.GetConn()
	if conn == nil {
		return fmt.Errorf("连接未建立")
	}

	fullPath := c.config.BasePath + path
	err := conn.Delete(fullPath, -1)
	if err != nil && err != zk.ErrNoNode {
		return err
	}

	// 从重建列表中移除
	c.ephemeralNodesMu.Lock()
	delete(c.ephemeralNodes, fullPath)
	c.ephemeralNodesMu.Unlock()

	return nil
}
