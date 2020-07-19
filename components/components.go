package components

import (
	"database/sql"
	"github.com/spf13/viper"
)

type Components struct {
	Settings        *viper.Viper
	DB *sql.DB
}

