// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package memory

import (
	"context"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/szabba/devroute/backend/auth"
)

// A UserRepository "persists" users in memory.
type UserRepository struct {
	lock   sync.Mutex
	events map[auth.UserID][]auth.UserEvent
}

var _ auth.UserRepository = &UserRepository{}

// NewUserRepository creates an in-memory repository of users.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		events: make(map[auth.UserID][]auth.UserEvent),
	}
}

// New tries to create a versioned user with no events applied.
func (repo *UserRepository) New(_ context.Context) (*auth.VersionedUser, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	rawID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	versioned := &auth.VersionedUser{}
	versioned.ID = auth.UserID(rawID)
	return versioned, nil
}

// Load tries to restore a user with the given ID.
func (repo *UserRepository) Load(_ context.Context, id auth.UserID) (*auth.VersionedUser, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	savedEvts := repo.events[id]
	if len(savedEvts) == 0 {
		return nil, auth.ErrUserNotFound
	}

	version := auth.UserVersion(len(savedEvts))
	user := &auth.User{}
	for _, evt := range savedEvts {
		evt.ApplyTo(user)
	}

	return &auth.VersionedUser{User: *user, Version: version}, nil
}

// SaveEvents extends a user's history with the given events.
func (repo *UserRepository) SaveEvents(_ context.Context, user *auth.VersionedUser, evts ...auth.UserEvent) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	savedEvts := repo.events[user.ID]
	if len(savedEvts) != int(user.Version) {
		return auth.ErrConcurrentModification
	}

	repo.events[user.ID] = append(savedEvts, evts...)

	return nil
}
