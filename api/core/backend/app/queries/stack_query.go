package queries

import (
	"github.com/darshan-raul/accio/core/app/models"
	"github.com/jmoiron/sqlx"
)

// StackQueries struct for queries from Stack model.
type StackQueries struct {
	*sqlx.DB
}

// GetStacks method for getting all stacks.
func (q *StackQueries) GetStacks() ([]models.Stack, error) {
	// Define stacks variable.
	stacks := []models.Stack{}

	// Define query string.
	query := `SELECT id,name FROM stacks`

	// Send query to database.
	err := q.Select(&stacks, query)
	if err != nil {
		// Return empty object and error.
		return stacks, err
	}

	// Return query result.
	return stacks, nil
}

// GetStack method for getting one stack by given ID.
func (q *StackQueries) GetStack(id int) (models.Stack, error) {
	// Define stack variable.
	stack := models.Stack{}

	// Define query string.
	query := `SELECT id,name FROM stacks WHERE id = $1`

	// Send query to database.
	err := q.Get(&stack, query, id)
	if err != nil {
		// Return empty object and error.
		return stack, err
	}

	// Return query result.
	return stack, nil
}


// CreateStack method for creating book by given Stack object.
func (q *StackQueries) CreateStack(p *models.Stack) error {
	// Define query string.
	query := `INSERT INTO stacks (name,project_id) VALUES ($1, $2)`

	// Send query to database.
	_, err := q.Exec(query, p.Name, p.ProjectId)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}


// DeleteStack method for delete stack by given ID.
func (q *StackQueries) DeleteStack(id int) error {
	// Define query string.
	query := `DELETE FROM stacks WHERE id = $1`

	// Send query to database.
	_, err := q.Exec(query, id)
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}