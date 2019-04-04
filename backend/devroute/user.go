// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devroute

import (
	"errors"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const _HashStrength = bcrypt.DefaultCost

// ErrUsernameTaken indicates that a specified username is taken by a different, existing user.
var ErrUsernameTaken = errors.New("username taken")

// A UserID identifies a user of the system.
type UserID uuid.UUID

// A UserVersion is a version of a user.
// To fully identify a particular user state both the ID and version are necessary.
type UserVersion int

// A VersionedUser desrcibes a user with a
type VersionedUser struct {
	Version UserVersion
	User
}

// A User of the system.
type User struct {
	ID             UserID
	Username       string
	HashedPassword []byte
}

// IsPasswordCorrect checks whether the password passed in matches that of the user.
func (user *User) IsPasswordCorrect(password []byte) bool {
	return bcrypt.CompareHashAndPassword(user.HashedPassword, password) == nil
}

// A UserService manages users.
type UserService struct {
	reserver   UsernameReserver
	repository UserRepository
}

// NewUserService creates a new service to manage users.
func NewUserService(reserver UsernameReserver, repository UserRepository) *UserService {
	return &UserService{
		reserver:   reserver,
		repository: repository,
	}
}

// Create tries to create a new user.
func (svc *UserService) Create(request CreateUser) (*User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword(request.Password, _HashStrength)
	err := svc.reserver.ReserveUsername(request.Username)
	if err != nil {
		return nil, err
	}
	evt := &UserRegistered{
		Username:       request.Username,
		HashedPassword: hashedPassword,
	}
	versioned, err := svc.repository.New()
	if err != nil {
		return nil, err
	}
	err = svc.repository.SaveEvents(versioned, evt)
	if err != nil {
		return nil, err
	}
	evt.ApplyTo(&versioned.User)
	return &versioned.User, nil
}

// Get tries to retrieve an existing user.
func (svc *UserService) Get(id UserID) (*User, error) {
	versioned, err := svc.repository.Load(id)
	if err != nil {
		return nil, err
	}
	return &versioned.User, nil
}

// CreateUser is the parameter to (*UserService).Create.
type CreateUser struct {
	Username string
	Password []byte
}

// A UsernameReserver manages a set of usernames.
// A valid implementation must ensure a username is not used more than once.
type UsernameReserver interface {
	ReserveUsername(name string) error
}

// A UserRepository is responsible for the persistence of users.
type UserRepository interface {
	New() (*VersionedUser, error)
	Load(id UserID) (*VersionedUser, error)
	SaveEvents(user *VersionedUser, evts ...UserEvent) error
}

// A UserEvent is a fact about something that has happend to a user in the past.
type UserEvent interface {
	ApplyTo(user *User)
}

var _ = []UserEvent{&UserRegistered{}}

// UserRegistered is an event that represents a user registering within the system.
type UserRegistered struct {
	Username       string
	HashedPassword []byte
}

// ApplyTo implements UserEvent.
func (evt *UserRegistered) ApplyTo(user *User) {
	user.Username = evt.Username
	user.HashedPassword = evt.HashedPassword
}
