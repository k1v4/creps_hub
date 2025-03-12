package entity

type User struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Password      []byte `json:"-"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Username      string `json:"username"`
	AccessLevelId int    `json:"-"`
}
