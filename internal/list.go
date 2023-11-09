package internal

func (l *List) PushFront(r *Room) {
	newNode := &Node{Room: r}
	if l.Front == nil {
		l.Len = 1
		l.Front = newNode
		l.Back = newNode
	} else {
		l.Len++
		newNode.Next = l.Front
		l.Front = newNode
	}
}

func (l *List) PushBack(r *Room) {
	newNode := &Node{Room: r}
	if l.Front == nil {
		l.Len = 1
		l.Front = newNode
		l.Back = newNode
	} else {
		l.Len++
		l.Back.Next = newNode
		l.Back = newNode
	}
}

func (l *List) RemoveFront() {
	if l.Front != nil {
		l.Len--
		l.Front = l.Front.Next
		if l.Front == nil {
			l.Back = nil
		}
	}
}

func (l *List) ToArray(lenArr int) []*Room {
	if l.Front == nil || lenArr < 1 {
		return nil
	}
	res := make([]*Room, lenArr)
	cur := l.Front
	for i := range res {
		res[i] = cur.Room
		cur = cur.Next
		if cur == nil {
			break
		}
	}
	return res
}
