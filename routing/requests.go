package routing

import "go.mongodb.org/mongo-driver/bson/primitive"

// auth

type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// channels

type CreateChannelRequest struct {
	Name           string             `form:"name" json:"name" binding:"required"`
	OrganizationID primitive.ObjectID `form:"organization_id" json:"organization_id" binding:"required"`
	Private        bool               `form:"private" json:"private"`
}

type LeaveChannelRequest struct {
	ChannelID primitive.ObjectID `form:"channel_id" json:"channel_id" binding:"required"`
}

type JoinChannelRequest struct {
	ChannelID primitive.ObjectID `form:"channel_id" json:"channel_id" binding:"required"`
}

// users

type UserSearchRequest struct {
	Name           string `form:"username" json:"username" binding:"required"`
	OrganizationID string `form:"organization_id" json:"organization_id" binding:"required"`
}

// organizations

type CreateOrganizationRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

type AddUserToOrganizationRequest struct {
	OrganizationID string `form:"organization_id" json:"organization_id" binding:"required"`
	Username       string `form:"username" json:"username" binding:"required"`
}
