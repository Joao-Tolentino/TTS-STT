# Developer & Technical Documentation

This document provides a technical guide to the **TTS-STT** application's architecture and execution.

---

## System Architecture

The application combines the Fyne GUI loop with background goroutines for audio I/O using PortAudio and local command-line integrations.

```mermaid
graph TD
    UI[Fyne GUI] -->|Input Text| TTS[htgo-tts]
    TTS -->|Play Native| AudioOut
    
    UI -->|Start Recording| PortAudio[portaudio.Read]
    PortAudio -->|Stream into RAM| MutexBuffer[[]int Buffer]
    UI -->|Stop Recording| WAVEncoder[go-audio/wav]
    WAVEncoder -->|Save| TempWAV[temp/recording.wav]
    
    UI -->|Choose File| STT
    TempWAV --> STT
    
    STT[stt.go] -->|Convert 16kHz| FFMPEG[exec.Command ffmpeg]
    FFMPEG -->|Transcribe| Whisper[exec.Command whisper-cpp]
    Whisper -->|Read .txt| Output[Dialog pop-up]
```

---

## Directory Structure & File Roles

```
.
├── main.go             # Fyne GUI definitions, event listeners, and recording logic
├── stt.go              # ffmpeg format conversions and whisper-cpp execution wrapper
├── tts.go              # htgo-tts logic wrapper configured for Portuguese
├── models/             # Directory containing the whisper-cpp ggml weights
├── go.mod / go.sum     # Go dependency tracking
├── README.md           # General overview
└── Documentation.md    # Technical documentation
```

---

## Workflow

The execution flow of TTS-STT:
1. **Initialization**: Starts the `fyne.io/fyne/v2` app instance and creates layouts with inputs and buttons.
2. **Audio Capture**: Using `sync.Mutex` and goroutines, `portaudio` continuously reads microphone data into a global integer array while recording. Upon stopping, `go-audio/wav` flushes it to a disk file inside `temp/`.
3. **STT Inference**: The audio (whether picked manually or recorded) is pushed through `stt.go`. It shells out to `ffmpeg` to force 16000Hz mono PCM. It then shells out to `cmd.exe /c whisper-cpp -m models/ggml-tiny.bin` to transcribe the audio, retrieving the text and spawning a Fyne `dialog.ShowInformation` alert.
4. **TTS Generation**: Strings pass into `htgotts.Speech` with `voices.Portuguese`, compiling and natively playing the speech immediately.

---

## Launcher Compilation Guide

Because this application relies on CGO bindings for PortAudio and GLFW (Fyne), you must have a C compiler (like GCC/MinGW) installed on Windows.

### Compilation or Execution Commands

Execute the following commands in order within your terminal:

```powershell
# Restore modules
go mod tidy

# Build the GUI executable (hiding the command line window on Windows)
go build -ldflags="-H windowsgui" -o tts-stt.exe .

# Alternatively, run directly
go run .
```
