package main

import (
	"testing"
)

func TestInitUser(t *testing.T) {
	var (
		 userMap userRecords
		 sampleUserID  = "sampleID"
		 samplePermutation = []int{1,2,3}
		 startingCurrentPos = 0
		 startPermCount = 0
	)
	userMap = InitMemMap().(userRecords)

	userMap.InitUser(sampleUserID, samplePermutation)

	user, ok := userMap.m[sampleUserID]
	if !ok {
		t.Errorf("user with id: %s was not created", sampleUserID)
		t.FailNow()
	}
	for i := range samplePermutation {
		if samplePermutation[i] != user.permutation[i] {
			t.Errorf("user's permutation is invalid. expected %v, got %v", samplePermutation, user.permutation)
			break
		}
	}
	if user.currentPos != startingCurrentPos {
		t.Errorf("user's currentPos is invalid. expected %v, got %v",startingCurrentPos, user.currentPos)
	}
	if user.permsCount != startPermCount {
		t.Errorf("user's permsCount is invalid. expected %v, got %v",startPermCount, user.permsCount)
	}
}

func TestUserExists(t *testing.T) {
	var (
		userMap userRecords
		userID = "sampleID"
		user = userRecord{permutation:[]int{}, currentPos:0, permsCount:0}
	)

	userMap = InitMemMap().(userRecords)
	userMap.m[userID] = user

	if !userMap.UserExists(userID) {
		t.Errorf("user with id %s has to exist, but it isnt", userID)
	}
}

func TestNextForUser(t *testing.T) {
	var (
		userMap userRecords

		userID = "sampleID"

		initPerm = []int{1,2}
		initPermLengthFactorial = 2

		firstIterationPerm = []int{1,2}
		secondIterationPerm = []int{2,1}

		user = userRecord{permutation:initPerm, currentPos:0, permsCount:0}
	)

	userMap = InitMemMap().(userRecords)
	_, ok := userMap.NextForUser(userID)
	if ok {
		t.Errorf("user with id %s has to not exist, but it exists", userID)
		t.FailNow()
	}

	userMap.m[userID] = user

	for i := 0; i < initPermLengthFactorial+1; i++ {
		out, ok := userMap.NextForUser(userID)
		if !ok {
			t.Errorf("user with id %s has to exist, but it isnt", userID)
			t.FailNow()
		}
		switch i {
		case 0:
			if len(initPerm) != len(out) {
				t.Errorf("permutation changed lenght, expected %d, got %d", len(initPerm), len(out))
				t.FailNow()
			}
			for j := range firstIterationPerm {
				if firstIterationPerm[j] != out[j] {
					t.Errorf("invalid permutation value, expected %v, got %v", firstIterationPerm, out)
				}
			}
		case 1:
			if len(initPerm) != len(out) {
				t.Errorf("permutation changed lenght, expected %d, got %d", len(initPerm), len(out))
				t.FailNow()
			}
			for j := range secondIterationPerm {
				if secondIterationPerm[j] != out[j] {
					t.Errorf("invalid permutation value, expected %v, got %v", secondIterationPerm, out)
				}
			}
		case initPermLengthFactorial:
			if len(out) != 0 {
				t.Errorf("invalid expired permutation lenght, expected 0, got %v", len(out))
			}
		}
	}
}

func TestFactorial(t *testing.T) {
	var (
		initials = []int{0, 1, 2, 3, 4}
		factorials = []int{1, 1, 2, 6, 24}
	)

	for i := range initials {
		if f := factorial(initials[i]); f != factorials[i] {
			t.Errorf("failed factorial for %d: expected: %d, got %d", initials[i], factorials[i], f)
		}
	}
}



func TestCopySlice(t *testing.T) {
	var (
		initialSlice = []int{1,2,3,4,5}
		finalSlice = []int{1,2,3,4,5}
	)

	testedSlice := copySlice(initialSlice)

	for i := range initialSlice {
		if testedSlice[i] != initialSlice[i] {
			t.Errorf("initial slice: %v, expected: %v, got: %v", initialSlice, finalSlice, testedSlice)
			break
		}
	}

	initialSlice[0] = 0

	for i := range finalSlice {
		if testedSlice[i] != finalSlice[i] {
			t.Error("slice copy was failed, still poiting to initial slice")
			break
		}
	}
}

func TestSwapWithNext(t *testing.T) {
	var (
		initial = []int{1,2}
		initialSingle = []int{1}
		initialEmpty = []int{}
		bigNumber = len(initial) + 1
		lastPos = len(initial) - 1

		swapped = []int{2,1}
		swapBig = []int{2,1}
		swapSingle = []int{1}
		swapEmpty = []int{}
	)

	worked := []int{1,2}
	position := 0
	newPosition := swapWithNext(worked, position)
	if newPosition != position +1 {
		t.Errorf("swap failed, expected new position: %v, got %v", position + 1, newPosition)
	}
	for i := range worked {
		if worked[i] != swapped[i] {
			t.Errorf("swap failed, expected: %v, got %v", swapped, worked)
			break
		}
	}

	worked = []int{1,2}
	position = lastPos
	newPosition = swapWithNext(worked, position)
	if newPosition != lastPos {
		t.Errorf("swap failed, expected new position: %v, got %v", lastPos, newPosition)
	}

	for i := range worked {
		if worked[i] != swapped[i] {
			t.Errorf("swap failed, expected: %v, got %v", swapped, worked)
			break
		}
	}

	// check swap if `position` bigger than slice len
	worked = []int{1,2}
	swapWithNext(worked, bigNumber)
	for i := range worked {
		if worked[i] != swapBig[i] {
			t.Errorf("swap failed, expected: %v, got %v", swapBig, worked)
			break
		}
	}

	// check swaps for slices with len 1
	swapWithNext(initialSingle, 0)
	for i := range initialSingle {
		if swapSingle[i] != swapSingle[i] {
			t.Errorf("swap failed, expected: %v, got %v", swapSingle, initialSingle)
			break
		}
	}

	// check swaps with slices with len 0
	swapWithNext(initialEmpty, 0)
	if len(initialEmpty) != 0 {
			t.Errorf("swap failed, expected: %v, got %v", swapEmpty, initialEmpty)
		}

}