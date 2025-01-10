package driver

const (
	// MimeTypeImage output base64 mine-type.
	MimeTypeImage = "image/png"

	// MimeTypeAudio output base64 mine-type.
	MimeTypeAudio = "audio/wav"

	sampleRate = 8000 // Hz
	
	imageStringDpi = 72.0

	// OptionShowHollowLine shows hollow line
	OptionShowHollowLine = 2
	// OptionShowSlimeLine shows slime line
	OptionShowSlimeLine = 4
	// OptionShowSineLine shows sine line
	OptionShowSineLine = 8
)

var endingBeepSound []byte

func init() {
	endingBeepSound = changeSpeed(beepSound, 1.4)
}
