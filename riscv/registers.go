package riscv

type Reg struct {
	id int
}

type Registers struct {
	reg [32]uint32
	pc  uint32
}

const (
	reg_zero int = 0
	reg_ra   int = 1
	reg_sp   int = 2
	reg_gp   int = 3
	reg_tp   int = 4
	reg_t0   int = 5
	reg_t1   int = 6
	reg_t2   int = 7
	reg_s0   int = 8
	reg_fp   int = 8
	reg_s1   int = 9
	reg_a0   int = 10
	reg_a1   int = 11
	reg_a2   int = 12
	reg_a3   int = 13
	reg_a4   int = 14
	reg_a5   int = 15
	reg_a6   int = 16
	reg_a7   int = 17
	reg_s2   int = 18
	reg_s3   int = 19
	reg_s4   int = 20
	reg_s5   int = 21
	reg_s6   int = 22
	reg_s7   int = 23
	reg_s8   int = 24
	reg_s9   int = 25
	reg_s10  int = 26
	reg_s11  int = 27
	reg_t3   int = 28
	reg_t4   int = 29
	reg_t5   int = 30
	reg_t6   int = 31
)
