package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/scottjbarr/gameoflife-go"
	"github.com/wsxiaoys/terminal/color"
)

type Game struct {
	Generation int
	Rows       []Row
	SleepTime  int
}

func (g *Game) Get(col int, row int) *Cell {
	return &g.Rows[row].Cells[col]
}

type Row struct {
	Cells []Cell
}

type Cell struct {
	Col       int
	Row       int
	Value     int
	NextValue int
	Color     string
}

// offset values tht are used when finding neighbours
var offsets [8][]int

func init() {
	// offsets define the x (column), and y (row) offset from origin (0, 0)
	offsets[0] = []int{-1, -1}
	offsets[1] = []int{0, -1}
	offsets[2] = []int{1, -1}
	offsets[3] = []int{-1, 0}
	offsets[4] = []int{1, 0}
	offsets[5] = []int{-1, 1}
	offsets[6] = []int{0, 1}
	offsets[7] = []int{1, 1}
}

// Draw all Cell structures
func (game *Game) Draw() {
	clearScreen()

	fmt.Printf("Generation : %v\n", game.Generation)

	for _, r := range game.Rows {
		for _, c := range r.Cells {
			if c.Value == 1 {
				color.Print(fmt.Sprintf("%v ", c.Color))
			} else {
				color.Print("@c.")
			}
		}
		fmt.Printf("\n")
	}

	game.Sleep()
}

// Return the number of live neighbours a Cell has
func (g *Game) NeighbourCount(c *Cell) int {
	alive := 0

	for _, offset := range offsets {
		col := c.Col + offset[0]
		row := c.Row + offset[1]

		if col >= 0 && col < len(g.Rows[0].Cells) {
			if row >= 0 && row < len(g.Rows) {
				neighbour := g.Rows[row].Cells[col]
				if neighbour.IsAlive() {
					alive += 1
				}
			}
		}
	}

	return alive
}

func (g *Game) Sleep() {
	time.Sleep(time.Millisecond * time.Duration(g.SleepTime))
}

// Returns true if the Cell is currently alive, false otherwise
func (c *Cell) IsAlive() bool {
	return c.Value == 1
}

// Kill the Cell in the next generation
func (c *Cell) Die() {
	c.NextValue = 0
}

// Keep a Cell alive in the next generation
func (c *Cell) Live() {
	c.NextValue = 1

	// Use a green block
	c.Color = "@{G}"
}

// Bring the Cell back to life in the next generation
func (c *Cell) Spawn() {
	c.NextValue = 1

	// Use a white block
	c.Color = "@{W}"
}

// Copy next values of the Cell's to the current values.
func (g *Game) PrepareValues() {
	for row := 0; row < len(g.Rows); row++ {
		for col := 0; col < len(g.Rows[0].Cells); col++ {
			c := g.Get(col, row)
			c.Value = c.NextValue
		}
	}
}

func NewGame(filename string, sleepTime int) *Game {

	// load the game data
	data, err := loader.ReadLifeData(filename)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	rowCount := len(data)

	// colCount is the maximum string length of a row of data
	maxLength := 0

	for _, line := range data {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}

	colCount := maxLength

	rows := make([]Row, rowCount)

	for row := 0; row < len(rows); row++ {
		rows[row].Cells = make([]Cell, colCount)

		for col := 0; col < colCount; col++ {
			c := *NewCell(col, row)

			if row < len(data) && col < len(data[row]) {
				if string(data[row][col]) == "*" {
					// this is a live cell
					c.Live()
					c.Value = 1
				}
			}

			rows[row].Cells[col] = c
		}
	}

	game := &Game{Rows: rows, SleepTime: sleepTime}

	return game
}

// Create a new Cell
func NewCell(col int, row int) *Cell {
	return &Cell{Value: 0, Col: col, Row: row}
}

// Update all Cells
func (game *Game) Tick() {
	game.Generation += 1

	for row := 0; row < len(game.Rows); row++ {
		for col := 0; col < len(game.Rows[0].Cells); col++ {

			c := &game.Rows[row].Cells[col]

			n := game.NeighbourCount(c)

			if c.IsAlive() {
				if n < 2 {
					// live cell with less than 2 live neighbour
					c.Die()
				} else if n == 2 || n == 3 {
					// live cell with 2 or 3 live neighbours
					c.Live()
				} else if n > 3 {
					// live cell with 3+ live neighbours
					c.Die()
				}
			} else if !c.IsAlive() && n == 3 {
				// dead cell with 3 neighbours
				c.Spawn()
			}
		}
	}

	game.PrepareValues()
}

// Clear the terminal
func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func main() {
	iterations := flag.Int("iterations", 100, "Iterations")
	sleepTime := flag.Int("sleep", 100, "Pause between frames (ms)")
	filename := flag.String("file", "", "Name of life data file")

	flag.Parse()

	game := NewGame(*filename, *sleepTime)

	for i := 0; i < *iterations; i++ {
		game.Draw()
		game.Tick()
	}
}
