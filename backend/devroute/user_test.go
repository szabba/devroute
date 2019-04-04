// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devroute_test

import (
	"testing"

	"github.com/szabba/assert"

	"github.com/szabba/devroute/backend/devroute"
	"github.com/szabba/devroute/backend/devroute/memory"
)

func TestNewUserCanBeCreated(t *testing.T) {
	// given
	reserver := memory.NewUsernameReserver()
	repository := memory.NewUserRepository()
	service := devroute.NewUserService(reserver, repository)

	request := devroute.CreateUser{
		Username: "ulrich-ser",
		Password: []byte("shazam"),
	}

	// when
	user, err := service.Create(request)

	// then
	assert.That(user != nil, t.Fatalf, "user must not be nil")
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
	assert.That(user.Username == request.Username, t.Fatalf, "got %q, want %q", user.Username, request.Username)
	assert.That(user.IsPasswordCorrect(request.Password), t.Fatalf, "password %q should be valid", request.Password)
}

func TestNewlyCreatedUserCanBeRetrieved(t *testing.T) {
	// given
	reserver := memory.NewUsernameReserver()
	repository := memory.NewUserRepository()
	service := devroute.NewUserService(reserver, repository)

	created, err := service.Create(devroute.CreateUser{
		Username: "ulrich-ser",
		Password: []byte("shazam"),
	})
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)

	// when
	found, err := service.Get(created.ID)

	// then
	assert.That(found != nil, t.Fatalf, "found user must not be nil")
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)
}

func TestTwoUsersWithTheSameNameCannotBeCreated(t *testing.T) {
	// given
	reserver := memory.NewUsernameReserver()
	repository := memory.NewUserRepository()
	service := devroute.NewUserService(reserver, repository)

	request := devroute.CreateUser{
		Username: "ulrich-ser",
		Password: []byte("shazam"),
	}

	_, err := service.Create(request)
	assert.That(err == nil, t.Fatalf, "unexpected error: %s", err)

	// when
	second, err := service.Create(request)

	// then
	assert.That(second == nil, t.Fatalf, "there should be no second user: %#v", second)
	assert.That(err == devroute.ErrUsernameTaken, t.Fatalf, "got error %q, want %q", err, devroute.ErrUsernameTaken)
}
