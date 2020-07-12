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
	GroupName string   `form:"name" json:"name"`
	UserIds   []string `json:"users" form:"users" binding:"required"`
}

type UpdateGroupRequest struct {
	Id            string   `form:"id" json:"id"`
	GroupName     string   `form:"name" json:"name"`
	AddUserIds    []string `json:"add_users" form:"add_users"`
	RemoveUserIds []string `json:"remove_users" form:"remove_users"`
}

// users

type UserSearchRequest struct {
	Name string `form:"username" json:"username" binding:"required"`
}

// organizations

type CreateOrganizationRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}
