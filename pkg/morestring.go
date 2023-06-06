package pkg

import "fmt"

func GetLastKString(k int, slice []string) []string {
	finalList := []string{}
	sliceLength := len(slice)

	if sliceLength < k {
		panic(fmt.Sprintf("slice with length %d is less than k value of %d", sliceLength, k))
	}

	startAppendingIdx := sliceLength - k

	for i, v := range slice {
		if i >= startAppendingIdx {
			finalList = append(finalList, v)
		}
	}

	return finalList
}

func GetFirstKString(k int, slice []string) []string {
	finalList := []string{}
	sliceLength := len(slice)

	if sliceLength < k {
		panic(fmt.Sprintf("slice with length %d is less than k value of %d", sliceLength, k))
	}

	endAppendingIdx := k - 1

	for i, v := range slice {
		if i <= endAppendingIdx {
			finalList = append(finalList, v)
		}
	}

	return finalList
}
