package configurator

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfiguratorOption func(cfg any) error

func Setup(cfg any, opts ...ConfiguratorOption) error {
	err := isValidType(cfg)
	if err != nil {
		return err
	}
	err = setupDefaultConfig(cfg)
	if err != nil {
		return err
	}
	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func Dump(cfg any, filePath string) error {
	err := isValidType(cfg)
	if err != nil {
		return err
	}

	// create parent directories if not exist
	rootPath := path.Dir(filePath)
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		if err := os.MkdirAll(rootPath, os.FileMode(0755)); err != nil {
			return fmt.Errorf("%w: %s", ErrFailedToDumpConfig, err.Error())
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToDumpConfig, err.Error())
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	enc.SetIndent(2)

	err = enc.Encode(cfg)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToDumpConfig, err.Error())
	}
	return nil
}

func isValidType(cfg any) error {
	cfgValue := reflect.ValueOf(cfg)
	if cfgValue.Kind() != reflect.Pointer || reflect.Indirect(reflect.ValueOf(cfgValue)).Kind() != reflect.Struct {
		return fmt.Errorf("%w: %s", ErrUnsupportedConfigType, "config muse be defined as a struct and passed as a pointer")
	}

	return nil
}

func setupDefaultConfig(cfg any) error {
	err := iterateConfig(cfg, "")
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToCreateDefaultConfig, err.Error())
	}
	return nil
}

func iterateConfig(input any, fieldPrefix string) error {
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

		if !field.CanAddr() || !field.CanInterface() {
			return fmt.Errorf("%w: %s", ErrUnsupportedConfigType, fPrefix)
		}

		switch field.Kind() {
		case reflect.Struct:
			err := iterateConfig(field.Addr().Interface(), fPrefix)
			if err != nil {
				return err
			}
		// supported types
		case reflect.Bool, reflect.String, reflect.Int, reflect.Int64, reflect.Float64:
			value := fieldStruct.Tag.Get("default")
			if len(value) == 0 {
				return fmt.Errorf("%w: %s", ErrDefaultNotFound, fPrefix)
			}
			if value == Empty {
				value = ""
			}
			if err := yaml.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
				return fmt.Errorf("%w: %s", ErrUnmarshal, fPrefix)
			}
		case reflect.Slice:
			value := fieldStruct.Tag.Get("default")
			var values []string

			if len(value) == 0 {
				return fmt.Errorf("%w: %s", ErrDefaultNotFound, fPrefix)
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
