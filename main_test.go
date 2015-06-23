package main

import (
	// "fmt"
	"reflect"
	"testing"
)

// Test helper. Thanks again, @keighl
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// ***.
// *...
// .*..
// ....
func buildGlider() *Game {
	g := NewGame(3, 3)
	// rows := g.Rows

	// row 0
	g.Rows[0].Cells[0].Value = 1
	g.Rows[0].Cells[1].Value = 1
	g.Rows[0].Cells[2].Value = 1

	// row 1
	g.Rows[1].Cells[0].Value = 1

	// row 2
	g.Rows[2].Cells[1].Value = 1

	return g
}

func TestIsAlive(t *testing.T) {
	g := buildGlider()

	c := g.Get(2, 0)
	expect(t, true, c.IsAlive())
}

func TestIsAliveFalse(t *testing.T) {
	g := buildGlider()

	c := g.Get(0, 2)
	expect(t, false, c.IsAlive())
}

func TestGameGet(t *testing.T) {
	g := buildGlider()
	expect(t, g.Rows[2].Cells[0], g.Get(0, 2))
}

func TestNeighbourCount1(t *testing.T) {
	g := buildGlider()
	g.Draw()
	c := g.Get(0, 2)
	expect(t, 2, g.NeighbourCount(&c))
}

func TestNeighbourCount2(t *testing.T) {
	g := buildGlider()

	c := g.Get(2, 0)
	expect(t, 1, g.NeighbourCount(&c))
}
