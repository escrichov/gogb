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
		"../screenshots/tests/blargg/cpu_instrs.png",
		300000000,
	)
}

func TestBlargg_01_special(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/01-special.gb",
		"../screenshots/tests/blargg/cpu_instrs/01-special.png",
		10000000,
	)
}

func TestBlargg_02_interrups(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/02-interrupts.gb",
		"../screenshots/tests/blargg/cpu_instrs/02-interrupts.png",
		10000000,
	)
}

func TestBlargg_03_op_sp_hl(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/03-op sp,hl.gb",
		"../screenshots/tests/blargg/cpu_instrs/03-op sp,hl.png",
		50000000,
	)
}

func TestBlargg_04_op_r_imm(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/04-op r,imm.gb",
		"../screenshots/tests/blargg/cpu_instrs/04-op r,imm.png",
		20000000,
	)
}

func TestBlargg_05_op_rp(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/05-op rp.gb",
		"../screenshots/tests/blargg/cpu_instrs/05-op rp.png",
		20000000,
	)
}

func TestBlargg_06_ld_rr(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/06-ld r,r.gb",
		"../screenshots/tests/blargg/cpu_instrs/06-ld r,r.png",
		10000000,
	)
}

func TestBlargg_07_jr_jp_call_ret_rst(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/07-jr,jp,call,ret,rst.gb",
		"../screenshots/tests/blargg/cpu_instrs/07-jr,jp,call,ret,rst.png",
		20000000,
	)
}

func TestBlargg_08_misc_instrs(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/08-misc instrs.gb",
		"../screenshots/tests/blargg/cpu_instrs/08-misc instrs.png",
		40000000,
	)
}

func TestBlargg_09_op_rr(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/09-op r,r.gb",
		"../screenshots/tests/blargg/cpu_instrs/09-op r,r.png",
		40000000,
	)
}

func TestBlargg_10_bit_ops(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/10-bit ops.gb",
		"../screenshots/tests/blargg/cpu_instrs/10-bit ops.png",
		60000000,
	)
}

func TestBlargg_11_op_a_hl(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cpu_instrs/individual/11-op a,(hl).gb",
		"../screenshots/tests/blargg/cpu_instrs/11-op a,(hl).png",
		100000000,
	)
}

func TestBlargg_cgb_sound_01_registers(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/01-registers.gb",
		"../screenshots/tests/cgb_sound/01-registers.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_02_len_ctr(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/02-len ctr.gb",
		"../screenshots/tests/cgb_sound/02-len ctr.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_03_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/03-trigger.gb",
		"../screenshots/tests/cgb_sound/03-trigger.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_04_sweep(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/04-sweep.gb",
		"../screenshots/tests/cgb_sound/04-sweep.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_05_sweep_details(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/05-sweep details.gb",
		"../screenshots/tests/cgb_sound/05-sweep details.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_06_overflow_on_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/06-overflow on trigger.gb",
		"../screenshots/tests/cgb_sound/06-overflow on trigger.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_07_len_sweep_period_sync(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/07-len sweep period sync.gb",
		"../screenshots/tests/cgb_sound/07-len sweep period sync.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_08_len_ctr_during_power(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/08-len ctr during power.gb",
		"../screenshots/tests/cgb_sound/08-len ctr during power.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_09_wave_read_while_on(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/09-wave read while on.gb",
		"../screenshots/tests/cgb_sound/09-wave read while on.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_10_wave_trigger_while_on(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/10-wave trigger while on.gb",
		"../screenshots/tests/cgb_sound/10-wave trigger while on.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_11_regs_after_power(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/11-regs after power.gb",
		"../screenshots/tests/cgb_sound/11-regs after power.png",
		100000000,
	)
}

func TestBlargg_cgb_sound_12_wave(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/rom_singles/12-wave.gb",
		"../screenshots/tests/cgb_sound/12-wave.png",
		100000000,
	)
}

