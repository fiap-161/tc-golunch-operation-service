package entity

type Admin struct {
	Id       string
	Email    string
	Password string
}

func (a Admin) Build(password string) Admin {
	return Admin{
		Id:       a.Id,
		Email:    a.Email,
		Password: password,
	}
}
