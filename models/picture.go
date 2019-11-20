package models

import (
	"awise-socialNetwork/helpers"
	"time"
)

// Picture table models
type Picture struct {
	ID        int       `json:"id"`
	Path      string    `json:"path"`
	Origin    string    `json:"origin"`
	Size      string    `json:"size"`
	Active    int       `json:"active"`
	Delete    int       `json:"delete"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// FindPicture for find one picture by id
func FindPicture(id int) (*Picture, error) {
	picture := Picture{}
	result, err := db.Query("SELECT id, path, origin, size, active, `delete`, created_at, updated_at FROM tbl_pictures WHERE id = ?", id)
	if err != nil {
		return &picture, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&picture.ID, &picture.Path, &picture.Origin, &picture.Size, &picture.Active, &picture.CreatedAt, &picture.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &picture, nil
}

// FindAllPictures for find all avatars in the database
func FindAllPictures() ([]*Picture, error) {
	avatars := []*Picture{}
	result, err := db.Query("SELECT id, path, origin, size, active, `delete`, created_at, updated_at FROM tbl_pictures")
	if err != nil {
		return avatars, err
	}
	defer result.Close()
	for result.Next() {
		picture := Picture{}
		err := result.Scan(&picture.ID, &picture.Path, &picture.Origin, &picture.Size, &picture.Active, &picture.CreatedAt, &picture.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		avatars = append(avatars, &picture)
	}
	return avatars, nil
}

// Update a user
func (a *Picture) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_pictures SET path = ?, origin = ?, size = ?, active = ?, `delete` = ?, updated_at = ? WHERE id = ?")
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

// NewPicture create a new picture
func NewPicture(path string, origin string, size string) (*Picture, error) {
	picture := &Picture{}
	stmt, err := db.Prepare("INSERT INTO tbl_account(path, origin, size, active, `delete`, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return picture, err
	}
	defer stmt.Close()

	utc := helpers.GetUtc()

	result, err := stmt.Exec(path, origin, size, 0, 0, utc, utc)
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
