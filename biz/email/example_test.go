package email_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yi-nology/common/biz/email"
)

// ExampleSMTPSender_Send 发送简单邮件示例
func ExampleSMTPSender_Send() {
	// 创建配置
	config := &email.Config{
		Host:      "smtp.example.com",
		Port:      587,
		Username:  "user@example.com",
		Password:  "password",
		UseTLS:    true,
		FromEmail: "sender@example.com",
		FromName:  "发件人名称",
	}

	// 创建发送器
	sender, err := email.NewSMTPSender(config)
	if err != nil {
		panic(err)
	}
	defer sender.Close()

	// 构建邮件
	message := &email.Message{
		To:      []string{"recipient@example.com"},
		Subject: "测试邮件",
		Text:    "这是一封测试邮件的纯文本内容。",
	}

	// 发送邮件
	ctx := context.Background()
	if err := sender.Send(ctx, message); err != nil {
		panic(err)
	}

	fmt.Println("邮件发送成功")
}

// ExampleSendHTML 发送HTML邮件示例
func ExampleSendHTML() {
	config := email.DefaultConfig()
	config.Host = "smtp.example.com"
	config.Username = "user@example.com"
	config.Password = "password"
	config.FromEmail = "sender@example.com"

	sender, err := email.NewSMTPSender(config)
	if err != nil {
		panic(err)
	}
	defer sender.Close()

	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; }
			.header { background-color: #4CAF50; color: white; padding: 10px; }
			.content { padding: 20px; }
		</style>
	</head>
	<body>
		<div class="header">
			<h1>欢迎使用邮件SDK</h1>
		</div>
		<div class="content">
			<p>这是一封HTML格式的测试邮件。</p>
			<p>支持丰富的格式和样式。</p>
		</div>
	</body>
	</html>
	`

	message := &email.Message{
		To:       []string{"recipient@example.com"},
		Cc:       []string{"cc@example.com"},
		Subject:  "HTML邮件测试",
		HTML:     htmlContent,
		Priority: email.PriorityHigh,
	}

	ctx := context.Background()
	if err := sender.Send(ctx, message); err != nil {
		panic(err)
	}

	fmt.Println("HTML邮件发送成功")
}

// ExampleSendWithAttachment 发送带附件的邮件示例
func ExampleSendWithAttachment() {
	config := &email.Config{
		Host:      "smtp.example.com",
		Port:      587,
		Username:  "user@example.com",
		Password:  "password",
		FromEmail: "sender@example.com",
	}

	sender, err := email.NewSMTPSender(config)
	if err != nil {
		panic(err)
	}
	defer sender.Close()

	// 创建附件
	attachmentData := strings.NewReader("这是附件的内容")
	attachment := &email.Attachment{
		Filename:    "document.txt",
		ContentType: "text/plain",
		Data:        attachmentData,
	}

	message := &email.Message{
		To:          []string{"recipient@example.com"},
		Subject:     "带附件的邮件",
		Text:        "请查看附件内容。",
		Attachments: []*email.Attachment{attachment},
	}

	ctx := context.Background()
	if err := sender.Send(ctx, message); err != nil {
		panic(err)
	}

	fmt.Println("带附件的邮件发送成功")
}

// ExampleSMTPSender_SendBatch 批量发送邮件示例
func ExampleSMTPSender_SendBatch() {
	config := email.DefaultConfig()
	config.Host = "smtp.example.com"
	config.Username = "user@example.com"
	config.Password = "password"
	config.FromEmail = "sender@example.com"

	sender, err := email.NewSMTPSender(config)
	if err != nil {
		panic(err)
	}
	defer sender.Close()

	// 构建多封邮件
	messages := []*email.Message{
		{
			To:      []string{"user1@example.com"},
			Subject: "通知1",
			Text:    "这是第一封通知邮件",
		},
		{
			To:      []string{"user2@example.com"},
			Subject: "通知2",
			Text:    "这是第二封通知邮件",
		},
		{
			To:      []string{"user3@example.com"},
			Subject: "通知3",
			Text:    "这是第三封通知邮件",
		},
	}

	// 批量发送
	ctx := context.Background()
	if err := sender.SendBatch(ctx, messages); err != nil {
		panic(err)
	}

	fmt.Println("批量邮件发送成功")
}

// ExampleMessage_Validate 验证邮件消息示例
func ExampleMessage_Validate() {
	message := &email.Message{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "测试邮件",
		Text:    "邮件内容",
	}

	if err := message.Validate(); err != nil {
		fmt.Println("邮件验证失败:", err)
		return
	}

	fmt.Println("邮件验证通过")
}

// ExampleDefaultConfig 使用默认配置示例
func ExampleDefaultConfig() {
	config := email.DefaultConfig()

	// 修改必要的配置
	config.Host = "smtp.gmail.com"
	config.Username = "your-email@gmail.com"
	config.Password = "your-app-password"
	config.FromEmail = "your-email@gmail.com"
	config.FromName = "Your Name"

	fmt.Printf("默认配置:\n")
	fmt.Printf("  端口: %d\n", config.Port)
	fmt.Printf("  使用TLS: %v\n", config.UseTLS)
	fmt.Printf("  连接池大小: %d\n", config.PoolSize)
	fmt.Printf("  最大重试次数: %d\n", config.MaxRetries)
	fmt.Printf("  重试间隔: %v\n", config.RetryInterval)
	fmt.Printf("  连接超时: %v\n", config.ConnectTimeout)
	fmt.Printf("  发送超时: %v\n", config.SendTimeout)
}

// ExamplePriority 设置邮件优先级示例
func ExamplePriority() {
	message := &email.Message{
		To:       []string{"recipient@example.com"},
		Subject:  "紧急通知",
		Text:     "这是一封高优先级的邮件",
		Priority: email.PriorityHigh,
	}

	fmt.Printf("邮件优先级: %d\n", message.Priority)
	// Output: 邮件优先级: 5
}

// ExampleMessage_scheduleExample 计划发送邮件示例
func ExampleMessage_scheduleExample() {
	sendTime := time.Now().Add(1 * time.Hour)

	message := &email.Message{
		To:      []string{"recipient@example.com"},
		Subject: "计划发送的邮件",
		Text:    "这封邮件将在1小时后发送",
		SendAt:  &sendTime,
	}

	fmt.Printf("邮件计划发送时间: %s\n", message.SendAt.Format(time.RFC3339))
}
