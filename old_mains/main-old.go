package old_mains

//
//import (
//    emulator_old "emulator-go/emulator-old"
//    "fmt"
//    "os"
//)
//
//func main() {
//
//    bootRomData, err := os.ReadFile("./assets/roms/bootroms/DMG_ROM.bin")
//    if err != nil {
//        fmt.Println("Boot Rom file not found:", err)
//        return
//    }
//
//    gameRomData, err := os.ReadFile("./assets/roms/tetris.bin")
//    if err != nil {
//        fmt.Println("Game Rom file not found:", err)
//        return
//    }
//
//    var e emulator_old.Emulator
//
//    if err := e.Initialize(bootRomData, gameRomData); err != nil {
//        panic(err)
//    }
//
//    //imgui.RunImgui()
//
//    e.Run()
//
//    e.Destroy()
//}
