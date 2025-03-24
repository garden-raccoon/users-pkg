package models

import (
	proto "github.com/garden-raccoon/users-pkg/protocols/users"

	"github.com/gofrs/uuid"
)

type User struct {
	UserUUID  uuid.UUID
	Username  string
	Email     string
	Phone     string
	FirstName string
	LastName  string
	Addresses []string
}

type UpdateUserRequest struct {
	UserUUID  uuid.UUID
	Username  *string
	Email     *string
	Phone     *string
	FirstName *string
	LastName  *string
	Avatar    *string
	Addresses []string
}

// UserFromProto is
func UserFromProto(pb *proto.User) *User {
	return &User{
		UserUUID: uuid.FromBytesOrNil(pb.UserUuid),
		Email:    pb.Email,
		Username: pb.Username,
		Phone:    pb.Phone,
	}
}

func (u User) Proto() *proto.User {
	user := &proto.User{
		UserUuid: u.UserUUID.Bytes(),
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
	}
	return user
}

// Proto is
func Proto(u UpdateUserRequest) *proto.UpdateUserRequest {
	fields := &proto.UpdateUserRequest{UserUuid: u.UserUUID.Bytes()}

	if u.Email != nil {
		fields.Email = *u.Email
	}

	if u.Username != nil {
		fields.Username = *u.Username
	}

	if u.FirstName != nil {
		fields.FirstName = *u.FirstName
	}

	if u.LastName != nil {
		fields.LastName = *u.LastName
	}
	if u.Phone != nil {
		fields.Phone = *u.Phone
	}
	if u.Avatar != nil {
		fields.Avatar = *u.Avatar
	}
	if u.Addresses != nil {
		fields.Addresses = u.Addresses
	}
	return fields
}

// UpdateUserRequestFromProto is
func UpdateUserRequestFromProto(pb *proto.UpdateUserRequest) *UpdateUserRequest {
	req := &UpdateUserRequest{
		UserUUID: uuid.FromBytesOrNil(pb.UserUuid),
	}
	if pb.Email != "" {
		req.Email = &pb.Email
	}

	if pb.Username != "" {
		req.Username = &pb.Username
	}

	if pb.FirstName != "" {
		req.FirstName = &pb.FirstName
	}

	if pb.LastName != "" {
		req.LastName = &pb.LastName
	}
	if pb.Avatar != "" {
		req.Avatar = &pb.Avatar
	}
	if pb.Phone != "" {
		req.Phone = &pb.Phone
	}
	if pb.Addresses != nil {
		req.Addresses = pb.Addresses
	}
	return req
}
