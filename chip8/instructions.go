package chip8

type Instruction int

const (
	SYS       Instruction = iota // 0nnn
	CLS                          // 00E0
	RET                          // 00EE
	JPAddr                       // 1nnn
	CALL                         // 2nnn
	SEVxByte                     // 3xkk
	SNEVxByte                    // 4xkk
	SEVxVy                       // 5xy0
	LDVxByte                     // 6xkk
	ADDVxByte                    // 7xkk
	LDVxVy                       // 8xy0
	OR                           // 8xy1
	AND                          // 8xy2
	XOR                          // 8xy3
	ADDVxVy                      // 8xy4
	SUB                          // 8xy5
	SHR                          // 8xy6
	SUBN                         // 8xy7
	SHL                          // 8xyE
	SNEVxVy                      // 9xy0
	LDIAddr                      // Annn
	JPV0Addr                     // Bnnn
	RND                          // Cxkk
	DRW                          // Dxyn
	SKP                          // Ex9E
	SKNP                         // ExA1
	LDVxDT                       // Fx07
	LDVxK                        // Fx0A
	LDDTVx                       // Fx15
	LDSTVx                       // Fx18
	ADDIVx                       // Fx1E
	LDFVx                        // Fx29
	LDBVx                        // Fx33
	LDIVx                        // Fx55
	LDVxI                        // Fx65
	UNKNOWN
)
