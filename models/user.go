package models

import (
	proto "github.com/garden-raccoon/users-pkg/protocols/users"

	"github.com/gofrs/uuid"
)

type User struct {
	UserUUID         uuid.UUID
	Username         string
	Email            string
	Phone            string
	FirstName        string
	LastName         string
	Avatar           string
	ValidationStatus int
	Addresses        []*Address
}
type Address struct {
	Street string
	City   string
	Gps    string
}

type SignUpResponse struct {
	UserUUID uuid.UUID
}
type UpdateUserRequest struct {
	UserUUID         uuid.UUID
	ValidationStatus *int
	Username         *string
	Email            *string
	Phone            *string
	FirstName        *string
	LastName         *string
	Avatar           *string
	Addresses        []*Address
}

// UserFromProto is

func (u User) Proto() *proto.User {
	user := &proto.User{
		UserUuid:  u.UserUUID.Bytes(),
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Avatar:    u.Avatar,
		Addresses: convertToProtoAddresses(u.Addresses),
		Phone:     u.Phone,
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
		fields.Addresses = convertToProtoAddresses(u.Addresses)

	}

	if u.ValidationStatus != nil {
		fields.ValidationStatus = int64(*u.ValidationStatus)
	}

	return fields
}
func addressToProto(address *Address) *proto.Address {
	return &proto.Address{
		City:   address.City,
		Street: address.Street,
		Gps:    address.Gps,
	}
}
func addressFromProto(address *proto.Address) *Address {
	return &Address{
		City:   address.City,
		Street: address.Street,
		Gps:    address.Gps,
	}
}
func convertToProtoAddresses(a []*Address) *proto.Addresses {
	protoAddresses := &proto.Addresses{}
	for _, address := range a {
		protoAddresses.Addresses = append(protoAddresses.Addresses, addressToProto(address))
	}
	return protoAddresses
}
func convertFromProtoAddresses(pb *proto.Addresses) []*Address {
	var addresses []*Address
	for _, a := range pb.Addresses {
		addresses = append(addresses, addressFromProto(a))
	}
	return addresses
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
		req.Addresses = convertFromProtoAddresses(pb.Addresses)
	}
	if pb.ValidationStatus != 0 {
		req.ValidationStatus = int64ToPtrInt(pb.ValidationStatus)
	}
	return req
}
func UserFromProto(pb *proto.User) *User {
	req := &User{
		UserUUID: uuid.FromBytesOrNil(pb.UserUuid),
	}
	if pb.Email != "" {
		req.Email = pb.Email
	}

	if pb.Username != "" {
		req.Username = pb.Username
	}

	if pb.FirstName != "" {
		req.FirstName = pb.FirstName
	}

	if pb.LastName != "" {
		req.LastName = pb.LastName
	}
	if pb.Avatar != "" {
		req.Avatar = pb.Avatar
	}
	if pb.Phone != "" {
		req.Phone = pb.Phone
	}

	if pb.Addresses != nil {
		req.Addresses = convertFromProtoAddresses(pb.Addresses)
	}
	return req
}
func int64ToPtrInt(status int64) *int {
	toIntPtr := int(status)
	return &toIntPtr
}
