package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReservationNumber() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input reservation number: ")

	numStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	numStr = strings.TrimSpace(numStr)

	numInt, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return numInt, nil
}
