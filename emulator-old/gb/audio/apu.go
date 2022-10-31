package apu

import (
	"math"
)

// TODO: Audio APU

const (
	sampleRate = 44100
	twoPi      = 2 * math.Pi
	perSample  = 1 / float64(sampleRate)

	cpuTicksPerSample    = float64(4194304) / sampleRate
	maxFrameBufferLength = 5000
)

// APU is the GameBoy's audio processing unit. Audio comprises four
// channels, each one controlled by a set of registers.
//
// Channels 1 and 2 are both Square channels, channel 3 is an arbitrary
// waveform channel which can be set in RAM, and channel 4 outputs noise.
type APU struct {
	playing bool

	memory      [52]byte
	waveformRam []byte

	chn1, chn2, chn3, chn4 *Channel
	tickCounter            float64
	lVol, rVol             float64

	audioBuffer chan [2]byte
}

// Init the sound emulation for a Gameboy.
func (a *APU) Init(sound bool) {
}

// Starts a goroutine which plays the sound
func (a *APU) playSound(bufferSeconds int) {

}

func (a *APU) Buffer(cpuTicks int, speed int) {

}

var soundMask = []byte{
	/* 0xFF10 */ 0xFF, 0xC0, 0xFF, 0x00, 0x40,
	/* 0xFF15 */ 0x00, 0xC0, 0xFF, 0x00, 0x40,
	/* 0xFF1A */ 0x80, 0x00, 0x60, 0x00, 0x40,
	/* 0xFF20 */ 0x00, 0x3F, 0xFF, 0xFF, 0x40,
	/* 0xFF24 */ 0xFF, 0xFF, 0x80,
}

var channel3Volume = map[byte]float64{0: 0, 1: 1, 2: 0.5, 3: 0.25}

var squareLimits = map[byte]float64{
	0: -0.25, // 12.5% ( _-------_-------_------- )
	1: -0.5,  // 25%   ( __------__------__------ )
	2: 0,     // 50%   ( ____----____----____---- ) (normal)
	3: 0.5,   // 75%   ( ______--______--______-- )
}

// Read returns a value from the APU.
func (a *APU) Read(address uint16) byte {
	return 0x00
}

// Write a value to the APU registers.
func (a *APU) Write(address uint16, value byte) {

}

// WriteWaveform writes a value to the waveform ram.
func (a *APU) WriteWaveform(address uint16, value byte) {

}

// ToggleSoundChannel toggles a sound channel for debugging.
func (a *APU) ToggleSoundChannel(channel int) {

}

func (a *APU) LogSoundState() {

}

// Extract some envelope variables from a byte.
func (a *APU) extractEnvelope(val byte) (volume, direction, sweep byte) {
	return
}
