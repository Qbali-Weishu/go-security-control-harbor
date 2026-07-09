package domain

type conflictCarrier interface {
	GetConflicts() []Conflict
}

func (c Component) GetConflicts() []Conflict {
	return c.Conflicts
}
