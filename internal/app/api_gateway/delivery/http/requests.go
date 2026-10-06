package http

type Generate1Request struct {
	NotificationType string `json:"notification_type"`
	Text             string `json:"text"`
	UserId           int    `json:"user_id"`
}

type GenerateNRequest struct {
	NotificationTypes []string `json:"notification_types"`
	Number            int      `json:"number"`
	UserIds           []int    `json:"user_ids"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
