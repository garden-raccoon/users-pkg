package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/garden-raccoon/users-pkg/models"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
	"time"

	proto "github.com/garden-raccoon/users-pkg/protocols/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type IUserAPI interface {

	// CheckAuth is
	CheckAuth(token []byte) (*models.User, error)

	GetUserRoles(token []byte) (roles []string, err error)

	GetIsAdmin(token []byte) (bool, error)
	// UserByUUID is
	UserByUUID(userUUID uuid.UUID) (*models.User, error)

	UserByPhone(phone string) (*models.User, error)

	UserByEmail(email string) (*models.User, error)

	UpdateUser(user *models.UpdateUserRequest) (*models.User, error)

	UpdateAddress(user *models.UpdateAddressRequest) (*models.User, error)
	AddAddresses(user *models.AddAddressRequest) (*models.User, error)
	DeleteAddress(addrUuid, userUuid uuid.UUID) error
	// SignUpByEmail is
	SignUpByEmail(email string, password []byte) (*models.SignUpResponse, error)

	// SignInByEmail is
	SignInByEmail(email string, password []byte) ([]byte, error)

	// SignUpByPhone is
	SignUpByPhone(phone string, password []byte) (*models.SignUpResponse, error)

	ResetPassword(phone string) error

	UpdatePassword(userUUID uuid.UUID, password []byte) error

	// SignInByPhone is
	SignInByPhone(phone string, password []byte) ([]byte, error)

	RequestPasswordChange(phone string) error
	VerifyUserByPhone(phone, checkPhrase string) ([]byte, error)

	VerifyPasswordChange(phone, checkPhrase string, password []byte) ([]byte, error)
	HealthCheck() error

	// Close GRPC Api connection
	Close() error
}

// UsersAPI is profile-service GRPC UsersAPI
// structure with client Connection
type UsersAPI struct {
	addr    string
	timeout time.Duration
	mu      sync.Mutex
	*grpc.ClientConn
	proto.UserServiceClient
	grpc_health_v1.HealthClient
}

// New create new Users IEmployerAPI instance
func New(addr string, timeout time.Duration) (IUserAPI, error) {
	api := &UsersAPI{timeout: timeout}

	if err := api.initConn(addr); err != nil {
		return nil, fmt.Errorf("create Users UsersAPI:  %w", err)
	}
	api.HealthClient = grpc_health_v1.NewHealthClient(api.ClientConn)

	api.UserServiceClient = proto.NewUserServiceClient(api.ClientConn)
	return api, nil
}

func (api *UsersAPI) RequestPasswordChange(phone string) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	req := &proto.PasswordChangeRequest{Phone: phone}
	if _, err := api.UserServiceClient.RequestPasswordChange(ctx, req); err != nil {
		return fmt.Errorf("request password change failed: %w", err)
	}
	return nil
}

func (api *UsersAPI) UpdateUser(user *models.UpdateUserRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	protoUser := models.Proto(user)
	resp, err := api.UserServiceClient.UpdateUser(ctx, protoUser)
	if err != nil {
		return nil, fmt.Errorf("updateUser api request: %w", err)
	}
	return models.UserFromProto(resp), nil
}
func (api *UsersAPI) UpdateAddress(user *models.UpdateAddressRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	protoUser := models.UpdateAddressToProto(user)
	resp, err := api.UserServiceClient.UpdateAddress(ctx, protoUser)
	if err != nil {
		return nil, fmt.Errorf("UpdateAddress api request: %w", err)
	}
	return models.UserFromProto(resp), nil
}
func (api *UsersAPI) AddAddresses(user *models.AddAddressRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	protoUser := models.AddAddressToProto(user)
	resp, err := api.UserServiceClient.AddAddresses(ctx, protoUser)
	if err != nil {
		return nil, fmt.Errorf("AddAddress api request: %w", err)
	}
	return models.UserFromProto(resp), nil
}

func (api *UsersAPI) GetUserRoles(token []byte) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	protoToken := &proto.TokenRequest{Token: token}
	resp, err := api.UserServiceClient.GetUserRoles(ctx, protoToken)
	if err != nil {
		return nil, fmt.Errorf("GetUserRoles api request: %w", err)
	}
	return resp.UserRoles, nil
}
func (api *UsersAPI) GetIsAdmin(token []byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	protoToken := &proto.TokenRequest{Token: token}
	resp, err := api.UserServiceClient.GetIsAdmin(ctx, protoToken)
	if err != nil {
		return false, fmt.Errorf("GetIsAdmin api request: %w", err)
	}
	return resp.IsAdmin, nil
}

// initConn initialize connection to Grpc servers
func (api *UsersAPI) initConn(addr string) (err error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	connParams := grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  100 * time.Millisecond,
			Multiplier: 1.2,
			MaxDelay:   1 * time.Second,
		},
		MinConnectTimeout: 5 * time.Second,
	})
	api.ClientConn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp), connParams)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	return
}

