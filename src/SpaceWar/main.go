package main


import (
	"github.com/gen2brain/raylib-go/raylib"
	//"math/rand"
	//"time"
	//"strconv"
	"fmt"
	"math"
)


var DEFAULT_SCREEN_WIDTH int32 = 1500;
var DEFAULT_SCREEN_HEIGHT int32 = 750;
var PLAYER_VELOCITY float32 = 0.004;
var BULLET_VELOCITY float32 = 0.01;
var TICKS_DELAY_SHOOT int = 50;

var BASIC_ENEMY_VELOCITY float32 = 0.002;

// -------------------------------------------------------------------------------------------------------------------------------
//
//                                            D I S P L A Y A B L E   S T R U C T S 
//
// -------------------------------------------------------------------------------------------------------------------------------

type DisplayObject interface {
	GetTexture() rl.Texture2D
	GetHeight() int32
	GetWidth() int32
	Tick(pos int)
	Die(pos int)
	Draw()
	Move()
}

// Implementation of Generic Displayables Objects
type DisplayableObject struct {
	texture rl.Texture2D 
	marginLeft float32
	marginTop float32
	rotation float32
	scaleX float32
	scaleY float32
	height int32
	image string
	width int32
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

func (x *Player) GetTexture() rl.Texture2D {
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
		b := &Bullet{dp: DisplayableObject{marginLeft: x.dp.marginLeft, marginTop: x.dp.marginTop}, radius: 4.0}
		AddDisplayableObject(b);
		AddBullet(b);
		x.ticksLastShoot = 0;
	}
}

func (x *Player) Die(pos int) {
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
	radius float32
	wantToDie bool
}

func (x Bullet) GetWidth() int32 {
	return x.dp.GetWidth();
}

func (x Bullet) GetHeight() int32 {
	return x.dp.GetHeight()
}

func (x *Bullet) GetTexture() rl.Texture2D {
	return x.dp.GetTexture();
}

func (x Bullet) GetCenter() rl.Vector2 {
	var v rl.Vector2;
	v.X = float32(x.GetWidth());
	v.Y = float32(x.GetHeight() + int32(40 * int32(rl.GetScreenHeight()) / DEFAULT_SCREEN_WIDTH) - 8);
	return v;
}

func (x Bullet) Draw() {
	rl.DrawCircleV(x.GetCenter(), x.radius, rl.Red)
}

func (x *Bullet) Move() {
	x.dp.marginLeft = float32(math.Min(1.0, float64(x.dp.marginLeft + BULLET_VELOCITY)));
}

func (x *Bullet) Die(pos int) {
	if (x.dp.marginLeft == 1.0 || x.wantToDie) {
		RemoveDisplayableObject(pos);
		RemoveBullet(x);
	}
}

func (x *Bullet) Tick(pos int) {
	x.Move();
	x.Draw();
	x.Die(pos);
}



// Implementation of the Enemy Displayable Struct

type Enemy struct {
	dp DisplayableObject
}

func (x Enemy) GetWidth() int32 {
	return x.dp.GetWidth();
}

func (x Enemy) GetHeight() int32 {
	return x.dp.GetHeight()
}

func (x *Enemy) GetTexture() rl.Texture2D {
	return x.dp.GetTexture();
}

func (x *Enemy) Draw() {
	var v rl.Vector2;
	v.X = float32(x.GetWidth());
	v.Y = float32(x.GetHeight());
	rl.DrawTextureEx(x.GetTexture(), v, x.dp.rotation, 1.0, rl.White)
}

func (x *Enemy) Move() {
	return
}

func (x *Enemy) Tick(pos int) {
	x.Move();
	x.Draw();
	x.Die(pos);
}

func (x *Enemy) Die(pos int) {
	if (x.dp.marginLeft == -0.2 || x.CollideWithBullet()) {
		RemoveDisplayableObject(pos);
	}
}

