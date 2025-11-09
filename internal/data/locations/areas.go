// File: internal/data/locations/districts.go
package locations

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
)

/****************************************************************************************
 *										Declarations									*
 ***************************************************************************************/
// Coordinates represents geographical coordinates.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AreaType represents the different types of areas.
type AreaType string

// Area type constants.
const (
	AreaTypeCity    AreaType = "city"
	AreaTypeTown    AreaType = "town"
	AreaTypeVillage AreaType = "village"
)

// Area represents a geographical area within a district.
type Area struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	DistrictID  int         `json:"district_id"`
	AreaType    AreaType    `json:"area_type"`
	Coordinates Coordinates `json:"coordinates"`
	IsActive    *bool       `json:"is_active"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// Areas is a slice of Area.
type Areas []Area

// AreaFilter represents filtering options for querying areas.
type AreaFilter struct {
	Name       string
	DistrictID *int
	AreaType   *AreaType
	IsActive   *bool
	Default    filters.Filters
}

// AreasModel represents the model for areas.
type AreaModel struct {
	DB *sql.DB
}

// ValidateCoordinates validates the latitude and longitude values.
func ValidateCoordinates(v *validator.Validator, c Coordinates) {
	if c.Latitude != 0 || c.Longitude != 0 {
		v.Check(v.Matches(fmt.Sprintf("%f", c.Latitude), validator.LatitudeRX), "latitude", "must be a valid latitude")
		v.Check(v.Matches(fmt.Sprintf("%f", c.Longitude), validator.LongitudeRX), "longitude", "must be a valid longitude")
	}
}

// ValidateArea validates the fields of an Area.
func ValidateArea(v *validator.Validator, a *Area) {
	v.Check(a.Name != "", "name", "must be provided")
	v.Check(len(a.Name) <= 255, "name", "must not be more than 255 bytes long")
	v.Check(a.DistrictID > 0, "district_id", "must be provided and greater than zero")
	v.Check(a.AreaType == AreaTypeCity || a.AreaType == AreaTypeTown || a.AreaType == AreaTypeVillage, "area_type", "must be a valid area type")
	ValidateCoordinates(v, a.Coordinates)
}

/****************************************************************************************
 *										Methods											*
 ***************************************************************************************/

// Insert adds a new area to the database.
func (m *AreaModel) Insert(a Area) error {
	query := `
		INSERT INTO areas (name, district_id, area_type, latitude, longitude, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	args := []any{a.Name, a.DistrictID, a.AreaType, a.Coordinates.Latitude, a.Coordinates.Longitude, a.IsActive}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "areas_name_key"):
			return errors.ErrDuplicateValue("name")
		case errors.IsForeignKeyViolation(err):
			return errors.ErrInvalidDistrictID
		default:
			return err
		}
	}

	return nil
}

// Update modifies an existing area in the database.
func (m *AreaModel) Update(a *Area) error {
	query := `
		UPDATE areas
		SET name = $1, district_id = $2, area_type = $3, latitude = $4, longitude = $5, is_active = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	args := []any{a.Name, a.DistrictID, a.AreaType, a.Coordinates.Latitude, a.Coordinates.Longitude, a.IsActive, a.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&a.UpdatedAt)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "areas_name_key"):
			return errors.ErrDuplicateName
		case errors.IsForeignKeyViolation(err):
			return errors.ErrInvalidDistrictID
		case errors.IsEditConflict(err):
			return errors.ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete Permanently removes an area from the database.
func (m *AreaModel) Delete(id int) error {
	query := `
		DELETE FROM areas
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
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

// Get retrieves a specific area by its ID.
func (m *AreaModel) GetByID(id int) (*Area, error) {
	query := `
		SELECT id, name, district_id, area_type, latitude, longitude, is_active, created_at, updated_at
		FROM areas
		WHERE id = $1
	`

	var a Area

	scan := []any{
		&a.ID,
		&a.Name,
		&a.DistrictID,
		&a.AreaType,
		&a.Coordinates.Latitude,
		&a.Coordinates.Longitude,
		&a.IsActive,
		&a.CreatedAt,
		&a.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(scan...)
	if err != nil {
		switch {
		case errors.ErrNoRows(err):
			return nil, errors.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &a, nil
}

// GetAll retrieves all areas matching the provided filter criteria.
func (m *AreaModel) GetAll(filter *AreaFilter) (Areas, filters.MetaData, error) {
	query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, name, district_id, area_type, latitude, longitude, is_active, created_at, updated_at
        FROM areas
        WHERE ($1 = '' OR LOWER(name) ILIKE LOWER('%%' || $1 || '%%'))
        AND ($2::bigint IS NULL OR district_id = $2)
        AND ($3::text IS NULL OR area_type = $3)
        AND ($4::boolean IS NULL OR is_active = $4)
        ORDER BY %s %s, id ASC
        LIMIT $5 OFFSET $6`, filter.Default.SortColumn(), filter.Default.SortDirection())

	args := []any{
		filter.Name,
		filter.DistrictID,
		filter.AreaType,
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
	areas := Areas{}

	for rows.Next() {
		var a Area
		scan := []any{
			&totalRecords,
			&a.ID,
			&a.Name,
			&a.DistrictID,
			&a.AreaType,
			&a.Coordinates.Latitude,
			&a.Coordinates.Longitude,
			&a.IsActive,
			&a.CreatedAt,
			&a.UpdatedAt,
		}

		err := rows.Scan(scan...)
		if err != nil {
			return nil, filters.EmptyMetaData, err
		}

		areas = append(areas, a)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.EmptyMetaData, err
	}

	metadata := filters.CalculateMetaData(totalRecords, filter.Default.Page, filter.Default.PageSize)

	return areas, metadata, nil
}
