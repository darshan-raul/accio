package queries

import (

	"github.com/darshan-raul/accio/core/app/models"
	"github.com/jmoiron/sqlx"
)

// ResourceTypeQueries struct for queries from ResourceType model.
type ResourceTypeQueries struct {
	*sqlx.DB
}

// GetResourceTypes method for getting all books.
func (q *ResourceTypeQueries) GetResourceTypes() ([]models.ResourceType, error) {
	// Define books variable.
	resourectypes := []models.ResourceType{}

	// Define query string.
	query := `SELECT * FROM resourcetypes`

	// Send query to database.
	err := q.Select(&resourectypes, query)
	if err != nil {
		// Return empty object and error.
		return resourectypes, err
	}

	// Return query result.
	return resourectypes, nil
}

// GetResourceTypeByName method for getting one resourcetype by given Name.
func (q *ResourceTypeQueries) GetResourceTypeByName(name string) (models.ResourceType, error) {
	// Define resourcetype variable.
	resourcetype := models.ResourceType{}

	// Define query string.
	query := `SELECT id,name FROM resourcetypes WHERE name = $1`

	// Send query to database.
	err := q.Get(&resourcetype, query, name)
	if err != nil {
		// Return empty object and error.
		return resourcetype, err
	}

	// Return query result.
	return resourcetype, nil
}