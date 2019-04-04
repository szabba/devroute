// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package memory

import (
	"sync"

	"github.com/szabba/devroute/backend/devroute"
)

// A UsernameReserver manages username reservations in memory.
type UsernameReserver struct {
	lock          sync.Mutex
	reservedNames map[string]bool
}

var _ devroute.UsernameReserver = &UsernameReserver{}

// NewUsernameReserver creates an in-memory record of memory reservations.
func NewUsernameReserver() *UsernameReserver {
	return &UsernameReserver{
		reservedNames: map[string]bool{},
	}
}

// ReserveUsername tries to reserve a username.
// It will fail if the username was reserved before.
func (reserver *UsernameReserver) ReserveUsername(name string) error {
	if reserver.reservedNames[name] {
		return devroute.ErrUsernameTaken
	}
	reserver.reservedNames[name] = true
	return nil
}
