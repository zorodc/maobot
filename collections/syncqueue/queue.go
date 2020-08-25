/* Implements a synchronized queue. */

package syncqueue

import "sync"

type Queue struct {
	sync.Mutex
	items []interface{}
	idx     uint // start index
}

// Return the number of elements in the queue.
func (this *Queue) Count() uint {
	this.Lock()
	defer this.Unlock()

	return uint(len(this.items)) - this.idx
}

func (this *Queue) Clear() {
	this.Lock()
	defer this.Unlock()

	this.items = nil
}

func (this *Queue) Empty() bool {
	return this.Count() == 0
}

// Get the front (leftmost) item, or nil if the queue is empty.
func (this *Queue) Front() interface{} {
	if this.Empty() {
		return nil
	}

	this.Lock()
	defer this.Unlock()
	
	return this.items[this.idx]
}

// Get the back (rightmost) item, or nil if the queue is empty.
func (this *Queue) Back() interface{} {
	if this.Empty() {
		return nil
	}

	this.Lock()
	defer this.Unlock()

	return this.items[len(this.items)-1]
}

// Prepends the last popped item.
func (this *Queue) Rollback() {
	this.Lock()
	defer this.Unlock()

	if this.idx > 0 {
		this.idx--
	}
}

// Push an item to the front (left) of the queue.
func (this *Queue) Prepend(elem interface{}) {
	this.Lock()
	defer this.Unlock()

	if this.idx > 0 {
		this.idx--
		this.items[this.idx] = elem
	} else {
		// Perform a CONS operation. Expensive.
		this.items = append([]interface{}{elem}, this.items)
	}
}

// Push an item to the back (right) of the queue.
func (this *Queue) Append(elem interface{}) {
	this.Lock()
	defer this.Unlock()

	this.items = append(this.items, elem)
}

// Pop an item from the front (left) of the queue, or nil if the queue is empty.
func (this *Queue) PopFront() interface{} {
	if this.Empty() {
		return nil
	}

	this.Lock()
	defer this.Unlock()

	old_idx := this.idx
	this.idx++
	return this.items[old_idx]
}

// Pop an item from the back (right) of the queue, or nil if the queue is empty.
func (this *Queue) PopBack() interface{} {
	if this.Empty() {
		return nil
	}

	this.Lock()
	defer this.Unlock()

	old_items := this.items
	// Cut off the last element.
	this.items = this.items[0:len(this.items)-1]
	return old_items[len(old_items)-1]
}
