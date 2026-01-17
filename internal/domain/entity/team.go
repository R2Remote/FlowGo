package entity

type Team struct {
	BaseEntity
	Name        string
	Description string
	OwnerId     uint64
}
