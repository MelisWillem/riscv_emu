#!/bin/bash
# Click on view->serial
# there is an option -no-graphic that seems to work, but can't exit the console
qemu-system-riscv32 -machine virt -bios hello.elf
