package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/schmorrison/Zoho"
	"github.com/schmorrison/Zoho/crm"
)

func main() {
	z := zoho.New()
	z.SetZohoDomain("eu")

	// to start oAuth2 flow
	scopes := []zoho.ScopeString{
		zoho.BuildScope(zoho.Crm, zoho.ModulesScope, zoho.AllMethod, zoho.NoOp),
	}

	// The authorization request will provide a link that must be clicked on or pasted into a browser.
	// Sometimes it will show the consent screen, upon consenting it will redirect to the redirectURL (currently the server doesn't return a value to the browser once getting the code)
	// The redirectURL provided here must match the URL provided when generating the clientID/secret
	// if the provided redirectURL is a localhost domain, the function will create a server on that port (use non-privileged port), and wait for the redirect to occur.
	// if the redirect provides the authorization code in the URL parameter "code", then the server catches it and provides it to the function for generating AccessToken and RefreshToken

	if err := z.AuthorizationCodeRequest(os.Getenv("ZOHO_CLIENT_ID"), os.Getenv("ZOHO_CLIENT_SECRET"), scopes, "http://localhost:8080/oauthredirect"); err != nil {
		log.Fatal(err)
	}

	c := crm.New(z)

	go func() {
		for {
			time.Sleep(time.Minute)
			err := z.RefreshTokenRequest()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// The API for getting module records is bound to change once the returned data types can be defined.
	// The returned JSON values are subject to change given that custom fields are an instrinsic part of zoho. (see brainstorm above)

	continueScan := true
	page := 0
	for continueScan {
		out, err := c.ListRecords(&crm.Account{}, crm.AccountsModule, map[string]zoho.Parameter{
			"page": zoho.Parameter(fmt.Sprintf("%d", page)),
		})
		if err != nil {
			log.Fatal(err)
		}
		data := out.(*crm.Account)

		for _, entry := range data.Data {
			newEntry := newData{
				AantalPersonen: -1,
				ID:             entry.ID,
			}
			upsert := crm.UpdateRecordData{
				Trigger: []string{"workflow"},
				Data:    []interface{}{newEntry},
			}
			out, err := c.UpdateRecord(upsert, crm.AccountsModule, entry.ID)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Triggered %s %s\n", entry.ID, out.Data[0].Message)
			time.Sleep(3 * time.Second)
		}
		if len(data.Data) < 200 {
			continueScan = false
		}
		page++
	}
}

type newData struct {
	AantalPersonen int    `json:"Aantal_Personen"`
	ID             string `json:"id"`
}
