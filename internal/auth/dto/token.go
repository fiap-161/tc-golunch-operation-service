package dto

type LoginRequestDTO struct {
	UserID   string `json:"user_id" binding:"required"`
	UserType string `json:"user_type" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponseDTO struct {
	Token string `json:"token"`
}
