package playerHandler

import (
	"github.com/alvaromfcunha/lol-elo-police/pkg/lol"
	"gorm.io/gorm"
)

type PlayerHandler struct {
	Db     *gorm.DB
	LolApi lol.LolApi
}
