// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserFactory struct {
	repo     UserRepository
	reserver UsernameReserver
}

// A UsernameReserver manages a set of usernames.
// A valid implementation must ensure a username is not used more than once.
type UsernameReserver interface {
	ReserveUsername(ctx context.Context, id UserID, name string) error
	GetIDByName(ctx context.Context, name string) (*UserID, error)
}

func NewUserFactory(repo UserRepository, reserver UsernameReserver) *UserFactory {
	return &UserFactory{repo, reserver}
}

func (fact *UserFactory) Create(ctx context.Context, username, password string) (*User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), _HashStrength)
	versioned, err := fact.repo.New(ctx)
	if err != nil {
		return nil, err
	}

	err = fact.reserver.ReserveUsername(ctx, versioned.ID, username)
	if err != nil {
		log.Printf("cannot resevrve name %s: %s", username, err)
		return nil, ErrUsernameTaken
	}
	evt := &UserRegistered{
		Username:       username,
		HashedPassword: hashedPassword,
	}
	err = fact.repo.SaveEvents(ctx, versioned, evt)
	if err != nil {
		return nil, err
	}
	evt.ApplyTo(&versioned.User)
	return &versioned.User, nil
}
