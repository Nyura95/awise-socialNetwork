package models

import (
	"awise-socialNetwork/helpers"
	"time"
)

// AccountAvatar table models
type AccountAvatar struct {
	ID        int       `json:"id"`
	IDAccount int       `json:"id_account"`
	IDAvatar  int       `json:"id_avatar"`
	Delete    int       `json:"delete"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FindAccountAvatar for find one accountAvatar by id
func FindAccountAvatar(id int) (*AccountAvatar, error) {
	accountAvatar := AccountAvatar{}
	result, err := db.Query("SELECT id, id_account, id_avatar, `delete`, created_at, updated_at FROM tbl_account_avatars WHERE id = ?", id)
	if err != nil {
		return &accountAvatar, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&accountAvatar.ID, &accountAvatar.IDAccount, &accountAvatar.IDAvatar, &accountAvatar.Delete, &accountAvatar.CreatedAt, &accountAvatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &accountAvatar, nil
}

// FindAllAccountAvatars for find all accountAvatars in the database
func FindAllAccountAvatars() ([]*AccountAvatar, error) {
	accountAvatars := []*AccountAvatar{}
	result, err := db.Query("SELECT id, id_account, id_avatar, `delete`, created_at, updated_at FROM tbl_account_avatars")
	if err != nil {
		return accountAvatars, err
	}
	defer result.Close()
	for result.Next() {
		accountAvatar := AccountAvatar{}
		err := result.Scan(&accountAvatar.ID, &accountAvatar.IDAccount, &accountAvatar.IDAvatar, &accountAvatar.Delete, &accountAvatar.CreatedAt, &accountAvatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		accountAvatars = append(accountAvatars, &accountAvatar)
	}
	return accountAvatars, nil
}

// Update a user
func (a *AccountAvatar) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_account_avatars SET id_account = ?, id_avatar = ?, `delete` = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.IDAccount, a.IDAvatar, a.Delete, time.UTC, a.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewAccountAvatar create a new accountAvatar
func NewAccountAvatar(idaccount int, idavatar int) (*AccountAvatar, error) {
	accountAvatar := &AccountAvatar{}
	stmt, err := db.Prepare("INSERT INTO tbl_account_avatars(id_account, id_avatar, `delete`, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return accountAvatar, err
	}
	defer stmt.Close()

	utc := helpers.GetUtc()

	result, err := stmt.Exec(idaccount, idavatar, 0, utc, utc)
	if err != nil {
		return accountAvatar, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return accountAvatar, err
	}

	accountAvatar, _ = FindAccountAvatar(int(ID))

	return accountAvatar, nil
}
