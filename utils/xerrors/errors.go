package xerrors

import (
	"fmt"
	"runtime"
)

// 自定义错误类型
type Error struct {
	Code    int64  `json:"code,omitempty"`  // 错误码
	Message string `json:"error,omitempty"` // 错误消息
	Cause   error  `json:"-"`               // 原始错误
	Stack   string `json:"-"`               // 堆栈信息
}

// 实现 error 接口
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}
	return e.Message
}

// 实现 fmt.Stringer 接口
func (e *Error) String() string {
	return fmt.Sprintf("code:%d msg:%s err:%+v", e.Code, e.Message, e.Cause)
}

// 实现 errors.Wrapper 接口
func (e *Error) Unwrap() error {
	return e.Cause
}

// 检查错误是否与目标错误相同
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// 创建新的自定义错误
func New(code int64, msg string, cause error) *Error {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	stack := fmt.Sprintf("%s\n\t%s:%d", f.Name(), file, line)

	return &Error{
		Code:    code,
		Message: msg,
		Cause:   cause,
		Stack:   stack,
	}
}

// 包装原始错误为自定义错误
func Wrap(code int64, msg string, cause error) error {
	return New(code, msg, cause)
}

// 检查错误链中是否存在指定类型的错误
func IsType(err error, code int64) bool {
	if err == nil {
		return false
	}

	customErr, ok := err.(*Error)
	if !ok {
		return false
	}

	return customErr.Code == code
}
