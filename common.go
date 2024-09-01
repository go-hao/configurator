package configurator

import "errors"

var (
	ErrUnsupportedConfigType       = errors.New("unsupported config type")
	ErrFailedToDumpConfig          = errors.New("failed to dump config")
	ErrFailedToCreateDefaultConfig = errors.New("failed to create default config")
	ErrFailedToUpdateFromFile      = errors.New("failed to update config from file")
	ErrFailedToUpdateFromEnv       = errors.New("failed to update config from env")
	ErrFailedToUpdateFromRemote    = errors.New("failed to update config from remote")
	ErrDefaultNotFound             = errors.New("default value is required")
	ErrInvalidSliceFormat          = errors.New("invalid slice format")
	ErrUnmarshal                   = errors.New("unmarshal error")
	ErrConfigurationError          = errors.New("configuration error")
)

const Empty string = "EMPTY"
