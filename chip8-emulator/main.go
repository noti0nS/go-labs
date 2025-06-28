package main

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chip8 struct {
	/*
	   Memory Map:
	   +---------------+= 0xFFF (4095) End of Chip-8 RAM
	   |               |
	   |               |
	   |               |
	   |               |
	   |               |
	   | 0x200 to 0xFFF|
	   |     Chip-8    |
	   | Program / Data|
	   |     Space     |
	   |               |
	   |               |
	   |               |
	   +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
	   |               |
	   |               |
	   |               |
	   +---------------+= 0x200 (512) Start of most Chip-8 programs
	   | 0x000 to 0x1FF|
	   | Reserved for  |
	   |  interpreter  |
	   +---------------+= 0x000 (0) Start of Chip-8 RAM
	*/
	Memory          [0xFFF]uint8
	Stack           [0xF]uint16
	Registers       [0xF]uint8
	AddressRegister uint16
	DelayRegister   uint8
	SoundRegister   uint8
	PC              uint16 // Program Counter
	SP              uint8  // Stack Pointer
}

const (
	displayWidth  = 64
	displayHeight = 32
	scale         = 20
)

var sprites [16][5]uint8 = [16][5]uint8{
	{0xF0, 0x90, 0x90, 0x90, 0xF0}, // 0
	{0x20, 0x60, 0x20, 0x20, 0x70}, // 1
	{0xF0, 0x10, 0xF0, 0x80, 0xF0}, // 2
	{0xF0, 0x10, 0xF0, 0x10, 0xF0}, // 3
	{0x90, 0x90, 0xF0, 0x10, 0x10}, // 4
	{0xF0, 0x80, 0xF0, 0x10, 0xF0}, // 5
	{0xF0, 0x80, 0xF0, 0x90, 0xF0}, // 6
	{0xF0, 0x10, 0x20, 0x40, 0x40}, // 7
	{0xF0, 0x90, 0xF0, 0x90, 0xF0}, // 8
	{0xF0, 0x90, 0xF0, 0x10, 0xF0}, // 9
	{0xF0, 0x90, 0xF0, 0x90, 0x90}, // A
	{0xE0, 0x90, 0xE0, 0x90, 0xE0}, // B
	{0xF0, 0x80, 0x80, 0x80, 0xF0}, // C
	{0xE0, 0x90, 0x90, 0x90, 0xE0}, // D
	{0xF0, 0x80, 0x80, 0x80, 0xF0}, // E
	{0xF0, 0x80, 0xF0, 0x80, 0x80}, // F
}

func main() {
	emu := &Chip8{}
	if err := loadRom(emu); err != nil {
		fmt.Println(err)
		return
	}

	rl.InitWindow(displayWidth*scale, displayHeight*scale, "CHIP-8")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Black)
		rl.DrawText("CHIP-8 Emulator", 32*scale, 16*scale, 12, rl.RayWhite)

		rl.EndDrawing()
	}

}

func loadRom(emu *Chip8) error {
	rom, err := os.ReadFile("roms/IBM.ch8")
	if err != nil {
		return fmt.Errorf("it wasn't possible open the rom's file due to the following error: %w", err)
	}
	for i, byte := range rom {
		emu.Memory[0x200+i] = byte
	}

	return nil
}
