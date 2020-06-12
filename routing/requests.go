package routing

// auth

type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// group

type CreateGroupRequest struct {
	GroupName string   `form:"name"`
	UserIds   []string `json:"users" form:"users" binding:"required"`
}
