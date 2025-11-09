// File: internal/data/cattle/cattle.go
package cattle

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/data/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
)

/****************************************************************************************
 *										Declarations									*
 ***************************************************************************************/
type Sex string

const (
	Male    Sex = "male"
	Female  Sex = "female"
	Unknown Sex = "unknown"
)

// Cattle represents a cattle entity.
type Cattle struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	BreedID   int    `json:"breed_id"`
	TagNumber string `json:"tag_number"`
	Sex       Sex    `json:"sex"`
	AgeMonths int    `json:"age_months"`
	WeightKg  int    `json:"weight_kg"`
	// For simplicity, vaccinations are represented as a comma-separated string.
	// In future, will be normalized into a separate table.
	Vaccinations string `json:"vaccinations"`
	// HealthRecords are represented as a comma-separated string.
	// In future, will be normalized into a separate table.
	MedicalHistory string `json:"medical_history"`
	IsPregnant     *bool  `json:"is_pregnant"`
	IsCastrated    *bool  `json:"is_castrated"`
	// For Simplicity IsActive is for soft deletion, sold or deceased cattle.
	IsActive  *bool  `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Cattles is a slice of Cattle.
type Cattles []Cattle

type CattleFilter struct {
	OwnerID     *int
	BreedID     *int
	TagNumber   string
	Sex         *Sex
	AgeMonths   *int
	WeightKg    *int
	IsPregnant  *bool
	IsCastrated *bool
	IsActive    *bool
	CreatedAt   string
	UpdatedAt   string
	Default     filters.Filters
}

// CattleModel represents the model for cattle.
type CattleModel struct {
	DB *sql.DB
}

/****************************************************************************************
 *										Methods											*
 ***************************************************************************************/

// Insert inserts a new cattle record into the database.
func (m *CattleModel) Insert(c *Cattle) error {
	query := `
		INSERT INTO cattle (
			owner_id, breed_id, tag_number, sex, age_months, weight_kg,
			vaccinations, medical_history, is_pregnant, is_castrated, is_active,
			created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query,
		c.OwnerID, c.BreedID, c.TagNumber, c.Sex, c.AgeMonths, c.WeightKg,
		c.Vaccinations, c.MedicalHistory, c.IsPregnant, c.IsCastrated, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "cattle_tag_number_key"):
			return errors.ErrDuplicateValue("tag_number")
		case errors.IsForeignKeyViolation(err):
			return errors.ErrForeignKeyViolation
		default:
			return err
		}
	}
	return nil
}

// Update updates an existing cattle record in the database.
func (m *CattleModel) Update(c *Cattle) error {
	query := `
		UPDATE cattle
		SET
			owner_id = $1,
			breed_id = $2,
			tag_number = $3,
			sex = $4,
			age_months = $5,
			weight_kg = $6,
			vaccinations = $7,
			medical_history = $8,
			is_pregnant = $9,
			is_castrated = $10,
			is_active = $11,
			updated_at = NOW()
		WHERE id = $12
		RETURNING updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query,
		c.OwnerID, c.BreedID, c.TagNumber, c.Sex, c.AgeMonths, c.WeightKg,
		c.Vaccinations, c.MedicalHistory, c.IsPregnant, c.IsCastrated, c.IsActive,
		c.ID,
	).Scan(&c.UpdatedAt)
	if err != nil {
		switch {
		case errors.IsUniqueViolation(err, "cattle_tag_number_key"):
			return errors.ErrDuplicateValue("tag_number")
		case errors.IsEditConflict(err):
			return errors.ErrEditConflict
		case errors.IsForeignKeyViolation(err):
			return errors.ErrForeignKeyViolation
		default:
			return err
		}
	}
	return nil
}

// Delete permanently deletes a cattle record from the database.
func (m *CattleModel) Delete(id int) error {
	query := `DELETE FROM cattle WHERE id = $1`
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

// GetByID retrieves a cattle record by its ID.
func (m *CattleModel) GetByID(id int) (*Cattle, error) {
	query := `
		SELECT
			id, owner_id, breed_id, tag_number, sex, age_months, weight_kg,
			vaccinations, medical_history, is_pregnant, is_castrated, is_active,
			created_at, updated_at
		FROM cattle
		WHERE id = $1
	`
	var c Cattle
	scan := []any{
		&c.ID, &c.OwnerID, &c.BreedID, &c.TagNumber, &c.Sex, &c.AgeMonths, &c.WeightKg,
		&c.Vaccinations, &c.MedicalHistory, &c.IsPregnant, &c.IsCastrated, &c.IsActive,
		&c.CreatedAt, &c.UpdatedAt,
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
	return &c, nil
}

// GetAll retrieves all cattle records from the database with optional filtering.
func (m *CattleModel) GetAll(filter *CattleFilter) (Cattles, filters.MetaData, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),
			id, owner_id, breed_id, tag_number, sex, age_months, weight_kg,
			vaccinations, medical_history, is_pregnant, is_castrated, is_active,
			created_at, updated_at
		FROM cattle
		WHERE
			($1::int IS NULL OR owner_id = $1) AND
			($2::int IS NULL OR breed_id = $2) AND
			($3::text IS NULL OR LOWER(tag_number) LIKE LOWER('%%' || $3 || '%%')) AND
			($4::text IS NULL OR sex = $4) AND
			($5::int IS NULL OR age_months = $5) AND
			($6::int IS NULL OR weight_kg = $6) AND
			($7::boolean IS NULL OR is_pregnant = $7) AND
			($8::boolean IS NULL OR is_castrated = $8) AND
			($9::boolean IS NULL OR is_active = $9)
		ORDER BY %s %s, id ASC
		LIMIT $10 OFFSET $11`, filter.Default.SortColumn(), filter.Default.SortDirection())

	args := []any{
		filter.OwnerID,
		filter.BreedID,
		filter.TagNumber,
		filter.Sex,
		filter.AgeMonths,
		filter.WeightKg,
		filter.IsPregnant,
		filter.IsCastrated,
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
	cattles := Cattles{}
	for rows.Next() {
		var c Cattle
		scan := []any{
			&c.ID, &c.OwnerID, &c.BreedID, &c.TagNumber, &c.Sex, &c.AgeMonths, &c.WeightKg,
			&c.Vaccinations, &c.MedicalHistory, &c.IsPregnant, &c.IsCastrated, &c.IsActive,
			&c.CreatedAt, &c.UpdatedAt,
		}
		err := rows.Scan(scan...)
		if err != nil {
			return nil, filters.EmptyMetaData, err
		}
		cattles = append(cattles, c)
	}
	if err = rows.Err(); err != nil {
		return nil, filters.EmptyMetaData, err
	}

	metaData := filters.CalculateMetaData(totalRecords, filter.Default.Page, filter.Default.PageSize)
	return cattles, metaData, nil
}
