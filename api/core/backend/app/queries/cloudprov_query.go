package queries

import (

	"github.com/darshan-raul/accio/core/app/models"
	"github.com/jmoiron/sqlx"
)

// CloudProvQueries struct for queries from CloudProv model.
type CloudProvQueries struct {
	*sqlx.DB
}

// GetCloudProvs method for getting all books.
func (q *CloudProvQueries) GetCloudProvs() ([]models.CloudProv, error) {
	// Define books variable.
	cloudprovs := []models.CloudProv{}

	// Define query string.
	query := `SELECT * FROM cloudproviders`

	// Send query to database.
	err := q.Select(&cloudprovs, query)
	if err != nil {
		// Return empty object and error.
		return cloudprovs, err
	}

	// Return query result.
	return cloudprovs, nil
}

// // GetCloudProvsByAuthor method for getting all books by given author.
// func (q *CloudProvQueries) GetCloudProvsByAuthor(author string) ([]models.CloudProv, error) {
// 	// Define books variable.
// 	books := []models.CloudProv{}

// 	// Define query string.
// 	query := `SELECT * FROM books WHERE author = $1`

// 	// Send query to database.
// 	err := q.Get(&books, query, author)
// 	if err != nil {
// 		// Return empty object and error.
// 		return books, err
// 	}

// 	// Return query result.
// 	return books, nil
// }

// // GetCloudProv method for getting one book by given ID.
// func (q *CloudProvQueries) GetCloudProv(id uuid.UUID) (models.CloudProv, error) {
// 	// Define book variable.
// 	book := models.CloudProv{}

// 	// Define query string.
// 	query := `SELECT * FROM books WHERE id = $1`

// 	// Send query to database.
// 	err := q.Get(&book, query, id)
// 	if err != nil {
// 		// Return empty object and error.
// 		return book, err
// 	}

// 	// Return query result.
// 	return book, nil
// }

// // CreateCloudProv method for creating book by given CloudProv object.
// func (q *CloudProvQueries) CreateCloudProv(b *models.CloudProv) error {
// 	// Define query string.
// 	query := `INSERT INTO books VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

// 	// Send query to database.
// 	_, err := q.Exec(query, b.ID, b.CreatedAt, b.UpdatedAt, b.UserID, b.Title, b.Author, b.CloudProvStatus, b.CloudProvAttrs)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// // UpdateCloudProv method for updating book by given CloudProv object.
// func (q *CloudProvQueries) UpdateCloudProv(id uuid.UUID, b *models.CloudProv) error {
// 	// Define query string.
// 	query := `UPDATE books SET updated_at = $2, title = $3, author = $4, book_status = $5, book_attrs = $6 WHERE id = $1`

// 	// Send query to database.
// 	_, err := q.Exec(query, id, b.UpdatedAt, b.Title, b.Author, b.CloudProvStatus, b.CloudProvAttrs)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }

// // DeleteCloudProv method for delete book by given ID.
// func (q *CloudProvQueries) DeleteCloudProv(id uuid.UUID) error {
// 	// Define query string.
// 	query := `DELETE FROM books WHERE id = $1`

// 	// Send query to database.
// 	_, err := q.Exec(query, id)
// 	if err != nil {
// 		// Return only error.
// 		return err
// 	}

// 	// This query returns nothing.
// 	return nil
// }
