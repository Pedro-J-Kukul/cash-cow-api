// File: internal/data/cattle/breeds.go
package cattle

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/data/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
)

/****************************************************************************************
 *										Declarations									*
 ***************************************************************************************/
// Breed represents a cattle breed.
type Breed struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Breeds is a slice of Breed.
type Breeds []Breed

// BreedModel represents the model for cattle breeds.
type BreedModel struct {
	DB *sql.DB
}

// BreedFilter represents filtering options for querying cattle breeds.
type BreedFilter struct {
	Name     string
	IsActive *bool
	Default  filters.Filters
}

// Validate validates the fields of a Breed.
func (f *BreedFilter) Validate(v *validator.Validator, b *Breed) {
	v.Check(b.Name != "", "name", "must be provided")
	v.Check(len(b.Name) <= 255, "name", "must not be more than 255 characters long")
}

/****************************************************************************************
 *										Methods										*
 ***************************************************************************************/
//  Insert inserts a new cattle breed into the database.
func (m *BreedModel) Insert(b *Breed) error {
	query := `
		INSERT INTO cattle_breeds (name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, b.Name, b.Description, b.IsActive).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "cattle_breeds_name_key"):
			return errors.ErrDuplicateValue("name")
		default:
			return err
		}
	}
	return nil
}

// Update updates an existing cattle breed in the database.
func (m *BreedModel) Update(b *Breed) error {
	query := `
		UPDATE cattle_breeds
		SET name = $1, description = $2, is_active = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, b.Name, b.Description, b.IsActive, b.ID).Scan(&b.UpdatedAt)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "cattle_breeds_name_key"):
			return errors.ErrDuplicateValue("name")
		case errors.IsEditConflict(err):
			return errors.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete Permanently deletes a cattle breed from the database.
func (m *BreedModel) Delete(id int) error {
	query := `
		DELETE FROM cattle_breeds
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.IsForeignKeyViolation(err):
			return errors.ErrForeignKeyViolation
		default:
			return err
		}
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.ErrRecordNotFound
	}
	return nil
}

// GetByField retrieves a cattle breed by a specified field and value.
func (m *BreedModel) GetByID(id int) (*Breed, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM cattle_breeds
		WHERE id = $1
	`
	var b Breed
	scan := []any{
		&b.ID,
		&b.Name,
		&b.Description,
		&b.IsActive,
		&b.CreatedAt,
		&b.UpdatedAt,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(scan...)
	if err != nil {
		switch {
		case errors.IsEditConflict(err):
			return nil, errors.ErrEditConflict
		case errors.ErrNoRows(err):
			return nil, errors.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &b, nil
}

// GetAll retrieves all cattle breeds from the database.
func (m *BreedModel) GetAll(filter *BreedFilter) (Breeds, filters.MetaData, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, description, is_active, created_at, updated_at
		FROM cattle_breeds
		WHERE ($1 = '' OR LOWER(name) LIKE LOWER('%%' || $1 || '%%'))
		AND ($2::boolean IS NULL OR is_active = $2)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filter.Default.SortColumn(), filter.Default.SortDirection())

	args := []any{
		filter.Name,
		filter.IsActive,
		filter.Default.Limit(),
		filter.Default.Offset(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.EmptyMetaData, err
	}
	defer rows.Close()

	totalRecords := 0
	breeds := Breeds{}
	for rows.Next() {
		var count int
		var b Breed
		scan := []any{
			&count,
			&b.ID,
			&b.Name,
			&b.Description,
			&b.IsActive,
			&b.CreatedAt,
			&b.UpdatedAt,
		}
		err := rows.Scan(scan...)
		if err != nil {
			return nil, filters.EmptyMetaData, err
		}
		totalRecords = count
		breeds = append(breeds, b)
	}
	if err = rows.Err(); err != nil {
		return nil, filters.EmptyMetaData, err
	}

	metaData := filters.CalculateMetaData(totalRecords, filter.Default.Page, filter.Default.PageSize)
	return breeds, metaData, nil
}
