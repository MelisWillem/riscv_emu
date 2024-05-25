#!/bin/bash
TC_PATH="~/Downloads/riscv32-elf-ubuntu-22.04-llvm-nightly-2024.03.01-nightly/riscv/bin"
AS="${TC_PATH}/riscv32-unknown-elf-as"
LD="${TC_PATH}/riscv32-unknown-elf-ld"

bash -c "$AS -march=rv32i -mabi=ilp32 -o hello_world.o -c hello_world.s"
bash -c "$LD -T link.ld --no-dynamic-linker -m elf32lriscv -static -nostdlib -s -o hello.elf hello_world.o"
