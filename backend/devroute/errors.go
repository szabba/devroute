// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package devroute

import (
	"errors"
)

var (
	// ErrNotFound indicates that an entity could not be found.
	ErrNotFound = errors.New("not found")

	// ErrConcurrentModification indicates an entity was modified by another process.
	// The failing request might or might not be retryable.
	ErrConcurrentModification = errors.New("concurrent modification")
)
