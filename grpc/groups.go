package grpc

import (
	"context"
	"github.com/bradenrayhorn/switchboard-core/repositories"
	"github.com/bradenrayhorn/switchboard-protos/groups"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupsServer struct{}

func (g GroupsServer) GetGroups(c context.Context, r *groups.GetGroupsRequest) (*groups.GetGroupsResponse, error) {
	userId, err := primitive.ObjectIDFromHex(r.UserId)
	if err != nil {
		return nil, err
	}
	accessibleGroups, err := repositories.Group.GetGroups(userId)
	if err != nil {
		return nil, err
	}
	var groupIds []string
	for _, group := range accessibleGroups {
		groupIds = append(groupIds, group.ID.Hex())
	}
	return &groups.GetGroupsResponse{GroupIds: groupIds}, nil
}
