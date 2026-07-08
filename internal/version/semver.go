package version

import (
	"fmt"
	"strconv"
	"strings"
)

// Version 表示语义版本号
type Version struct {
	Parts [3]int
}

// Parse 解析语义版本字符串
func Parse(raw string) (Version, error) {
	fields := strings.Split(raw, ".")
	if len(fields) != 3 {
		return Version{}, fmt.Errorf("期望包含三个数字部分的语义版本: %q", raw)
	}
	var v Version
	for i, part := range fields {
		n, err := strconv.Atoi(part)
		if err != nil {
			return Version{}, fmt.Errorf("无效的数字段 %q 在 %q 中: %w", part, raw, err)
		}
		v.Parts[i] = n
	}
	return v, nil
}

// Compare 比较两个版本字符串
// 返回: -1 (left < right), 0 (left == right), 1 (left > right)
func Compare(left, right string) (int, error) {
	l, err := Parse(left)
	if err != nil {
		return 0, err
	}
	r, err := Parse(right)
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(l.Parts); i++ {
		if l.Parts[i] < r.Parts[i] {
			return -1, nil
		}
		if l.Parts[i] > r.Parts[i] {
			return 1, nil
		}
	}
	return 0, nil
}

// InRange 检查候选版本是否在指定范围内
// min 是包含的，maxExclusive 是排除的
func InRange(candidate, min, maxExclusive string) (bool, error) {
	if min != "" {
		cmp, err := Compare(candidate, min)
		if err != nil {
			return false, err
		}
		if cmp < 0 {
			return false, nil
		}
	}
	if maxExclusive != "" {
		cmp, err := Compare(candidate, maxExclusive)
		if err != nil {
			return false, err
		}
		if cmp >= 0 {
			return false, nil
		}
	}
	return true, nil
}
