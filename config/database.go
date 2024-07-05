package config

import (
	"os"

	"github.com/caesar-rocks/orm"
)

func ProvideDatabase(env *EnvironmentVariables) *orm.Database {
	return orm.NewDatabase(&orm.DatabaseConfig{
		DBMS:  orm.DBMS(env.DBMS),
		DSN:   env.DSN,
		Debug: os.Getenv("DEBUG_DATABASE") == "true",
	})
}
