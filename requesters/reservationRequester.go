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

const reservationsMenu = `Reservations menu:
	1 -- view your reservations 
	2 -- update your reservation
	0 -- go to main menu
`

const reservationsKey = "reservations"

func (r *Requester) ProcessReservationsActions() error {
	var (
		menuItem int
		err      error
	)

	r.cache.Set(reservationsKey, make([]uuid.UUID, 0))

	for {
		fmt.Printf("\n\n%s", reservationsMenu)

		if menuItem, err = input.MenuItem(); err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
			continue
		}

		switch menuItem {
		case 1:
			if err = r.ViewReservations(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.UpdateReservation(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) ViewReservations() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	var reservationsID []uuid.UUID
	if err := r.cache.Get(reservationsKey, &reservationsID); err != nil {
		return err
	}

	request := HTTPRequest{
		Method: http.MethodGet,
		URL:    r.baseURL + "/api/reservations",
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

	var reservations []*jsonmodels.ReservationModel
	if err = json.Unmarshal(response.Body, &reservations); err != nil {
		return err
	}

	printReservations(reservations)
	copyReservationIDsToArray(&reservationsID, reservations)
	r.cache.Set(reservationsKey, reservationsID)

	return nil
}

func (r *Requester) UpdateReservation() error {
	var tokens dto.ReaderTokensDTO
	if err := r.cache.Get(tokensKey, &tokens); err != nil {
		return err
	}

	var reservationsID []uuid.UUID
	if err := r.cache.Get(reservationsKey, &reservationsID); err != nil {
		return err
	}

	num, err := input.ReservationNumber()
	if err != nil {
		return err
	}

	if num > len(reservationsID) || num < 0 {
		return errors.New("reservation number out of range")
	}

	reservationID := reservationsID[num]

	request := HTTPRequest{
		Method: http.MethodPut,
		URL:    r.baseURL + fmt.Sprintf("/api/reservations/%s", reservationID.String()),
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
		},
		Body:    reservationID,
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

	fmt.Printf("\n\nReservation successfully updated!\n")

	return nil
}

func printReservations(reservations []*jsonmodels.ReservationModel) {
	t := table.NewWriter()
	t.SetTitle("Reservations")
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatTitle
	t.AppendHeader(table.Row{"No.", "Issue Date", "Return Date", "State"})

	for i, r := range reservations {
		t.AppendRow(table.Row{i, r.IssueDate.Format("2006-01-02"), r.ReturnDate.Format("2006-01-02"), r.State})
	}
	fmt.Println(t.Render())
}

func copyReservationIDsToArray(reservationIDs *[]uuid.UUID, reservations []*jsonmodels.ReservationModel) {
	for _, reservation := range reservations {
		*reservationIDs = append(*reservationIDs, reservation.ID)
	}
}
