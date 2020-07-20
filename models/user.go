package models

var (
	UserStatusPending = "pending"
	UserStatusComplete = "complete"
)

type User struct {
	ID string
	Email string
	Password string
	Name string
	Phone string
	Address string
	Status string
}