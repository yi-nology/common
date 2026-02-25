package email

import "errors"

var (
	// 邮件验证错误
	ErrInvalidFrom   = errors.New("invalid from address")
	ErrNoRecipients  = errors.New("no recipients specified")
	ErrEmptySubject  = errors.New("empty subject")
	ErrEmptyBody     = errors.New("empty email body")
	ErrInvalidConfig = errors.New("invalid email configuration")

	// 发送错误
	ErrSendFailed     = errors.New("failed to send email")
	ErrConnectFailed  = errors.New("failed to connect to SMTP server")
	ErrAuthFailed     = errors.New("SMTP authentication failed")
	ErrTimeout        = errors.New("email sending timeout")
	ErrTooManyRetries = errors.New("exceeded maximum retry attempts")

	// 附件错误
	ErrAttachmentTooLarge = errors.New("attachment size exceeds limit")
	ErrInvalidAttachment  = errors.New("invalid attachment")
)
