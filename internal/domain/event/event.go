package event

import "time"

// DomainEvent 领域事件接口
type DomainEvent interface {
	EventType() string
	OccurredOn() time.Time
}

// BaseEvent 基础事件
type BaseEvent struct {
	OccurredAt time.Time
}

// OccurredOn 返回事件发生时间
func (e *BaseEvent) OccurredOn() time.Time {
	return e.OccurredAt
}
