package main

//Imports
import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Main function - GUI
func main() {
	// Instantiate the Fyne Framework app and window
	a := app.New()
	w := a.NewWindow("TTS-STT")

	// Labels
	greet := widget.NewLabel("Welcome to the GO Text-to-Speech and Speech-to-Text App!")
	input_lbl := widget.NewLabel("Input the text to be said below!")

	// Entries
	input := widget.NewEntry()
	input.SetPlaceHolder("Write here ...")

	// Set the window content
	w.SetContent(container.NewVBox(
		greet,
		input_lbl,
		input,
		widget.NewButton("Send to TTS!", func() {
			var text_to_say = input.Text
			tts(text_to_say)
		}),
	))

	// Show the window and start the app event loop
	w.ShowAndRun()
}