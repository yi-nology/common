package xerrors

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	// 创建原始错误
	originalErr := errors.New("this is an original error")

	// 使用 Wrap 函数创建自定义错误
	err := Wrap(100, "this is a custom error", originalErr)

	// 断言错误类型
	customErr, ok := err.(*Error)
	if !ok {
		t.Fatalf("expected *Error, got %T", err)
	}

	// 检查错误码
	if customErr.Code != 100 {
		t.Fatalf("expected code 100, got %d", customErr.Code)
	}

	// 检查错误消息
	expectedMsg := "this is a custom error: this is an original error"
	if customErr.Error() != expectedMsg {
		t.Fatalf("expected message %q, got %q", expectedMsg, customErr.Error())
	}

	// 检查原始错误
	if !errors.Is(err, originalErr) {
		t.Fatalf("expected original error in cause")
	}

	// 检查堆栈信息
	if customErr.Stack == "" {
		t.Fatalf("expected stack trace, got none")
	}
}

func TestIsType(t *testing.T) {
	// 创建原始错误
	originalErr := errors.New("this is an original error")

	// 使用 Wrap 函数创建自定义错误
	err := Wrap(100, "this is a custom error", originalErr)

	// 使用 IsType 函数检查错误类型
	if !IsType(err, 100) {
		t.Fatalf("expected IsType to return true for code 100")
	}

	// 使用 IsType 函数检查错误类型
	if IsType(err, 200) {
		t.Fatalf("expected IsType to return false for code 200")
	}
}
