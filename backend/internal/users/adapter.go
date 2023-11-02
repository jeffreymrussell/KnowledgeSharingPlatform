package users

import (
	"KnowledgeSharingPlatform/internal"
	"database/sql"
	"errors"
)

type SQLiteAdapter struct {
	DB *sql.DB
}

func (s *SQLiteAdapter) SaveUser(user User) error {
	_, err := s.DB.Exec(`INSERT INTO users (username, password_hash, email) VALUES (?, ?, ?)`,
		user.Username, user.PasswordHash, user.Email)

	return err
}

func (s *SQLiteAdapter) GetUserByUsername(username string) (User, error) {
	var user User
	row := s.DB.QueryRow(`SELECT id, username, password_hash, email, created_at, updated_at FROM users WHERE username = ?`, username)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, internal.ErrNotFound
		}
		return User{}, err
	}
	return user, nil
}
