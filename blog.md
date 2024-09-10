# Implementing CHIP-8 emulator in Go
## Preface:
Needed a decent sized project to be done in 3-4 days. Not too trivial or small, no web project, not so long that takes a week or more, which is more or less keyword for always in progress.

### Why emulator:
As close to hardware as possible, gives insight into hardware details and some tricks/quirks used in software. Allows better understading of software behaviour.

### Why CHIP-8:
Perfect for 2-3 days project. Not complicated like NES and subsequent systems, yet gives full taste of emulated systems. Also, a great gateway to other/serious emulators like x86, ARM and RISC-V

### Why Go:
A static, compiled and yet simple language. No need to worry about memory management. Python like syntax. Also, haven't done any project in Go for a while. Wanter to get my hands dirty again.

## CHIP-8 hardware overview:
## Intro
A virtual hardware/programming language. It has instructions, supports 16 bit memory and other peripherals, which can be mapped as per implementer. Memory is unified program and data memory, dedicated stack, attached keyboard, display, sound and timer. Instructions for all peripherals. 
## Memory:
CHIP-8 can address upto 4 KB memory. First 512 bytes are reserved for interpreter, as per convention. So most of the ROMs assume code to start at 0x200 (512). Sprite memory access assumes no such convention, but to keep it uniform, I ensured all addresses, except font sprites, lie above 0x200 boundry. 

## Registers:
Language has 16 general purpose 8 bit registers, named V0-V15. Many instructions use V15 as a flag register and usually avoided in general purpose operations. There is I register which is only 16 bit register, typically used for storing address. There are also regsiters TD for timer and SD for sound. 

There are some shadow registers like program counter (PC) and stack pointer (SP) which are not directly accessible by instructions.

## Stack:
CHIP-8 provisions 16 entry stack with each entry 16 bit wide. It is typically used for saving return address in function call. Threre are no explicit push and pop instructions.

## Peripherals: Keyboard, Display, Sound, Timer, Sprites:
### Keyboard: 
CHIP-8 based system originally used 16 key keyboard, although instructions allow for upto 8 bit keycodes. Most of the ROMs though assume 16 keys only.

### Display: 
CHIP-8 based systems used 64x32 1 bit display, which is used by most of the ROMs
#### Sprites:
CHIP-8 uses sprites for drawing. Basically each sprite is in the memory representing a drawing block of 8 columns and 1-15 rows. The drawing instruction simply copies this block to display memory.
#### Font sprites:
These are fixed sprites of size 5x8, representing characters 0-9 and A-F.

### Sound:
CHIP-8 supports just one tone and instruction allows setting a sound timer. Tone starts whenever non-zero value is written to this timer and remians on till timer reaches 0.

### Timer:
CHIP-8 supports 8 bit timer. There is one instruction to set timer value and other to read current timer value.

## Iteration-0: Raylib-go minimal example
I decided to use [Raylib-go](https://github.com/gen2brain/raylib-go) for this project. No specific reason for choosing raylib over SDL, just wanted to explore it. From my experience on this project, I will continue to use it in similar projects.

Before starting actual emulator implementation, I wanted to overview raylib and its golang binding. I started with its hello world code. Next, I finalized window dimensions to be 1280x640 with 64x32 pixels, each pixel represented by a square of size 18x18 and spacing of 2 between pixels. 

Display buffer is implemented as 32 element array of uint64. Essentially, each 64 bit unsigned int represents a row on display. This works since CHIP-8 display has just 1 bit color. To render this display buffer, I've used 2 nested loops, drawing pixels column by column, from left to right and top to bottom.

## Iteration-1: minimum functioning system
### Loading ROM:
First thing we need to do is find a Chip-8 ROM and load it. I used test ROMs from [Timendus Chip 8 test suite](https://github.com/Timendus/chip8-test-suite/). Chip-8 ROM does not have any header or section. It starts right with instructions. Only contraint is loading address of ROM should be 0x200, since jump instructions expect it. 

I created a Machine struct in code and added 4KB uint8 array to represent RAM. In this RAM I read ROM file from slice [0x200:]. This loads instructions at right address and we don't have to do any additional calculation. 

First ROM I want to run correctly is Chip-8 splash screen. It needs following instructions to be implemented:
- 00E0: clear the screen
- 6xnn: load immediate value in normal register x
- Annn: Load immediate value in index (address) register
- Dxyn: Draw spirte
- 1nnn: Jump

Before implementing these instructions, we need to add registers and other fields to processor. I added general purpose register array, index register and program counter register.

Now only tricky part in above instructions is draw sprite. Dxyn gives x, y and n, i.e. draw sprite at location (Vx, Vy) in display buffer. x and y are register numbers instead of actual coordinates. `n` indicates number of bytes to draw. This makes sprite sizes from 1x8 to 15x8 possible.

Implementation is straightforward. Only different thing I did after encountering few bugs was to iterate and shift mask left to right.


## Iteration-2: Add just one more instruction for IBM logo
- Add instruction
## Iteration-3: Using test ROM and clearing all tests
- Grouping cases in switch
## Iteration-4: Flags
- Order of setting registers
## Iteration-5: Keyboard
- Key to scancode mapping
## Iteration-6: Timer
- Register and update
## Iteration-7: Sound
- Stub implementation
## Iteration-8: Random number generator
- Fixed number, tetris problem, go rand module
## Iteration-9: Threading
- Graphics in goroutine, segfaults, machine in gorouting
## Iteration-10: Music
- Raylib raw waves, playing mp3
## References:
- [Go raylib](https://github.com/gen2brain/raylib-go/)
- [Timendus Chip 8 test suite](https://github.com/Timendus/chip8-test-suite/)
- [Kripod's Chip 8 ROM collection](https://github.com/kripod/chip8-roms/)
- [Cowgod's Chip-8 reference](http://devernay.free.fr/hacks/chip8/)
- [Game sound](https://elements.envato.com/vibrant-game-game-key-1-LQA9WPL)
