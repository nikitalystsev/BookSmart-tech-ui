package input

import (
	"bufio"
	"fmt"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"os"
	"strconv"
	"strings"
)

func Review() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input review: ")

	review, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	review = strings.TrimSpace(review)

	return review, nil
}

func Rating() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input rating: ")

	ratingStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	ratingStr = strings.TrimSpace(ratingStr)

	ratingInt, err := strconv.Atoi(ratingStr)
	if err != nil {
		return 0, err
	}

	return ratingInt, nil
}

func RatingParams() (dto.RatingInputDTO, error) {
	var (
		ratingDTO dto.RatingInputDTO
		err       error
	)

	if ratingDTO.Review, err = Review(); err != nil {
		return dto.RatingInputDTO{}, err
	}
	if ratingDTO.Rating, err = Rating(); err != nil {
		return dto.RatingInputDTO{}, err
	}

	return ratingDTO, nil
}
