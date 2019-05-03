// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth_test

import (
	"context"
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/devroute/backend/auth"
	"github.com/szabba/devroute/backend/auth/memory"
)

func TestNewUserCanBeCreated(t *testing.T) {
	// given
	reserver := memory.NewUsernameMapping()
	repository := memory.NewUserRepository()
	factory := auth.NewUserFactory(repository, reserver)

	ctx := context.Background()
	username := "ulrich-ser"
	password := "shazam"

	// when
	user, err := factory.Create(ctx, username, password)

	// then
	assert.That(user != nil, t.Fatalf, "user must not be nil")
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
	assert.That(user.Username == username, t.Fatalf, "got %q, want %q", user.Username, username)
	assert.That(user.IsPasswordCorrect(password), t.Fatalf, "password %q should be valid", password)
}

func TestNewlyCreatedUserCanBeRetrieved(t *testing.T) {
	// given
	reserver := memory.NewUsernameMapping()
	repository := memory.NewUserRepository()
	factory := auth.NewUserFactory(repository, reserver)

	ctx := context.Background()
	created, err := factory.Create(ctx, "ulrich-ser", "shazam")
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)

	// when
	found, err := repository.Load(ctx, created.ID)

	// then
	assert.That(found != nil, t.Fatalf, "found user must not be nil")
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
}

func TestTwoUsersWithTheSameNameCannotBeCreated(t *testing.T) {
	// given
	reserver := memory.NewUsernameMapping()
	repository := memory.NewUserRepository()
	factory := auth.NewUserFactory(repository, reserver)

	ctx := context.Background()
	_, err := factory.Create(ctx, "ulrich-ser", "shazam")
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)

	// when
	second, err := factory.Create(ctx, "ulrich-ser", "shazam")

	// then
	assert.That(second == nil, t.Fatalf, "there should be no second user: %#v", second)
	assert.That(err == auth.ErrUsernameTaken, t.Fatalf, "got error %q, want %q", err, auth.ErrUsernameTaken)
}
