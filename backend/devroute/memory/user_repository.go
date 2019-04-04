// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package memory

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/szabba/devroute/backend/devroute"
)

// A UserRepository "persists" users in memory.
type UserRepository struct {
	lock   sync.Mutex
	events map[devroute.UserID][]devroute.UserEvent
}

var _ devroute.UserRepository = &UserRepository{}

// NewUserRepository creates an in-memory repository of users.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		events: make(map[devroute.UserID][]devroute.UserEvent),
	}
}

// New tries to create a versioned user with no events applied.
func (repo *UserRepository) New() (*devroute.VersionedUser, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	rawID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	versioned := &devroute.VersionedUser{}
	versioned.ID = devroute.UserID(rawID)
	return versioned, nil
}

// Load tries to restore a user with the given ID.
func (repo *UserRepository) Load(id devroute.UserID) (*devroute.VersionedUser, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	savedEvts := repo.events[id]
	if len(savedEvts) == 0 {
		return nil, devroute.ErrNotFound
	}

	version := devroute.UserVersion(len(savedEvts))
	user := &devroute.User{}
	for _, evt := range savedEvts {
		evt.ApplyTo(user)
	}

	return &devroute.VersionedUser{User: *user, Version: version}, nil
}

// SaveEvents extends a user's history with the given events.
func (repo *UserRepository) SaveEvents(user *devroute.VersionedUser, evts ...devroute.UserEvent) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	savedEvts := repo.events[user.ID]
	if len(savedEvts) != int(user.Version) {
		return devroute.ErrConcurrentModification
	}

	repo.events[user.ID] = append(savedEvts, evts...)

	return nil
}
