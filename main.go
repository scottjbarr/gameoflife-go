package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/wsxiaoys/terminal/color"
)

type Game struct {
	Generation int
	Rows       []Row
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
	fmt.Println("Generation : ", game.Generation)

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

func NewGame(rowCount int, colCount int) *Game {
	rows := make([]Row, rowCount)

	for row := 0; row < len(rows); row++ {
		rows[row].Cells = make([]Cell, colCount)

		for col := 0; col < colCount; col++ {
			rows[row].Cells[col] = *NewCell(col, row)
		}
	}

	game := &Game{Rows: rows}

	// Hard coding a Glider until I add the loader
	//
	//     ..*
	//     *.*
	//     .**

	// row 0
	rows[0].Cells[2].Value = 1

	// row 1
	rows[1].Cells[0].Value = 1
	rows[1].Cells[2].Value = 1

	// row 2
	rows[2].Cells[1].Value = 1
	rows[2].Cells[2].Value = 1

	return game
}

// Create a new Cell
func NewCell(col int, row int) *Cell {
	return &Cell{Value: 0, Col: col, Row: row}
}

// Update all Cells
func (game *Game) Tick() { // game *Game) {
	game.Generation += 1

	for row := 0; row < len(game.Rows); row++ {
		for col := 0; col < len(game.Rows[0].Cells); col++ {

			c := &game.Rows[row].Cells[col]

			n := game.NeighbourCount(c)

			if c.IsAlive() {
				if n < 2 {
					// live cell with less than 2 live neighbour
					c.Die() // NextValue = 0
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

// Clear the terminal.
func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

// Sleep for the given number of seconds.
func sleep(seconds int) {
	time.Sleep(time.Millisecond * time.Duration(seconds))
}

func main() {
	colCount := flag.Int("cols", 40, "Grid columns")
	rowCount := flag.Int("rows", 20, "Grid rows")
	iterations := flag.Int("iterations", 70, "Iterations")
	sleepTime := flag.Int("sleep", 100, "Pause between frames (ms)")

	flag.Parse()

	game := NewGame(*rowCount, *colCount)

	for i := 0; i < *iterations; i++ {
		clearScreen()
		game.Draw()
		sleep(*sleepTime)
		game.Tick()
		fmt.Printf("\n")
	}
}
