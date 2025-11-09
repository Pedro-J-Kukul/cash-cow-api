// File: internal/data/locations/regions.go
package locations

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

// Region represents a geographical region.
type Region struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// RegionsModel represents the model for regions.
type RegionModel struct {
	DB *sql.DB
}

// Regions is a slice of Region.
type Regions []Region

// RegionFilter represents filtering options for querying regions.
type RegionFilter struct {
	Name    string
	Code    string
	Default filters.Filters
}

// ValidateRegion validates the fields of a Region.
func ValidateRegion(v *validator.Validator, r *Region) {
	v.Check(r.Name != "", "name", "must be provided")
	v.Check(len(r.Name) <= 255, "name", "must not be more than 255 bytes long")
	v.Check(r.Code != "", "code", "must be provided")
	v.Check(len(r.Code) <= 10, "code", "must not be more than 10 bytes long")
}

/****************************************************************************************
 *										Methods											*
 ***************************************************************************************/

// Insert adds a new region to the database.
func (m *RegionModel) Insert(r Region) error {
	query := `
		INSERT INTO regions (name, code)
		VALUES ($1, $2)
	`
	args := []any{r.Name, r.Code}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "code"):
			return errors.ErrDuplicateValue("code")
		case errors.IsUniqueViolation(err, "name"):
			return errors.ErrDuplicateValue("name")
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
		return errors.ErrInsertFailed
	}
	return nil
}

// Update modifies an existing region in the database.
func (m *RegionModel) Update(r *Region) error {
	query := `
		UPDATE regions
		SET name = $1, code = $2
		WHERE id = $3
	`
	args := []any{r.Name, r.Code, r.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "code"):
			return errors.ErrDuplicateCode
		case errors.IsUniqueViolation(err, "name"):
			return errors.ErrDuplicateName
		case errors.IsEditConflict(err):
			return errors.ErrEditConflict
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
		return sql.ErrNoRows
	}

	return nil
}

// Delete permanently removes a region from the database.
func (m *RegionModel) Delete(id int) error {
	query := `
		DELETE FROM regions
		WHERE id = $1
	`
	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
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
		return sql.ErrNoRows
	}

	return nil
}

// GetByID retrieves a region by its ID.
func (m *RegionModel) GetByID(id int) (*Region, error) {
	query := `
		SELECT id, name, code
		FROM regions
		WHERE id = $1
	`
	args := []any{id}

	var r Region

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// dryscan pattern
	scan := []any{
		&r.ID,
		&r.Name,
		&r.Code,
	}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(scan...)
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

	return &r, nil
}

// GetAll retrieves all regions from the database.
func (m *RegionModel) GetAll(r *RegionFilter) (Regions, filters.MetaData, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, name, code
		FROM regions
		WHERE (LOWER(name) ILIKE LOWER('%%' || $1 || '%%') OR $1 = '')
		AND (LOWER(code) ILIKE LOWER('%%' || $2 || '%%') OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, r.Default.SortColumn(), r.Default.SortDirection())

	args := []any{r.Name, r.Code, r.Default.Limit(), r.Default.Offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, filters.MetaData{}, err
	}
	defer rows.Close()

	totalRecords := 0
	regions := Regions{}

	for rows.Next() {
		var r Region

		// dryscan pattern
		scan := []any{
			&totalRecords,
			&r.ID,
			&r.Name,
			&r.Code,
		}

		err := rows.Scan(scan...)
		if err != nil {
			return nil, filters.EmptyMetaData, err
		}

		regions = append(regions, r)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.EmptyMetaData, err
	}

	metaData := filters.CalculateMetaData(totalRecords, r.Default.Page, r.Default.PageSize)

	return regions, metaData, nil
}
