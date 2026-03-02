package main

// Imports
import (
	htgotts "github.com/hegedustibor/htgo-tts"
	handlers "github.com/hegedustibor/htgo-tts/handlers"
	voices "github.com/hegedustibor/htgo-tts/voices"
)

// Text-to-Speech logic
func tts(text string) {
	// Speech configurations: output Folder, Language, Audio handler
	speech := htgotts.Speech{
		Folder: "audio", 
		Language: voices.Portuguese,
		Handler: &handlers.Native{},
	}
	
	// Attempts to Play the audio from the text
	if err := speech.Speak(text) 

	// Catches any possible error and passes to terminal
	err !=nil {
		panic(err)
	}
}