package models

import (
	"awise-socialNetwork/helpers"
	"time"
)

// Avatar table models
type Avatar struct {
	ID        int       `json:"id"`
	Path      string    `json:"path"`
	Origin    string    `json:"origin"`
	Size      string    `json:"size"`
	Active    int       `json:"active"`
	Delete    int       `json:"delete"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// FindAvatar for find one avatar by id
func FindAvatar(id int) (*Avatar, error) {
	avatar := Avatar{}
	result, err := db.Query("SELECT id, path, origin, size, active, `delete`, created_at, updated_at FROM tbl_avatars WHERE id = ?", id)
	if err != nil {
		return &avatar, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&avatar.ID, &avatar.Path, &avatar.Origin, &avatar.Size, &avatar.Active, &avatar.CreatedAt, &avatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &avatar, nil
}

// FindAllAvatars for find all avatars in the database
func FindAllAvatars() ([]*Avatar, error) {
	avatars := []*Avatar{}
	result, err := db.Query("SELECT id, path, origin, size, active, `delete`, created_at, updated_at FROM tbl_avatars")
	if err != nil {
		return avatars, err
	}
	defer result.Close()
	for result.Next() {
		avatar := Avatar{}
		err := result.Scan(&avatar.ID, &avatar.Path, &avatar.Origin, &avatar.Size, &avatar.Active, &avatar.CreatedAt, &avatar.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		avatars = append(avatars, &avatar)
	}
	return avatars, nil
}

// Update a user
func (a *Avatar) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_avatars SET path = ?, origin = ?, size = ?, active = ?, `delete` = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Path, a.Origin, a.Size, a.Active, a.Delete, time.UTC, a.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewAvatar create a new avatar
func NewAvatar(path string, origin string, size string) (*Avatar, error) {
	avatar := &Avatar{}
	stmt, err := db.Prepare("INSERT INTO tbl_account(path, origin, size, active, `delete`, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return avatar, err
	}
	defer stmt.Close()

	utc := helpers.GetUtc()

	result, err := stmt.Exec(path, origin, size, 0, 0, utc, utc)
	if err != nil {
		return avatar, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return avatar, err
	}

	avatar, _ = FindAvatar(int(ID))

	return avatar, nil
}
