# TTS-STT (Text-To-Speech / Speech-To-Text)

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](LICENSE)
[![Platform: Windows](https://img.shields.io/badge/Platform-Windows-0078D6.svg?logo=windows&logoColor=white)](#)

A native GUI desktop application written in Go utilizing the Fyne framework. It provides offline voice-to-text inference utilizing `whisper-cpp` and dynamic text-to-speech rendering natively.

---

## Features

- **Fyne GUI**: A complete, cross-platform graphical user interface featuring file pickers, text inputs, and dynamic recording buttons.
- **Microphone Recording**: Directly captures microphone input into memory using `portaudio` and saves the buffer as a `.wav` file on-demand.
- **Offline Speech-To-Text (STT)**: Utilizes the `ggml-tiny.bin` model via the `whisper-cpp` command-line executable. Automatically converts inputs to 16kHz mono using `ffmpeg`.
- **Text-To-Speech (TTS)**: Translates written text inputs into spoken audio using the `htgo-tts` library, pre-configured for the Portuguese language.

---

## Quick Start

1. Clone or download the repository.
2. **Prerequisites**: Ensure you have installed Go, `ffmpeg` (available in system PATH), and the `whisper-cpp` CLI executable.
3. Fetch Go dependencies with `go mod tidy`.
4. Execute `go run .` to launch the GUI window.

---

## Configuration Details

This application depends on external executables (`ffmpeg` and `whisper-cpp`) existing in your system's PATH. It expects a Whisper model file to be present locally at `models/ggml-tiny.bin`.

---

## Usage Guidelines

- **TTS**: Type your text into the input field and click "Send to TTS!". It will generate the audio locally in the `audio` directory and automatically play it.
- **STT (File)**: Click "Choose File" to select an existing audio file. It will be converted by ffmpeg and transcribed by Whisper.
- **STT (Recording)**: Click "Start Recording" to use your microphone. When done, click "Stop Recording". The application saves it to the `temp` folder and immediately runs inference on it.

---

## Technical Documentation

For developers interested in directory structures, code architecture, or compilation guidelines, please refer to the **[Documentation.md](Documentation.md)** file.

---

## License

This project is licensed under the **GNU Affero General Public License Version 3 (AGPLv3)**. See the LICENSE file for details.
