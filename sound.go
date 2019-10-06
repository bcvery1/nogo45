package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var (
	mainTrackStreamer beep.Streamer

	sounds = make(map[Track]*beep.Buffer)

	hasErrored = false
)

type Track string

const (
	hurtSound        Track = "hurt"
	attackSound      Track = "attack1"
	coinPickupSound  Track = "coinpickup"
	deniedSound      Track = "denied"
	projectile1Sound Track = "projectile1"
	projectile2Sound Track = "projectile2"
)

func SetupAudio() {
	mainTrack, err := os.Open(filepath.Join(binPath, "assets/audio/mainTrack.wav"))
	if err != nil {
		panic(err)
		hasErrored = true
		return
	}

	streamer, format, err := wav.Decode(mainTrack)
	if err != nil {
		panic(err)
		hasErrored = true
		return
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		panic(err)
		hasErrored = true
		return
	}

	mainTrackStreamer = &effects.Volume{
		Streamer: beep.Loop(-1, streamer),
		Base:     2,
		Volume:   -3,
	}
	speaker.Play(mainTrackStreamer)

	loadSound(hurtSound)
	loadSound(attackSound)
	loadSound(coinPickupSound)
	loadSound(deniedSound)
	loadSound(projectile1Sound)
	loadSound(projectile2Sound)
}

func loadSound(sound Track) {
	filename := filepath.Join(binPath, "assets", "audio", fmt.Sprintf("%s.wav", sound))
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
		return
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		panic(err)
		return
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	_ = streamer.Close()

	sounds[sound] = buffer
}

func PlaySound(sound Track) {
	if hasErrored {
		return
	}

	b, ok := sounds[sound]
	if !ok {
		fmt.Printf("not found %v\n", sound)
		return
	}

	speaker.Play(b.Streamer(0, b.Len()))
}
