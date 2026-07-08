package domain

// conflictCarrier 是具有冲突的类型接口
type conflictCarrier interface {
	GetConflicts() []Conflict
}

// GetConflicts 返回组件的冲突列表
func (c Component) GetConflicts() []Conflict {
	return c.Conflicts
}
