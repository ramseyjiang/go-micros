package viperconf

import (
	"context"
	"encoding/json"

	"github.com/ramseyjiang/go-micros/shared/apierror"
	"github.com/ramseyjiang/go-micros/shared/srvlogs/v2"
	"github.com/spf13/viper"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const SecretPrefixGCP = "SECRET_GCP:"

var gcpSecretManagerClient *secretmanager.Client

func GetGCPSecret(varname string, name string) (string, error) {
	// Create/Get the client.
	if gcpSecretManagerClient == nil {
		var clientErr error
		gcpSecretManagerClient, clientErr = secretmanager.NewClient(context.Background())
		if clientErr != nil {
			return name, apierror.NewAPIError(clientErr, 500, "", "Error initialising GCP Secret manager client for (%s): %v", varname, clientErr)
		}
	}

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, accessErr := gcpSecretManagerClient.AccessSecretVersion(context.Background(), req)
	if accessErr != nil {
		srvlogs.Errorf("gcpSecretManagerClient.AccessSecretVersion (%s) ERROR: %v", varname, accessErr)
		return name, apierror.NewAPIError(accessErr, 500, "", "Error reading GCP secret for %s", varname)
	}

	return string(result.Payload.Data), nil
}

func GetGCPSecretFullConfig(thisViper *viper.Viper, cfgLocation string, setDefault bool) error {
	newValue, getErr := GetGCPSecret("GetGCPSecretFullConfig", cfgLocation)
	if getErr != nil {
		return apierror.NewAPIError(getErr, 500, "", "Error loading config from GCP secret (%s): %v", cfgLocation, getErr)
	}

	tmpmap := make(map[string]interface{})
	jErr := json.Unmarshal([]byte(newValue), &tmpmap)
	if jErr != nil {
		return apierror.NewAPIError(jErr, 500, "", "Error unmarshalling config from GCP secret (%s): %v", cfgLocation, jErr)
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