func TestBlargg_cgb_sound(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/cgb_sound/cgb_sound.gb",
		"../screenshots/tests/cgb_sound/cgb_sound.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_01_registers(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/01-registers.gb",
		"../screenshots/tests/dmg_sound/01-registers.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_02_len_ctr(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/02-len ctr.gb",
		"../screenshots/tests/dmg_sound/02-len ctr.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_03_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/03-trigger.gb",
		"../screenshots/tests/dmg_sound/03-trigger.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_04_sweep(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/04-sweep.gb",
		"../screenshots/tests/dmg_sound/04-sweep.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_05_sweep_details(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/05-sweep details.gb",
		"../screenshots/tests/dmg_sound/05-sweep details.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_06_overflow_on_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/06-overflow on trigger.gb",
		"../screenshots/tests/dmg_sound/06-overflow on trigger.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_07_len_sweep_period_sync(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/07-len sweep period sync.gb",
		"../screenshots/tests/dmg_sound/07-len sweep period sync.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_08_len_ctr_during_power(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/08-len ctr during power.gb",
		"../screenshots/tests/dmg_sound/08-len ctr during power.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_09_wave_read_while_on(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/09-wave read while on.gb",
		"../screenshots/tests/dmg_sound/09-wave read while on.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_10_wave_trigger_while_on(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/10-wave trigger while on.gb",
		"../screenshots/tests/dmg_sound/10-wave trigger while on.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_11_regs_after_power(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/11-regs after power.gb",
		"../screenshots/tests/dmg_sound/11-regs after power.png",
		100000000,
	)
}

func TestBlargg_dmg_sound_12_wave_write_while_on(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/rom_singles/12-wave write while on.gb",
		"../screenshots/tests/dmg_sound/12-wave write while on.png",
		100000000,
	)
}

func TestBlargg_dmg_sound(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/dmg_sound/dmg_sound.gb",
		"../screenshots/tests/dmg_sound/dmg_sound.png",
		100000000,
	)
}

func TestBlargg_instr_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/instr_timing/instr_timing.gb",
		"../screenshots/tests/instr_timing/instr_timing.png",
		100000000,
	)
}

func TestBlargg_interrupt_time(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/interrupt_time/interrupt_time.gb",
		"../screenshots/tests/interrupt_time/interrupt_time.png",
		100000000,
	)
}

func TestBlargg_mem_timing_01_read_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing/individual/01-read_timing.gb",
		"../screenshots/tests/mem_timing/01-read_timing.png",
		100000000,
	)
}

func TestBlargg_mem_timing_02_write_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing/individual/02-write_timing.gb",
		"../screenshots/tests/mem_timing/02-write_timing.png",
		100000000,
	)
}

func TestBlargg_mem_timing_03_modify_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing/individual/03-modify_timing.gb",
		"../screenshots/tests/mem_timing/03-modify_timing.png",
		100000000,
	)
}

func TestBlargg_mem_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing/mem_timing.gb",
		"../screenshots/tests/mem_timing/mem_timing.png",
		100000000,
	)
}

func TestBlargg_mem_timing_2_01_read_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing-2/rom_singles/01-read_timing.gb",
		"../screenshots/tests/mem_timing-2/01-read_timing.png",
		100000000,
	)
}

func TestBlargg_mem_timing_2_02_write_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing-2/rom_singles/02-write_timing.gb",
		"../screenshots/tests/mem_timing-2/02-write_timing.png",
		100000000,
	)
}

func TestBlargg_mem_timing_2_03_modify_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing-2/rom_singles/03-modify_timing.gb",
		"../screenshots/tests/mem_timing-2/03-modify_timing.png",
		100000000,
	)
}

func TestBlargg_mem_timing_2(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/mem_timing-2/mem_timing.gb",
		"../screenshots/tests/mem_timing-2/mem_timing.png",
		100000000,
	)
}

func TestBlargg_oam_bug_1_lcd_sync(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/1-lcd_sync.gb",
		"../screenshots/tests/blargg/oam_bug/1-lcd_sync.png",
		100000000,
	)
}

func TestBlargg_oam_bug_2_causes(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/2-causes.gb",
		"../screenshots/tests/blargg/oam_bug/2-causes.png",
		100000000,
	)
}

func TestBlargg_oam_bug_3_non_causes(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/3-non_causes.gb",
		"../screenshots/tests/blargg/oam_bug/3-non_causes.png",
		100000000,
	)
}

func TestBlargg_oam_bug_4_scanline_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/4-scanline_timing.gb",
		"../screenshots/tests/blargg/oam_bug/4-scanline_timing.png",
		100000000,
	)
}

func TestBlargg_oam_bug_5_timing_bug(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/5-timing_bug.gb",
		"../screenshots/tests/blargg/oam_bug/5-timing_bug.png",
		100000000,
	)
}

func TestBlargg_oam_bug_6_timing_no_bug(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/6-timing_no_bug.gb",
		"../screenshots/tests/blargg/oam_bug/6-timing_no_bug.png",
		100000000,
	)
}

