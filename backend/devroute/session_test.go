// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devroute_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/szabba/assert"
	"github.com/szabba/devroute/backend/devroute"
)

func TestSessionTokenRetainsTheUserIDItWasCreatedWith(t *testing.T) {
	// given
	rawUserID, _ := uuid.NewV4()
	userID := devroute.UserID(rawUserID)
	issuedAt := time.Unix(0, 0)
	expiresAfter := time.Minute

	// when
	token := devroute.NewSessionToken(userID, issuedAt, expiresAfter)

	// then
	assert.That(token.ID() == userID, t.Fatalf, "got user ID %q, want %q", token.ID(), userID)
}

func TestSessionTokenHasExpectedValidityRange(t *testing.T) {
	// given
	rawUserID, _ := uuid.NewV4()
	userID := devroute.UserID(rawUserID)
	issuedAt := time.Unix(0, 0)
	expiresAfter := time.Minute

	lastInvalidBefore := issuedAt.Add(-time.Duration(1))
	firstValid := issuedAt
	middleValid := issuedAt.Add(expiresAfter / 2)
	lastValid := issuedAt.Add(expiresAfter)
	firstInvalidAfter := lastValid.Add(time.Duration(1))

	// when
	token := devroute.NewSessionToken(userID, issuedAt, expiresAfter)

	// then
	t.Logf("token %v", token)
	assert.That(!token.IsValidAt(lastInvalidBefore), t.Errorf, "should be invalid at %v", lastInvalidBefore)
	assert.That(token.IsValidAt(firstValid), t.Errorf, "should be valid at %v", firstValid)
	assert.That(token.IsValidAt(middleValid), t.Errorf, "should be valid at %v", middleValid)
	assert.That(token.IsValidAt(lastValid), t.Errorf, "should be valid at %v", lastValid)
	assert.That(!token.IsValidAt(firstInvalidAfter), t.Errorf, "\tshould be invalid at %v", firstInvalidAfter)
}
