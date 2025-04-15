package models

type User struct {
	Id       uint    `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name"`
	Lastname string  `json:"lastname"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Todos    []Todos `gorm:"foreignKey:UserID" json:"Todos"` // <- relacja 1:N
}

type Todos struct {
	Id     uint   `json:"id"`
	UserID uint   `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status Status `json:"status"`
}
type Status string

const (
	Done       Status = "done"
	InProgress Status = "inProgress"
	Todo       Status = "todo"
)
