package viperconf

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ramseyjiang/go-micros/shared/apierror"
	"github.com/spf13/viper"
)

const PrefixFile = "FILE:"
const PrefixFileNoTrim = "FILE_NOTRIM:"

func LoadFromFileNoTrim(varname string, filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", apierror.NewAPIError(err, 500, "", "Error read file (%s) for (%s): %v", filename, varname, err)
	}

	return string(data), nil
}

func LoadFromFile(varname string, filename string) (string, error) {
	strdata, err := LoadFromFileNoTrim(varname, filename)
	if err != nil {
		return "", apierror.NewAPIError(err, 500, "", "")
	}

	return strings.TrimSpace(strdata), nil
}

func LoadFromFileFullConfig(thisViper *viper.Viper, cfgLocation string, setDefault bool) error {

	newValue, getErr := LoadFromFileNoTrim("LoadFromFileFullConfig", cfgLocation)
	if getErr != nil {
		return apierror.NewAPIError(getErr, 500, "", "Error loading config from GCP secret: %v", getErr)
	}

	tmpmap := make(map[string]interface{})
	jErr := json.Unmarshal([]byte(newValue), &tmpmap)
	if jErr != nil {
		return apierror.NewAPIError(jErr, 500, "", "Error unmarshalling config from GCP secret: %v", jErr)
	}

	for k, v := range tmpmap {
		if setDefault {
			thisViper.SetDefault(k, v)
		} else {
			thisViper.Set(k, v)
		}
	}

	return nil
}