func (x Enemy) CollideWithBullet() bool{
	haveCollide := false;
	var rec rl.Rectangle;
	rec.X = float32(x.GetWidth());
	rec.Y = float32(x.GetHeight() - x.dp.height);
	rec.Height = float32(x.dp.height);
	rec.Width = float32(x.dp.width);
	// rl.DrawRectangleRec(rec, rl.Red);   --- For debuging purposes
	for i:=0; i<len(bulletsArray); i++ {
		bullet := bulletsArray[i];
		v := bullet.GetCenter();
		if (rl.CheckCollisionCircleRec(v, bullet.radius, rec)) {
			haveCollide = true;
			bullet.wantToDie = true;
			break;
		}
	}
	return haveCollide;
}

// Implementation of the Linear Enemy Displayable Struct

type LinearEnemy struct {
	enemy Enemy
}

func (x LinearEnemy) GetWidth() int32 {
	return x.enemy.GetWidth();
}

func (x LinearEnemy) GetHeight() int32 {
	return x.enemy.GetHeight()
}

func (x *LinearEnemy) GetTexture() rl.Texture2D {
	return x.enemy.GetTexture();
}

func (x *LinearEnemy) Draw() {
	x.enemy.Draw();
}

func (x *LinearEnemy) Move() {
	x.enemy.dp.marginLeft = float32(math.Max(-0.2, float64(x.enemy.dp.marginLeft - BASIC_ENEMY_VELOCITY)));
}

func (x *LinearEnemy) Tick(pos int) {
	x.Move();
	x.Draw();
	x.Die(pos);
}

func (x *LinearEnemy) Die(pos int) {
	x.enemy.Die(pos);
}

// -------------------------------------------------------------------------------------------------------------------------------
//
//                                                   E N E M Y   G E N E R A T O R
//
// -------------------------------------------------------------------------------------------------------------------------------

func GenerateLinearEnemy(mt float32, ml float32) {
	dp := DisplayableObject{marginTop: mt, marginLeft: ml, rotation: -90.0, image: "assets/ship1.png"};
	GenerateLinearEnemyDp(dp);
}

func GenerateLinearEnemyDp(dp DisplayableObject) {
	enemy := &LinearEnemy{enemy: Enemy{dp: dp}};
	AddDisplayableObject(enemy);
}

// -------------------------------------------------------------------------------------------------------------------------------
//
//                                     M A N I P U L A T I N G   D I S P L A Y   O B J E C T S
//
// -------------------------------------------------------------------------------------------------------------------------------


// Manipulating generics display objects
var displayArray []DisplayObject;
var newDisplayArray []DisplayObject;

func AddDisplayableObject[V DisplayObject](dp V) {
	displayArray = append(displayArray, dp);
}

func RemoveDisplayableObject(i int) {
	newDisplayArray[i] = newDisplayArray[len(newDisplayArray)-1]
    newDisplayArray = newDisplayArray[:len(newDisplayArray)-1]
}

// Manipulating Bullets
var bulletsArray []*Bullet;

func AddBullet(b *Bullet) {
	bulletsArray = append(bulletsArray, b);
}

func RemoveBullet(b *Bullet) {
	newBulletsArray := bulletsArray;
	for i:=0; i<len(bulletsArray); i++ {
		if(bulletsArray[i] == b) {
			newBulletsArray[i] = newBulletsArray[len(newBulletsArray)-1]
    		newBulletsArray = newBulletsArray[:len(newBulletsArray)-1]
		}
	}
	bulletsArray = newBulletsArray;
}

// -------------------------------------------------------------------------------------------------------------------------------
//
//                                                 L E V E L   M A N I P U L A T I O N
//
// -------------------------------------------------------------------------------------------------------------------------------

type Level interface {
	Tick()
	isEnded() bool
	End()
}

type Wave interface {
	makeWave();
	canActivateWave(ticks int) bool;
}

type DefinedLevel struct {
	ticks int;
	waves []Wave;
	waveNumber int;
}

