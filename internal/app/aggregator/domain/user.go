package domain

type User struct {
	Id     int
	Email  string
	Active bool
}

type Settings struct {
	Id        int
	UserId    int
	AvatarUrl string
	FirstName string
	LastName  string
}
