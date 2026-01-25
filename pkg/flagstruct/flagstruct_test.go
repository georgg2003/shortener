package flagstruct_test

import (
	"flag"
	"testing"

	"github.com/georgg2003/shortener/pkg/flagstruct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFromFlags(t *testing.T) {
	type Config struct {
		Name string `flag:"name"`
		Env  string `flag:"env"`
	}

	cfg := Config{}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	err := flagstruct.ReadFromFlags(fs, &cfg)
	require.NoError(t, err)

	err = fs.Parse([]string{"-name=app", "-env=prod"})
	require.NoError(t, err)

	assert.Equal(t, "app", cfg.Name)
	assert.Equal(t, "prod", cfg.Env)
}
