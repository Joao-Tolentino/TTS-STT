package main

//Imports
import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Audio recording variables and global window for pop up
var (
	recording bool
	stream    *portaudio.Stream
	buffer    []int
	mutex     sync.Mutex
	audio_bttn *widget.Button

	appWindow fyne.Window  // Global window
)

// Start recording using the system microphone
func startRecording() {
	// Initialization
	mutex.Lock()
	defer mutex.Unlock()

	portaudio.Initialize()

	in := make([]int16, 64)

	s, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), &in)
	if err != nil {
		panic(err)
	}

	buffer = []int{}
	stream = s
	stream.Start()

	recording = true

	// Record the audio stream
	go func() {
		for recording {
			stream.Read()

			mutex.Lock()
			for _, v := range in {
				buffer = append(buffer, int(v))
			}
			mutex.Unlock()
		}
	}()
}

// End the recording and save the file
func stopRecording() {
    mutex.Lock()
    defer mutex.Unlock()

    recording = false
    stream.Stop()
    stream.Close()
    portaudio.Terminate()

    wavFile := saveWav() 

    // Run STT in background
    go func() {
        outDir := "transcriptions"
        os.MkdirAll(outDir, os.ModePerm)

        txtFile, transcription, err := stt("models/ggml-tiny.bin", wavFile, outDir)
        if err != nil {
            dialog.ShowError(err, appWindow)
            return
        }

    	dialog.ShowInformation("Recording Transcription", transcription, appWindow)
        fmt.Println("Saved transcription to:", txtFile)
    }()
}

// Save the recorded audio in a temp folder and send to Speech to Text
func saveWav() string {
	os.MkdirAll("temp", os.ModePerm)

	// Name the file using timestamp
	filename := fmt.Sprintf("temp/recording_%d.wav", time.Now().Unix())

	// create the recording file
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Save the audio
	enc := wav.NewEncoder(file, 44100, 16, 1, 1)

	intBuf := &audio.IntBuffer{
		Data:           buffer,
		Format:         &audio.Format{NumChannels: 1, SampleRate: 44100},
		SourceBitDepth: 16,
	}

	if err := enc.Write(intBuf); err != nil {
		panic(err)
	}

	enc.Close()

	return filename 
}


// Main function - GUI
func main() {
	// Instantiate the Fyne Framework app and window
	a := app.New()
	w := a.NewWindow("TTS-STT")

	appWindow = w // pass the window for the global variable

	// Labels
	greet := widget.NewLabel("Welcome to the GO Text-to-Speech and Speech-to-Text App!")
	input_lbl := widget.NewLabel("Input the text to be said below!")
	file_lbl := widget.NewLabel("Select audio file to transcribe.")
	filePathLabel := widget.NewLabel("No file selected")
	record_lbl := widget.NewLabel("Record your audio to transcribe!")

	// Entries
	input := widget.NewEntry()
	input.SetPlaceHolder("Write here ...")

	// Handle file selection for the audio file used in stt
	fileSelectedCallback := func(r fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if r == nil {
			filePathLabel.SetText("Selection cancelled")
			return
		}

		// Close the reader
		defer r.Close()

		// Update the label with the selected file's path
		filePath := r.URI().Path()
		filePathLabel.SetText(filePath)

		// Run STT in a goroutine to avoid blocking the GUI
    	go func() {
			outDir := "transcriptions"
			os.MkdirAll(outDir, os.ModePerm)

			txtFile, transcription, err := stt("models/ggml-tiny.bin", filePath, outDir)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// Show the result in a popup on the main thread
			dialog.ShowInformation("Transcription Result", transcription, w)

			fmt.Println("Saved transcription to:", txtFile)
		}()
	}

	// Buttons
	// The send for the text input
	input_bttn := widget.NewButton("Send to TTS!", func() {
						var text_to_say = input.Text
						tts(text_to_say)
					});

	// File operation
	file_bttn := widget.NewButton("Choose File", func() {
		dialog.ShowFileOpen(fileSelectedCallback, w)
	})

	// Start/Stop dor the audio recording
	tempButton := widget.NewButton("Start Recording", nil)
	audio_bttn = tempButton // pass to global

	// Make the button start and stop recording while tracking states with the text
	audio_bttn.OnTapped = func() {
		
		// Start the recording and change the text
		if !recording {
			startRecording()
			audio_bttn.SetText("Stop Recording")
		} else { // Stops recording and change the text
			stopRecording()
			audio_bttn.SetText("Start Recording")
		}
	}

	// Set the window content
	w.SetContent(container.NewVBox(
		greet,
		input_lbl,
		input,
		input_bttn,
		file_lbl,
		file_bttn,
		filePathLabel,
		record_lbl,
		audio_bttn,
	))

	// Show the window and start the app event loop
	w.ShowAndRun()
}
