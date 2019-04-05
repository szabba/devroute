// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devroute

import (
	"fmt"
	"time"
)

type SessionToken struct {
	id           UserID
	issuedAt     time.Time
	expiresAfter time.Duration
}

func NewSessionToken(id UserID, issuedAt time.Time, expiresAfter time.Duration) SessionToken {
	return SessionToken{
		id:           id,
		issuedAt:     issuedAt.UTC(),
		expiresAfter: expiresAfter,
	}
}

func (token SessionToken) String() string {
	return fmt.Sprintf(
		"%x, valid from %v to %v",
		token.id,
		token.issuedAt,
		token.issuedAt.Add(token.expiresAfter))
}

func (token SessionToken) ID() UserID { return token.id }

func (token SessionToken) IsValidAt(t time.Time) bool {
	min := token.issuedAt
	max := token.issuedAt.Add(token.expiresAfter)
	notBefore := min.Before(t) || min.Equal(t)
	notAfter := max.After(t) || max.Equal(t)
	return notBefore && notAfter
}
