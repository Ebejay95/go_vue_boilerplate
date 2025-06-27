package handlers

import (
	"context"
	"database/sql"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "myapp/proto/user"
)

type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func NewUserGRPCServer(db *sql.DB) *UserGRPCServer {
	return &UserGRPCServer{db: db}
}

func (s *UserGRPCServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, email, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to query users: %v", err)
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var id int32
		var name, email string
		var createdAt, updatedAt time.Time

		err := rows.Scan(&id, &name, &email, &createdAt, &updatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan user: %v", err)
		}

		user := &pb.User{
			Id:        id,
			Name:      name,
			Email:     email,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}
		users = append(users, user)
	}

	return &pb.GetUsersResponse{Users: users}, nil
}

func (s *UserGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	now := time.Now()
	var id int32

	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Name, req.Email, now, now,
	).Scan(&id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	user := &pb.User{
		Id:        id,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	return &pb.CreateUserResponse{User: user}, nil
}

func (s *UserGRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	now := time.Now()

	result, err := s.db.ExecContext(ctx,
		"UPDATE users SET name = $1, email = $2, updated_at = $3 WHERE id = $4",
		req.Name, req.Email, now, req.Id,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	var createdAt time.Time
	err = s.db.QueryRowContext(ctx,
		"SELECT created_at FROM users WHERE id = $1",
		req.Id,
	).Scan(&createdAt)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch updated user: %v", err)
	}

	user := &pb.User{
		Id:        req.Id,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(now),
	}

	return &pb.UpdateUserResponse{User: user}, nil
}

func (s *UserGRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	result, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}
