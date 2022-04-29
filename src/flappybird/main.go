package main


import (
	"github.com/gen2brain/raylib-go/raylib"
	//"math/rand"
	//"time"
	//"strconv"
	"fmt"
	//"math"
)

// ---------------------

type DisplayObject interface {
	GetWidth()
	GetHeight()
	GetTexture()
}

type MyTexture struct {
	texture rl.Texture2D
}

type Player struct {
	height float64
	width float64
	image string
	texture rl.Texture2D 
}

func (x Player) GetWidth() int32 {
	return int32(x.width * float64(rl.GetScreenWidth()))
}

func (x Player) GetHeight() int32 {
	return int32(x.height * float64(rl.GetScreenHeight()))
}

func (x Player) GetTexture() rl.Texture2D {
	return rl.LoadTextureFromImage(rl.LoadImage(x.image))
}

// ---------------------




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

	player := Player{height: 0.05, width: 0.05, image: "assets/player.png"}

	for !rl.WindowShouldClose() {
        rl.BeginDrawing();
		rl.ClearBackground(rl.RayWhite);

		rl.DrawTexture(player.GetTexture(), player.GetWidth(), player.GetHeight(), rl.White)
		rl.EndDrawing();
        //time.Sleep(50000000);
	}

	rl.CloseWindow();

}