func (api *UsersAPI) ResetPassword(phone string) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	_, err := api.UserServiceClient.ResetPassword(ctx, &proto.ResetPasswordRequest{Phone: phone})
	if err != nil {
		return fmt.Errorf("resetPassword api request: %w", err)
	}

	return nil
}
func (api *UsersAPI) UpdatePassword(userUUID uuid.UUID, password []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	req := &proto.UpdatePasswordRequest{
		UserUuid: userUUID.Bytes(),
		Password: password,
	}

	_, err := api.UserServiceClient.UpdatePassword(ctx, req)
	if err != nil {
		return fmt.Errorf("updatePassword api request: %w", err)
	}
	return nil
}

func (api *UsersAPI) CheckAuth(token []byte) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	protoToken := &proto.TokenRequest{Token: token}
	resp, err := api.UserServiceClient.CheckAuth(ctx, protoToken)
	if err != nil {
		return nil, fmt.Errorf("checkAuth api request: %w", err)
	}

	return models.UserFromProto(resp), nil
}

// SignUpByEmail is
func (api *UsersAPI) SignUpByEmail(email string, password []byte) (*models.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	opts := &proto.SignUpRequest{
		LoginType: &proto.SignUpRequest_Email{Email: email},
		Password:  password,
	}
	resp, err := api.UserServiceClient.SignUp(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("signUp by email request has been failed: %w", err)
	}
	signUpResp := &models.SignUpResponse{UserUUID: uuid.FromBytesOrNil(resp.UserUuid)}
	return signUpResp, nil
}

// SignInByEmail is
func (api *UsersAPI) SignInByEmail(email string, password []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	opts := &proto.SignInRequest{
		LoginType: &proto.SignInRequest_Email{Email: email},
		Password:  password,
	}
	resp, err := api.UserServiceClient.SignIn(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("signIn by email api request: %w", err)
	}

	return resp.Token, nil
}

// SignUpByPhone is
func (api *UsersAPI) SignUpByPhone(phone string, password []byte) (*models.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	opts := &proto.SignUpRequest{
		LoginType: &proto.SignUpRequest_Phone{Phone: phone},
		Password:  password,
	}
	resp, err := api.UserServiceClient.SignUp(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("signUp by phone request has been failed: %w", err)
	}
	signUpResp := &models.SignUpResponse{UserUUID: uuid.FromBytesOrNil(resp.UserUuid)}

	return signUpResp, nil
}

func (api *UsersAPI) VerifyUserByPhone(phone, checkPhrase string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	opts := &proto.VerifyByPhoneRequest{
		Phone:       phone,
		CheckPhrase: checkPhrase,
	}
	resp, err := api.UserServiceClient.VerifyUserByPhone(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("verifying user by phone failed: %w", err)
	}

	return resp.Token, nil
}

func (api *UsersAPI) VerifyPasswordChange(phone, checkPhrase string, password []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	opts := &proto.VerifyPasswordChangeRequest{
		Phone:       phone,
		CheckPhrase: checkPhrase,
		Password:    password,
	}
	resp, err := api.UserServiceClient.VerifyPasswordChange(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("verifying password change failed: %w", err)
	}
	return resp.Token, nil
}

// SignInByPhone is
func (api *UsersAPI) SignInByPhone(phone string, password []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	opts := &proto.SignInRequest{
		LoginType: &proto.SignInRequest_Phone{Phone: phone},
		Password:  password,
	}
	resp, err := api.UserServiceClient.SignIn(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("signIn by phone api request: %w", err)
	}

	return resp.Token, nil
}

func (api *UsersAPI) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	api.mu.Lock()
	defer api.mu.Unlock()

	resp, err := api.HealthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "userapi"})
	if err != nil {
		return fmt.Errorf("healthcheck error: %w", err)
	}

	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("node is %s", errors.New("service is unhealthy"))
	}
	return nil
}

func (api *UsersAPI) UserByUUID(userUUID uuid.UUID) (*models.User, error) {
	opts := &proto.UserGetter{
		Getter: &proto.UserGetter_UserUuid{
			UserUuid: userUUID.Bytes(),
		},
	}
	return api.getUser(opts)
}
func (api *UsersAPI) UserByPhone(phone string) (*models.User, error) {
	opts := &proto.UserGetter{
		Getter: &proto.UserGetter_Phone{
			Phone: phone,
		},
	}
	return api.getUser(opts)
}
func (api *UsersAPI) UserByEmail(email string) (*models.User, error) {
	opts := &proto.UserGetter{
		Getter: &proto.UserGetter_Email{
			Email: email,
		},
	}
	return api.getUser(opts)
}

func (api *UsersAPI) getUser(opts *proto.UserGetter) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()

	resp, err := api.UserServiceClient.UserBy(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("userapi request failed: %w", err)
	}

	return models.UserFromProto(resp), nil
}
func (api *UsersAPI) DeleteAddress(addrUuid, userUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), api.timeout)
	defer cancel()
	req := &proto.DeleteAddressRequest{
		AddressUuid: addrUuid.Bytes(),
		UserUuid:    userUuid.Bytes(),
	}
	_, err := api.UserServiceClient.DeleteAddress(ctx, req)
	if err != nil {
		return fmt.Errorf("DeleteAddress api request: %w", err)
	}
	return nil
}
