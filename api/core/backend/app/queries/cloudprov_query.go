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


// GetCloudProvByName method for getting one cloudprov by given Name.
func (q *CloudProvQueries) GetCloudProvBySlug(slug string) (models.CloudProv, error) {
	// Define cloudprov variable.
	cloudprov := models.CloudProv{}

	// Define query string.
	query := `SELECT id,name FROM cloudproviders WHERE slug = $1`

	// Send query to database.
	err := q.Get(&cloudprov, query, slug)
	if err != nil {
		// Return empty object and error.
		return cloudprov, err
	}

	// Return query result.
	return cloudprov, nil
}