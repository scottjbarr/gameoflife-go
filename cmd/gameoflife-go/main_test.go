package main

import (
	"reflect"
	"testing"
)

// Test helper. Thanks again, @keighl
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

// Create a Game that has a Glider in it.
//
//     ..*
//     *.*
//     .**
//
func buildGlider() *Game {
	g := NewGame(4, 4, 50)

	g.Get(2, 0).Live()

	g.Get(0, 1).Live()
	g.Get(2, 1).Live()

	g.Get(1, 2).Live()
	g.Get(2, 2).Live()

	g.PrepareValues()

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

	expect(t, g.Rows[2].Cells[0], *g.Get(0, 2))
}

func TestNeighbourCount1(t *testing.T) {
	g := buildGlider()
	c := g.Get(0, 2)

	expect(t, 2, g.NeighbourCount(c))
}

func TestNeighbourCount2(t *testing.T) {
	g := buildGlider()
	c := g.Get(2, 0)

	expect(t, 1, g.NeighbourCount(c))
}