func TestBlargg_oam_bug_7_timing_effect(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/7-timing_effect.gb",
		"../screenshots/tests/blargg/oam_bug/7-timing_effect.png",
		100000000,
	)
}

func TestBlargg_oam_bug_8_instr_effect(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/rom_singles/8-instr_effect.gb",
		"../screenshots/tests/blargg/oam_bug/8-instr_effect.png",
		100000000,
	)
}

func TestBlargg_oam_bug(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/oam_bug/oam_bug.gb",
		"../screenshots/tests/blargg/oam_bug/oam_bug.png",
		100000000,
	)
}

func TestBlargg_halt_bug(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/gb-test-roms/halt_bug.gb",
		"../screenshots/tests/blargg/interrupt_time/interrupt_time.png",
		100000000,
	)
}

func TestMooneye_acceptance_bits_mem_oam(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/bits/mem_oam.gb",
		"../screenshots/tests/mooneye/bits/mem_oam.png",
		1000000,
	)
}

func TestMooneye_acceptance_bits_reg_f(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/bits/reg_f.gb",
		"../screenshots/tests/mooneye/bits/reg_f.png",
		1000000,
	)
}

func TestMooneye_acceptance_bits_unused_hwio_gs(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/bits/unused_hwio-GS.gb",
		"../screenshots/tests/mooneye/bits/unused_hwio-GS.png",
		1000000,
	)
}

func TestMooneye_acceptance_instr_daa(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/instr/daa.gb",
		"../screenshots/tests/mooneye/instr/daa.png",
		10000000,
	)
}

func TestMooneye_acceptance_interrupts_ie_push(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/interrupts/ie_push.gb",
		"../screenshots/tests/mooneye/interrupts/ie_push.png",
		10000000,
	)
}

func TestMooneye_acceptance_oam_dma_basic(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/oam_dma/basic.gb",
		"../screenshots/tests/mooneye/oam_dma/basic.png",
		10000000,
	)
}

func TestMooneye_acceptance_oam_dma_reg_read(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/oam_dma/reg_read.gb",
		"../screenshots/tests/mooneye/oam_dma/reg_read.png",
		10000000,
	)
}

func TestMooneye_acceptance_oam_dma_sources_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/oam_dma/sources-GS.gb",
		"../screenshots/tests/mooneye/oam_dma/sources-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_hblank_ly_scx_timing_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/hblank_ly_scx_timing-GS.gb",
		"../screenshots/tests/mooneye/ppu/hblank_ly_scx_timing-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_intr_1_2_timing_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/intr_1_2_timing-GS.gb",
		"../screenshots/tests/mooneye/ppu/intr_1_2_timing-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_intr_2_0_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/intr_2_0_timing.gb",
		"../screenshots/tests/mooneye/ppu/intr_2_0_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_intr_2_mode0_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/intr_2_mode0_timing.gb",
		"../screenshots/tests/mooneye/ppu/intr_2_mode0_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_intr_2_mode0_timing_sprites(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/intr_2_mode0_timing_sprites.gb",
		"../screenshots/tests/mooneye/ppu/intr_2_mode0_timing_sprites.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_intr_2_mode3_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/intr_2_mode3_timing.gb",
		"../screenshots/tests/mooneye/ppu/intr_2_mode3_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_intr_2_oam_ok_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/intr_2_oam_ok_timing.gb",
		"../screenshots/tests/mooneye/ppu/intr_2_oam_ok_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_lcdon_timing_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/lcdon_timing-GS.gb",
		"../screenshots/tests/mooneye/ppu/lcdon_timing-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_lcdon_write_timing_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/lcdon_write_timing-GS.gb",
		"../screenshots/tests/mooneye/ppu/lcdon_write_timing-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_stat_irq_blocking(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/stat_irq_blocking.gb",
		"../screenshots/tests/mooneye/ppu/stat_irq_blocking.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_stat_lyc_onoff(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/stat_lyc_onoff.gb",
		"../screenshots/tests/mooneye/ppu/stat_lyc_onoff.png",
		10000000,
	)
}

func TestMooneye_acceptance_ppu_vblank_stat_intr_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ppu/vblank_stat_intr-GS.gb",
		"../screenshots/tests/mooneye/ppu/vblank_stat_intr-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_timer_div_write(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/div_write.gb",
		"../screenshots/tests/mooneye/timer/div_write.png",
		10000000,
	)
}

func TestMooneye_acceptance_timer_rapid_toggle(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/rapid_toggle.gb",
		"../screenshots/tests/mooneye/timer/rapid_toggle.png",
		100000000,
	)
}

