package juice

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"
)

// Environment defines a environment.
// It contains a database connection configuration.
type Environment struct {
	// DataSource is a string in a driver-specific format.
	DataSource string

	// Driver is a driver for
	Driver string

	// MaxIdleConnNum is a maximum number of idle connections.
	MaxIdleConnNum int

	// MaxOpenConnNum is a maximum number of open connections.
	MaxOpenConnNum int

	// MaxConnLifetime is a maximum lifetime of a connection.
	MaxConnLifetime int

	// MaxIdleConnLifetime is a maximum lifetime of an idle connection.
	MaxIdleConnLifetime int

	// attrs is a map of attributes.
	attrs map[string]string
}

// setAttr sets a value of the attribute.
func (e *Environment) setAttr(key, value string) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}
	e.attrs[key] = value
}

// Attr returns a value of the attribute.
func (e *Environment) Attr(key string) string {
	return e.attrs[key]
}

// ID returns a identifier of the environment.
func (e *Environment) ID() string {
	return e.Attr("id")
}

// provider is a environment value provider.
// It provides a value of the environment variable.
func (e *Environment) provider() EnvValueProvider {
	return GetEnvValueProvider(e.Attr("provider"))
}

// Connect returns a database connection.
func (e *Environment) Connect() (*sql.DB, error) {
	db, err := sql.Open(e.Driver, e.DataSource)
	if err != nil {
		return nil, err
	}
	if e.MaxIdleConnNum > 0 {
		db.SetMaxIdleConns(e.MaxIdleConnNum)
	}
	if e.MaxOpenConnNum > 0 {
		db.SetMaxOpenConns(e.MaxOpenConnNum)
	}
	if e.MaxConnLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(e.MaxConnLifetime) * time.Second)
	}
	if e.MaxIdleConnLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(e.MaxIdleConnLifetime) * time.Second)
	}
	return db, nil
}

// Environments is a collection of environments.
type Environments struct {
	// Default is an identifier of the default environment.
	Default string

	// envs is a map of environments.
	// The key is an identifier of the environment.
	envs map[string]*Environment
}

// DefaultEnv returns the default environment.
func (e *Environments) DefaultEnv() (*Environment, error) {
	return e.Use(e.Default)
}

// Use returns the environment specified by the identifier.
func (e *Environments) Use(id string) (*Environment, error) {
	env, exists := e.envs[id]
	if !exists {
		return nil, errors.New("environment not found")
	}
	return env, nil
}

// EnvValueProvider defines a environment value provider.
type EnvValueProvider interface {
	Get(key string) (string, error)
}

// envValueProviderLibraries is a map of environment value providers.
var envValueProviderLibraries = map[string]EnvValueProvider{}

// ValueProvider is a default environment value provider.
type ValueProvider struct{}

// Get returns a value of the environment variable.
func (p ValueProvider) Get(key string) (string, error) {
	return key, nil
}

// OsEnvValueProvider is a environment value provider that uses os.Getenv.
type OsEnvValueProvider struct{}

// Get returns a value of the environment variable.
// It uses os.Getenv.
func (p OsEnvValueProvider) Get(key string) (string, error) {
	var err error
	key = formatRegexp.ReplaceAllStringFunc(key, func(find string) string {
		value := os.Getenv(formatRegexp.FindStringSubmatch(find)[1])
		if len(value) == 0 {
			err = fmt.Errorf("environment variable %s not found", find)
		}
		return value
	})
	return key, err
}

// RegisterEnvValueProvider registers a environment value provider.
// The key is a name of the provider.
// The value is a provider.
// It allows to override the default provider.
func RegisterEnvValueProvider(name string, provider EnvValueProvider) {
	envValueProviderLibraries[name] = provider
}

// defaultEnvValueProvider is a default environment value provider.
var defaultEnvValueProvider = &ValueProvider{}

// GetEnvValueProvider returns a environment value provider.
func GetEnvValueProvider(key string) EnvValueProvider {
	if provider, exists := envValueProviderLibraries[key]; exists {
		return provider
	}
	return defaultEnvValueProvider
}

func init() {
	// Register the default environment value provider.
	RegisterEnvValueProvider("env", &OsEnvValueProvider{})
}
