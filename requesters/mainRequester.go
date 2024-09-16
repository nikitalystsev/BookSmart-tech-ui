package requesters

import (
	"fmt"
	"github.com/nikitalystsev/BookSmart-tech-ui/input"
	myCache "github.com/nikitalystsev/BookSmart-tech-ui/pkg/cache"
	"os"
	"time"
)

const mainMenu = `Main menu:
	1 -- sign up
	2 -- sign in as reader
	3 -- sign in as administrator
	4 -- view books catalog
	0 -- exit program
`

type Requester struct {
	cache           myCache.ICache
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	baseURL         string
}

func NewRequester(
	accessTokenTTL,
	refreshTokenTTL time.Duration,
	port string,
) *Requester {
	return &Requester{
		cache:           myCache.NewCache(),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		baseURL:         "http://localhost:" + port,
	}
}

func (r *Requester) Run() {
	for {
		fmt.Printf("\n\n%s", mainMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
			continue
		}

		switch menuItem {
		case 1:
			if err = r.SignUp(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.ProcessReaderActions(); err != nil {
				continue
			}
		case 3:
			if err = r.ProcessAdminActions(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 4:
			if err = r.ProcessBookCatalogActions(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			os.Exit(0)
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
