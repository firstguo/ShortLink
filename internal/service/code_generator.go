package service

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Base62 字符集
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base62Len   = int64(len(base62Chars))
)

// CodeGenerator 短码生成器
type CodeGenerator struct {
	mu         sync.Mutex
	epoch      int64 // 起始时间戳（毫秒）
	workerID   int64 // 工作机器ID
	sequence   int64 // 序列号
	lastTime   int64 // 上次生成ID的时间
	codeLength int   // 短码长度
}

// NewCodeGenerator 创建短码生成器
func NewCodeGenerator(epoch time.Time, workerID int64, codeLength int) (*CodeGenerator, error) {
	if workerID < 0 || workerID > 0xFFF {
		return nil, fmt.Errorf("workerID must be between 0 and 4095")
	}

	if codeLength < 4 || codeLength > 10 {
		return nil, fmt.Errorf("codeLength must be between 4 and 10")
	}

	return &CodeGenerator{
		epoch:      epoch.UnixMilli(),
		workerID:   workerID,
		sequence:   0,
		lastTime:   0,
		codeLength: codeLength,
	}, nil
}

// Generate 生成短码
func (g *CodeGenerator) Generate() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixMilli()

	// 如果当前时间等于上次生成时间，序列号+1
	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & 0xFFF // 12位序列号，最大值4095
		if g.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for now <= g.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 时间变化，序列号重置为0
		g.sequence = 0
	}

	g.lastTime = now

	// 生成 Snowflake ID
	// 结构：(timestamp - epoch) << 22 | workerID << 12 | sequence
	id := ((now - g.epoch) << 22) | (g.workerID << 12) | g.sequence

	// 转换为 Base62 短码
	return ToBase62(id, g.codeLength)
}

// ToBase62 将整数转换为 Base62 编码的字符串
func ToBase62(num int64, length int) string {
	if num == 0 {
		return stringOfChar('0', length)
	}

	var result []byte
	for num > 0 {
		result = append(result, base62Chars[num%base62Len])
		num /= base62Len
	}

	// 反转字符串
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// 如果长度不足，左侧补 '0'
	if len(result) < length {
		padding := stringOfChar('0', length-len(result))
		result = append([]byte(padding), result...)
	}

	return string(result)
}

// FromBase62 将 Base62 编码的字符串转换为整数
func FromBase62(s string) int64 {
	var num int64
	for _, c := range s {
		num = num*base62Len + int64(indexOf(base62Chars, byte(c)))
	}
	return num
}

// stringOfChar 生成指定长度的重复字符
func stringOfChar(c byte, length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = c
	}
	return string(result)
}

// indexOf 查找字节在字符串中的位置
func indexOf(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
