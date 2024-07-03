package main

import (
	"fmt"
	"math"
	"time"
)

type rectangle struct {
	length int
	width  int
}

type square struct {
	length int
}

// new shape example (task 2)
type circle struct {
	radius int
}

func (r rectangle) area() int {
	return r.length * r.width
}

func (s square) area() int {
	return s.length * s.length
}

// new shape example (task 2)
func (c circle) area() int {
	return int(math.Pi * float64(c.radius) * float64(c.radius))
}

func (r *rectangle) scale(x int) {
	r.length = r.length * x
	r.width = r.width * x
}

func (s *square) scale(x int) {
	s.length = s.length * x
}

// new shape example (task 2)
func (c *circle) scale(x int) {
	c.radius = c.radius * x
}

type shape interface {
	area() int
}

type shapeExt interface {
	shape
	scale(int)
}

func sumArea(x, y shape) int {
	return x.area() + y.area()
}

func sumAreaScaleBefore(n int, x, y shapeExt) int {
	x.scale(n)
	y.scale(n)
	return x.area() + y.area()
}

func test() {
	var r rectangle = rectangle{1, 2}
	var s square = square{3}
	x1 := r.area() + s.area()
	fmt.Printf("%d \n", x1)
	x2 := sumArea(r, s)
	fmt.Printf("%d \n", x2)
	pt := &r
	x3 := pt.area()
	fmt.Printf("%d \n", x3)
	x4 := sumAreaScaleBefore(3, &r, &s)
	fmt.Printf("%d \n", x4)
}

func testNewShape() {
	var r rectangle = rectangle{1, 2}
	var c circle = circle{3}
	x1 := r.area() + c.area()
	fmt.Printf("%d \n", x1)
	x2 := sumArea(r, c)
	fmt.Printf("%d \n", x2)
	pt := &r
	x3 := pt.area()
	fmt.Printf("%d \n", x3)
	x4 := sumAreaScaleBefore(3, &r, &c)
	fmt.Printf("%d \n", x4)
}

// Introducing unique function names for overloaded methods

func area_Rec(r rectangle) int {
	return r.length * r.width
}

func area_Sq(s square) int {
	return s.length * s.length
}

// new shape example (task 2)
func area_Circle(c circle) int {
	return int(math.Pi * float64(c.radius) * float64(c.radius))
}

// "value" method implies "pointer" method
func area_RecPtr(r *rectangle) int {
	return area_Rec(*r)
}

func area_SqPtr(s *square) int {
	return area_Sq(*s)
}

// new shape example (task 2)
func area_CirclePtr(c *circle) int {
	return area_Circle(*c)
}

func scale_RecPtr(r *rectangle, x int) {
	r.length = r.length * x
	r.width = r.width * x
}

func scale_SqPtr(s *square, x int) {
	s.length = s.length * x
}

// new shape example (task 2)
func scale_CirclePtr(c *circle, x int) {
	c.radius = c.radius * x
}

// Run-time method lookup

func area_Lookup(x interface{}) int {
	var y int

	switch v := x.(type) {
	case square:
		y = area_Sq(v)
	case rectangle:
		y = area_Rec(v)
	// new shape example (task 2)
	case circle:
		y = area_Circle(v)
	}
	return y

}

func sumArea_Lookup(x, y interface{}) int {
	return area_Lookup(x) + area_Lookup(y)
}

// expanded with circle example (task 2)
func test_Lookup() {
	var r rectangle = rectangle{1, 2}
	var s square = square{3}
	var c circle = circle{3}
	x1 := area_Rec(r) + area_Sq(s) + area_Circle(c)
	fmt.Printf("%d \n", x1)
	x2 := sumArea_Lookup(r, s)
	// rectangle <= interface{}
	// square <= interface{}
	x3 := sumArea_Lookup(r, c)
	// rectangle <= interface{}
	// circle <= interface{}
	fmt.Printf("%d \n", x2)
	fmt.Printf("%d \n", x3)
}

// Dictionary translation

type shape_Value struct {
	val  interface{}
	area func(interface{}) int
}

type shapeExt_Value struct {
	val   interface{}
	area  func(interface{}) int
	scale func(interface{}, int)
}

// shapExt <= shape
func fromShapeExtToShape(x shapeExt_Value) shape_Value {
	return shape_Value{x.val, x.area}
}

func sumArea_Dict(x, y shape_Value) int {
	return x.area(x.val) + y.area(y.val)
}

func sumAreaScaleBefore_Dict(n int, x, y shapeExt_Value) int {
	x.scale(x.val, n)
	y.scale(y.val, n)
	return x.area(x.val) + y.area(y.val)
}

