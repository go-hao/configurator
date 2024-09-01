package configurator

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func WithUpdateFromFile(filePath string) ConfiguratorOption {
	return func(cfg any) error {
		configFile, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrFailedToUpdateFromFile, err.Error())
		}
		err = yaml.Unmarshal(configFile, cfg)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrFailedToUpdateFromFile, err.Error())
		}
		return nil
	}
}

func WithUpdateFromEnv(envPrefix string) ConfiguratorOption {
	return func(cfg any) error {
		err := iterateEnvConfig(cfg, "", envPrefix)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrFailedToUpdateFromEnv, err.Error())
		}
		return nil
	}
}

func iterateEnvConfig(input any, fieldPrefix string, keyPrefix string) error {
	inputValue := reflect.Indirect(reflect.ValueOf(input))
	inputType := inputValue.Type()

	// loop through all fields
	for i := 0; i < inputType.NumField(); i++ {
		field := inputValue.Field(i)
		fieldStruct := inputType.Field(i)

		fPrefix := inputType.Field(i).Name

		if len(fieldPrefix) != 0 {
			fPrefix = fmt.Sprintf("%s.%s", fieldPrefix, fPrefix)
		}

		kPrefix := fieldStruct.Tag.Get("yaml")
		if len(kPrefix) == 0 {
			kPrefix = inputType.Field(i).Name
		}

		if len(keyPrefix) != 0 {
			kPrefix = fmt.Sprintf("%s_%s", keyPrefix, kPrefix)
			kPrefix = strings.Replace(strings.ToUpper(kPrefix), "-", "__", -1)
		}

		if !field.CanAddr() || !field.CanInterface() {
			return fmt.Errorf("%w: %s", ErrUnsupportedConfigType, fPrefix)
		}

		switch field.Kind() {
		case reflect.Struct:
			err := iterateEnvConfig(field.Addr().Interface(), fPrefix, kPrefix)
			if err != nil {
				return err
			}
		// supported types
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int64, reflect.Float64:
			value := os.Getenv(kPrefix)
			if len(value) == 0 {
				continue
			}
			if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
				return fmt.Errorf("%w: %s", ErrUnmarshal, fPrefix)
			}
		case reflect.Slice:
			value := os.Getenv(kPrefix)
			var values []string

			if len(value) == 0 {
				continue
			}

			if !strings.HasPrefix(value, "[") || !strings.HasSuffix(value, "]") {
				return fmt.Errorf("%w: %s", ErrInvalidSliceFormat, fPrefix)
			}

			if value == Empty || value == "[]" {
				values = []string{}
			} else {
				value, _ = strings.CutPrefix(value, "[")
				value, _ = strings.CutSuffix(value, "]")
				value = strings.Replace(value, Empty, "", -1)
				values = strings.Split(value, ",")
			}

			marshaled, _ := yaml.Marshal(values)
			if err := yaml.Unmarshal(marshaled, field.Addr().Interface()); err != nil {
				return fmt.Errorf("%w: %s", ErrUnmarshal, fPrefix)
			}
		default:
			return fmt.Errorf("%w: %s", ErrUnsupportedConfigType, fPrefix)
		}
	}
	return nil
}
