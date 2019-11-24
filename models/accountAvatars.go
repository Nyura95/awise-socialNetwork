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
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FindAccountAvatar for find one accountAvatar by id
func FindAccountAvatar(id int) (*AccountAvatar, error) {
	accountAvatar := AccountAvatar{}
	result, err := db.Query("SELECT id, id_account, id_avatar, `delete`, active, created_at, updated_at FROM tbl_account_avatars WHERE id = ?", id)
	if err != nil {
		return &accountAvatar, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&accountAvatar.ID, &accountAvatar.IDAccount, &accountAvatar.IDAvatar, &accountAvatar.Delete, &accountAvatar.Active, &accountAvatar.CreatedAt, &accountAvatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &accountAvatar, nil
}

// FindAccountAvatarByIDAccountActive for find one accountAvatar active by id
func FindAccountAvatarByIDAccountActive(id int) (*AccountAvatar, error) {
	accountAvatar := AccountAvatar{}
	result, err := db.Query("SELECT id, id_account, id_avatar, `delete`, active, created_at, updated_at FROM tbl_account_avatars WHERE id_account = ? AND active = 1 LIMIT 1", id)
	if err != nil {
		return &accountAvatar, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&accountAvatar.ID, &accountAvatar.IDAccount, &accountAvatar.IDAvatar, &accountAvatar.Delete, &accountAvatar.Active, &accountAvatar.CreatedAt, &accountAvatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &accountAvatar, nil
}

// FindAllAccountAvatars for find all accountAvatars in the database
func FindAllAccountAvatars() ([]*AccountAvatar, error) {
	accountAvatars := []*AccountAvatar{}
	result, err := db.Query("SELECT id, id_account, id_avatar, `delete`, active, created_at, updated_at FROM tbl_account_avatars")
	if err != nil {
		return accountAvatars, err
	}
	defer result.Close()
	for result.Next() {
		accountAvatar := AccountAvatar{}
		err := result.Scan(&accountAvatar.ID, &accountAvatar.IDAccount, &accountAvatar.IDAvatar, &accountAvatar.Delete, &accountAvatar.Active, &accountAvatar.CreatedAt, &accountAvatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		accountAvatars = append(accountAvatars, &accountAvatar)
	}
	return accountAvatars, nil
}

// FindAllAccountAvatarsByIDAccount for find all accountAvatars in the database by id_account
func FindAllAccountAvatarsByIDAccount(idAccount int) ([]*AccountAvatar, error) {
	accountAvatars := []*AccountAvatar{}
	result, err := db.Query("SELECT id, id_account, id_avatar, `delete`, active, created_at, updated_at FROM tbl_account_avatars WHERE id_account = ?", idAccount)
	if err != nil {
		return accountAvatars, err
	}
	defer result.Close()
	for result.Next() {
		accountAvatar := AccountAvatar{}
		err := result.Scan(&accountAvatar.ID, &accountAvatar.IDAccount, &accountAvatar.IDAvatar, &accountAvatar.Delete, &accountAvatar.Active, &accountAvatar.CreatedAt, &accountAvatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		accountAvatars = append(accountAvatars, &accountAvatar)
	}
	return accountAvatars, nil
}

// Update a user
func (a *AccountAvatar) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_account_avatars SET id_account = ?, id_avatar = ?, `delete` = ?, active = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.IDAccount, a.IDAvatar, a.Delete, a.Active, helpers.GetUtc(), a.ID)
	if err != nil {
		return err
	}

	return nil
}

// DisabledAllAvatarByIDAccount remove all avatars for an account
func DisabledAllAvatarByIDAccount(idAccount int) error {
	stmt, err := db.Prepare("UPDATE tbl_account_avatars SET active = 0, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(helpers.GetUtc(), idAccount)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllAvatarByIDAccount remove all avatars for an account
func DeleteAllAvatarByIDAccount(idAccount int) error {
	stmt, err := db.Prepare("UPDATE tbl_account_avatars SET `delete` = 1,  updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(1, helpers.GetUtc(), idAccount)
	if err != nil {
		return err
	}

	return nil
}

// NewAccountAvatar create a new accountAvatar
func NewAccountAvatar(idaccount int, idavatar int) (*AccountAvatar, error) {
	accountAvatar := &AccountAvatar{}
	stmt, err := db.Prepare("INSERT INTO tbl_account_avatars(id_account, id_avatar, `delete`, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return accountAvatar, err
	}
	defer stmt.Close()

	utc := helpers.GetUtc()

	result, err := stmt.Exec(idaccount, idavatar, 0, 1, utc, utc)
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
