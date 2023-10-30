# GoGB

GoBoy is a Nintendo GameBoy emulator written in go.
This emulator was built as a development exercise.
It is still a work in progress.

## Dependencies

```
github.com/veandco/go-sdl2/sdl
github.com/veandco/go-sdl2/img
github.com/veandco/go-sdl2/ttf
```

These dependencies require sdl2 library to be installed, you can install it in your system following this guide:

https://github.com/veandco/go-sdl2#requirements

## Installation

```
go get github.com/escrichov/gogb
```

## Usage

```
gogb [rom_file]
```

## Controls

| Gameboy                  | Emulator   |
|--------------------------|------------|
| Up, Down, Left, Right    | Arrow Keys |
| Start                    | Enter      |
| Select                   | Tab        |
| A                        | z          |
| B                        | x          |

### Other emulator functionality

| Emulator | Functionality              |
|----------|----------------------------|
| Escape   | Exit                       |
| v        | V-Sync on/off              |
| t        | Show fps                   |
| f        | Full screen                |
| k        | Take snapshot snapshot.png |
| s        | Save state to save.bess    |
| r        | Reset                      |
| p        | Pause                      |
