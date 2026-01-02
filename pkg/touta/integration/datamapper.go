package integration

import (
	"fmt"

	"github.com/toutaio/toutago-datamapper/engine"
)

// DataMapperConfig holds configuration for creating a datamapper instance.
type DataMapperConfig struct {
	ConfigPath string // Path to datamapper YAML configuration file
}

// NewDataMapper creates a new datamapper from a configuration file.
// The datamapper uses a configuration-driven approach where mappings are
// defined in YAML files.
//
// Example:
//
//	mapper, err := NewDataMapper(&DataMapperConfig{
//	    ConfigPath: "config/datamapper.yaml",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer mapper.Close()
//
//	// Use the mapper
//	user := &User{ID: "1", Name: "John"}
//	err = mapper.Insert(ctx, "User", user)
//
// For detailed configuration format, see:
// https://github.com/toutaio/toutago-datamapper
func NewDataMapper(config *DataMapperConfig) (*engine.Mapper, error) {
	if config.ConfigPath == "" {
		return nil, fmt.Errorf("datamapper config path is required")
	}
	
	return engine.NewMapper(config.ConfigPath)
}
