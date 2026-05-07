package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

// SMTPSender SMTP邮件发送器
type SMTPSender struct {
	config *Config
	auth   smtp.Auth
}

// NewSMTPSender 创建SMTP邮件发送器
func NewSMTPSender(config *Config) (*SMTPSender, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	var auth smtp.Auth
	if config.Username != "" && config.Password != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	return &SMTPSender{
		config: config,
		auth:   auth,
	}, nil
}

// Send 发送邮件
func (s *SMTPSender) Send(ctx context.Context, message *Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	// 设置默认发件人
	if message.From == "" {
		if s.config.FromEmail == "" {
			return ErrInvalidFrom
		}
		message.From = s.formatAddress(s.config.FromEmail, s.config.FromName)
	}

	// 构建邮件内容
	body, err := s.buildMessage(message)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// 发送邮件
	return s.sendWithRetry(ctx, message, body)
}

// SendBatch 批量发送邮件
func (s *SMTPSender) SendBatch(ctx context.Context, messages []*Message) error {
	for _, msg := range messages {
		if err := s.Send(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭连接
func (s *SMTPSender) Close() error {
	// SMTP连接是短连接,不需要特殊关闭
	return nil
}

// sendWithRetry 发送邮件(带重试)
func (s *SMTPSender) sendWithRetry(ctx context.Context, message *Message, body []byte) error {
	var lastErr error

	for i := 0; i < s.config.MaxRetries; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(s.config.RetryInterval):
			}
		}

		if err := s.send(ctx, message, body); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	if lastErr != nil {
		return fmt.Errorf("%w: %v", ErrTooManyRetries, lastErr)
	}
	return ErrSendFailed
}

// send 实际发送邮件
func (s *SMTPSender) send(ctx context.Context, message *Message, body []byte) error {
	// 创建带超时的连接
	conn, err := s.dial(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectFailed, err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	// 如果支持STARTTLS,启用TLS
	if s.config.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName:         s.config.Host,
				InsecureSkipVerify: s.config.InsecureSkipTLS,
			}
			if err := client.StartTLS(tlsConfig); err != nil {
				return err
			}
		}
	}

	// 认证
	if s.auth != nil {
		if err := client.Auth(s.auth); err != nil {
			return fmt.Errorf("%w: %v", ErrAuthFailed, err)
		}
	}

	// 设置发件人
	if err := client.Mail(s.extractEmail(message.From)); err != nil {
		return err
	}

	// 设置收件人
	recipients := append(message.To, message.Cc...)
	recipients = append(recipients, message.Bcc...)
	for _, addr := range recipients {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := w.Write(body); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return client.Quit()
}

// dial 创建连接
func (s *SMTPSender) dial(ctx context.Context) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	dialer := &net.Dialer{
		Timeout:   s.config.ConnectTimeout,
		KeepAlive: s.config.KeepAlivePeriod,
	}

	return dialer.DialContext(ctx, "tcp", addr)
}

// buildMessage 构建邮件内容
func (s *SMTPSender) buildMessage(message *Message) ([]byte, error) {
	var buf strings.Builder

	// 邮件头
	s.writeHeaders(&buf, message)

	// 邮件体
	if message.HasAttachment() {
		return s.buildMIMEMessage(&buf, message)
	}

	if message.IsHTML() {
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(message.HTML)
	} else {
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		buf.WriteString("\r\n")
		buf.WriteString(message.Text)
	}

	return []byte(buf.String()), nil
}

// writeHeaders 写入邮件头
func (s *SMTPSender) writeHeaders(buf *strings.Builder, message *Message) {
	buf.WriteString(fmt.Sprintf("From: %s\r\n", message.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(message.To, ", ")))

	if len(message.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(message.Cc, ", ")))
	}

	if message.ReplyTo != "" {
		buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", message.ReplyTo))
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", s.encodeSubject(message.Subject)))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))

	// 优先级
	if message.Priority != PriorityNormal {
		buf.WriteString(fmt.Sprintf("X-Priority: %d\r\n", message.Priority))
	}

	// 自定义邮件头
	for key, value := range message.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
}

// buildMIMEMessage 构建MIME邮件(带附件)
func (s *SMTPSender) buildMIMEMessage(buf *strings.Builder, message *Message) ([]byte, error) {
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

	// 邮件正文
	if message.IsHTML() {
		s.writeMIMEPart(writer, "text/html", message.HTML)
	} else {
		s.writeMIMEPart(writer, "text/plain", message.Text)
	}

	// 附件
	for _, attachment := range message.Attachments {
		if err := s.writeAttachment(writer, attachment); err != nil {
			return nil, err
		}
	}

	writer.Close()
	return []byte(buf.String()), nil
}

// writeMIMEPart 写入MIME部分
func (s *SMTPSender) writeMIMEPart(writer *multipart.Writer, contentType, content string) error {
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", contentType+"; charset=UTF-8")
	header.Set("Content-Transfer-Encoding", "quoted-printable")

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(content))
	return err
}

// writeAttachment 写入附件
func (s *SMTPSender) writeAttachment(writer *multipart.Writer, attachment *Attachment) error {
	header := textproto.MIMEHeader{}

	if attachment.ContentType == "" {
		attachment.ContentType = "application/octet-stream"
	}

	header.Set("Content-Type", attachment.ContentType)
	header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s",
		mime.QEncoding.Encode("UTF-8", filepath.Base(attachment.Filename))))
	header.Set("Content-Transfer-Encoding", "base64")

	part, err := writer.CreatePart(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, attachment.Data)
	return err
}

// formatAddress 格式化邮箱地址
func (s *SMTPSender) formatAddress(email, name string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}

// extractEmail 提取邮箱地址
func (s *SMTPSender) extractEmail(addr string) string {
	if idx := strings.Index(addr, "<"); idx >= 0 {
		if end := strings.Index(addr[idx:], ">"); end >= 0 {
			return addr[idx+1 : idx+end]
		}
	}
	return addr
}

// encodeSubject 编码主题
func (s *SMTPSender) encodeSubject(subject string) string {
	return mime.QEncoding.Encode("UTF-8", subject)
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Host == "" {
		return fmt.Errorf("%w: host is required", ErrInvalidConfig)
	}
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("%w: invalid port", ErrInvalidConfig)
	}
	return nil
}
