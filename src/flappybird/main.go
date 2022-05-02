package main


import (
	"github.com/gen2brain/raylib-go/raylib"
	//"math/rand"
	"time"
	//"strconv"
	"fmt"
	"math"
)


var DEFAULT_SCREEN_WIDTH int32 = 1500;
var DEFAULT_SCREEN_HEIGHT int32 = 750;
var PLAYER_VELOCITY float32 = 0.004;
var BULLET_VELOCITY float32 = 0.01;
var TICKS_DELAY_SHOOT int = 50;

// ---------------------

type DisplayObject interface {
	GetWidth() int32
	GetHeight() int32
	GetTexture() rl.Texture2D
	Draw()
	Move()
	Die(pos int)
	Tick(pos int)
}

// Implementation of Generic Displayables Objects
type DisplayableObject struct {
	height int32
	width int32
	marginTop float32
	marginLeft float32
	rotation float32
	image string
	scaleX float32
	scaleY float32
	texture rl.Texture2D 
}

func (x DisplayableObject) GetWidth() int32 {
	return int32(x.marginLeft * float32(rl.GetScreenWidth()) + x.rotation / 360)
}

func (x DisplayableObject) GetHeight() int32 {
	return int32(x.marginTop * float32(rl.GetScreenHeight()))
}

func (x *DisplayableObject) GetTexture() rl.Texture2D {
	img := rl.LoadImage(x.image)
	x.scaleX = float32(int32(rl.GetScreenHeight()) / DEFAULT_SCREEN_HEIGHT);
	x.scaleY = float32(int32(rl.GetScreenWidth()) / DEFAULT_SCREEN_WIDTH);
	x.height = int32(40 * x.scaleX)
	x.width = int32(40 * x.scaleY)
	rl.ImageResize(img, x.height, x.width)
	return rl.LoadTextureFromImage(img)
}

func (x *DisplayableObject) Draw() {
	return
}

func (x *DisplayableObject) Move() {
	return
}

// Implementation of the Player Displayable Struct
type Player struct {
	dp DisplayableObject
	ticksLastShoot int
}

func (x Player) GetWidth() int32 {
	return x.dp.GetWidth();
}

func (x Player) GetHeight() int32 {
	return x.dp.GetHeight()
}

func (x Player) GetTexture() rl.Texture2D {
	return x.dp.GetTexture();
}

func (x *Player) Draw() {
	var v rl.Vector2;
	v.X = float32(x.GetWidth());
	v.Y = float32(x.GetHeight());
	rl.DrawTextureEx(x.GetTexture(), v, x.dp.rotation, 1.0, rl.White)
}

func (x *Player) Move() {
	if rl.IsKeyDown(rl.KeyUp){
		x.dp.marginTop = float32(math.Max(0.0, float64(x.dp.marginTop - PLAYER_VELOCITY)));
	}
	if rl.IsKeyDown(rl.KeyDown){
		x.dp.marginTop = float32(math.Min(1.0, float64(x.dp.marginTop + PLAYER_VELOCITY)));
	}
	if rl.IsKeyDown(rl.KeySpace) {
		x.Shoot()
	}
}

func (x *Player) Shoot() {
	if (x.ticksLastShoot >= TICKS_DELAY_SHOOT) {
		b := &Bullet{dp: DisplayableObject{marginLeft: x.dp.marginLeft, marginTop: x.dp.marginTop}}
		AddDisplayableObject(b);
		x.ticksLastShoot = 0;
	}
}

func (x Player) Die(pos int) {
	return
}

func (x *Player) Tick(pos int) {
	x.Move();
	x.Draw();
	x.Die(pos);
	x.ticksLastShoot++;
}

// Implementation of the Bullet Displayable Struct

type Bullet struct {
	dp DisplayableObject
}

func (x Bullet) GetWidth() int32 {
	return x.dp.GetWidth();
}

func (x Bullet) GetHeight() int32 {
	return x.dp.GetHeight()
}

func (x Bullet) GetTexture() rl.Texture2D {
	return x.dp.GetTexture();
}

func (x Bullet) Draw() {
	rl.DrawCircle(x.GetWidth(), x.GetHeight() + int32(40 * int32(rl.GetScreenHeight()) / DEFAULT_SCREEN_WIDTH) - 8, 4.0, rl.Red)
}

func (x *Bullet) Move() {
	x.dp.marginLeft = float32(math.Min(1.0, float64(x.dp.marginLeft + BULLET_VELOCITY)));
}

func (x Bullet) Die(pos int) {
	if (x.dp.marginLeft == 1.0) {
		RemoveDisplayableObject(pos);
	}
}

func (x *Bullet) Tick(pos int) {
	x.Move();
	x.Draw();
	x.Die(pos);
}

// ---------------------

var displayArray []DisplayObject;
var newDisplayArray []DisplayObject;


// Initialize basic white window with fullscreen
func InitWindow() {
	rl.InitWindow(DEFAULT_SCREEN_WIDTH, DEFAULT_SCREEN_HEIGHT, "FlappyApples");

	rl.MaximizeWindow();
	rl.SetTargetFPS(120);
	
	monitor := rl.GetCurrentMonitor();
	screenWidth := rl.GetMonitorWidth(monitor);
    screenHeight := rl.GetMonitorHeight(monitor);
	rl.SetWindowSize(screenWidth, screenHeight);
	rl.SetWindowPosition(0, 0);
}

func AddDisplayableObject[V DisplayObject](dp V) {
	displayArray = append(displayArray, dp);
}

func RemoveDisplayableObject(i int) {
	newDisplayArray[i] = newDisplayArray[len(newDisplayArray)-1]
    newDisplayArray = newDisplayArray[:len(newDisplayArray)-1]
}


func main() {
    fmt.Println(" Initializing the game....");

	InitWindow();

	player := &Player{dp: DisplayableObject{marginTop: 0.05, marginLeft: 0.09, rotation: 90.0, image: "assets/player.png"}};

	for !rl.WindowShouldClose() {
        rl.BeginDrawing();
		rl.ClearBackground(rl.RayWhite);
		player.Tick(-1);
		newDisplayArray = displayArray
		for i := len(displayArray)-1; i >= 0 ; i-- {
			displayArray[i].Tick(i)
		}
		displayArray = newDisplayArray
		rl.EndDrawing();
        time.Sleep(1000000);
	}

	rl.CloseWindow();

}