// expanded with circle example (task 2)
func test_Dict() {
	var r rectangle = rectangle{1, 2}
	var s square = square{3}
	var c circle = circle{3}

	// 1. Plain method calls
	x1 := area_Rec(r) + area_Sq(s) + area_Circle(c)
	fmt.Printf("%d \n", x1)
	x2 := sumArea(r, s)
	x3 := sumArea(r, c)
	fmt.Printf("%d \n", x2)
	fmt.Printf("%d \n", x3)
	pt := &r
	x4 := area_Rec(*pt)
	fmt.Printf("%d \n", x4)

	area_Rec_Wrapper := func(v interface{}) int {
		return area_Rec(v.(rectangle))
	}

	area_Sq_Wrapper := func(v interface{}) int {
		return area_Sq(v.(square))
	}

	// new shape example (task 2)
	area_Circle_Wrapper := func(v interface{}) int {
		return area_Circle(v.(circle))
	}

	rDictShape := shape_Value{r, area_Rec_Wrapper}
	sDictShape := shape_Value{s, area_Sq_Wrapper}
	// new shape example (task 2)
	cDictShape := shape_Value{c, area_Circle_Wrapper}

	x5 := sumArea_Dict(rDictShape, sDictShape)
	x6 := sumArea_Dict(rDictShape, cDictShape)
	fmt.Printf("%d \n", x5)
	fmt.Printf("%d \n", x6)

	area_RecPtr_Wrapper := func(v interface{}) int {
		return area_RecPtr(v.(*rectangle))
	}

	area_SqPtr_Wrapper := func(v interface{}) int {
		return area_SqPtr(v.(*square))
	}

	// new shape example (task 2)
	area_CirclePtr_Wrapper := func(v interface{}) int {
		return area_CirclePtr(v.(*circle))
	}

	scale_RecPtr_Wrapper := func(v interface{}, x int) {
		scale_RecPtr(v.(*rectangle), x)
	}

	scale_SqPtr_Wrapper := func(v interface{}, x int) {
		scale_SqPtr(v.(*square), x)
	}

	// new shape example (task 2)
	scale_CirclePtr_Wrapper := func(v interface{}, x int) {
		scale_CirclePtr(v.(*circle), x)
	}

	// Construct the appropriate interface values
	rDictShapeExt := shapeExt_Value{&r, area_RecPtr_Wrapper, scale_RecPtr_Wrapper}
	sDictShapeExt := shapeExt_Value{&s, area_SqPtr_Wrapper, scale_SqPtr_Wrapper}
	// new shape example (task 2)
	cDictShapeExt := shapeExt_Value{&c, area_CirclePtr_Wrapper, scale_CirclePtr_Wrapper}

	x7 := sumAreaScaleBefore_Dict(3, rDictShapeExt, sDictShapeExt)
	x8 := sumAreaScaleBefore_Dict(3, rDictShapeExt, cDictShapeExt)

	fmt.Printf("%d \n", x7)
	fmt.Printf("%d \n", x8)

	x9 := sumArea_Dict(fromShapeExtToShape(rDictShapeExt), fromShapeExtToShape(sDictShapeExt))
	x10 := sumArea_Dict(fromShapeExtToShape(rDictShapeExt), fromShapeExtToShape(cDictShapeExt))

	fmt.Printf("%d \n", x9)
	fmt.Printf("%d \n", x10)
}

func measureTime(fn func()) time.Duration {
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		fn()
	}
	return time.Since(start)
}

func iterationsRT(iterations int, r, s shape) {
	for i := 0; i < iterations; i++ {
		_ = sumArea(r, s)
	}
}

func iterationsDT(iterations int, rDictShape, sDictShape shape_Value) {
	for i := 0; i < iterations; i++ {
		_ = sumArea_Dict(rDictShape, sDictShape)
	}
}

func main() {

	//test()
	//testNewShape()
	//test_Lookup()
	//test_Dict()

	var r rectangle = rectangle{1, 2}
	var s square = square{3}

	/***** Measuring normal runtime calculation *****/
	rtTime := measureTime(func() { iterationsRT(1000, r, s) })
	fmt.Printf("rtTime: %v\n", rtTime)

	/***** Measuring Dictionary calculation *****/
	area_Rec_Wrapper := func(v interface{}) int {
		return area_Rec(v.(rectangle))
	}

	area_Sq_Wrapper := func(v interface{}) int {
		return area_Sq(v.(square))
	}

	//Directory Instances
	rDictShape := shape_Value{r, area_Rec_Wrapper}
	sDictShape := shape_Value{s, area_Sq_Wrapper}

	dtTime := measureTime(func() { iterationsDT(1000, rDictShape, sDictShape) })
	fmt.Printf("dtTime: %v\n", dtTime)

}
