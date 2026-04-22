// Package flags
package flags

import (
	"flag"
	"fmt"
)

func Parse() (dbPath, configPath *string, err error) {
	dbPath = flag.String("db", "", "path to sqlite database")
	configPath = flag.String("cf", "", "path to config json file")
	flag.Parse()

	if *dbPath == "" {
		return nil, nil, fmt.Errorf("error: db path is required.\nUse '--help' to know more")
	}

	if *configPath == "" {
		return nil, nil, fmt.Errorf("error: config path is required.\nUse '--help' to know more")
	}

	return dbPath, configPath, nil
}
