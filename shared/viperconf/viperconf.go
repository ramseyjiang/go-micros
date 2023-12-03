package viperconf

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/ramseyjiang/go-micros/shared/apierror"
	"github.com/ramseyjiang/go-micros/shared/helpers"
	"github.com/ramseyjiang/go-micros/shared/srvlog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ViperSetParamsFunc func()

type LogInterface interface {
	Println(v ...interface{})
	Printf(format string, args ...interface{})
}

var (
	ErrRemoteConfig = errors.New("Error loading remote config")
	ErrLocalConfig  = errors.New("Error loading local config")
)

func loadRemoteConfig(thisViper *viper.Viper, configLocation string, setDefaults bool) error {
	upperConfigLocation := strings.ToUpper(configLocation)
	switch {

	case strings.HasPrefix(upperConfigLocation, SecretPrefixGCP):
		return GetGCPSecretFullConfig(viper.GetViper(), configLocation[len(SecretPrefixGCP):], setDefaults)

	case strings.HasPrefix(upperConfigLocation, PrefixFileNoTrim):
		return LoadFromFileFullConfig(viper.GetViper(), configLocation[len(PrefixFileNoTrim):], setDefaults)

	case strings.HasPrefix(upperConfigLocation, PrefixFile):
		return LoadFromFileFullConfig(viper.GetViper(), configLocation[len(PrefixFile):], setDefaults)

	case strings.HasPrefix(upperConfigLocation, secretPrefixFirestore):
		return GetFirestoreFullConfig(viper.GetViper(), configLocation[len(secretPrefixFirestore):], setDefaults)

	}

	return apierror.NewAPIError(nil, 400, "loadRemoteConfig", "Invalid/Unsupported remote config location")
}

var defaultProtectedConfigItems = []string{"api-key", "jwt-passphrase", "db-loc", "sql-loc", "sql-location", "nosql-loc", "nosql-location", "oracle-db", "oracle-dsn", "vault-token", "zendesk-auth"}
var defaultProtectedSubstrings = []string{"password", "passwd", "apikey", "api-key", "secret", "-key", "-passphrase", "-token"}
var ProtectedConfigItems []string

func ShowConfig() {
	srvlog.GlobalLogger.Debug("Running with flags:")

	keys := viper.AllKeys()

	sort.Strings(keys)

	for _, cfgVarName := range keys {
		if fmt.Sprintf("%v", viper.Get(cfgVarName)) != "" {
			if helpers.StringInSliceCaseInsensitive(cfgVarName, ProtectedConfigItems) {
				srvlog.GlobalLogger.Debugf("  %s: ***OMITTED (protected custom key)***", cfgVarName)
				continue
			}

			foundStr := false
			for _, thisStr := range defaultProtectedConfigItems {
				if strings.HasSuffix(strings.ToLower(cfgVarName), thisStr) {
					foundStr = true
					srvlog.GlobalLogger.Debugf("  %s: ***OMITTED (protected key)***", cfgVarName)
					break
				}
			}
			for _, thisStr := range defaultProtectedSubstrings {
				if strings.Contains(strings.ToLower(cfgVarName), thisStr) {
					foundStr = true
					srvlog.GlobalLogger.Debugf("  %s: ***OMITTED (protected sub-string)***", cfgVarName)
					break
				}
			}
			if foundStr {
				continue
			}
		}
		srvlog.GlobalLogger.Debugf("  %s: %v", cfgVarName, viper.Get(cfgVarName))
	}
}

func SetupViper(logger LogInterface, applicationName string, setParamFunction ViperSetParamsFunc) error {
	return SetupViperV2(logger, applicationName, setParamFunction, nil)
}

