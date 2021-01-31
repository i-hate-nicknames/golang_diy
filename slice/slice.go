package slice

type IntSlice interface {
	Get(index int) (int, error)
	Set(index, val int) error
	Slice(start, end int) error
	Len() int
	Cap() int
	Append(other IntSlice) IntSlice
	Copy() IntSlice
}
