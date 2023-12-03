package viperconf

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// New creates a new global viper configuration item, using the type of the defaultValue
func New(configVarName string, defaultValue interface{}, configDescription string) {
	viper.SetDefault(configVarName, defaultValue)

	switch val := defaultValue.(type) {
	case string:
		pflag.String(configVarName, viper.GetString(configVarName), configDescription)
	case bool:
		pflag.Bool(configVarName, viper.GetBool(configVarName), configDescription)
	case float64:
		pflag.Float64(configVarName, viper.GetFloat64(configVarName), configDescription)
	case int64:
		pflag.Int64(configVarName, viper.GetInt64(configVarName), configDescription)
	case int32:
		pflag.Int32(configVarName, viper.GetInt32(configVarName), configDescription)
	case int:
		pflag.Int(configVarName, viper.GetInt(configVarName), configDescription)
	case []int:
		pflag.IntSlice(configVarName, viper.GetIntSlice(configVarName), configDescription)
	case []string:
		pflag.StringSlice(configVarName, viper.GetStringSlice(configVarName), configDescription)
	case time.Duration:
		pflag.Duration(configVarName, viper.GetDuration(configVarName), configDescription)
	case map[string]string:
		pflag.StringToString(configVarName, viper.GetStringMapString(configVarName), configDescription)
	case uint64:
		pflag.Uint64(configVarName, viper.GetUint64(configVarName), configDescription)
	case uint32:
		pflag.Uint32(configVarName, viper.GetUint32(configVarName), configDescription)
	case uint:
		pflag.Uint(configVarName, viper.GetUint(configVarName), configDescription)
	default:
		_ = val
		panic("Unknown/unsupported type for config variable " + configVarName)
	}

	viper.BindPFlag(configVarName, pflag.Lookup(configVarName))
}

// NewWithEnv creates a new global viper configuration item, using the type of the defaultValue and binds it to the supplied envVarName
func NewWithEnv(configVarName string, defaultValue interface{}, configDescription string, envVarName string) {

	New(configVarName, defaultValue, configDescription)
	viper.MustBindEnv(configVarName, envVarName)
}
