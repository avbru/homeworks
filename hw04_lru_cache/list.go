package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *item
	Back() *item
	PushFront(v interface{}) *item
	PushBack(v interface{}) *item
	Remove(i *item)
	MoveToFront(i *item)
}

type item struct {
	value interface{}
	next  *item
	prev  *item
}

type list struct {
	length    int
	firstNode *item
	lastNode  *item
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *item {
	return l.firstNode
}

func (l *list) Back() *item {
	return l.lastNode
}

func (l *list) PushBack(v interface{}) *item {
	newItem := &item{value: v, next: nil, prev: nil}

	if l.length == 0 {
		l.firstNode, l.lastNode = newItem, newItem
	} else {
		newItem.next = l.lastNode
		l.lastNode.prev = newItem
		l.lastNode = newItem
	}

	l.length++
	return newItem
}

func (l *list) PushFront(v interface{}) *item {
	newItem := &item{value: v, next: nil, prev: nil}

	if l.length == 0 {
		l.firstNode, l.lastNode = newItem, newItem
	} else {
		newItem.prev = l.firstNode
		l.firstNode.next = newItem
		l.firstNode = newItem
	}

	l.length++
	return newItem
}

func (l *list) Remove(i *item) {
	if i.prev != nil {
		i.prev.next = i.next
	} else {
		l.lastNode = i.next
	}

	if i.next != nil {
		i.next.prev = i.prev
	} else {
		l.firstNode = i.prev
	}

	l.length--
}

func (l *list) MoveToFront(i *item) {
	l.PushFront(i.value)
	l.Remove(i)
}
