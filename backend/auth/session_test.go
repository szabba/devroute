// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/szabba/assert"

	"github.com/szabba/devroute/backend/auth"
	"github.com/szabba/devroute/backend/auth/memory"
)

func init() {
	log.SetFlags(log.Flags() | log.Llongfile)
}

func TestSessionCanBeCreatedGivenValidCredentials(t *testing.T) {
	// given
	config := auth.SessionConfig{
		Issuer:   "issuer",
		ValidFor: 25 * time.Minute,
		Secret:   []byte("i kill giants"),
	}
	reserver := memory.NewUsernameMapping()
	repository := memory.NewUserRepository()

	userFactory := auth.NewUserFactory(repository, reserver)

	factory := auth.NewSessionFactory(reserver, repository, time.Now, config)

	ctx := context.Background()
	name, password := "user", "agent of asgard"

	_, err := userFactory.Create(ctx, name, password)
	assert.That(err == nil, t.Fatalf, "unexpected error creating user: %s", err)

	// when
	token, err := factory.Create(ctx, name, password)

	// then
	assert.That(token != "", t.Errorf, "expected non-empty token")
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
}

func TestSessionCannotBeCreatedForANonexistendUser(t *testing.T) {
	// given
	config := auth.SessionConfig{
		Issuer:   "issuer",
		ValidFor: 25 * time.Minute,
		Secret:   []byte("i kill giants"),
	}
	reserver := memory.NewUsernameMapping()
	repository := memory.NewUserRepository()

	factory := auth.NewSessionFactory(reserver, repository, time.Now, config)

	ctx := context.Background()
	name, password := "user", "agent of asgard"

	// when
	token, err := factory.Create(ctx, name, password)

	// then
	assert.That(token == "", t.Errorf, "expected token %s, got %s", "", token)
	assert.That(err == auth.ErrUserNotFound, t.Fatalf, "expected error %q, got %q", auth.ErrUserNotFound, err)
}

func TestSessionCannotBeCreatedUsingAnInvalidPasswords(t *testing.T) {
	// given
	config := auth.SessionConfig{
		Issuer:   "issuer",
		ValidFor: 25 * time.Minute,
		Secret:   []byte("i kill giants"),
	}
	reserver := memory.NewUsernameMapping()
	repository := memory.NewUserRepository()

	userFactory := auth.NewUserFactory(repository, reserver)

	factory := auth.NewSessionFactory(reserver, repository, time.Now, config)

	ctx := context.Background()
	name, password := "user", "agent of asgard"
	wrongPassword := "god of stories"

	_, err := userFactory.Create(ctx, name, password)
	assert.That(err == nil, t.Fatalf, "unexpected error creating user: %s", err)

	// when
	token, err := factory.Create(ctx, name, wrongPassword)

	// then
	assert.That(token == "", t.Errorf, "expected token %s, got %s", "", token)
	assert.That(err == auth.ErrInvalidCredentials, t.Fatalf, "expected error %q, got %q", auth.ErrInvalidCredentials, err)
}
