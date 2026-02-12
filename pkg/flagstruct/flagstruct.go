// Пакет, который упрощает парсинг конфига приложения из флагов.
// Пример конфига:
//
//	type Config struct {
//		ListenAddr       string `flag:"a" flag_usage:"listen addres"`
//		AuditURL         string `flag:"audit-url" flag_usage:"url to send audit logs"`
//		AuditStubEnabled bool   `mapstructure:"audit_stub_enabled" env:"AUDIT_STUB_ENABLED"`
//		DebugAddr        string `flag:"debug-addr"`
//	}
//
// TODO: Сейчас поддерживает только парсинг в строки. Нужно добавить поддержку других типов
package flagstruct

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
)

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
