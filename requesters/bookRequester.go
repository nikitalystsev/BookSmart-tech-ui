package requesters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-tech-ui/input"
	jsonmodels "github.com/nikitalystsev/BookSmart-web-api/core/models"
	"net/http"
	"time"
)

const catalogMenu = `Catalog's menu:
	1 -- view books
	2 -- next page
	3 -- view info about book
	4 -- add book to favorites
	5 -- reserve book
	6 -- view book ratings
	7 -- add book rating 
	0 -- go to main menu
`

const (
	pageLimit = 10

	booksKey      = "books"
	bookParamsKey = "bookParams"
)

func (r *Requester) ProcessBookCatalogActions() error {
	var (
		menuItem int
		err      error
	)

	r.cache.Set(bookParamsKey, dto.BookParamsDTO{Limit: pageLimit, Offset: 0})
	r.cache.Set(booksKey, make([]uuid.UUID, 0))

	for {
		fmt.Printf("\n\n%s", catalogMenu)

		if menuItem, err = input.MenuItem(); err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
			continue
		}

		switch menuItem {
		case 1:
			if err = r.viewFirstPage(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.viewNextPage(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 3:
			if err = r.ViewBook(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 4:
			if err = r.AddToFavorites(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 5:
			if err = r.ReserveBook(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 6:
			if err = r.viewBookRatings(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 7:
			if err = r.addNewBookRating(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
func (r *Requester) viewFirstPage() error {
	var bookParams dto.BookParamsDTO
	var bookPagesID []uuid.UUID

	isWithParams, err := input.IsWithParams()
	if err != nil {
		return err
	}

	if isWithParams {
		if bookParams, err = input.Params(); err != nil {
			return err
		}
	}

	bookParams.Limit = pageLimit
	bookParams.Offset = 0
	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + "/books",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		QueryParams: map[string]string{
			"title":           bookParams.Title,
			"author":          bookParams.Author,
			"publisher":       bookParams.Publisher,
			"copies_number":   fmt.Sprintf("%d", bookParams.CopiesNumber),
			"rarity":          bookParams.Rarity,
			"genre":           bookParams.Genre,
			"publishing_year": fmt.Sprintf("%d", bookParams.PublishingYear),
			"language":        bookParams.Language,
			"age_limit":       fmt.Sprintf("%d", bookParams.AgeLimit),
			"limit":           fmt.Sprintf("%d", bookParams.Limit),
			"offset":          fmt.Sprintf("%d", bookParams.Offset),
		},
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	var books []*jsonmodels.BookModel
	if err = json.Unmarshal(response.Body, &books); err != nil {
		return err
	}

	printBooks(books, 0)
	copyBookIDsToArray(&bookPagesID, books)
	r.cache.Set(booksKey, bookPagesID)
	r.cache.Set(
		bookParamsKey,
		dto.BookParamsDTO{
			Limit:  pageLimit,
			Offset: bookParams.Offset + pageLimit,
		},
	)

	return nil
}

func (r *Requester) viewNextPage() error {
	var bookParams dto.BookParamsDTO
	if err := r.cache.Get(bookParamsKey, &bookParams); err != nil {
		return err
	}

	var bookPagesID []uuid.UUID
	if err := r.cache.Get(booksKey, &bookPagesID); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + "/books",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		QueryParams: map[string]string{
			"title":           bookParams.Title,
			"author":          bookParams.Author,
			"publisher":       bookParams.Publisher,
			"copies_number":   fmt.Sprintf("%d", bookParams.CopiesNumber),
			"rarity":          bookParams.Rarity,
			"genre":           bookParams.Genre,
			"publishing_year": fmt.Sprintf("%d", bookParams.PublishingYear),
			"language":        bookParams.Language,
			"age_limit":       fmt.Sprintf("%d", bookParams.AgeLimit),
			"limit":           fmt.Sprintf("%d", bookParams.Limit),
			"offset":          fmt.Sprintf("%d", bookParams.Offset),
		},
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	var books []*jsonmodels.BookModel
	if err = json.Unmarshal(response.Body, &books); err != nil {
		return err
	}

	printBooks(books, bookParams.Offset)
	copyBookIDsToArray(&bookPagesID, books)
	r.cache.Set(booksKey, bookPagesID)
	r.cache.Set(
		bookParamsKey,
		dto.BookParamsDTO{
			Limit:  pageLimit,
			Offset: bookParams.Offset + pageLimit,
		},
	)

	return nil
}

func (r *Requester) ViewBook() error {
	var bookPagesID []uuid.UUID
	if err := r.cache.Get(booksKey, &bookPagesID); err != nil {
		return err
	}

	num, err := input.BookPagesNumber()
	if err != nil {
		return err
	}

	if num > len(bookPagesID) || num < 0 {
		return errors.New("book number out of range")
	}

	bookID := bookPagesID[num]

	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + fmt.Sprintf("/books/%s", bookID.String()),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	var book *jsonmodels.BookModel
	if err = json.Unmarshal(response.Body, &book); err != nil {
		return err
	}

	avgRating, err := r.getAvgRatingForBook(bookID)
	if err != nil {
		return err
	}

	printBook(book, avgRating, num)

	return nil

}

func (r *Requester) getAvgRatingForBook(bookID uuid.UUID) (float32, error) {
	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + "/ratings/avg",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		QueryParams: map[string]string{
			"book_id": bookID.String(),
		},
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return -1, err
	}

	if response.StatusCode == http.StatusNotFound {
		return -1, nil
	}

	if response.StatusCode != http.StatusOK {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return -1, err
		}
		return -1, errors.New(info)
	}

	var avgRating dto.AvgRatingDTO
	if err = json.Unmarshal(response.Body, &avgRating); err != nil {
		return -1, err
	}

	return avgRating.AvgRating, nil
}

func (r *Requester) AddToFavorites() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	var bookPagesID []uuid.UUID
	if err := r.cache.Get(booksKey, &bookPagesID); err != nil {
		return err
	}

	num, err := input.BookPagesNumber()
	if err != nil {
		return err
	}

	if num > len(bookPagesID) || num < 0 {
		return errors.New("book number out of range")
	}

	bookID := bookPagesID[num]

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/api/favorites",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
		},
		Body:    bookID,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("you are not authenticated")
	}

	if response.StatusCode != http.StatusCreated {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nBook successfully added to your favorites!\n")

	return nil
}

func (r *Requester) viewBookRatings() error {
	var bookPagesID []uuid.UUID
	if err := r.cache.Get(booksKey, &bookPagesID); err != nil {
		return err
	}

	num, err := input.BookPagesNumber()
	if err != nil {
		return err
	}

	if num > len(bookPagesID) || num < 0 {
		return errors.New("book number out of range")
	}

	bookID := bookPagesID[num]

	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + "/ratings",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		QueryParams: map[string]string{
			"book_id": bookID.String(),
		},
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	var ratings []*dto.RatingOutputDTO
	if err = json.Unmarshal(response.Body, &ratings); err != nil {
		return err
	}

	printRatings(ratings, num)

	return nil
}

func (r *Requester) addNewBookRating() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	var bookPagesID []uuid.UUID
	if err := r.cache.Get(booksKey, &bookPagesID); err != nil {
		return err
	}

	num, err := input.BookPagesNumber()
	if err != nil {
		return err
	}

	if num > len(bookPagesID) || num < 0 {
		return errors.New("book number out of range")
	}

	bookID := bookPagesID[num]

	ratingDTO, err := input.RatingParams()
	if err != nil {
		return err
	}
	ratingDTO.BookID = bookID

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/api/ratings",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
		},
		Body:    ratingDTO,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("you are not authenticated")
	}

	if response.StatusCode != http.StatusCreated {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nRating was successfully added!\n")

	return nil
}

func (r *Requester) ReserveBook() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	var bookPagesID []uuid.UUID
	if err := r.cache.Get(booksKey, &bookPagesID); err != nil {
		return err
	}

	num, err := input.BookPagesNumber()
	if err != nil {
		return err
	}

	if num > len(bookPagesID) || num < 0 {
		return errors.New("book number out of range")
	}

	bookID := bookPagesID[num]

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/api/reservations",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
		},
		Body:    bookID,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("you are not authenticated")
	}

	if response.StatusCode != http.StatusCreated {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nBook successfully reserved!\n")

	return nil
}

func printBook(book *jsonmodels.BookModel, avgRating float32, num int) {
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("Book №%d", num))
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatTitle

	t.AppendRow(table.Row{"Title", book.Title})
	t.AppendRow(table.Row{"Author", book.Author})
	t.AppendRow(table.Row{"Publisher", book.Publisher})
	t.AppendRow(table.Row{"Copies Number", book.CopiesNumber})
	t.AppendRow(table.Row{"Rarity", book.Rarity})
	t.AppendRow(table.Row{"Genre", book.Genre})
	t.AppendRow(table.Row{"Publishing Year", book.PublishingYear})
	t.AppendRow(table.Row{"Language", book.Language})
	t.AppendRow(table.Row{"Age Limit", book.AgeLimit})

	if avgRating == -1 {
		t.AppendRow(table.Row{"Avg Rating", "Has no rating"})
	} else {
		t.AppendRow(table.Row{"Avg Rating", fmt.Sprintf("%.1f", avgRating)})
	}

	fmt.Println(t.Render())
}

func printRatings(ratings []*dto.RatingOutputDTO, offset int) {
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("Отзывы на книгу №%d", offset))
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"No.", "Reader", "Review", "Rating"})

	for i, rating := range ratings {
		t.AppendRow(table.Row{offset + i, rating.Reader, rating.Review, rating.Rating})
	}
	fmt.Println(t.Render())
}

func printBooks(books []*jsonmodels.BookModel, offset int) {
	t := table.NewWriter()
	t.SetTitle(fmt.Sprintf("Страница книг №%d", offset/pageLimit+1))
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"No.", "Title", "Author"})

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "Author",
			WidthMax: 80,
		},
	})

	for i, book := range books {
		t.AppendRow(table.Row{offset + i, book.Title, book.Author})
	}
	fmt.Println(t.Render())
}

func copyBookIDsToArray(bookIDs *[]uuid.UUID, books []*jsonmodels.BookModel) {
	for _, book := range books {
		*bookIDs = append(*bookIDs, book.ID)
	}
}
