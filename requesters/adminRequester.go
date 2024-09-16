package requesters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-tech-ui/input"
	jsonmodels "github.com/nikitalystsev/BookSmart-web-api/core/models"
	"net/http"
	"time"
)

func (r *Requester) ProcessAdminActions() error {

	stopRefresh := make(chan struct{})

	if err := r.SignInAsAdmin(stopRefresh); err != nil {
		return err
	}

	for {
		fmt.Printf("\n\n%s", readerMainMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 1:
			err = r.ProcessAdminBookCatalogActions()
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = r.ProcessLibCardActions()
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err = r.ProcessReservationsActions()
			if err != nil {
				fmt.Println(err)
			}
		case 0:
			close(stopRefresh)
			fmt.Println("\n\nyou have successfully log out")
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) SignInAsAdmin(stopRefresh <-chan struct{}) error {
	readerSignInDTO, err := input.SignInParams()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/auth/admin/sign-in",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    readerSignInDTO,
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

	var tokens dto.ReaderTokensDTO
	if err = json.Unmarshal(response.Body, &tokens); err != nil {
		return err
	}

	r.cache.Set(tokensKey, tokens)

	fmt.Printf("\n\nAuthentication successful!\n")

	go r.Refreshing(r.accessTokenTTL, stopRefresh)

	return nil
}

const adminCatalogMenu = `Admin's Catalog menu:
	1 -- view books
	2 -- next page
	3 -- view info about book
	4 -- add book to favorites
	5 -- reserve book
	6 -- add new book
	7 -- delete book
	0 -- go to main menu
`

func (r *Requester) ProcessAdminBookCatalogActions() error {
	r.cache.Set(bookParamsKey, dto.BookParamsDTO{Limit: pageLimit, Offset: 0})
	r.cache.Set(booksKey, make([]uuid.UUID, 0))

	for {
		fmt.Printf("\n\n%s", adminCatalogMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
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
			if err = r.AddNewBook(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 7:
			if err = r.DeleteBook(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) AddNewBook() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	newBook, err := input.Book()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/api/admin/books",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
		},
		Body:    newBook,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nBook successfully created!\n")

	return nil
}

func (r *Requester) DeleteBook() error {
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

	if err = r.getReservationsByBook(bookID); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodDelete,
		URL:    r.baseURL + fmt.Sprintf("/api/admin/books/%s", bookID.String()),
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
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

	fmt.Printf("\n\nBook successfully deleted!\n")

	return nil
}

func (r *Requester) getReservationsByBook(bookID uuid.UUID) error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + "/api/admin/reservations",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
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

	if response.StatusCode == http.StatusNotFound {
		return nil
	}

	if response.StatusCode != http.StatusOK {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	var reservations []*jsonmodels.ReservationModel
	if err = json.Unmarshal(response.Body, &reservations); err != nil {
		return err
	}
	if len(reservations) > 0 {
		return errors.New("this book cannot be deleted, it is reserved")
	}

	return nil
}
