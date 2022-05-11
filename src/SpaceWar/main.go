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
	textureLoaded bool
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
	if (!x.textureLoaded) {
		img := rl.LoadImage(x.image)
		x.scaleX = float32(int32(rl.GetScreenHeight()) / DEFAULT_SCREEN_HEIGHT);
		x.scaleY = float32(int32(rl.GetScreenWidth()) / DEFAULT_SCREEN_WIDTH);
		x.height = int32(40 * x.scaleX)
		x.width = int32(40 * x.scaleY)
		rl.ImageResize(img, x.height, x.width)
		x.texture = rl.LoadTextureFromImage(img)
		x.textureLoaded = true
	}
	return x.texture;
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
	loadedSound rl.Sound
	haveLoadedSound bool
	ticksLastShoot int
	sound string
	core Core;
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
		b := &Bullet{dp: DisplayableObject{marginLeft: x.dp.marginLeft, marginTop: x.dp.marginTop}, radius: 4.0, core: x.core}
		x.core.AddDisplayableObject(b);
		x.core.AddBullet(b);
		if (!x.haveLoadedSound) {
			x.loadedSound = rl.LoadSound(x.sound);
			x.haveLoadedSound = true;
		}
		rl.PlaySoundMulti(x.loadedSound)
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
	core Core
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
		x.core.RemoveDisplayableObject(pos);
		x.core.RemoveBullet(x);
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
	dieSound rl.Sound
	core Core
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
	collide := x.CollideWithBullet();
	if (x.dp.marginLeft == -0.2 || collide) {
		x.core.RemoveDisplayableObject(pos);
		if (collide) {
			rl.PlaySoundMulti(x.dieSound)
		}
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
	bulletsArray := x.core.GetBulletsArray()
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

type EnemyFactory interface {
	GenerateLinearEnemy(mt float32, ml float32)
	GenerateLinearEnemyDp(dp DisplayableObject)
	Initialize()
}

type BasicEnemyFactory struct {
	dieSoundLinearEnemy rl.Sound;
}

func (x *BasicEnemyFactory) Initialize() {
	x.dieSoundLinearEnemy = rl.LoadSound("sounds/explosion01.wav");
}

func (x *BasicEnemyFactory)  GenerateLinearEnemy(mt float32, ml float32) {
	dp := DisplayableObject{marginTop: mt, marginLeft: ml, rotation: -90.0, image: "assets/ship1.png"};
	x.GenerateLinearEnemyDp(dp);
}

func (x *BasicEnemyFactory)  GenerateLinearEnemyDp(dp DisplayableObject) {
	enemy := &LinearEnemy{enemy: Enemy{dp: dp, dieSound: x.dieSoundLinearEnemy, core: coreGame}};
	coreGame.AddDisplayableObject(enemy);
}

// -------------------------------------------------------------------------------------------------------------------------------
//
//                                                 L E V E L   M A N I P U L A T I O N
//
// -------------------------------------------------------------------------------------------------------------------------------

type Level interface {
	SetCore(core Core)
	GetCore() Core
	isEnded() bool
	Initialize()
	Tick()
	End()
}

type Wave interface {
	canActivateWave(ticks int) bool;
	makeWave();
}

type DefinedLevel struct {
	waveNumber int;
	waves []Wave;
	sound string;
	ticks int;
	core Core;
}

func (x DefinedLevel) Initialize() {
	s := rl.LoadSound(x.sound);
	rl.PlaySoundMulti(s);
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

func (x *DefinedLevel) SetCore(core Core) {
	x.core = core;
}
func (x *DefinedLevel) GetCore() Core {
	return x.core;
}

type BasicWave struct {
	enemiesDp []DisplayableObject
	tickToLaunch int;
	core Core;
}

func (x BasicWave) canActivateWave(ticks int) bool {
	return x.tickToLaunch == ticks;
}

func (x BasicWave) makeWave() {
	enemyFactory := x.core.GetEnemyFactory()
	for i:=0; i<len(x.enemiesDp); i++ {
		enemyFactory.GenerateLinearEnemyDp(x.enemiesDp[i]);
	}
}

func (x *BasicWave) addNewEnemy(dp DisplayableObject) {
	x.enemiesDp = append(x.enemiesDp, dp);
}

// -------------------------------------------------------------------------------------------------------------------------------
//
//                                                G A M E   E N G I N E   &   C O R E
//
// -------------------------------------------------------------------------------------------------------------------------------

// Manipulating Bullets

type BulletController interface {
	GetBulletsArray() []*Bullet
	RemoveBullet(b *Bullet)
	AddBullet(b *Bullet)
}

type BasicBulletController struct {
	bulletsArray []*Bullet;
}

func (x *BasicBulletController) AddBullet(b *Bullet) {
	x.bulletsArray = append(x.bulletsArray, b);
}

func (x *BasicBulletController) RemoveBullet(b *Bullet) {
	newBulletsArray := x.bulletsArray;
	for i:=0; i<len(x.bulletsArray); i++ {
		if(x.bulletsArray[i] == b) {
			newBulletsArray[i] = newBulletsArray[len(newBulletsArray)-1]
    		newBulletsArray = newBulletsArray[:len(newBulletsArray)-1]
		}
	}
	x.bulletsArray = newBulletsArray;
}

func (x *BasicBulletController) GetBulletsArray() []*Bullet{
	return x.bulletsArray;
}

// Manipulating all Displayable Objects

type DisplayableObjectController interface {
	AddDisplayableObject(dp DisplayObject)
	RemoveDisplayableObject(i int)
	Tick();
}

type BasicDisplayableObjectController struct {
	displayArray []DisplayObject
	newDisplayArray []DisplayObject
}

func (x *BasicDisplayableObjectController) AddDisplayableObject(dp DisplayObject) {
	x.displayArray = append(x.displayArray, dp);
}

func (x *BasicDisplayableObjectController) RemoveDisplayableObject(i int) {
	x.newDisplayArray[i] = x.newDisplayArray[len(x.newDisplayArray)-1]
    x.newDisplayArray = x.newDisplayArray[:len(x.newDisplayArray)-1]
}

func (x *BasicDisplayableObjectController) Tick() {
	x.newDisplayArray = x.displayArray
	for i := len(x.displayArray)-1; i >= 0 ; i-- {
		x.displayArray[i].Tick(i)
	}
	x.displayArray = x.newDisplayArray
}

// Game Engine to manipulate the game

type GameEngine interface {
	SetLevel(level Level)
	IsGameEnded() bool
	InitializeGame()
	Tick();
}

type BasicGameEngine struct {
	player *Player
	level Level
	core Core
}

func (x *BasicGameEngine) InitializeGame() {
	x.player = &Player{dp: DisplayableObject{marginTop: 0.05, marginLeft: 0.09, rotation: 90.0, image: "assets/player.png"}, 
	sound: "sounds/laserfire01.ogg", core: x.core};

	x.level = GenerateLevel(x.core);
	x.level.Initialize();
}

func (x BasicGameEngine) IsGameEnded() bool {
	return false
}

func (x *BasicGameEngine) SetLevel(level Level) {
	x.level = level;
}

func (x *BasicGameEngine) Tick() {
	x.player.Tick(-1);
	x.level.Tick();
}

// Core of the Game

type Core interface {
	DisplayableObjectController;
	BulletController;
	GameEngine;

	GetEnemyFactory() EnemyFactory;
}

type CoreGame struct {
	displayableObjectController DisplayableObjectController;
	bulletController BulletController;
	enemyFactory EnemyFactory;
	gameEngine GameEngine;
}

func (x *CoreGame) InitializeGame() {
	x.displayableObjectController = &BasicDisplayableObjectController{};
	x.bulletController = &BasicBulletController{};
	x.gameEngine = &BasicGameEngine{core: x};
	x.enemyFactory = &BasicEnemyFactory{};

	rl.InitAudioDevice();
	InitWindow();

	x.enemyFactory.Initialize();
	x.gameEngine.InitializeGame();
}

// Manipulating Displayable objects
func (x *CoreGame) AddDisplayableObject(dp DisplayObject) {
	x.displayableObjectController.AddDisplayableObject(dp);
}

func (x *CoreGame) RemoveDisplayableObject(i int) {
	x.displayableObjectController.RemoveDisplayableObject(i);
}

//Manipulating Bullets
func (x *CoreGame) AddBullet(b *Bullet) {
	x.bulletController.AddBullet(b);
}

func (x *CoreGame) RemoveBullet(b *Bullet) {
	x.bulletController.RemoveBullet(b);
}

func (x *CoreGame) GetBulletsArray() []*Bullet{
	return x.bulletController.GetBulletsArray();
}

// Manipulating Game
func (x CoreGame) IsGameEnded() bool {
	return x.gameEngine.IsGameEnded();
}

func (x *CoreGame) SetLevel(level Level) {
	x.gameEngine.SetLevel(level);
}

func (x *CoreGame) Tick() {
	rl.BeginDrawing();
	rl.ClearBackground(rl.RayWhite);
	x.displayableObjectController.Tick();
	x.gameEngine.Tick();
	rl.EndDrawing();
}

// Other functionality
func (x CoreGame) GetEnemyFactory() EnemyFactory {
	return x.enemyFactory;
}


// -------------------------------------------------------------------------------------------------------------------------------
//
//                                                             M A I N
//
// -------------------------------------------------------------------------------------------------------------------------------

var coreGame Core;

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

func GenerateLevel(core Core) Level {
	level := &DefinedLevel{sound: "sounds/battle_02.mp3", core: core}

	wave1 := &BasicWave{tickToLaunch: 1, core: level.GetCore()};
	wave1.addNewEnemy(DisplayableObject{marginTop: 0.05, marginLeft: 1.00, rotation: -90.0, image: "assets/ship1.png"});
	wave1.addNewEnemy(DisplayableObject{marginTop: 0.65, marginLeft: 1.20, rotation: -90.0, image: "assets/ship1.png"});
	wave1.addNewEnemy(DisplayableObject{marginTop: 0.45, marginLeft: 1.30, rotation: -90.0, image: "assets/ship1.png"});
	level.AddWave(wave1);

	wave2 := &BasicWave{tickToLaunch: 700, core: level.GetCore()};
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.08, marginLeft: 1.00, rotation: -90.0, image: "assets/ship1.png"});
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.75, marginLeft: 1.10, rotation: -90.0, image: "assets/ship1.png"});
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.45, marginLeft: 1.35, rotation: -90.0, image: "assets/ship1.png"});
	wave2.addNewEnemy(DisplayableObject{marginTop: 0.70, marginLeft: 1.35, rotation: -90.0, image: "assets/ship1.png"});
	level.AddWave(wave2);

	wave3 := &BasicWave{tickToLaunch: 1400, core: level.GetCore()};
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.13, marginLeft: 1.00, rotation: -90.0, image: "assets/ship1.png"});
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.27, marginLeft: 1.05, rotation: -90.0, image: "assets/ship1.png"});
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.13, marginLeft: 1.10, rotation: -90.0, image: "assets/ship1.png"});
	wave3.addNewEnemy(DisplayableObject{marginTop: 0.27, marginLeft: 1.15, rotation: -90.0, image: "assets/ship1.png"});
	level.AddWave(wave3);

	return level;
}

func main() {
    fmt.Println(" Initializing the game....");

	coreGame = &CoreGame{};
	coreGame.InitializeGame();

	for !rl.WindowShouldClose() && !coreGame.IsGameEnded() {
		coreGame.Tick();
	}

	rl.CloseWindow();
}
