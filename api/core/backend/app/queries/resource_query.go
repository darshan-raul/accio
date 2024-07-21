package queries

import (
	"github.com/darshan-raul/accio/core/app/models"
	"github.com/jmoiron/sqlx"
)

// ResourceQueries struct for queries from Resource model.
type ResourceQueries struct {
	*sqlx.DB
}

// GetResources method for getting all resources.
func (q *ResourceQueries) GetResources() ([]models.Resource, error) {
	// Define resources variable.
	resources := []models.Resource{}

	// Define query string.
	query := `SELECT id,name FROM resources`

	// Send query to database.
	err := q.Select(&resources, query)
	if err != nil {
		// Return empty object and error.
		return resources, err
	}

	// Return query result.
	return resources, nil
}

// GetResource method for getting one resource by given ID.
func (q *ResourceQueries) GetResource(id int) (models.Resource, error) {
	// Define resource variable.
	resource := models.Resource{}

	// Define query string.
	query := `SELECT id,name FROM resources WHERE id = $1`

	// Send query to database.
	err := q.Get(&resource, query, id)
	if err != nil {
		// Return empty object and error.
		return resource, err
	}

	// Return query result.
	return resource, nil
}


// CreateResource method for creating book by given Resource object.
func (q *ResourceQueries) CreateResource(r *models.Resource) error {
	// Define query string.
	query := `INSERT INTO resources (name,cloud_prov_id,stack_id,res_type_id,created_at) VALUES ($1,$2,$3,$4,$5)`

	// Send query to database.
	_, err := q.Exec(query, r.Name, r.CloudProvId, r.StackId, r.TypeId , r.CreatedAt)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}


// DeleteResource method for delete resource by given ID.
func (q *ResourceQueries) DeleteResource(id int) error {
	// Define query string.
	query := `DELETE FROM resources WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}