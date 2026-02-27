package subscriptions

type CreateRequestDto struct {
	FeedURL  string `json:"feed_url" binding:"required,url,max=2048"`
	Title    string `json:"title" binding:"omitempty,min=1,max=200"`
	SiteURL  string `json:"site_url" binding:"omitempty,url,max=2048"`
	Category string `json:"category" binding:"omitempty,min=1,max=50"`
}

type UpdateRequestDto struct {
	Title    *string `json:"title" binding:"omitempty,min=1,max=200"`
	Category *string `json:"category" binding:"omitempty,min=1,max=50"`
	IsActive *bool   `json:"is_active" binding:"omitempty"`
}

type ResponseDto struct {
	ID        string `json:"id"`
	FeedURL   string `json:"feed_url"`
	Title     string `json:"title"`
	SiteURL   string `json:"site_url"`
	Category  string `json:"category"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
