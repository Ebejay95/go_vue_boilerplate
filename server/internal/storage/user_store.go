package storage

import (
	"sync"
	"backend-grpc-server/internal/models"
)

// UserStore interface defines the contract for user storage operations
type UserStore interface {
	GetUser(id int32) (*models.User, bool)
	CreateUser(user *models.User) *models.User
	ListUsers() []*models.User
}

// InMemoryUserStore implements UserStore using in-memory storage
type InMemoryUserStore struct {
	users  map[int32]*models.User
	mutex  sync.RWMutex
	nextID int32
}

// NewInMemoryUserStore creates a new in-memory user store
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users:  make(map[int32]*models.User),
		nextID: 1,
	}
}

// GetUser retrieves a user by ID
func (s *InMemoryUserStore) GetUser(id int32) (*models.User, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

// CreateUser creates a new user
func (s *InMemoryUserStore) CreateUser(user *models.User) *models.User {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = s.nextID
	s.users[s.nextID] = user
	s.nextID++

	return user
}

// ListUsers returns all users
func (s *InMemoryUserStore) ListUsers() []*models.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return users
}