package flagstruct

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
)

// TODO support another types

func ReadFromFlags(fs *flag.FlagSet, c any) error {
	v := reflect.ValueOf(c)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return errors.New("c must be pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		flagName, ok := field.Tag.Lookup("flag")
		if !ok {
			continue
		}

		if field.Type.Kind() != reflect.String {
			return fmt.Errorf("field %s must be string", field.Name)
		}

		fieldPtr := v.Field(i).Addr().Interface().(*string)

		fs.StringVar(fieldPtr, flagName, *fieldPtr, field.Tag.Get("flag_usage"))
	}

	return nil
}
