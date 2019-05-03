// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignedToken struct {
	Raw    *jwt.Token
	Signed string
}

type TokenIssuer struct {
	name     string
	secret   string
	now      func() time.Time
	validFor time.Duration
}

func NewTokenIssuer(name, secret string, now func() time.Time, validFor time.Duration) *TokenIssuer {
	return &TokenIssuer{name, secret, now, validFor}
}

func (issuer *TokenIssuer) Issue(user User) (SignedToken, error) {
	claims := &jwt.StandardClaims{
		Issuer:    issuer.name,
		Subject:   user.ID.String(),
		ExpiresAt: issuer.now().Add(issuer.validFor).Unix(),
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := raw.SignedString(issuer.secret)
	return SignedToken{raw, signed}, nil
}