func SetupViperV2(logger LogInterface, applicationName string, setParamFunction ViperSetParamsFunc, runWhenConfigChanges ViperSetParamsFunc) error {
	var confErr error

	if logger == nil {
		logger = srvlog.GlobalLogger
	}

	if logger == nil {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	viper.SetConfigName(applicationName)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	remoteRPCConfig := os.Getenv("REMOTE_RPC_CONFIG")
	if remoteRPCConfig != "" {
		remoteCfgErr := loadRemoteConfig(viper.GetViper(), remoteRPCConfig, true)
		if remoteCfgErr != nil {
			log.Printf("Remote RPC Config ERROR: %v", remoteCfgErr)
			os.Exit(1)
			// return remoteCfgErr
		}
		srvlog.Infof("Successfully loaded REMOTE_RPC_CONFIG: %s", remoteRPCConfig)
		// Temporary - REMOVE THIS IN FUTURE !!!!
		ConfigLoaded = true
	} else {
		srvlog.Debugf("WARNING: Falling back to legacy rpcConfig mode")
		// Temporary - REMOVE THIS IN FUTURE !!!!
		InitConfig()
	}

	// Inherit RPC Config params as defaults
	if RPCViper != nil {
		for _, cfgVarName := range RPCViper.AllKeys() {
			if !RPCViper.IsSet(cfgVarName) {
				continue
			}
			viper.SetDefault(cfgVarName, RPCViper.Get(cfgVarName))
		}
	}

	remoteConfig := os.Getenv("REMOTE_CONFIG")
	if remoteConfig != "" {
		remoteCfgErr := loadRemoteConfig(viper.GetViper(), remoteConfig, false)
		if remoteCfgErr != nil {
			log.Printf("Remote Config ERROR: %v", remoteCfgErr)
			os.Exit(1)
		}
		srvlog.Infof("Successfully loaded REMOTE_CONFIG: %s", remoteConfig)
	} else {
		ex, _ := os.Executable()
		viper.AddConfigPath(filepath.Dir(ex))
		viper.AddConfigPath("/opt/" + applicationName)
		viper.AddConfigPath(".")
		confErr = viper.ReadInConfig()
		if confErr != nil {
			if _, castOK := confErr.(viper.ConfigFileNotFoundError); !castOK {
				srvlog.Infof("ERROR Reading config: %v", confErr)
				if os.Getenv("CONFIG_OPTIONAL") == "" {
					return ErrLocalConfig
				}
			}
		} else {
			srvlog.Debug("Successfully loaded config from file")
		}
	}

	if viper.Get("environment") == nil {
		viper.SetDefault("environment", "")
	}
	if pflag.Lookup("environment") == nil {
		pflag.String("environment", viper.GetString("environment"), "Environment")
		viper.BindPFlag("environment", pflag.Lookup("environment"))
	}

	if viper.Get("config") == nil {
		viper.SetDefault("config", "")
	}
	if pflag.Lookup("config") == nil {
		pflag.String("config", viper.GetString("config"), "Configuration file location")
		viper.BindPFlag("config", pflag.Lookup("config"))
	}

	if viper.Get("debug") == nil {
		viper.SetDefault("debug", false)
	}
	if pflag.Lookup("debug") == nil {
		pflag.Bool("debug", viper.GetBool("debug"), "Debug mode")
		viper.BindPFlag("debug", pflag.Lookup("debug"))
	}

	if viper.Get("grpc-auth-jwt-file") == nil {
		viper.SetDefault("grpc-auth-jwt-file", "")
	}
	if pflag.Lookup("grpc-auth-jwt-file") == nil {
		pflag.String("grpc-auth-jwt-file", viper.GetString("grpc-auth-jwt-file"), "GRPC authentication credentials file")
		viper.BindPFlag("grpc-auth-jwt-file", pflag.Lookup("grpc-auth-jwt-file"))
	}

	if viper.Get("tracing-enabled") == nil {
		viper.SetDefault("tracing-enabled", false)
	}
	if pflag.Lookup("tracing-enabled") == nil {
		pflag.Bool("tracing-enabled", viper.GetBool("tracing-enabled"), "Enable tracing?")
		viper.BindPFlag("tracing-enabled", pflag.Lookup("tracing-enabled"))
	}

	if viper.Get("export-metrics") == nil {
		viper.SetDefault("export-metrics", false)
	}
	if pflag.Lookup("export-metrics") == nil {
		pflag.Bool("export-metrics", viper.GetBool("export-metrics"), "Export metrics?")
		viper.BindPFlag("export-metrics", pflag.Lookup("export-metrics"))
	}

	if viper.Get("datadog-tracing") == nil {
		viper.SetDefault("datadog-tracing", false)
	}
	if pflag.Lookup("datadog-tracing") == nil {
		pflag.Bool("datadog-tracing", viper.GetBool("datadog-tracing"), "Enable Datadog tracing?")
		viper.BindPFlag("datadog-tracing", pflag.Lookup("datadog-tracing"))
	}

	if viper.Get("datadog-trace-addr") == nil {
		viper.SetDefault("datadog-trace-addr", "datadog:8126")
	}
	if pflag.Lookup("datadog-trace-addr") == nil {
		pflag.String("datadog-trace-addr", viper.GetString("datadog-trace-addr"), "Datadog Trace IP:Port")
		viper.BindPFlag("datadog-trace-addr", pflag.Lookup("datadog-trace-addr"))
	}

	if viper.Get("gcp-project") == nil {
		viper.SetDefault("gcp-project", os.Getenv("GCP_PROJECT"))
	}
	if pflag.Lookup("gcp-project") == nil {
		pflag.String("gcp-project", viper.GetString("gcp-project"), "Google Cloud Project ID")
		viper.BindPFlag("gcp-project", pflag.Lookup("gcp-project"))
	}

	setParamFunction()

	pflag.Parse()

	viper.OnConfigChange(func(e fsnotify.Event) {
		srvlog.Infof("Config file changed: %s", e.Name)
		if runWhenConfigChanges != nil {
			runWhenConfigChanges()
		}
	})

	if viper.GetString("config") != "" {
		srvlog.Infof("Loading config from: %s", viper.GetString("config"))
		viper.SetConfigFile(viper.GetString("config"))
		confErr := viper.ReadInConfig()
		if confErr != nil {
			srvlog.Infof("Config load error: %v", confErr)
			return ErrLocalConfig
		}
		viper.WatchConfig()
	} else {
		if confErr != nil {
			srvlog.Debug("WARNING: No config file to read! - Not auto reloading config file changes")
		} else {
			viper.WatchConfig()
		}
		if viper.ConfigFileUsed() != "" {
			srvlog.Infof("Loaded default config from: %s", viper.ConfigFileUsed())
		}
	}

	ParseViperConfigForRemoteVars(viper.GetViper())

	setProxyFlags(viper.GetViper())

	return nil
}

// ParseViperConfigForRemoteVars parses through params and tried to fetch them from remote locations
func ParseViperConfigForRemoteVars(viperCfg *viper.Viper) {
	for _, cfgVarName := range viperCfg.AllKeys() {
		if strings.ToUpper(cfgVarName) == "REMOTE_RPC_CONFIG" ||
			strings.ToUpper(cfgVarName) == "REMOTE_CONFIG" {
			continue
		}
		cfgVarValue := viperCfg.Get(cfgVarName)
		cfgVarValueString, isString := cfgVarValue.(string)
		if !isString {
			// only look at config parameters that are strings
			continue
		}
		switch {
		case strings.HasPrefix(cfgVarValueString, SecretPrefixGCP):
			secretName := cfgVarValueString[len(SecretPrefixGCP):]
			if newValue, getErr := GetGCPSecret(cfgVarName, secretName); getErr == nil {
				viperCfg.Set(cfgVarName, newValue)
			}

		case strings.HasPrefix(cfgVarValueString, PrefixFileNoTrim):
			secretName := cfgVarValueString[len(PrefixFileNoTrim):]
			if newValue, getErr := LoadFromFileNoTrim(cfgVarName, secretName); getErr == nil {
				viperCfg.Set(cfgVarName, newValue)
			}

		case strings.HasPrefix(cfgVarValueString, PrefixFile):
			secretName := cfgVarValueString[len(PrefixFile):]
			if newValue, getErr := LoadFromFile(cfgVarName, secretName); getErr == nil {
				viperCfg.Set(cfgVarName, newValue)
			}
		}

	}
}

// setProxyFlags allows for setting the HTTP_PROXY without having to set it on the ENV, but rather in the settings.
// This is for test debugging in for example vscode where setting an env variable isn't immediately obvious and follows go defaults for proxying.
// i.e. the default HTTP client will respect these settings
func setProxyFlags(viperCfg *viper.Viper) {

	// INFO FROM GO STANDARD LIBRARY
	// ---------
	// ProxyFromEnvironment returns the URL of the proxy to use for a
	// given request, as indicated by the environment variables
	// HTTP_PROXY, HTTPS_PROXY and NO_PROXY (or the lowercase versions
	// thereof). HTTPS_PROXY takes precedence over HTTP_PROXY for https
	// requests.
	//
	// The environment values may be either a complete URL or a
	// "host[:port]", in which case the "http" scheme is assumed.
	// An error is returned if the value is a different form.
	//
	// A nil URL and nil error are returned if no proxy is defined in the
	// environment, or a proxy should not be used for the given request,
	// as defined by NO_PROXY.
	//
	// As a special case, if req.URL.Host is "localhost" (with or without
	// a port number), then a nil URL and nil error will be returned.

	httpProxy := viperCfg.GetString("HTTP_PROXY")
	if httpProxy == "" {
		httpProxy = viperCfg.GetString("http_proxy")
	}
	if httpProxy != "" {
		err := os.Setenv("HTTP_PROXY", httpProxy)
		if err != nil {
			srvlog.Info("HTTP_PROXY is set. The default HTTP client will be proxying requests via ", httpProxy)
		}
	}
	httpsProxy := viperCfg.GetString("HTTPS_PROXY")
	if httpsProxy == "" {
		httpsProxy = viperCfg.GetString("https_proxy")
	}
	if httpsProxy != "" {
		err := os.Setenv("HTTPS_PROXY", httpsProxy)
		if err != nil {
			srvlog.Info("HTTPS_PROXY is set. The default HTTP client will be proxying https requests via ", httpsProxy)
		}
	}
	noProxy := viperCfg.GetString("NO_PROXY")
	if noProxy == "" {
		noProxy = viperCfg.GetString("no_proxy")
	}
	if noProxy != "" {
		err := os.Setenv("NO_PROXY", noProxy)
		if err != nil {
			srvlog.Info("NO_PROXY is set. The default HTTP client will ignore the following endpoints when proxying ", httpsProxy)
		}
	}
}
