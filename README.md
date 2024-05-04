# riscv_emu

A simple qemu emulator written in golang.

## tools
There are 2 executables in this project

1. Dumper: Takes in an elf-file and dumps the instructions in there instruction format. Used to test the decoding of instructions
2. Emulator: Takes in and elf-file and emulates the program, the elf-file should be an executable and not a library.

