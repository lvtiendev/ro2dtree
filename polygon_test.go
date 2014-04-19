package ro2dtree

import (
	"testing"
)

func TestPolygonHasAPropertyBoundingBox(t *testing.T) {
	points := Points{
		NewPoint(0, 0),
		NewPoint(22, 10),
		NewPoint(10, 15),
		NewPoint(19, 20),
		NewPoint(1, 20),
		NewPoint(0, 0),
	}
	box := NewPolygon(points).Box()
	expectPoint(t, box.TopLeft, 0, 0)
	expectPoint(t, box.BottomRight, 22, 20)
}

func TestPolygonHasProperCentroid(t *testing.T) {
	points := Points{
		NewPoint(0, 0),
		NewPoint(10, 10),
		NewPoint(5, 8),
		NewPoint(0, 0),
	}
	p := NewPolygon(points)
	expected := NewPoint(5, 6)
	if p.centroid != expected {
		t.Errorf("Expecting a centroid of %v got %v", expected, p.centroid)
	}
}

func TestPolygonHitWhenOutside(t *testing.T) {
	p := buildRectangle(0, 0, 10, 10)
	expectFalse(t, p.Hit(0, 0))
	expectFalse(t, p.Hit(10, 10))
	expectFalse(t, p.Hit(0, 10))
	expectFalse(t, p.Hit(10, 0))
	expectFalse(t, p.Hit(11, 5))
	expectFalse(t, p.Hit(5, 11))
	expectFalse(t, p.Hit(25, 25))
}

func TestPolygonHitWhenInside(t *testing.T) {
	p := buildRectangle(0, 0, 10, 10)
	expectTrue(t, p.Hit(1, 1))
	expectTrue(t, p.Hit(9, 4))
}

func buildRectangle(x1, y1, x2, y2 float64) *SimplePolygon {
	points := Points{
		NewPoint(x1, y1),
		NewPoint(x2, y1),
		NewPoint(x2, y2),
		NewPoint(x1, y2),
		NewPoint(x1, y1),
	}
	return NewPolygon(points)
}

func expectFalse(t *testing.T, actual bool) {
	if actual {
		t.Error("Expecting false, got true")
	}
}

func expectTrue(t *testing.T, actual bool) {
	if !actual {
		t.Error("Expecting true, got false")
	}
}

func expectPoint(t *testing.T, point Point, x, y float64) {
	if point.X != x {
		t.Errorf("expecting point.X to be %v got %v", x, point.X)
	}
	if point.Y != y {
		t.Errorf("expecting point.Y to be %v got %v", y, point.Y)
	}
}