package main

import "fmt"

// Example

type rectangle struct {
	length int
	width  int
}

type square struct {
	length int
}

func (r rectangle) area() int {
	return r.length * r.width
}

func (s square) area() int {
	return s.length * s.length
}

func (r *rectangle) scale(x int) {
	r.length = r.length * x
	r.width = r.width * x
}

func (s *square) scale(x int) {
	s.length = s.length * x
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
	// Implicit conversion to
	// (*pt).area()
	//
	// Hence, any "value" receiver also implies the corresponding "pointer" receiver.

	fmt.Printf("%d \n", x3)

	//  x3 := sumAreaScaleBefore(3, r, s)
	//
	// "rectangle does not implement shapeExt (scale method has pointer receiver)"
	// same applies to square

	x4 := sumAreaScaleBefore(3, &r, &s)

	fmt.Printf("%d \n", x4)

}

// Introducing unique function names for overloaded methods

func area_Rec(r rectangle) int {
	return r.length * r.width
}

func area_Sq(s square) int {
	return s.length * s.length
}

// "value" method implies "pointer" method
func area_RecPtr(r *rectangle) int {
	return area_Rec(*r)
}

func area_SqPtr(s *square) int {
	return area_Sq(*s)
}

func scale_RecPtr(r *rectangle, x int) {
	r.length = r.length * x
	r.width = r.width * x
}

func scale_SqPtr(s *square, x int) {
	s.length = s.length * x
}

// Run-time method lookup

func area_Lookup(x interface{}) int {
	var y int

	switch v := x.(type) {
	case square:
		y = area_Sq(v)
	case rectangle:
		y = area_Rec(v)
	}
	return y

}

func sumArea_Lookup(x, y interface{}) int {
	return area_Lookup(x) + area_Lookup(y)
}

func test_Lookup() {

	var r rectangle = rectangle{1, 2}
	var s square = square{3}

	x1 := area_Rec(r) + area_Sq(s)

	fmt.Printf("%d \n", x1)

	x2 := sumArea_Lookup(r, s)
	// rectangle <= interface{}
	// square <= interface{}

	fmt.Printf("%d \n", x2)
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

func test_Dict() {
	var r rectangle = rectangle{1, 2}
	var s square = square{3}

	// 1. Plain method calls

	x1 := area_Rec(r) + area_Sq(s)

	fmt.Printf("%d \n", x1)

	x2 := sumArea(r, s)

	fmt.Printf("%d \n", x2)

	pt := &r

	x3 := area_Rec(*pt)
	// Implicit conversion from pointer to value

	fmt.Printf("%d \n", x3)

	//  x3 := sumAreaScaleBefore(3, r, s)
	//
	// "rectangle does not implement shapeExt (scale method has pointer receiver)"
	// same applies to square

	// 2. Calling sumArea

	// Wrapper functions are needed for the following reason.
	// (a) area_Rec has type func(rectangle) int
	// (b) We need to store area_Rec in the "area" dictionary entry which has type func(interface{}) int
	// (c) We cast area_Rec to the approrpriate type

	area_Rec_Wrapper := func(v interface{}) int {
		return area_Rec(v.(rectangle))

	}

	area_Sq_Wrapper := func(v interface{}) int {
		return area_Sq(v.(square))

	}

	rDictShape := shape_Value{r, area_Rec_Wrapper}

	sDictShape := shape_Value{s, area_Sq_Wrapper}

	x4 := sumArea_Dict(rDictShape, sDictShape)

	fmt.Printf("%d \n", x4)

	// 3. Calling sumAreaScaleBefore

	area_RecPtr_Wrapper := func(v interface{}) int {
		return area_RecPtr(v.(*rectangle))

	}

	area_SqPtr_Wrapper := func(v interface{}) int {
		return area_SqPtr(v.(*square))

	}

	scale_RecPtr_Wrapper := func(v interface{}, x int) {
		scale_RecPtr(v.(*rectangle), x)

	}

	scale_SqPtr_Wrapper := func(v interface{}, x int) {
		scale_SqPtr(v.(*square), x)

	}

	// Construct the appropriate interface values
	rDictShapeExt := shapeExt_Value{&r, area_RecPtr_Wrapper, scale_RecPtr_Wrapper}

	sDictShapeExt := shapeExt_Value{&s, area_SqPtr_Wrapper, scale_SqPtr_Wrapper}

	x5 := sumAreaScaleBefore_Dict(3, rDictShapeExt, sDictShapeExt)

	fmt.Printf("%d \n", x5)

	// 4. Calling sumArea with a shapeExt value

	x6 := sumArea_Dict(fromShapeExtToShape(rDictShapeExt), fromShapeExtToShape(sDictShapeExt))

	fmt.Printf("%d \n", x6)
}

func main() {

	test()

	test_Lookup()

	test_Dict()

}
