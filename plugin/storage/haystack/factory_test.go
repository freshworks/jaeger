package haystack

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactory(t *testing.T) {
	f := NewFactory()
	assert.Empty(t, f)

	flagSet := &flag.FlagSet{}
	f.AddFlags(flagSet)
	assert.NotEmpty(t, flagSet)

	v := &viper.Viper{}
	f.InitFromViper(v)
	assert.NotEmpty(t, flagSet)

}
