package input

import (
	"bufio"
	"fmt"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"os"
	"strconv"
	"strings"
)

func IsWithParams() (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Would you like to enter search parameters?(Y/N): ")

	isWithParams, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	isWithParams = strings.TrimSpace(isWithParams)
	if isWithParams != "n" && isWithParams != "N" {
		return true, nil
	}

	return false, nil
}

func Title() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input title: ")

	title, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	title = strings.TrimSpace(title)

	return title, nil
}

func Author() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input author: ")

	author, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	author = strings.TrimSpace(author)

	return author, nil
}

func Publisher() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input publisher: ")

	publisher, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	publisher = strings.TrimSpace(publisher)

	return publisher, nil
}

func Rarity() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input rarity: ")

	rarity, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	rarity = strings.TrimSpace(rarity)

	return rarity, nil
}

func Genre() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input genre: ")

	genre, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	genre = strings.TrimSpace(genre)

	return genre, nil
}

func PublishingYear() (uint, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input publishing year: ")

	yearStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	yearStr = strings.TrimSpace(yearStr)

	yearInt, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, err
	}

	year := uint(yearInt)

	return year, nil
}

func Language() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input language: ")

	language, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	language = strings.TrimSpace(language)

	return language, nil
}

func AgeLimit() (uint, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input age limit: ")

	ageLimitStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	ageLimitStr = strings.TrimSpace(ageLimitStr)

	ageLimitInt, err := strconv.Atoi(ageLimitStr)
	if err != nil {
		return 0, err
	}

	ageLimit := uint(ageLimitInt)

	return ageLimit, nil
}

func CopiesNumber() (uint, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input book's copies number: ")

	copiesNumStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	copiesNumStr = strings.TrimSpace(copiesNumStr)

	copiesNumInt, err := strconv.Atoi(copiesNumStr)
	if err != nil {
		return 0, err
	}

	copiesNum := uint(copiesNumInt)

	return copiesNum, nil
}

func Params() (dto.BookParamsDTO, error) {
	var params dto.BookParamsDTO
	var err error

	params.Title, err = Title()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}
	params.Author, err = Author()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}
	params.Publisher, err = Publisher()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}
	params.Rarity, err = Rarity()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}
	params.Genre, err = Genre()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}
	params.PublishingYear, err = PublishingYear()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}
	params.Language, err = Language()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}
	params.AgeLimit, err = AgeLimit()
	if err != nil {
		return dto.BookParamsDTO{Limit: 5, Offset: 0}, err
	}

	return params, nil
}

func BookPagesNumber() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input book pages number: ")

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

func Book() (dto.BookDTO, error) {
	var book dto.BookDTO
	var err error

	book.Title, err = Title()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.Author, err = Author()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.Publisher, err = Publisher()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.CopiesNumber, err = CopiesNumber()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.Rarity, err = Rarity()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.Genre, err = Genre()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.PublishingYear, err = PublishingYear()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.Language, err = Language()
	if err != nil {
		return dto.BookDTO{}, err
	}
	book.AgeLimit, err = AgeLimit()
	if err != nil {
		return dto.BookDTO{}, err
	}

	return book, nil
}
