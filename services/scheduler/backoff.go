// Copyright 2016-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package scheduler

import (
	"math"
	"time"
)

// type BackoffFunc func(attempts int) time.Duration

// exponentialBackoff is the default backoff function. It performs
// exponential backoff.
func exponentialBackoff(attempts int) time.Duration {
	if attempts == 0 {
		return time.Duration(0)
	}
	return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond
}
