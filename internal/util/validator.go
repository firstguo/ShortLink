package util

import (
	"net/url"
	"strings"
)

// ValidateURL 验证 URL 是否合法
func ValidateURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return &ValidationError{Field: "original_url", Message: "URL cannot be empty"}
	}

	// 限制 URL 长度
	if len(rawURL) > 2048 {
		return &ValidationError{Field: "original_url", Message: "URL length cannot exceed 2048 characters"}
	}

	// 解析 URL
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return &ValidationError{Field: "original_url", Message: "Invalid URL format"}
	}

	// 验证协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return &ValidationError{Field: "original_url", Message: "URL scheme must be http or https"}
	}

	// 验证主机
	if parsedURL.Host == "" {
		return &ValidationError{Field: "original_url", Message: "URL must contain a host"}
	}

	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
