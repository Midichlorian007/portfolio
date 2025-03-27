package repository

const (
	createTableQry = `
	CREATE TABLE IF NOT EXISTS 
    users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
          )`

	GetUserQry = `SELECT id, name FROM users WHERE id = ?`

	UpdateUserQry = `UPDATE users SET name = ? WHERE id = ?`

	CreateUserQry = `INSERT INTO users (name) VALUES (?)`

	DeleteUserQry = `DELETE FROM users WHERE id = ?`
)