func (x *DefinedLevel) Tick() {
	x.ticks++;
	if (x.waveNumber < len(x.waves) && x.waves[x.waveNumber].canActivateWave(x.ticks)) {
		x.waves[x.waveNumber].makeWave();
		x.waveNumber++;
	}
}

func (x DefinedLevel) isEnded() bool {
	return false;
}

func (x DefinedLevel) End() {

}

func (x *DefinedLevel) AddWave(w Wave) {
	x.waves = append(x.waves, w);
}

type BasicWave struct {
	tickToLaunch int;
	enemiesDp []DisplayableObject
}

func (x BasicWave) canActivateWave(ticks int) bool {
	return x.tickToLaunch == ticks;
}

func (x BasicWave) makeWave() {
	for i:=0; i<len(x.enemiesDp); i++ {
		GenerateLinearEnemyDp(x.enemiesDp[i]);
	}
}

func (x *BasicWave) addNewEnemy(dp DisplayableObject) {
	x.enemiesDp = append(x.enemiesDp, dp);
}

// -------------------------------------------------------------------------------------------------------------------------------
//
//                                                             M A I N
//
// -------------------------------------------------------------------------------------------------------------------------------


// Initialize basic white window with fullscreen
func InitWindow() {
	rl.InitWindow(DEFAULT_SCREEN_WIDTH, DEFAULT_SCREEN_HEIGHT, "Space War");

	rl.MaximizeWindow();
	rl.SetTargetFPS(100);
	
	monitor := rl.GetCurrentMonitor();
	screenWidth := rl.GetMonitorWidth(monitor);
    screenHeight := rl.GetMonitorHeight(monitor);
	rl.SetWindowSize(screenWidth, screenHeight);
	rl.SetWindowPosition(0, 0);
}

func GenerateLevel() Level {
	level := &DefinedLevel{}

	wave1 := &BasicWave{tickToLaunch: 1};
	wave1.addNewEnemy(DisplayableObject{marginTop: 0.05, marginLeft: 1.00, rotation: -90.0, image: "assets/ship1.png"});
	wave1.addNewEnemy(DisplayableObject{marginTop: 0.65, marginLeft: 1.20, rotation: -90.0, image: "assets/ship1.png"});
	wave1.addNewEnemy(DisplayableObject{marginTop: 0.45, marginLeft: 1.30, rotation: -90.0, image: "assets/ship1.png"});
	level.AddWave(wave1);

	wave2 := &BasicWave{tickToLaunch: 700};
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.08, marginLeft: 1.00, rotation: -90.0, image: "assets/ship1.png"});
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.75, marginLeft: 1.10, rotation: -90.0, image: "assets/ship1.png"});
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.45, marginLeft: 1.35, rotation: -90.0, image: "assets/ship1.png"});
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.70, marginLeft: 1.35, rotation: -90.0, image: "assets/ship1.png"});
	level.AddWave(wave2);

	wave3 := &BasicWave{tickToLaunch: 1400};
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.13, marginLeft: 1.00, rotation: -90.0, image: "assets/ship1.png"});
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.27, marginLeft: 1.05, rotation: -90.0, image: "assets/ship1.png"});
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.13, marginLeft: 1.10, rotation: -90.0, image: "assets/ship1.png"});
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.27, marginLeft: 1.15, rotation: -90.0, image: "assets/ship1.png"});
	level.AddWave(wave3);

	return level;
}

func main() {
    fmt.Println(" Initializing the game....");

	InitWindow();

	player := &Player{dp: DisplayableObject{marginTop: 0.05, marginLeft: 0.09, rotation: 90.0, image: "assets/player.png"}};

	level := GenerateLevel();

	for !rl.WindowShouldClose() && !level.isEnded() {
        rl.BeginDrawing();
		rl.ClearBackground(rl.RayWhite);
		player.Tick(-1);
		newDisplayArray = displayArray
		for i := len(displayArray)-1; i >= 0 ; i-- {
			displayArray[i].Tick(i)
		}
		displayArray = newDisplayArray
		level.Tick();
		rl.EndDrawing();
	}

	rl.CloseWindow();

}
