package service

import (
	"context"
	"owner-service/internal/core/port"
	"owner-service/internal/core/service/user"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var _ user.UserServer = (*grpcServer)(nil)

type Config struct {
}

type grpcServer struct {
	user.UnimplementedUserServer
	config   *Config
	userRepo port.UserRepository
}

func NewGRPCServer(config *Config, userRepo port.UserRepository, opts ...grpc.ServerOption) (*grpc.Server, error) {

	logger := zap.L().Named("grpc_server")
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(
			func(duration time.Duration) zapcore.Field {
				return zap.Int64(
					"grpc.time_ns",
					duration.Nanoseconds(),
				)
			},
		),
	}

	opts = append(opts,
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, zapOpts...),
		)),
	)

	gsrv := grpc.NewServer(opts...)
	srv, err := newgrpcServer(config, userRepo)
	if err != nil {
		return nil, err
	}

	user.RegisterUserServer(gsrv, srv)
	return gsrv, nil
}

func newgrpcServer(config *Config, userRepo port.UserRepository) (srv *grpcServer, err error) {
	srv = &grpcServer{
		config:   config,
		userRepo: userRepo,
	}
	return srv, nil
}

func (s *grpcServer) Get(ctx context.Context, req *user.UserRequest) (*user.UserResponse, error) {
	// Extract the logger from the context
	logger := zap.L().Named("grpc_server")
	logger.Info("Received Get request", zap.String("user_id", req.UserId))

	// TODO: authorize user through a mechanism
	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		logger.Error("Failed to parse user ID", zap.Error(err))
		return nil, err
	}

	// Get user from repository
	retUser, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		logger.Error("Failed to get user", zap.Error(err))
		return nil, err
	}

	usr := user.UserResponse{
		Id:        retUser.ID.Hex(),
		FirstName: retUser.FirstName,
		LastName:  retUser.LastName,
		Email:     retUser.Email,
		Password:  retUser.Password,
		Phone:     retUser.Phone,
		IsActive:  retUser.IsActive,
		CreatedAt: retUser.CreatedAt.String(),
		UpdatedAt: retUser.UpdatedAt.String(),
	}

	if role, err := user.StringToUserRole(retUser.Role.String()); err == nil {
		usr.Role = role
	}

	return &usr, nil
}
