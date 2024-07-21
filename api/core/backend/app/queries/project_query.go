package queries

import (
	"fmt"

	"github.com/darshan-raul/accio/core/app/models"
	"github.com/jmoiron/sqlx"
)

// ProjectQueries struct for queries from Project model.
type ProjectQueries struct {
	*sqlx.DB
}

// GetProjects method for getting all projects.
func (q *ProjectQueries) GetProjects() ([]models.Project, error) {
	// Define projects variable.
	projects := []models.Project{}

	// Define query string.
	query := `SELECT id,name,created_at FROM projects`
	fmt.Println("get projects called")
	// Send query to database.
	err := q.Select(&projects, query)
	if err != nil {
		// Return empty object and error.
		return projects, err
	}

	// Return query result.
	return projects, nil
}

// GetProject method for getting one project by given ID.
func (q *ProjectQueries) GetProject(id int) (models.Project, error) {
	// Define project variable.
	project := models.Project{}

	// Define query string.
	query := `SELECT id,name FROM projects WHERE id = $1`

	// Send query to database.
	err := q.Get(&project, query, id)
	if err != nil {
		// Return empty object and error.
		return project, err
	}

	// Return query result.
	return project, nil
}


// GetProject method for getting one project by given ID.
func (q *ProjectQueries) GetProjectByName(name string) (models.Project, error) {
	// Define project variable.
	project := models.Project{}

	// Define query string.
	query := `SELECT id,name FROM projects WHERE name = $1`

	// Send query to database.
	err := q.Get(&project, query, name)
	if err != nil {
		// Return empty object and error.
		return project, err
	}

	// Return query result.
	return project, nil
}

// CreateProject method for creating book by given Project object.
func (q *ProjectQueries) CreateProject(p *models.Project) error {
	// Define query string.
	query := `INSERT INTO projects (name,created_at) VALUES ($1, $2)`

	// Send query to database.
	_, err := q.Exec(query, p.Name, p.CreatedAt)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}


// DeleteProject method for delete project by given ID.
func (q *ProjectQueries) DeleteProject(id int) error {
	// Define query string.
	query := `DELETE FROM projects WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}