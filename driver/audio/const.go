package audio

const (
	// MimeTypeAudio output base64 mine-type.
	MimeTypeAudio = "audio/wav"

	sampleRate = 8000 // Hz
)

var endingBeepSound []byte

func init() {
	endingBeepSound = changeSpeed(beepSound, 1.4)
}
