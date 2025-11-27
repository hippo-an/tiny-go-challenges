package config

import "errors"

var (
	ErrNameRequired       = errors.New("project name is required")
	ErrModulePathRequired = errors.New("module path is required")
	ErrInvalidFramework   = errors.New("invalid framework selection")
	ErrInvalidDatabase    = errors.New("invalid database selection")
	ErrOutputDirExists    = errors.New("output directory already exists")
)
