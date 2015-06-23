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
	Tick int
	Rows []Row
}

func (g *Game) Get(col int, row int) Cell {
	return g.Rows[row].Cells[col]
}

type Row struct {
	Cells []Cell
}

type Cell struct {
	Col       int
	Row       int
	Value     int
	NextValue int
}

// neighbour offsets
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

func (game *Game) Draw() {
	fmt.Println("game.tick : ", game.Tick)

	for _, r := range game.Rows {
		for _, c := range r.Cells {
			if c.Value == 1 {
				color.Print("@{G} ")
			} else {
				color.Print("@c.")
			}
		}
		fmt.Printf("\n")
	}
}

func (g *Game) NeighbourCount(c *Cell) int {
	alive := 0

	// check if cell to left is alive
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

func (c *Cell) IsAlive() bool {
	return c.Value == 1
}

func (c *Cell) Die() {
	c.NextValue = 0
}

func (c *Cell) Live() {
	c.NextValue = 1
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

	// The glider
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

func NewCell(col int, row int) *Cell {
	return &Cell{Value: 0, Col: col, Row: row}
}

func Tick(game *Game) {
	game.Tick += 1

	// update all rows
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
				c.Live()
			}
		}
	}

	// copy new value to value, and reset new value
	for row := 0; row < len(game.Rows); row++ {
		for col := 0; col < len(game.Rows[0].Cells); col++ {
			c := &game.Rows[row].Cells[col]
			c.Value = c.NextValue
		}
	}
}

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func sleep(sleepTime int) {
	time.Sleep(time.Millisecond * 200)
}

func main() {
	iterations := flag.Int("iterations", 10, "Iterations")
	rowCount := flag.Int("rows", 4, "Grid rows")
	colCount := flag.Int("cols", 4, "Grid columns")
	sleepTime := flag.Int("sleep", 100, "Pause between frames (ms)")

	flag.Parse()

	// load data file, pass data file to NewGame
	game := NewGame(*rowCount, *colCount)

	for i := 0; i < *iterations; i++ {
		clearScreen()
		game.Draw()
		sleep(*sleepTime)
		Tick(game)
		fmt.Printf("\n")
	}
}
