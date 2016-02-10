// Copyright 2015 Hajime Hoshi
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

// +build !js,!windows

package audio

import (
	"time"

	"golang.org/x/mobile/exp/audio"
)

var players = map[*audio.Player]struct{}{}

func play(src ReadSeekCloser, sampleRate int) error {
	// TODO: audio.NewPlayer interprets WAV header, which we don't want.
	// Use OpenAL or native API instead.
	p, err := audio.NewPlayer(src, audio.Stereo16, int64(sampleRate))
	if err != nil {
		return err
	}
	players[p] = struct{}{}
	return p.Play()
}

func initialize() {
	audioEnabled = true
	go func() {
		for {
			deleted := []*audio.Player{}
			for p, _ := range players {
				if p.State() == audio.Stopped {
					p.Close()
					deleted = append(deleted, p)
				}
			}
			for _, p := range deleted {
				delete(players, p)
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
}