package main

import (
	"bytes"
	_ "embed"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

//go:embed resources/oam-theme.wav
var introThemeData []byte

//go:embed resources/sweet_puppy_paws.wav
var sweetPuppyPawsData []byte

const (
	sampleRate         = 44100
	introFadeStartTime = 42 * time.Second // Start fading after 42 seconds
)

var (
	audioContext   *audio.Context
	introPlayer    *audio.Player
	introVolume    float64 = 1.0
	introFading    bool
	introFadeSpeed float64 = 0.0033 // Fade out over ~5 seconds (at 60 FPS)
)

// InitAudio initializes the audio context
func InitAudio() {
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}
}

// PlayIntroTheme starts playing the intro theme music
func PlayIntroTheme() {
	if audioContext == nil {
		InitAudio()
	}

	// Don't restart if already playing
	if introPlayer != nil && introPlayer.IsPlaying() {
		return
	}

	// Decode WAV (16-bit PCM)
	stream, err := wav.Decode(audioContext, bytes.NewReader(introThemeData))
	if err != nil {
		log.Printf("Failed to decode intro theme: %v", err)
		return
	}

	// Create player
	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		log.Printf("Failed to create audio player: %v", err)
		return
	}

	introPlayer = player
	introVolume = 1.0
	introFading = false
	introPlayer.SetVolume(introVolume)
	introPlayer.Play()
	log.Println("Intro theme started playing")
}

// StartIntroFadeOut begins fading out the intro theme
func StartIntroFadeOut() {
	if introPlayer != nil && introPlayer.IsPlaying() {
		introFading = true
		log.Println("Starting intro theme fade out")
	}
}

// UpdateIntroAudio should be called every frame to handle fade effects
func UpdateIntroAudio() {
	if introPlayer == nil {
		return
	}

	// Start fading after 42 seconds of playback
	if !introFading && introPlayer.IsPlaying() {
		position := introPlayer.Position()
		if position >= introFadeStartTime {
			StartIntroFadeOut()
		}
	}

	if introFading && introVolume > 0 {
		introVolume -= introFadeSpeed
		if introVolume <= 0 {
			introVolume = 0
			introPlayer.Pause()
			introFading = false
			log.Println("Intro theme faded out completely")
		}
		introPlayer.SetVolume(introVolume)
	}
}

// StopIntroTheme immediately stops the intro theme
func StopIntroTheme() {
	if introPlayer != nil {
		introPlayer.Pause()
		introPlayer = nil
		introFading = false
		introVolume = 1.0
	}
}

// IsIntroThemePlaying returns true if the intro theme is currently playing
func IsIntroThemePlaying() bool {
	return introPlayer != nil && introPlayer.IsPlaying()
}

// PlaySweetPuppyPawsSound plays the sweet puppy paws audio file
func PlaySweetPuppyPawsSound() {
	if audioContext == nil {
		InitAudio()
	}

	stream, err := wav.Decode(audioContext, bytes.NewReader(sweetPuppyPawsData))
	if err != nil {
		log.Printf("Failed to decode sweet_puppy_paws.wav: %v", err)
		return
	}

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		log.Printf("Failed to create player for sweet_puppy_paws.wav: %v", err)
		return
	}

	player.Play()
}
