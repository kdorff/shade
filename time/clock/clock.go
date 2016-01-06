// Copyright 2016 Richard Hawkins
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
// Package clock TODO doc

package clock

import "time"

type Clock struct {
	Timer *time.Timer
}

func New() (*Clock, error) {
	c := Clock{
		Timer: time.NewTimer(time.Millisecond),
	}

	return &c, nil
}

// Clock TODO doc
func (c *Clock) Tick(limit int) {
	<-c.Timer.C
	// TODO: this is wrong, it should only restrict that the func is run at most _limit_ times a second.
	// But for now I think it will work.
	c.Timer.Reset(time.Millisecond / (1000 / time.Duration(limit)))
}
