package requesters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-tech-ui/input"
	"net/http"
	"time"
)

const readerMainMenu = `Main menu:
	1 -- go to books catalog 
	2 -- go to library card
	3 -- go to your reservations
	0 -- log out
`

const tokensKey = "tokens"

func (r *Requester) ProcessReaderActions() error {
	var (
		menuItem int
		err      error
	)
	stopRefresh := make(chan struct{})

	if err = r.SignIn(stopRefresh); err != nil {
		fmt.Printf("\n\n%s\n", err.Error())
		return err
	}

	for {
		fmt.Printf("\n\n%s", readerMainMenu)

		if menuItem, err = input.MenuItem(); err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
			continue
		}

		switch menuItem {
		case 1:
			if err = r.ProcessBookCatalogActions(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.ProcessLibCardActions(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 3:
			if err = r.ProcessReservationsActions(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			close(stopRefresh)
			r.cache.Clear()
			fmt.Println("\n\nyou have successfully log out")
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) SignUp() error {
	readerSignUpDTO, err := input.SignUpParams()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/auth/sign-up",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    readerSignUpDTO,
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

	fmt.Printf("\n\nRegistration completed successfully!\n")

	return nil
}

func (r *Requester) SignIn(stopRefresh <-chan struct{}) error {
	readerSignInDTO, err := input.SignInParams()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/auth/sign-in",
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

func (r *Requester) Refresh() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/auth/refresh",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    tokens.RefreshToken,
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

	if err = json.Unmarshal(response.Body, &tokens); err != nil {
		return err
	}

	r.cache.Set(tokensKey, tokens)

	//fmt.Printf("\n\nSuccessful refresh tokens!\n")

	return nil
}

func (r *Requester) Refreshing(interval time.Duration, stopRefresh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.Refresh(); err != nil {
				fmt.Printf("\n\nerror refreshing tokens: %v\n", err)
			}
		case <-stopRefresh:
			return
		}
	}
}
