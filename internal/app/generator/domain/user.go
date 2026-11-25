package domain

type Settings struct {
	Id        int
	UserId    int
	AvatarUrl string
	FirstName string
	LastName  string
}

type NotificationActor struct {
	Id   int
	Name string
}
