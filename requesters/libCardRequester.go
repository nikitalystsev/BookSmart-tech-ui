package requesters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-tech-ui/input"
	jsonmodels "github.com/nikitalystsev/BookSmart-web-api/core/models"
	"log"
	"net/http"
	"time"
)

const libCardMenu = `Library card menu:
	1 -- create library card
	2 -- update library card
	3 -- view info library card
	0 -- go to main menu
`

func (r *Requester) ProcessLibCardActions() error {
	var (
		menuItem int
		err      error
	)
	for {
		fmt.Printf("\n\n%s", libCardMenu)

		if menuItem, err = input.MenuItem(); err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
			continue
		}

		switch menuItem {
		case 1:
			if err = r.CreateLibCard(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.UpdateLibCard(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 3:
			if err = r.ViewLibCard(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) CreateLibCard() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodPost,
		URL:    r.baseURL + "/api/lib-cards",
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

	if response.StatusCode != http.StatusCreated {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nSuccessfully created library card!\n")

	return nil
}

func (r *Requester) UpdateLibCard() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodPut,
		URL:    r.baseURL + "/api/lib-cards",
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

	fmt.Printf("\n\nSuccessfully updated library card!\n")

	return nil
}

func (r *Requester) ViewLibCard() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + "/api/lib-cards",
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

	var libCard *jsonmodels.LibCardModel
	err = json.Unmarshal(response.Body, &libCard)
	if err != nil {
		log.Fatal(err)
	}

	printLibCard(libCard)

	return nil

}

func printLibCard(libCard *jsonmodels.LibCardModel) {
	t := table.NewWriter()
	t.SetTitle("Library card")
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatTitle

	issueDateStr := libCard.IssueDate.Format("02.01.2006")

	statusStr := "Inactive"
	if libCard.ActionStatus {
		statusStr = "Active"
	}

	t.AppendRow(table.Row{"Number", libCard.LibCardNum})
	t.AppendRow(table.Row{"Validity", libCard.Validity})
	t.AppendRow(table.Row{"Issue date", issueDateStr})
	t.AppendRow(table.Row{"Status", statusStr})

	fmt.Println(t.Render())
}
