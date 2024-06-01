# riscv_emu

A simple qemu emulator written in golang.

## Tools
There are 2 executables in this project

1. Dumper: Takes in an elf-file and dumps the instructions in there instruction format. Used to test the decoding of instructions
2. Emulator: Takes in and elf-file and emulates the program, the elf-file should be an executable and not a library.

## Examples

### Dumper

Example run:
``` go run ./tools/dumper --file=./elf_files/hello.elf ```

Output:
```
2024/06/01 17:37:41 Elf file with machine string=EM_RISCV 
2024/06/01 17:37:41 section[0] has name- with sectionType=SHT_NULL(0) 
2024/06/01 17:37:41 section[1] has name-.text with sectionType=SHT_PROGBITS(1) 
2024/06/01 17:37:41 section[2] has name-.riscv.attributes with sectionType=SHT_LOPROC+3(1879048195) 
2024/06/01 17:37:41 section[3] has name-.shstrtab with sectionType=SHT_STRTAB(3) 
2024/06/01 17:37:41 Printing out section[1] 
2024/06/01 17:37:41 Decoding instructions of base instruction set
2024/06/01 17:37:41 prog[0]=680513 
2024/06/01 17:37:41 decoded as IInstr{imm=104, rs1=0, func3=0, rd=10, opcode=19} 
2024/06/01 17:37:41 prog[4]=1005b7 
2024/06/01 17:37:41 decoded as IInstr{imm=65536, rd=11, opcode=55 
2024/06/01 17:37:41 prog[8]=0a58023 
2024/06/01 17:37:41 decoded as IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35} 
2024/06/01 17:37:41 prog[12]=650513 
2024/06/01 17:37:41 decoded as IInstr{imm=101, rs1=0, func3=0, rd=10, opcode=19} 
2024/06/01 17:37:41 prog[16]=0a58023 
2024/06/01 17:37:41 decoded as IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35} 
2024/06/01 17:37:41 prog[20]=6c0513 
2024/06/01 17:37:41 decoded as IInstr{imm=108, rs1=0, func3=0, rd=10, opcode=19} 
2024/06/01 17:37:41 prog[24]=0a58023 
2024/06/01 17:37:41 decoded as IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35} 
2024/06/01 17:37:41 prog[28]=6c0513 
2024/06/01 17:37:41 decoded as IInstr{imm=108, rs1=0, func3=0, rd=10, opcode=19} 
2024/06/01 17:37:41 prog[32]=0a58023 
2024/06/01 17:37:41 decoded as IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35} 
2024/06/01 17:37:41 prog[36]=6f0513 
2024/06/01 17:37:41 decoded as IInstr{imm=111, rs1=0, func3=0, rd=10, opcode=19} 
2024/06/01 17:37:41 prog[40]=0a58023 
2024/06/01 17:37:41 decoded as IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35} 
2024/06/01 17:37:41 prog[44]=0006f 
2024/06/01 17:37:41 decoded as IInstr{(imm3=0, imm2=0, imm1=0, imm0=0) imm=0, rd=0, opcode=111} 
```

### Emulator

Example run:
``` go run ./tools/emulator/ -file=./elf_files/hello.elf -memory_size=1000000 -memory_offset=268000000 ```

Output:
```
2024/06/01 22:39:16 Registering base instruction set in decoder
2024/06/01 22:39:16 Executing body of section 1 with name .text 
2024/06/01 22:39:16 decoding instruction at byte offset 0
2024/06/01 22:39:16 executing instruction I=IInstr{imm=104, rs1=0, func3=0, rd=10, opcode=19}
2024/06/01 22:39:16 decoding instruction at byte offset 4
2024/06/01 22:39:16 executing instruction I=IInstr{imm=65536, rd=11, opcode=55
2024/06/01 22:39:16 decoding instruction at byte offset 8
2024/06/01 22:39:16 executing instruction I=IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35}
2024/06/01 22:39:16 decoding instruction at byte offset 12
2024/06/01 22:39:16 executing instruction I=IInstr{imm=101, rs1=0, func3=0, rd=10, opcode=19}
2024/06/01 22:39:16 decoding instruction at byte offset 16
2024/06/01 22:39:16 executing instruction I=IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35}
2024/06/01 22:39:16 decoding instruction at byte offset 20
2024/06/01 22:39:16 executing instruction I=IInstr{imm=108, rs1=0, func3=0, rd=10, opcode=19}
2024/06/01 22:39:16 decoding instruction at byte offset 24
2024/06/01 22:39:16 executing instruction I=IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35}
2024/06/01 22:39:16 decoding instruction at byte offset 28
2024/06/01 22:39:16 executing instruction I=IInstr{imm=108, rs1=0, func3=0, rd=10, opcode=19}
2024/06/01 22:39:16 decoding instruction at byte offset 32
2024/06/01 22:39:16 executing instruction I=IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35}
2024/06/01 22:39:16 decoding instruction at byte offset 36
2024/06/01 22:39:16 executing instruction I=IInstr{imm=111, rs1=0, func3=0, rd=10, opcode=19}
2024/06/01 22:39:16 decoding instruction at byte offset 40
2024/06/01 22:39:16 executing instruction I=IInstr{imm=10, rs1=11, func3=0, rd=0, opcode=35}
2024/06/01 22:39:16 decoding instruction at byte offset 44
2024/06/01 22:39:16 executing instruction I=IInstr{(imm3=0, imm2=0, imm1=0, imm0=0) imm=0, rd=0, opcode=111}
```

