package models

import (
	"awise-socialNetwork/helpers"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")
var jwtKeySecret = []byte("my_secret_key_secret")

// AccessToken table model
type AccessToken struct {
	ID           int       `json:"id"`
	IDAccount    int       `json:"id_account"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	FlagDelete   int       `json:"flag_delete"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// FindAccessToken for find one access_token by id
func FindAccessToken(id int) (*AccessToken, error) {
	accessToken := AccessToken{}
	result, err := db.Query("SELECT id, id_account, token, refresh_token, flag_delete, expired_at, created_at, updated_at FROM tbl_access_token WHERE id = ? LIMIT 1", id)
	if err != nil {
		return &accessToken, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&accessToken.ID, &accessToken.IDAccount, &accessToken.Token, &accessToken.RefreshToken, &accessToken.FlagDelete, &accessToken.ExpiredAt, &accessToken.CreatedAt, &accessToken.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &accessToken, nil
}

// FindAccessTokenByToken for find one access_token by token
func FindAccessTokenByToken(token string) (*AccessToken, error) {
	accessToken := AccessToken{}
	result, err := db.Query("SELECT id, id_account, token, refresh_token, flag_delete, expired_at, created_at, updated_at FROM tbl_access_token WHERE token = ? LIMIT 1", token)
	if err != nil {
		return &accessToken, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&accessToken.ID, &accessToken.IDAccount, &accessToken.Token, &accessToken.RefreshToken, &accessToken.FlagDelete, &accessToken.ExpiredAt, &accessToken.CreatedAt, &accessToken.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &accessToken, nil
}

// FindAllAccessToken for find all access_token in the database
func FindAllAccessToken() ([]*AccessToken, error) {
	accessTokens := []*AccessToken{}
	result, err := db.Query("SELECT id, id_account, token, refresh_token, flag_delete, expired_at, created_at, updated_at FROM tbl_access_token")
	if err != nil {
		return accessTokens, err
	}
	defer result.Close()
	for result.Next() {
		accessToken := AccessToken{}
		err := result.Scan(&accessToken.ID, &accessToken.IDAccount, &accessToken.Token, &accessToken.RefreshToken, &accessToken.FlagDelete, &accessToken.ExpiredAt, &accessToken.CreatedAt, &accessToken.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		accessTokens = append(accessTokens, &accessToken)
	}
	return accessTokens, nil
}

// Update a access_token
func (a *AccessToken) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_access_token SET id_account = ?, token = ?, refresh_token = ?, flag_delete = ?, expired_at = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.IDAccount, a.Token, a.RefreshToken, a.FlagDelete, a.ExpiredAt, helpers.GetUtc(), a.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllAccessTokenByIDAccount disable all accessToken for one account
func DeleteAllAccessTokenByIDAccount(IDAccount int) error {
	stmt, err := db.Prepare("UPDATE tbl_access_token SET flag_delete = 1, updated_at = ? WHERE id_account = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(helpers.GetUtc(), IDAccount)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAccessTokenByID disable accessToken by id
func DeleteAccessTokenByID(ID int) error {
	stmt, err := db.Prepare("UPDATE tbl_access_token SET flag_delete = 1, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(helpers.GetUtc(), ID)
	if err != nil {
		return err
	}

	return nil
}

// NewAccessToken create a new access_token
func NewAccessToken(IDAccount int, expiredAt time.Time) (*AccessToken, error) {
	stmt, err := db.Prepare("INSERT INTO tbl_access_token(id_account, token, refresh_token, flag_delete, expired_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	utc := helpers.GetUtc()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"IDAccount": IDAccount,
		"nbf":       time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, _ := token.SignedString(jwtKey)
	tokenStringSecret, _ := token.SignedString(jwtKeySecret)

	result, err := stmt.Exec(IDAccount, tokenString, tokenStringSecret, 0, expiredAt, utc, utc)
	if err != nil {
		return nil, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	accessToken, _ := FindAccessToken(int(ID))

	return accessToken, nil
}
