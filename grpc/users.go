package grpc

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userspb "github.com/swarnakumar/go-identity/proto/users/v1"
)

func (s *server) GenerateToken(ctx context.Context, in *userspb.GenerateTokenRequest) (*userspb.Token, error) {
	email := in.GetEmail()
	pwd := in.GetPassword()
	if email == "" || pwd == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Email and Password are BOTH mandatory!!!")

	}

	ok := s.db.Users.VerifyPassword(ctx, email, pwd)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid Credentials")
	}

	jwtClaims := map[string]interface{}{"user": email}
	token := s.jwt.Create(jwtClaims)
	return &userspb.Token{Token: token}, nil

}

func (s *server) VerifyToken(ctx context.Context, data *userspb.Token) (*userspb.TokenValidityMessage, error) {
	token := data.GetToken()
	if token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Token not passed!")
	}

	t, err := s.jwt.Get(token)
	if err != nil || t == nil {
		return &userspb.TokenValidityMessage{Valid: false}, nil
	}

	user, present := t.PrivateClaims()["user"]

	if !present || user == nil || user == "" {
		return &userspb.TokenValidityMessage{Valid: false}, nil
	}

	_, err = s.db.Users.GetByEmail(ctx, user.(string))
	return &userspb.TokenValidityMessage{Valid: err == nil}, nil
}

func (s *server) checkValidUser(ctx context.Context) (bool, error) {
	token := s.IsValidToken(ctx)
	if !token.present {
		return false, status.Errorf(codes.Unauthenticated, "Authentication Token NOT PASSED.")
	}

	if !token.valid {
		return false, status.Errorf(codes.Unauthenticated, "Auth Token is Invalid.")
	}

	user := s.GetUser(ctx)

	if user == nil {
		return false, status.Errorf(codes.Unauthenticated, "Unknown User")
	}

	return true, nil

}

func (s *server) Me(ctx context.Context, _ *emptypb.Empty) (*userspb.UserDetails, error) {
	userValid, err := s.checkValidUser(ctx)
	if !userValid {
		return nil, err
	}
	user := s.GetUser(ctx)

	return &userspb.UserDetails{
		Email:     user.Email,
		IsActive:  user.IsActive,
		IsAdmin:   user.IsAdmin,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		LastLogin: timestamppb.New(user.LastLogin.Time),
		CreatedBy: user.CreatedBy.String,
	}, nil
}

func (s *server) RefreshToken(ctx context.Context, _ *emptypb.Empty) (*userspb.Token, error) {
	userValid, err := s.checkValidUser(ctx)
	if !userValid {
		return nil, err
	}

	user := s.GetUser(ctx)

	jwtClaims := map[string]interface{}{"user": user.Email}
	token := s.jwt.Create(jwtClaims)
	return &userspb.Token{Token: token}, nil
}

func (s *server) ChangePassword(ctx context.Context, in *userspb.ChangePasswordRequest) (*emptypb.Empty, error) {
	userValid, err := s.checkValidUser(ctx)
	if !userValid {
		return nil, err
	}
	user := s.GetUser(ctx)

	ok := s.db.Users.VerifyPassword(ctx, user.Email, in.CurrentPwd)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unknown User")
	}

	_, err = s.db.Users.ChangePassword(ctx, user.Email, in.NewPwd, &user.Email)
	if err == nil {
		return nil, nil
	}

	return nil, status.Errorf(codes.InvalidArgument, "Unable to change password: %s", err.Error())
}
