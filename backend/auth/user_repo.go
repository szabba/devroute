// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import "context"

// A UserRepository is responsible for the persistence of users.
type UserRepository interface {
	New(ctx context.Context) (*VersionedUser, error)
	Load(ctx context.Context, id UserID) (*VersionedUser, error)
	SaveEvents(ctx context.Context, user *VersionedUser, evts ...UserEvent) error
}
