// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tokenverify

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zsjinwei/openim-tools/errs"
)

const HoursOneDay = 24
const minutesBefore = 5

type Claims struct {
	UserID     string
	PlatformID int // login platform
	jwt.RegisteredClaims
}

func BuildClaims(uid string, platformID int, ttl int64) Claims {
	now := time.Now()
	before := now.Add(-time.Minute * time.Duration(minutesBefore))
	return Claims{
		UserID:     uid,
		PlatformID: platformID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(ttl*HoursOneDay) * time.Hour)), // Expiration time
			IssuedAt:  jwt.NewNumericDate(now),                                                 // Issuing time
			NotBefore: jwt.NewNumericDate(before),                                              // Begin Effective time
		},
	}
}

func GetClaimFromToken(tokensString string, secretFunc jwt.Keyfunc) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokensString, &Claims{}, secretFunc)
	if err == nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
		return nil, errs.ErrTokenUnknown
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		return nil, mapValidationError(ve)
	}

	return nil, errs.ErrTokenUnknown
}

func mapValidationError(ve *jwt.ValidationError) error {
	if ve.Errors&jwt.ValidationErrorMalformed != 0 {
		return errs.ErrTokenMalformed
	} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
		return errs.ErrTokenExpired
	} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
		return errs.ErrTokenNotValidYet
	}
	return errs.ErrTokenUnknown
}