func TestMooneye_acceptance_timer_tim00(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim00.gb",
		"../screenshots/tests/mooneye/timer/tim00.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tim00_div_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim00_div_trigger.gb",
		"../screenshots/tests/mooneye/timer/tim00_div_trigger.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tim01(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim01.gb",
		"../screenshots/tests/mooneye/timer/tim01.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tim01_div_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim01_div_trigger.gb",
		"../screenshots/tests/mooneye/timer/tim01_div_trigger.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tim10(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim10.gb",
		"../screenshots/tests/mooneye/timer/tim10.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tim10_div_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim10_div_trigger.gb",
		"../screenshots/tests/mooneye/timer/tim10_div_trigger.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tim11(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim11.gb",
		"../screenshots/tests/mooneye/timer/tim11.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tim11_div_trigger(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tim11_div_trigger.gb",
		"../screenshots/tests/mooneye/timer/tim11_div_trigger.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tima_reload(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tima_reload.gb",
		"../screenshots/tests/mooneye/timer/tima_reload.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tima_write_reloading(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tima_write_reloading.gb",
		"../screenshots/tests/mooneye/timer/tima_write_reloading.png",
		1000000,
	)
}

func TestMooneye_acceptance_timer_tma_write_reloading(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/timer/tma_write_reloading.gb",
		"../screenshots/tests/mooneye/timer/tma_write_reloading.png",
		1000000,
	)
}

func TestMooneye_acceptance_serial_boot_sclk_align_dmgABCmgb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/serial/boot_sclk_align-dmgABCmgb.gb",
		"../screenshots/tests/mooneye/serial/boot_sclk_align-dmgABCmgb.png",
		1000000,
	)
}

func TestMooneye_acceptance_add_sp_e_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/add_sp_e_timing.gb",
		"../screenshots/tests/mooneye/add_sp_e_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_boot_div_S(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/boot_div-S.gb",
		"../screenshots/tests/mooneye/boot_div-S.png",
		1000000,
	)
}

func TestMooneye_acceptance_boot_div_dmg0(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/boot_div-dmg0.gb",
		"../screenshots/tests/mooneye/boot_div-dmg0.png",
		1000000,
	)
}

func TestMooneye_acceptance_boot_div2_S(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/boot_div2-S.gb",
		"../screenshots/tests/mooneye/boot_div2-S.png",
		1000000,
	)
}

func TestMooneye_acceptance_boot_hwio_S(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/boot_hwio-S.gb",
		"../screenshots/tests/mooneye/boot_hwio-S.png",
		1000000,
	)
}

func TestMooneye_acceptance_boot_hwio_dmg0(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/boot_hwio-dmg0.gb",
		"../screenshots/tests/mooneye/boot_hwio.png",
		1000000,
	)
}

func TestMooneye_acceptance_boot_regs_dmg0(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/boot_regs-dmg0.gb",
		"../screenshots/tests/mooneye/boot_regs-dmg0.png",
		10000000,
	)
}

func TestMooneye_acceptance_call_cc_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/call_cc_timing.gb",
		"../screenshots/tests/mooneye/call_cc_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_call_cc_timing2(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/call_cc_timing2.gb",
		"../screenshots/tests/mooneye/call_cc_timing2.png",
		1000000,
	)
}

func TestMooneye_acceptance_call_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/call_timing.gb",
		"../screenshots/tests/mooneye/call_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_call_timing2(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/call_timing2.gb",
		"../screenshots/tests/mooneye/call_timing2.png",
		10000000,
	)
}

func TestMooneye_acceptance_di_timing_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/di_timing-GS.gb",
		"../screenshots/tests/mooneye/di_timing-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_div_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/div_timing.gb",
		"../screenshots/tests/mooneye/div_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_ei_sequence(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ei_sequence.gb",
		"../screenshots/tests/mooneye/ei_sequence.png",
		10000000,
	)
}

func TestMooneye_acceptance_ei_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ei_timing.gb",
		"../screenshots/tests/mooneye/ei_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_halt_ime0_ei(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/halt_ime0_ei.gb",
		"../screenshots/tests/mooneye/halt_ime0_ei.png",
		1000000,
	)
}

func TestMooneye_acceptance_halt_ime0_nointr_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/halt_ime0_nointr_timing.gb",
		"../screenshots/tests/mooneye/halt_ime0_nointr_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_halt_ime1_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/halt_ime1_timing.gb",
		"../screenshots/tests/mooneye/halt_ime1_timing.png",
		10000000,
	)
}

