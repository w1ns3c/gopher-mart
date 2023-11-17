package utils

import (
	"errors"
	"fmt"
	"strconv"
)

// test data 220070014780441 | 2
var (
	errCovertion error = errors.New("can't convert to int")
)

func LuhnValidator(sequence string) bool {
	l := len(sequence)
	checkSign, err := strconv.Atoi(sequence[l-1:])
	if err != nil {
		fmt.Println(errCovertion)
		return false
	}

	sum := 0
	for i := l - 2; i >= 0; i-- {
		number, err := strconv.Atoi(string(sequence[i]))
		if err != nil {
			fmt.Println(errCovertion)
			return false
		}

		if (l-i)%2 == 0 {
			number *= 2
			number = number%10 + number/10
		}
		sum += number
	}

	return (sum+checkSign)%10 == 0
}
