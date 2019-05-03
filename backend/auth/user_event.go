// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

// A UserEvent is a fact about something that has happend to a user in the past.
type UserEvent interface {
	ApplyTo(user *User)
}

var _ = []UserEvent{&UserRegistered{}}

// UserRegistered is an event that represents a user registering within the system.
type UserRegistered struct {
	Username       string
	HashedPassword []byte
}

// ApplyTo implements UserEvent.
func (evt *UserRegistered) ApplyTo(user *User) {
	user.Username = evt.Username
	user.HashedPassword = evt.HashedPassword
}
