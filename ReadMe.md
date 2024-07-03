# P1 Translating methods, interfaces and structural subtyping
##### by Josiah Heinen (61728), Vasiliki Konstanti (76282), Christian (55313)

---

### 1. Compare the run-time performance of RT and DT (e.g. call “sumArea” in a loop and measure which version runs faster)
In order to measure the performance of the two approaches, we created the following functions:
```
func iterationsRT(iterations int, r, s shape) {
	for i := 0; i < iterations; i++ {
		_ = sumArea(r, s)
	}
}
```

```
func iterationsDT(iterations int, rDictShape, sDictShape shape_Value) {
	for i := 0; i < iterations; i++ {
		_ = sumArea_Dict(rDictShape, sDictShape)
	}
}
```
In the main method, we pass this function as a parameter to another function:
```
func measureTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}
```
The results were as follows:

![task1](https://github.com/ChristianHauff42/model-based-sw/assets/102160452/65ed4f4d-ee6d-4ad1-bee6-bd04f440222f)

---

### 2. Apply the RT and DT approach to one further example of your own choice.

As the example we created a new shape "circle" with methods analogous to rectangle and square:
```
type circle struct {
	radius int
}

func (c circle) area() int {
	return int(math.Pi * float64(c.radius) * float64(c.radius))
}

func (c *circle) scale(x int) {
	c.radius = c.radius * x
}

func area_Circle(c circle) int {
	return int(math.Pi * float64(c.radius) * float64(c.radius))
}

func area_CirclePtr(c *circle) int {
	return area_Circle(*c)
}

func scale_CirclePtr(c *circle, x int) {
	c.radius = c.radius * x
}
```
As such we also expanded the functions ```area_Lookup``` and the corresponding test functions.
The tests yields following results:

![task2](https://github.com/ChristianHauff42/model-based-sw/assets/102160452/cff1f4c1-db23-466f-8618-bed08cf3cb25)

---

### 3. Extend RT and DT to deal with type assertions.
We added type assertions to the functions ```sumAreaVariant```, ```sumAreaVariant_Lookup``` and ```sumAreaVariant_Dict```:
```
func sumAreaVariant(x, y shape) int {
    z, ok := y.(square)
    if !ok {
        fmt.Println("Type assertion failed")
        return -1
    }
    return x.area() + y.area() + z.length
}
```

The expanded test functions yield the following results:

![task3](https://github.com/ChristianHauff42/model-based-sw/assets/102160452/130cefe4-e7f6-4574-98f1-5f0cba065cbe)

---

### 4. Extend RT and DT to deal with type bounds.
We added the following extensions:
```
type node[T any] struct {
    val  T
    next *node[T]
}

type Show interface {
    show() string
}

func showNode[T Show](n *node[T]) string {
    var s string
    for n != nil {
        s = s + n.val.show() + " -> "
        n = n.next
    }
    s = s + "nil"
    return s
}

func (r rectangle) show() string {
    return fmt.Sprintf("Rectangle(%d, %d)", r.length, r.width)
}

func (s square) show() string {
    return fmt.Sprintf("Square(%d)", s.length)
}

func (c circle) show() string {
    return fmt.Sprintf("Circle(%d)", c.radius)
}
```

The added tests yield the following results:

![task4](https://github.com/ChristianHauff42/model-based-sw/assets/102160452/1f5a4d80-5206-42d2-9ada-54e97cac2f21)

---
