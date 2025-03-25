# Mango - Music Anywhere

Mango is a modern, cross-platform music player built with Go and Wails.

## Features

- Clean, minimalist interface
- Album browser with cover art
- Album and track details
- Playback controls
- Keyboard shortcuts for controlling playback
- Support for FLAC, MP3, WAV and Ogg Vorbis audio files

## Music Directory Structure

Mango expects your music library to be organized in a specific way:

```
Music/
├── Album1/
│   ├── 01 - Track1.flac
│   ├── 02 - Track2.flac
│   ├── 03 - Track3.flac
│   └── folder.jpg
├── Album2/
│   ├── 01 - Track1.flac
│   ├── 02 - Track2.flac
│   └── folder.jpg
└── Album3/
    ├── 01 - Track1.flac
    ├── 02 - Track2.flac
    └── folder.jpg
```

**Important notes:**
- Each album should be in its own directory
- Album directories should contain FLAC, MP3, WAV or Ogg Vorbis files
- Include a `folder.jpg` file in each album directory for cover art

## Keyboard Shortcuts

- `Space` or `K`: Play/Pause
- `J`: Previous track
- `L`: Next track

## Building from Source

### Prerequisites

- Go 1.23 or later
- Node.js and npm
- Wails CLI v2.10.0 or later

### Install Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Build Instructions

#### Clone the repository

```bash
git clone https://github.com/snosek/mango.git
cd mango
```

#### For all platforms

```bash
wails build
```

#### For development

```bash
wails dev
```

### Platform-specific Instructions

#### macOS

The built application will be in the `build/bin` directory.

```bash
# Build for macOS
wails build -platform darwin
```

#### Linux

Dependencies:
- Required packages: `libgtk-3-dev`, `libwebkit2gtk-4.0-dev`, `libasound2-dev`

```bash
# Install dependencies (Ubuntu/Debian)
sudo apt install libgtk-3-dev libwebkit2gtk-4.0-dev libasound2-dev

# Build for Linux
wails build -platform linux
```

#### Windows

Dependencies:
- Windows 10 or later
- WebView2 runtime

```bash
# Build for Windows
wails build -platform windows
```

## Running the Application

1. After building, run the executable in the `build/bin` directory
2. Click "Add Music Directory" to select your music folder
3. Browse your albums and enjoy your music!

## Acknowledgments

- [Wails](https://wails.io) for the Go/JavaScript framework
- [Beep](https://github.com/gopxl/beep) for audio playback
- [Taglib](https://github.com/wtolson/go-taglib) for handling music metadata
