// Listings
package listings

import (
	"context"
	"database/sql"
	"time"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/data/errors"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/filters"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/shared/validator"
)

/****************************************************************************************
 *										Declarations									*
 ***************************************************************************************/
type CattleClass string

const ( // Calves (Under 11 months)
	CattleClassCalf       CattleClass = "calf"
	CattleClassHeiferCalf CattleClass = "heifer_calf"
	CattleClassSteerCalf  CattleClass = "steer_calf"
	CattleClassBullCalf   CattleClass = "bull_calf"

	// Weaners (Typically 11–12 months or weaned)
	CattleClassWeaner       CattleClass = "weaner"
	CattleClassHeiferWeaner CattleClass = "heifer_weaner"
	CattleClassSteerWeaner  CattleClass = "steer_weaner"
	CattleClassBullWeaner   CattleClass = "bull_weaner"

	// Adult Cattle (Over 12 months)
	CattleClassYearling     CattleClass = "yearling"
	CattleClassHeifer       CattleClass = "heifer"
	CattleClassCow          CattleClass = "cow"
	CattleClassSpayedHeifer CattleClass = "spayed_heifer"
	CattleClassSpayedCow    CattleClass = "spayed_cow"
	CattleClassSteer        CattleClass = "steer"
	CattleClassBull         CattleClass = "bull"
)

// ListingPrices represents a listing for cattle sale.
type ListingPrice struct {
	ListingID   int64       `json:"listing_id"`
	CattleClass CattleClass `json:"cattle_class"`
	PricePerKg  int64       `json:"price_per_kg"`
	Quantity    int64       `json:"quantity"`
}

type ListingsPrices []ListingPrice

type ListingPricesModel struct {
	DB *sql.DB
}

type ListingPricesFilter struct {
	CattleClass *CattleClass
	Default     filters.Filters
}

// ValidateListingPrice validates the listing price fields.
func ValidateListingPrice(v *validator.Validator, lp *ListingPrice) {
	// Cattle class must be provided
	v.Check(lp.CattleClass != "", "cattle_class", "must be provided")

	// Price per kg must be greater than zero
	v.Check(lp.PricePerKg > 0, "price_per_kg", "must be greater than zero")
}

/****************************************************************************************
 *									Database Operations								*
 ***************************************************************************************/
//  Insert
func (lpm ListingPricesModel) Insert(lp *ListingPrice) error {
	query := `INSERT INTO listing_prices (listing_id, cattle_class, price_per_kg, quantity)
			  VALUES ($1, $2, $3. $4)`

	_, err := lpm.DB.Exec(query, lp.ListingID, lp.CattleClass, lp.PricePerKg, lp.Quantity)

	if err != nil {
		switch {
		case errors.IsForeignKeyViolation(err):
			return errors.ErrInvalidDistrictID
		default:
			return err
		}
	}
	return nil
}

// Update
func (lpm *ListingPricesModel) Update(lp *ListingPrice) error {
	query := `UPDATE listing_prices
			SET cattle_class = $1, price_per_kg = $2, quantity = $3
			WHERE listing_price = $4`

	args := []any{}
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

// Delete
// GetAll
// GetByID
