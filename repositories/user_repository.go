package repositories

import (
	"database/sql"
	"fmt"
	"go-api-demo/database"
	"go-api-demo/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

const (
	createUserQuery  = `INSERT INTO users (name, email) VALUES (?, ?)`
	getUserByIDQuery = `SELECT id, name, email, created_at FROM users WHERE id = ?`
	getAllUsersQuery = `SELECT id, name, email, created_at FROM users`
	updateUserQuery  = `UPDATE users SET name = ?, email = ? WHERE id = ?`
	deleteUserQuery  = `DELETE FROM users WHERE id = ?`
)

type SQLiteUserRepository struct {
	db *sql.DB // Uses the global DB connection we set up
}

func NewUserRespository() UserRepository { // Constructor like init()
	repo := &SQLiteUserRepository{
		db: database.DB,
	}

	if err := repo.ValidateSchema(); err != nil {
		panic(fmt.Sprintf("Schema validation failed %v", err))
	}

	if err := repo.TestQueries(); err != nil {
		panic(fmt.Sprintf("Query validation failed: %v", err))
	}

	return repo
}

// Extending SQLRepo
func (r *SQLiteUserRepository) ValidateSchema() error {
	// Check if table exists
	var tableName string

	// SELECT name - What We Want
	// We only want the "name" column from the results

	// sqlite_master: This is SQLite's internal "catalog" table
	// It's meta data and not the actual table so has no risk of crashing when querying
	// It contains information about ALL objects in your database

	// Example of what sqlite_master contains:
	// type        name        sql
	// ----------  ----------  ----------------------------------------
	// table       users       CREATE TABLE users (id INTEGER, ...)
	// table       posts       CREATE TABLE posts (id INTEGER, ...)
	// index       idx_email   CREATE INDEX idx_email ON users(email)

	//"Does a table named 'users' exist in this database?"
	err := r.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
	if err != nil {
		return fmt.Errorf("users table does not exist")
	}

	return nil
}

func (r *SQLiteUserRepository) TestQueries() error {
	queries := []string{
		createUserQuery,
		getUserByIDQuery,
		getAllUsersQuery,
		updateUserQuery,
		deleteUserQuery,
	}

	for _, query := range queries {
		// Try to prepare each query
		stmt, err := r.db.Prepare(query) // validates SQL syntax without executing it
		if err != nil {
			return fmt.Errorf("invalid query %s - Error: %v", query, err)
		}
		stmt.Close() // Clean up the prepared statement
	}
	return nil
}

func (r *SQLiteUserRepository) Create(user *models.User) error {
	result, err := r.db.Exec(createUserQuery, user.Name, user.Email)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}

	user.ID = int(id)
	return nil
}

func (r *SQLiteUserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{} // create empty struct
	err := r.db.QueryRow(getUserByIDQuery, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *SQLiteUserRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query(getAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *SQLiteUserRepository) Update(user *models.User) error {
	_, err := r.db.Exec(updateUserQuery, user.Name, user.Email, user.ID)
	return err
}

func (r *SQLiteUserRepository) Delete(id int) error {
	_, err := r.db.Exec(deleteUserQuery, id)
	return err
}
