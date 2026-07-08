package compat

import "compatgate/internal/domain"

func init() {
	// 确保 domain.Component 类型可用
	_ = domain.Component{}
}
