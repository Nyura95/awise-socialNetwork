package models

import (
	"awise-socialNetwork/helpers"
	"sync"
	"time"
)

// Account table models
type Account struct {
	ID        int    `json:"id"`
	IDAvatars string `json:"id_avatars"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`

	Score   int    `json:"score"`
	Level   int    `json:"level"`
	Credits int    `json:"credits"`
	Phone   string `json:"phone"`
	City    string `json:"city"`
	Country string `json:"country"`

	password  string
	IDScope   int       `json:"idScope"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// FindAccount for find one account by id
func FindAccount(id int) (*Account, error) {
	account := Account{}
	result, err := db.Query("SELECT id, id_avatars, first_name, last_name, username, email, bio, score, level, credits, phone, city, country, password, id_scope, created_at, updated_at FROM tbl_account WHERE id = ?", id)
	if err != nil {
		return &account, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&account.ID, &account.IDAvatars, &account.Firstname, &account.Lastname, &account.Username, &account.Email, &account.Bio, &account.Score, &account.Level, &account.Credits, &account.Phone, &account.City, &account.Country, &account.password, &account.IDScope, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &account, nil
}

// FindAccountByPassword for find one account by password
func FindAccountByPassword(password string) (*Account, error) {
	account := Account{}
	result, err := db.Query("SELECT id, id_avatars, first_name, last_name, username, bio, email, score, level, credits, phone, city, country, password, id_scope, created_at, updated_at FROM tbl_account WHERE password = ?", password)
	if err != nil {
		return &account, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&account.ID, &account.IDAvatars, &account.Firstname, &account.Lastname, &account.Username, &account.Bio, &account.Email, &account.Score, &account.Level, &account.Credits, &account.Phone, &account.City, &account.Country, &account.password, &account.IDScope, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &account, nil
}

// FindAllAccount for find all accounts in the database
func FindAllAccount() ([]*Account, error) {
	accounts := []*Account{}
	result, err := db.Query("SELECT id, id_avatars, first_name, last_name, username, bio, email, score, level, credits, phone, city, country, password, id_scope, created_at, updated_at FROM tbl_account")
	if err != nil {
		return accounts, err
	}
	defer result.Close()
	for result.Next() {
		account := Account{}
		err := result.Scan(&account.ID, &account.IDAvatars, &account.Firstname, &account.Lastname, &account.Username, &account.Bio, &account.Email, &account.Score, &account.Level, &account.Credits, &account.Phone, &account.City, &account.Country, &account.password, &account.IDScope, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

// FindAllAccountByEmailOrUsername check if an account already exist with a username or email
func FindAllAccountByEmailOrUsername(email string, username string) ([]*Account, error) {
	accounts := []*Account{}
	result, err := db.Query("SELECT id, id_avatars, first_name, last_name, username, bio, email, score, level, credits, phone, city, country, password, id_scope, created_at, updated_at FROM tbl_account WHERE email = ? OR username = ?", email, username)
	if err != nil {
		return accounts, err
	}
	defer result.Close()
	for result.Next() {
		account := Account{}
		err := result.Scan(&account.ID, &account.IDAvatars, &account.Firstname, &account.Lastname, &account.Username, &account.Bio, &account.Email, &account.Score, &account.Level, &account.Credits, &account.Phone, &account.City, &account.Country, &account.password, &account.IDScope, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

// CheckAccountExist check if account exist
func CheckAccountExist(IDAccounts ...int) bool {
	jobs := make(chan bool, len(IDAccounts))
	defer close(jobs)

	var wg sync.WaitGroup
	wg.Add(len(IDAccounts))

	for _, IDAccount := range IDAccounts {
		go func(IDAccount int) {
			defer wg.Done()
			account, err := FindAccount(IDAccount)
			if err != nil || account.ID == 0 {
				jobs <- false
			}
		}(IDAccount)
	}

	wg.Wait()

	return len(jobs) == 0
}

// Update a user
func (a *Account) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_account SET avatars = ?, firstname = ?, lastname = ?, username = ?, bio = ?, score = ?, level = ?, credits = ?, phone = ?, city = ?, country = ?, password = ?, id_scope = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.IDAvatars, a.Firstname, a.Lastname, a.Username, a.Bio, a.Score, a.Level, a.Credits, a.Phone, a.City, a.Country, a.password, a.IDScope, time.UTC, a.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewAccount create a new user
func NewAccount(username string, email string, password string) (*Account, error) {
	account := &Account{}
	stmt, err := db.Prepare("INSERT INTO tbl_account(id_avatars, first_name, last_name, username, bio, email, score, level, credits, phone, city, country, password, id_scope, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return account, err
	}
	defer stmt.Close()

	utc := helpers.GetUtc()

	result, err := stmt.Exec(0, "", "", username, "", email, 0, 0, 0, "", "", "", helpers.StringToMD5(username+":"+password), 1, utc, utc)
	if err != nil {
		return account, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return account, err
	}

	account, _ = FindAccount(int(ID))

	return account, nil
}
