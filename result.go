package ro2dtree

import (
	"sync/atomic"
)

var EmptyResult = newResult(nil, 0)

//todo stats
type ResultPool struct {
	misses   int64
	capacity int
	list     chan *Result
}

func newResultPool(count, capacity int) *ResultPool {
	pool := &ResultPool{
		capacity: capacity,
		list:     make(chan *Result, count),
	}
	for i := 0; i < cap(pool.list); i++ {
		pool.list <- newResult(pool, capacity)
	}
	return pool
}

func (pool *ResultPool) Checkout() *Result {
	select {
	case result := <-pool.list:
		return result
	default:
		atomic.AddInt64(&pool.misses, 1)
		return newResult(nil, pool.capacity)
	}
}

type Result struct {
	target   Point
	position int
	polygons Polygons
	groupMap map[int]int
	pool     *ResultPool
}

func newResult(pool *ResultPool, capacity int) *Result {
	return &Result{
		pool:     pool,
		polygons: make(Polygons, capacity),
		groupMap: make(map[int]int),
	}
}

func (r *Result) Add(polygon Polygon) bool {
	r.polygons[r.position] = polygon
	r.position++
	return r.position != len(r.polygons)
}

func (r *Result) AddUniqueByGroup(polygon Polygon) bool {
	oldPosition, present := r.groupMap[polygon.GroupId()]
	if present {
		oldPolygon := r.polygons[oldPosition]
		if r.Rank(oldPolygon) > r.Rank(polygon) {
			r.polygons[oldPosition] = polygon
		}
	} else {
		r.groupMap[polygon.GroupId()] = r.position
		r.polygons[r.position] = polygon
		r.position++
	}
	return r.position != len(r.polygons)
}

func (r *Result) Polygons() Polygons {
	return r.polygons[:r.position]
}

func (r *Result) Close() {
	if r.pool != nil {
		r.position = 0
		r.pool.list <- r
	}
}

func (r *Result) Len() int {
	return r.position
}

func (r *Result) Less(i, j int) bool {
	return r.Rank(r.polygons[i]) < r.Rank(r.polygons[j])
}

func (r *Result) Swap(i, j int) {
	r.polygons[i], r.polygons[j] = r.polygons[j], r.polygons[i]
}

func (r *Result) Rank(polygon Polygon) float64 {
	return polygon.Centroid().DistanceTo(r.target)
}
