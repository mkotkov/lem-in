package internal

type Result struct {
	AntsCount int
	Paths     []*List
}

type antFarm struct {
	FieldInfo  *FieldInfo
	AntsCount  int
	Start, End string
	Rooms      map[string]*Room
	StepsCount int
	Result     *Result
}

type Room struct {
	Name                string
	X, Y                int
	Paths               map[*Room]int
	ParentIn, ParentOut *Room
	VisitIn, VisitOut   bool
	Weight              [2]int
	Separated           bool
}

type FieldInfo struct {
	MODE             byte
	Start, End       bool
	IsStart, IsEnd   bool
	UsingCoordinates map[int]map[int]bool
}

type Node struct {
	Room *Room
	Next *Node
}

type List struct {
	Len   int
	Front *Node
	Back  *Node
}

type AntStruct struct {
	Num  int
	Path int
	Pos  int
	Next *AntStruct
}

type AntQueue struct {
	Front *AntStruct
	Back  *AntStruct
}

type WeightNode struct {
	Room   *Room
	Weight int
	Mark   bool
	Next   *WeightNode
}

type SortedQueue struct {
	Front *WeightNode
	Back  *WeightNode
}
