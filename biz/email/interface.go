package email

import (
	"context"
	"io"
	"time"
)

// Sender 邮件发送器接口
type Sender interface {
	// Send 发送邮件
	Send(ctx context.Context, message *Message) error

	// SendBatch 批量发送邮件
	SendBatch(ctx context.Context, messages []*Message) error

	// Close 关闭连接
	Close() error
}

// Message 邮件消息
type Message struct {
	// 基本信息
	From    string   // 发件人邮箱
	To      []string // 收件人邮箱列表
	Cc      []string // 抄送邮箱列表
	Bcc     []string // 密送邮箱列表
	ReplyTo string   // 回复邮箱

	// 邮件内容
	Subject string // 主题
	Text    string // 纯文本内容
	HTML    string // HTML内容

	// 附件
	Attachments []*Attachment

	// 邮件头
	Headers map[string]string

	// 优先级
	Priority Priority

	// 发送时间(可选,为空表示立即发送)
	SendAt *time.Time
}

// Attachment 附件
type Attachment struct {
	Filename    string    // 文件名
	ContentType string    // MIME类型
	Data        io.Reader // 文件内容
}

// Priority 邮件优先级
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityNormal Priority = 3
	PriorityHigh   Priority = 5
)

// Config 邮件配置
type Config struct {
	// SMTP配置
	Host     string // SMTP服务器地址
	Port     int    // SMTP端口
	Username string // 用户名
	Password string // 密码

	// TLS配置
	UseTLS          bool // 是否使用TLS
	InsecureSkipTLS bool // 是否跳过TLS证书验证

	// 发送配置
	FromEmail string // 默认发件人邮箱
	FromName  string // 默认发件人名称

	// 连接池配置
	PoolSize        int           // 连接池大小
	MaxRetries      int           // 最大重试次数
	RetryInterval   time.Duration // 重试间隔
	ConnectTimeout  time.Duration // 连接超时
	SendTimeout     time.Duration // 发送超时
	KeepAlive       bool          // 是否保持连接
	KeepAlivePeriod time.Duration // 保持连接周期
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Port:            587,
		UseTLS:          true,
		InsecureSkipTLS: false,
		PoolSize:        10,
		MaxRetries:      3,
		RetryInterval:   time.Second * 2,
		ConnectTimeout:  time.Second * 10,
		SendTimeout:     time.Second * 30,
		KeepAlive:       true,
		KeepAlivePeriod: time.Minute * 5,
	}
}

// Validate 验证邮件消息
func (m *Message) Validate() error {
	if m.From == "" {
		return ErrInvalidFrom
	}
	if len(m.To) == 0 {
		return ErrNoRecipients
	}
	if m.Subject == "" {
		return ErrEmptySubject
	}
	if m.Text == "" && m.HTML == "" {
		return ErrEmptyBody
	}
	return nil
}

// HasAttachment 是否有附件
func (m *Message) HasAttachment() bool {
	return len(m.Attachments) > 0
}

// IsHTML 是否是HTML邮件
func (m *Message) IsHTML() bool {
	return m.HTML != ""
}
