// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package memory

import (
	"context"
	"sync"

	"github.com/szabba/devroute/backend/auth"
)

// A UsernameMapping manages username reservations in memory.
type UsernameMapping struct {
	lock          sync.Mutex
	reservedNames map[string]*auth.UserID
}

var _ auth.UsernameReserver = &UsernameMapping{}

// NewUsernameMapping creates an in-memory record of memory reservations.
func NewUsernameMapping() *UsernameMapping {
	return &UsernameMapping{
		reservedNames: map[string]*auth.UserID{},
	}
}

// ReserveUsername tries to reserve a username.
// It will fail if the username was reserved before.
func (reserver *UsernameMapping) ReserveUsername(_ context.Context, id auth.UserID, name string) error {
	reserver.lock.Lock()
	defer reserver.lock.Unlock()

	if reserver.reservedNames[name] != nil {
		return auth.ErrUsernameTaken
	}
	reserver.reservedNames[name] = &id
	return nil
}

func (reserver *UsernameMapping) GetIDByName(_ context.Context, name string) (*auth.UserID, error) {
	reserver.lock.Lock()
	defer reserver.lock.Unlock()

	id := reserver.reservedNames[name]
	if id == nil {
		return nil, auth.ErrUserNotFound
	}
	return id, nil
}
