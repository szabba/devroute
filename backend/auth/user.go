// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import (
	"errors"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const _HashStrength = bcrypt.DefaultCost

var (
	// ErrUsernameTaken indicates that a specified username is taken by a different, existing user.
	ErrUsernameTaken = errors.New("username taken")

	ErrUserNotFound = errors.New("user not found")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrConcurrentModification = errors.New("concurrent modification")
)

// A UserID identifies a user of the system.
type UserID uuid.UUID

func ID(raw string) (UserID, error) {
	id, err := uuid.FromString(raw)
	return UserID(id), err
}

func (id UserID) String() string { return uuid.UUID(id).String() }

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
func (user *User) IsPasswordCorrect(password string) bool {
	return bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)) == nil
}
