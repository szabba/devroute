// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SessionFactory struct {
	reserver UsernameReserver
	repo     UserRepository
	now      func() time.Time
	config   SessionConfig
}

type SessionConfig struct {
	Issuer   string
	ValidFor time.Duration
	Secret   []byte
}

func NewSessionFactory(
	reserver UsernameReserver,
	repo UserRepository,
	now func() time.Time,
	config SessionConfig,
) *SessionFactory {
	return &SessionFactory{
		reserver: reserver,
		repo:     repo,
		now:      now,
		config:   config,
	}
}

func (fact *SessionFactory) Create(ctx context.Context, username, password string) (string, error) {
	user, err := fact.validate(ctx, username, password)
	if err != nil {
		return "", err
	}
	return fact.create(user)
}

func (fact *SessionFactory) validate(ctx context.Context, username, password string) (*User, error) {
	id, err := fact.reserver.GetIDByName(ctx, username)
	if err != nil {
		log.Printf("could not find ID for user named %s", username)
		return nil, err
	}

	versioned, err := fact.repo.Load(ctx, *id)
	if err != nil {
		log.Printf("failed to load user for ID %d", *id)
		return nil, err
	}

	if !versioned.IsPasswordCorrect(password) {
		log.Printf("invalid credentials for user named %s", username)
		return nil, ErrInvalidCredentials
	}

	return &versioned.User, nil
}

func (fact *SessionFactory) create(user *User) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    fact.config.Issuer,
		Subject:   user.ID.String(),
		ExpiresAt: fact.now().Add(fact.config.ValidFor).Unix(),
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := raw.SignedString(fact.config.Secret)
	if err != nil {
		log.Printf("failed signing token claims %#v (secret has length of %d)", claims, len(fact.config.Secret))
		return "", err
	}
	return signed, nil
}
