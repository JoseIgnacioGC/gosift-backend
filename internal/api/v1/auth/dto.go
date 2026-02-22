package auth

type RegisterRequestDto struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Timezone string `json:"timezone" binding:"omitempty,timezone"`
	Language string `json:"language" binding:"omitempty,bcp47_language_tag"`
}

type RegisterResponseDto struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Timezone  string `json:"timezone"`
	Language  string `json:"language"`
	CreatedAt string `json:"created_at"`
	Token     string `json:"token"`
}
