package storage

import (
	"testing"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresUserStore_CreateUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresUserStore(db)

	params := &models.CreateUserParams{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
		Role:  "user",
	}

	user, err := store.CreateUser(params)

	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, params.Name, user.Name)
	assert.Equal(t, params.Email, user.Email)
	assert.Equal(t, params.Age, user.Age)
	assert.Equal(t, params.Role, user.Role)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestPostgresUserStore_GetUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresUserStore(db)

	// Create a user first
	params := &models.CreateUserParams{
		Name:  "Jane Doe",
		Email: "jane@example.com",
		Age:   30,
		Role:  "admin",
	}

	createdUser, err := store.CreateUser(params)
	require.NoError(t, err)

	// Test getting the user
	user, exists := store.GetUser(createdUser.ID)

	assert.True(t, exists)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, createdUser.Name, user.Name)
	assert.Equal(t, createdUser.Email, user.Email)
}

func TestPostgresUserStore_GetUser_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresUserStore(db)

	user, exists := store.GetUser(99999)

	assert.False(t, exists)
	assert.Nil(t, user)
}

func TestPostgresUserStore_UpdateUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresUserStore(db)

	// Create a user first
	createParams := &models.CreateUserParams{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   25,
		Role:  "user",
	}

	createdUser, err := store.CreateUser(createParams)
	require.NoError(t, err)

	// Update the user
	updateParams := &models.UpdateUserParams{
		ID:    createdUser.ID,
		Name:  "John Smith",
		Email: "johnsmith@example.com",
		Age:   26,
		Role:  "admin",
	}

	updatedUser, err := store.UpdateUser(updateParams)

	require.NoError(t, err)
	assert.Equal(t, updateParams.ID, updatedUser.ID)
	assert.Equal(t, updateParams.Name, updatedUser.Name)
	assert.Equal(t, updateParams.Email, updatedUser.Email)
	assert.Equal(t, updateParams.Age, updatedUser.Age)
	assert.Equal(t, updateParams.Role, updatedUser.Role)
	assert.True(t, updatedUser.UpdatedAt.After(updatedUser.CreatedAt))
}

func TestPostgresUserStore_DeleteUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresUserStore(db)

	// Create a user first
	params := &models.CreateUserParams{
		Name:  "Delete Me",
		Email: "deleteme@example.com",
		Age:   25,
		Role:  "user",
	}

	createdUser, err := store.CreateUser(params)
	require.NoError(t, err)

	// Delete the user
	err = store.DeleteUser(createdUser.ID)
	require.NoError(t, err)

	// Verify user is deleted
	user, exists := store.GetUser(createdUser.ID)
	assert.False(t, exists)
	assert.Nil(t, user)
}

func TestPostgresUserStore_ListUsers(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresUserStore(db)

	// Create multiple users
	users := []models.CreateUserParams{
		{Name: "User 1", Email: "user1@example.com", Age: 25, Role: "user"},
		{Name: "User 2", Email: "user2@example.com", Age: 30, Role: "admin"},
		{Name: "User 3", Email: "user3@example.com", Age: 35, Role: "user"},
	}

	for _, params := range users {
		_, err := store.CreateUser(&params)
		require.NoError(t, err)
	}

	// Test listing users
	params := &models.ListUsersParams{
		Limit:  2,
		Offset: 0,
	}

	userList, total, err := store.ListUsers(params)

	require.NoError(t, err)
	assert.Equal(t, int32(3), total)
	assert.Len(t, userList, 2)
}
