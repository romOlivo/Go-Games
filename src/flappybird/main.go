package main


import (
	"github.com/gen2brain/raylib-go/raylib"
	//"math/rand"
	//"time"
	//"strconv"
	"fmt"
)

// Initialize basic white window with fullscreen
func InitWindow() {
	DEFAULT_SCREEN_WIDTH := int32(1500);
	DEFAULT_SCREEN_HEIGHT := int32(750);

	rl.InitWindow(DEFAULT_SCREEN_WIDTH, DEFAULT_SCREEN_HEIGHT, "FlappyApples");

	rl.MaximizeWindow();
	rl.SetTargetFPS(60);
	
	monitor := rl.GetCurrentMonitor();
	screenWidth := rl.GetMonitorWidth(monitor);
    screenHeight := rl.GetMonitorHeight(monitor);
	rl.SetWindowSize(screenWidth, screenHeight);
	rl.SetWindowPosition(0, 0);
}


func main() {
    fmt.Println(" Initializing the game....");

	InitWindow();

	for !rl.WindowShouldClose() {
        rl.BeginDrawing();
		rl.ClearBackground(rl.RayWhite);
		rl.EndDrawing();
        //time.Sleep(50000000);
	}

	rl.CloseWindow();

}
