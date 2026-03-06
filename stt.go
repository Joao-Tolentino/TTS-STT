package main

// Imports
import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Check the filetype, if is .wav whisper accepts as us
func IsWav(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".wav"
}

// Conversion for the audio files, whisper accepts only .wav
func ConvertToWhisperWav(inputPath string, outputPath string) error {

	// Use ffmpeg to make the conversion
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-ar", "16000", // sample rate
		"-ac", "1",     // mono
		"-c:a", "pcm_s16le",
		"-y",
		outputPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v | %s", err, string(out))
	}

	return nil
}

// Make the file check and converts if needed
func PrepareAudioForWhisper(input string) (string, error) {

	if IsWav(input) {
		return input, nil
	}

	// change the name of the file to .wav
	output := filepath.Join(
		filepath.Dir(input),
		filepath.Base(input)+".wav",
	)

	// apply the conversion to wav
	err := ConvertToWhisperWav(input, output)
	if err != nil {
		return "", err
	}

	return output, nil
}

// Read WAV PCM16 to float32
func ReadWavSamples(path string) ([]float32, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    // skip WAV header
    header := make([]byte, 44)
    if _, err := f.Read(header); err != nil {
        return nil, err
    }

	// Read the file
    data, err := io.ReadAll(f)
    if err != nil {
        return nil, err
    }

	// take samples from the audio for processing
    samples := make([]float32, len(data)/2)
    for i := 0; i < len(samples); i++ {
        val := int16(binary.LittleEndian.Uint16(data[i*2 : i*2+2]))
        samples[i] = float32(val) / 32768.0
    }

	//return the samples
    return samples, nil
}

// Speech-to-Text
func stt(modelPath string, audioPath string, outputDir string) (string, string, error) {
	// Prepare audio (convert if needed)
	audio, err := PrepareAudioForWhisper(audioPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to prepare audio: %w", err)
	}

	// Create output filename
	base := filepath.Base(audio)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	timestamp := time.Now().Unix()
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%s_%d.txt", name, timestamp))
	
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// arguments for the CLI, auto detection of the language is activated
	args := []string{
		"/c",
		"whisper-cpp",
		"-m", modelPath,
		audio,
		"-f", outputFile,
		"--output-txt",
		"-l", 
		"auto",
	}
	
	// Create command with cmd.exe
	cmd := exec.Command("cmd.exe", args...)
	
	// Run command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("whisper failed: %v\nOutput: %s", err, string(output))
	}

	// The CLI creates a .txt file with the same name as the audio file
	generatedTxt := strings.TrimSuffix(audio, ".wav") + ".txt"
	
	// Check if the file was created
	if _, err := os.Stat(generatedTxt); os.IsNotExist(err) {
		// Try to find any .txt file that might have been created
		files, err := filepath.Glob(filepath.Join(filepath.Dir(audio), "*.txt"))
		if err != nil || len(files) == 0 {
			return "", "", fmt.Errorf("no transcription file generated")
		}
		// Use the most recent .txt file
		generatedTxt = files[len(files)-1]
	}
	
	// Read the transcription
	transcriptionBytes, err := os.ReadFile(generatedTxt)
	if err != nil {
		return "", "", fmt.Errorf("could not read transcription file: %v", err)
	}
	
	transcription := strings.TrimSpace(string(transcriptionBytes))
	
	// Move/copy to output directory
	err = os.Rename(generatedTxt, outputFile)
	if err != nil {
		// If move fails, copy content to new file
		err = os.WriteFile(outputFile, transcriptionBytes, 0644)
		if err != nil {
			return "", "", fmt.Errorf("failed to save transcription: %w", err)
		}
		// Try to delete the original
		os.Remove(generatedTxt)
	}

	return outputFile, transcription, nil
}