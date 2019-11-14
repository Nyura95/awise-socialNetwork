package models

import (
	"awise-socialNetwork/helpers"
	"time"
)

// Picture table model
type Picture struct {
	ID        int       `json:"id"`
	Path      string    `json:"path"`
	Origin    string    `json:"origin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FindPicture for find one picture by id
func FindPicture(id int) (*Picture, error) {
	picture := Picture{}
	result, err := db.Query("SELECT id, path, origin, created_at, updated_at FROM tbl_pictures WHERE id = ? LIMIT 1", id)
	if err != nil {
		return &picture, err
	}
	defer result.Close()
	for result.Next() {
		err := result.Scan(&picture.ID, &picture.Path, &picture.Origin, &picture.CreatedAt, &picture.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
	}
	return &picture, nil
}

// FindAllPicture for find all access_token in the database
func FindAllPicture() ([]*Picture, error) {
	pictures := []*Picture{}
	result, err := db.Query("SELECT id, path, origin, created_at, updated_at FROM tbl_pictures")
	if err != nil {
		return pictures, err
	}
	defer result.Close()
	for result.Next() {
		picture := Picture{}
		err := result.Scan(&picture.ID, &picture.Path, &picture.Origin, &picture.CreatedAt, &picture.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		pictures = append(pictures, &picture)
	}
	return pictures, nil
}

// Update a picture
func (p *Picture) Update() error {
	stmt, err := db.Prepare("UPDATE tbl_pictures SET path = ?, origin = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Path, p.Origin, helpers.GetUtc(), p.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewPicture create a new access_token
func NewPicture(path string, origin string) (*Picture, error) {
	stmt, err := db.Prepare("INSERT INTO tbl_pictures(path, origin, created_at, updated_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	utc := helpers.GetUtc()

	result, err := stmt.Exec(path, origin, utc, utc)
	if err != nil {
		return nil, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	picture, _ := FindPicture(int(ID))

	return picture, nil
}