func TestMooneye_acceptance_halt_ime1_timing2_GS(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/halt_ime1_timing2-GS.gb",
		"../screenshots/tests/mooneye/halt_ime1_timing2-GS.png",
		10000000,
	)
}

func TestMooneye_acceptance_if_ie_registers(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/if_ie_registers.gb",
		"../screenshots/tests/mooneye/if_ie_registers.png",
		1000000,
	)
}

func TestMooneye_acceptance_intr_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/intr_timing.gb",
		"../screenshots/tests/mooneye/intr_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_jp_cc_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/jp_cc_timing.gb",
		"../screenshots/tests/mooneye/jp_cc_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_jp_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/jp_timing.gb",
		"../screenshots/tests/mooneye/jp_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_ld_hl_sp_e_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ld_hl_sp_e_timing.gb",
		"../screenshots/tests/mooneye/ld_hl_sp_e_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_oam_dma_restart(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/oam_dma_restart.gb",
		"../screenshots/tests/mooneye/oam_dma_restart.png",
		1000000,
	)
}

func TestMooneye_acceptance_oam_dma_start(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/oam_dma_start.gb",
		"../screenshots/tests/mooneye/oam_dma_start.png",
		1000000,
	)
}

func TestMooneye_acceptance_oam_dma_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/oam_dma_timing.gb",
		"../screenshots/tests/mooneye/oam_dma_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_pop_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/pop_timing.gb",
		"../screenshots/tests/mooneye/pop_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_push_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/push_timing.gb",
		"../screenshots/tests/mooneye/push_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_rapid_di_ei(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/rapid_di_ei.gb",
		"../screenshots/tests/mooneye/rapid_di_ei.png",
		1000000,
	)
}

func TestMooneye_acceptance_ret_cc_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/ret_cc_timing.gb",
		"../screenshots/tests/mooneye/ret_cc_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_reti_intr_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/reti_intr_timing.gb",
		"../screenshots/tests/mooneye/reti_intr_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_reti_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/reti_timing.gb",
		"../screenshots/tests/mooneye/reti_timing.png",
		1000000,
	)
}

func TestMooneye_acceptance_rst_timing(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/acceptance/rst_timing.gb",
		"../screenshots/tests/mooneye/rst_timing.png",
		1000000,
	)
}

func TestMooneye_emulator_only_mbc1_bits_bank1(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/bits_bank1.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/bits_bank1.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc1_bits_bank2(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/bits_bank2.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/bits_bank2.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc1_bits_mode(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/bits_mode.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/bits_mode.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc1_bits_ramg(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/bits_ramg.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/bits_ramg.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc1_multicart_rom_8Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/multicart_rom_8Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/multicart_rom_8Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc1_ram_256kb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/ram_256kb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/ram_256kb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc1_ram_64kb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/ram_64kb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/ram_64kb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc1_rom_16Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/rom_16Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/rom_16Mb.png",
		10000000,
	)
}
func TestMooneye_emulator_only_mbc1_rom_1Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/rom_1Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/rom_1Mb.png",
		10000000,
	)
}
func TestMooneye_emulator_only_mbc1_rom_2Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/rom_2Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/rom_2Mb.png",
		10000000,
	)
}
func TestMooneye_emulator_only_mbc1_rom_4Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/rom_4Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/rom_4Mb.png",
		10000000,
	)
}
func TestMooneye_emulator_only_mbc1_rom_512kb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/rom_512kb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/rom_512kb.png",
		10000000,
	)
}
func TestMooneye_emulator_only_mbc1_rom_8Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc1/rom_8Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc1/rom_8Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc2_bits_ramg(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc2/bits_ramg.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc2/bits_ramg.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc2_bits_romb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc2/bits_romb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc2/bits_romb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc2_bits_unused(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc2/bits_unused.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc2/bits_unused.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc2_ram(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc2/ram.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc2/ram.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc2_rom_1Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc2/rom_1Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc2/rom_1Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc2_rom_2Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc2/rom_2Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc2/rom_2Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc2_rom_512kb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc2/rom_512kb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc2/rom_512kb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_16Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_16Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_16Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_1Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_1Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_1Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_2Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_2Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_2Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_32Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_32Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_32Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_4Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_4Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_4Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_512kb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_512kb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_512kb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_64Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_64Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_64Mb.png",
		10000000,
	)
}

func TestMooneye_emulator_only_mbc5_rom_8Mb(t *testing.T) {
	RomTester(
		t,
		"../assets/roms/mooneye/emulator-only/mbc5/rom_8Mb.gb",
		"../screenshots/tests/mooneye/emulator-only/mbc5/rom_8Mb.png",
		10000000,
	)
}
