package main

import (
	"fmt"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	"log"
	"os"
	"path/filepath"
	"time"
)

func decodeFile(file *os.File) (beep.StreamSeekCloser, beep.Format, error) {
	fileType := filepath.Ext(file.Name())
	switch fileType {
	case ".mp3":
		return mp3.Decode(file)
	case ".wav":
		return wav.Decode(file)
	case ".flac":
		return flac.Decode(file)
	default:
		return nil, beep.Format{}, fmt.Errorf("Unknown file format: %v", fileType)
	}
}

func PlayFile(file *os.File) {

	streamer, format, err := decodeFile(file)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	sr := format.SampleRate * 10
	sampleRate := beep.Resample(1, format.SampleRate, sr, streamer)
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(sampleRate, beep.Callback(func() {
		done <- true
	})))
	<-done

}
