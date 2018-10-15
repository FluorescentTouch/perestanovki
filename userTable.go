package main

import "sync"

type(
	userRecords struct {
		m map[string]userRecord
		*sync.RWMutex
	}
	userRecord struct {
		permutation []int
		currentPos int
		permsCount int
	}
)

func InitMemMap() UserRecords {
	return userRecords{
		make(map[string]userRecord),
			new(sync.RWMutex),
	}
}

// InitUser will init new session for user
// will rewrite session if already exists
func (u userRecords) InitUser(userID string, perm []int) {
	u.Lock()
	defer u.Unlock()
	u.m[userID] = userRecord{
		permutation: perm,
		currentPos: 0,
		permsCount: 0,
	}
}

// UserExists checks if user exists in map
func (u userRecords) UserExists(userID string) bool {
	u.RLock()
	defer u.RUnlock()

	_, ok := u.m[userID]
	return ok
}

// NextForUser returns next permutation for user
// second argument replies if user exists in table
func (u userRecords) NextForUser(userID string) ([]int, bool) {
	u.Lock()
	defer u.Unlock()

	ur, ok := u.m[userID]
	if !ok {
		return []int{}, false
	}

	if ur.permsCount == factorial(len(ur.permutation)) {
		return []int{}, true
	}

	if ur.permsCount != 0 {
		newPos := swapWithNext(ur.permutation, ur.currentPos)
		ur.currentPos = newPos
	}

	ur.permsCount++
	u.m[userID] = ur

	outSlice := copySlice(ur.permutation)
	return outSlice, true
}

func factorial(num int) int {
	if num != 0 {
		return num * factorial(num-1)
	}
	return 1
}

func swapWithNext(slice []int, pos int) (newPosition int) {
	if len(slice) > 1 {
		if pos >= len(slice) - 1 {
			pos = 0
		}
		slice[pos], slice[pos+1] = slice[pos+1], slice[pos]
	}
	pos++
	return pos
}

func copySlice(slice []int) []int{
	out := make([]int, len(slice), len(slice))

	for i := range slice {
		out[i] = slice[i]
	}
	return out
}