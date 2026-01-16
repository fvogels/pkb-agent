package uid

var nextID = 0

func Generate() int {
	nextID++
	return nextID
}
