// Package models defines model for snippets
package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet - Define a Snippet type to hold the data for an individual snippet.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel - Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert - This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL statement to execute, backquotes to split string into two lines for readability
	stmt := `INSERT INTO snippets (title, content, created, expires)
					VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Exec() method on the embedded connection pool executes query
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// LastInsertId() gets the ID of newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, gets converted to an int type before returning.
	return int(id), nil
}

// Get - This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets
					WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Executes the statement
	row := m.DB.QueryRow(stmt, id)

	// Initialize a new zeroed Snippet struct.
	var s Snippet

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	// If everything went OK, then return the filled Snippet struct.
	return s, nil
}

// Latest - This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	// SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippets
					WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	var snippets []Snippet

	for rows.Next() {
		// Create a new zeroed Snippet struct.
		var s Snippet
		// Use rows.Scan() to copy the values from each field in the row
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return snippets, nil
}
