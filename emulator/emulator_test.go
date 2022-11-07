package emulator

import (
	"testing"
)

func TestCreateDestroyEmulator(t *testing.T) {
	e, err := NewEmulator(
		"../assets/roms/bgbtest.gb",
		"../saves/test.sav",
		"",
		"../assets/fonts/arial.ttf",
		false,
	)
	if err != nil {
		t.Fatalf(`Failed to create emulator: %v`, err)
	}
	defer e.Destroy()
}

func RomTester(t *testing.T, romFilename, originalScreenshotFilename string, runCycles uint64) {
	e, err := NewEmulator(
		romFilename,
		"../saves/test.sav",
		"",
		"../assets/fonts/arial.ttf",
		false)
	if err != nil {
		t.Fatalf(`Failed to create emulator: %v`, err)
	}
	defer e.Destroy()

	e.RunTest(runCycles)

	err = e.TakeSnapshot("../screenshots/test.png")
	if err != nil {
		t.Fatalf(`Failed to take snapshot: %v`, err)
	}

	filesEqual, err := FileCompare("../screenshots/test.png", originalScreenshotFilename)
	if err != nil {
		t.Fatalf(`Failed to compare files: %v`, err)
	}

	if !filesEqual {
		t.Fatalf(`Screenshot not equal to: %v`, originalScreenshotFilename)
	}
}

func TestRunEmulator(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/bgbtest.gb",
		"../screenshots/tests/bgbtest.png",
		1000000,
	)
}

func TestRomPokeBlue(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/pokeblue.gb",
		"../screenshots/tests/pokeblue.png",
		100000000,
	)
}

func TestRomPokeGreen(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/pokegreen.gb",
		"../screenshots/tests/pokegreen.png",
		100000000,
	)
}

func TestRomPokeRed(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/pokered.gb",
		"../screenshots/tests/pokered.png",
		100000000,
	)
}

func TestRomTetris(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/tetris.gb",
		"../screenshots/tests/tetris.png",
		10000000,
	)
}

func TestRomOpus5(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/opus5.gb",
		"../screenshots/tests/opus5.png",
		500000,
	)
}

func TestRomZelda(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/zelda.gb",
		"../screenshots/tests/zelda.png",
		10000000,
	)
}

func TestRomSupermarioland(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/supermarioland.gb",
		"../screenshots/tests/supermarioland.png",
		10000000,
	)
}

func TestBlarggCPUInstrs(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/cpu_instrs.gb",
		"../screenshots/tests/cpu_instrs.png",
		300000000,
	)
}

func TestBlargg_01_special(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/01-special.gb",
		"../screenshots/tests/cpu_instrs/01-special.png",
		10000000,
	)
}

func TestBlargg_02_interrups(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/02-interrupts.gb",
		"../screenshots/tests/cpu_instrs/02-interrupts.png",
		10000000,
	)
}

func TestBlargg_03_op_sp_hl(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/03-op sp,hl.gb",
		"../screenshots/tests/cpu_instrs/03-op sp,hl.png",
		50000000,
	)
}

func TestBlargg_04_op_r_imm(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/04-op r,imm.gb",
		"../screenshots/tests/cpu_instrs/04-op r,imm.png",
		20000000,
	)
}

func TestBlargg_05_op_rp(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/05-op rp.gb",
		"../screenshots/tests/cpu_instrs/05-op rp.png",
		20000000,
	)
}

func TestBlargg_06_ld_rr(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/06-ld r,r.gb",
		"../screenshots/tests/cpu_instrs/06-ld r,r.png",
		10000000,
	)
}

func TestBlargg_07_jr_jp_call_ret_rst(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/07-jr,jp,call,ret,rst.gb",
		"../screenshots/tests/cpu_instrs/07-jr,jp,call,ret,rst.png",
		20000000,
	)
}

func TestBlargg_08_misc_instrs(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/08-misc instrs.gb",
		"../screenshots/tests/cpu_instrs/08-misc instrs.png",
		40000000,
	)
}

func TestBlargg_09_op_rr(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/09-op r,r.gb",
		"../screenshots/tests/cpu_instrs/09-op r,r.png",
		40000000,
	)
}

func TestBlargg_10_bit_ops(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/10-bit ops.gb",
		"../screenshots/tests/cpu_instrs/10-bit ops.png",
		60000000,
	)
}

func TestBlargg_11_op_a_hl(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/11-op a,(hl).gb",
		"../screenshots/tests/cpu_instrs/11-op a,(hl).png",
		100000000,
	)
}
