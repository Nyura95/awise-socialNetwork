package models

import (
	"awise-socialNetwork/helpers"
	"time"
)

// Picture table models
type Picture struct {
	ID         int       `json:"id"`
	IDAccount  int       `json:"id_account"`
	Path       string    `json:"path"`
	PathBlured string    `json:"path_blured"`
	Origin     string    `json:"origin"`
	Size       string    `json:"size"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// FindPicture for find one picture by id
func FindPicture(id int) (*Picture, error) {
	picture := Picture{}
	result, err := db.Query("SELECT id, id_account, path, path_blured, origin, size, created_at, updated_at FROM tbl_pictures WHERE id = ?", id)
	if err != nil {
		return &picture, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&picture.ID, &picture.IDAccount, &picture.Path, &picture.PathBlured, &picture.Origin, &picture.Size, &picture.CreatedAt, &picture.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &picture, nil
}

// FindAllPictures for find all avatars in the database
func FindAllPictures() ([]*Picture, error) {
	avatars := []*Picture{}
	result, err := db.Query("SELECT id, id_account, path, path_blured, origin, size, created_at, updated_at FROM tbl_pictures")
	if err != nil {
		return avatars, err
	}
	defer result.Close()
	for result.Next() {
		picture := Picture{}
		err := result.Scan(&picture.ID, &picture.IDAccount, &picture.Path, &picture.PathBlured, &picture.Origin, &picture.Size, &picture.CreatedAt, &picture.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		avatars = append(avatars, &picture)
	}
	return avatars, nil
}

// Update a user
func (a *Picture) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_pictures SET path = ?, path_blured = ?, origin = ?, size = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Path, a.PathBlured, a.Origin, a.Size, helpers.GetUtc(), a.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewPicture create a new picture
func NewPicture(idaccount int, path string, pathBlured string, origin string, size string) (*Picture, error) {
	picture := &Picture{}
	stmt, err := db.Prepare("INSERT INTO tbl_pictures(id_account, path, path_blured, origin, size, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return picture, err
	}
	defer stmt.Close()

	utc := helpers.GetUtc()

	result, err := stmt.Exec(idaccount, path, pathBlured, origin, size, utc, utc)
	if err != nil {
		return picture, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return picture, err
	}

	picture, _ = FindPicture(int(ID))

	return picture, nil
}
