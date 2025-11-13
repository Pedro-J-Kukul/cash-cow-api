// File: internal/data/models.go
package data

import (
	"database/sql"

	"github.com/Pedro-J-Kukul/cash-cow-api/internal/data/cattle"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/data/locations"
	"github.com/Pedro-J-Kukul/cash-cow-api/internal/data/users"
)

/****************************************************************************************
 *										Declarations									*
 ***************************************************************************************/

// Models is a wrapper for all data models.
type Models struct {
	Cattle      cattle.CattleModel
	Breeds      cattle.BreedModel
	Users       users.UserModel
	Tokens      users.TokenModel
	Permissions users.PermissionModel
	Areas       locations.AreaModel
	Regions     locations.RegionModel
}

// NewModels initializes and returns a Models struct.

func NewModels(db *sql.DB) Models {
	return Models{
		Cattle:      cattle.CattleModel{DB: db},
		Breeds:      cattle.BreedModel{DB: db},
		Users:       users.UserModel{DB: db},
		Tokens:      users.TokenModel{DB: db},
		Permissions: users.PermissionModel{DB: db},
		Areas:       locations.AreaModel{DB: db},
		Regions:     locations.RegionModel{DB: db},
	}
}
