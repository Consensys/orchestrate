package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/tls"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const (
	component = "chain-registry.store.types"
	separator = "-"
)

type ConfigType int32

const (
	UNKNOWN        ConfigType = 0
	HTTPROUTER     ConfigType = 1
	HTTPMIDDLEWARE ConfigType = 2
	HTTPSERVICE    ConfigType = 3
	TCPROUTER      ConfigType = 4
	TCPSERVICE     ConfigType = 5
	TLSCERTIFICATE ConfigType = 6
	TLSOPTIONS     ConfigType = 7
	TLSSTORES      ConfigType = 8
)

type Config struct {
	tableName struct{} `sql:"config"` //nolint:unused,structcheck

	ID         int    `sql:",pk" json:"-"`
	Name       string `sql:",notnull"`
	TenantID   string
	ConfigType ConfigType
	Config     json.RawMessage
}

// UnmarshalJSONConfig unmarshal a config in JSON byte/string and checks that the config match the struct
func UnmarshalJSONConfig(config []byte, configStruct interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(config))
	dec.DisallowUnknownFields() // Force errors
	if err := dec.Decode(configStruct); err != nil {
		return errors.FromError(fmt.Errorf("bad config in json format - %v", err)).ExtendComponent(component)
	}
	return nil
}

// NewTraefikConfig returns a new traefik dynamic config with its fields initiated
func NewTraefikConfig() *dynamic.Configuration {
	return &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:     make(map[string]*dynamic.Router),
			Middlewares: make(map[string]*dynamic.Middleware),
			Services:    make(map[string]*dynamic.Service),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:  make(map[string]*dynamic.TCPRouter),
			Services: make(map[string]*dynamic.TCPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]tls.Store),
			Options: make(map[string]tls.Options),
		},
	}
}

// BuildConfiguration converts the config model to a traefik dynamic configuration
func BuildConfiguration(configs []*Config) (*dynamic.Configuration, error) {
	traefikConfig := NewTraefikConfig()

	for _, config := range configs {
		configStruct, err := GetConfigStruct(config.ConfigType)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		err = UnmarshalJSONConfig(config.Config, configStruct)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		// Prefix config names by tenant id
		name := strings.Join([]string{config.TenantID, config.Name}, separator)

		// TODO: check there are no duplicates in config names
		switch v := configStruct.(type) {
		case *dynamic.Router:
			traefikConfig.HTTP.Routers[name] = v
		case *dynamic.Middleware:
			traefikConfig.HTTP.Middlewares[name] = v
		case *dynamic.Service:
			traefikConfig.HTTP.Services[name] = v
		case *dynamic.TCPRouter:
			traefikConfig.TCP.Routers[name] = v
		case *dynamic.TCPService:
			traefikConfig.TCP.Services[name] = v
		case *tls.CertAndStores:
			return nil, errors.FromError(fmt.Errorf("%v config type is not allowed", v)).ExtendComponent(component)
		case *tls.Options:
			traefikConfig.TLS.Options[name] = *v
		case *tls.Store:
			traefikConfig.TLS.Stores[name] = *v
		default:
			return nil, errors.FromError(fmt.Errorf("unknown %d config type ", v)).ExtendComponent(component)
		}
	}

	return traefikConfig, nil
}

// TraefikConfigToStoreConfig converts a traefik dynamic config to a config model
func TraefikConfigToStoreConfig(traefikConfig *dynamic.Configuration, tenantID string) ([]*Config, error) {
	configs := make([]*Config, 0)

	v := reflect.ValueOf(*traefikConfig)
	// Loop over Configuration struct (HTTP, TCP or TLS)
	for i := 0; i < v.NumField(); i++ {
		element := reflect.Indirect(v.Field(i))
		if element.Kind() == reflect.Invalid {
			continue
		}

		// Loop over Configuration struct (HTTP_Routers, HTTP_Middlewares or HTTP_Services, TCP_Routers...)
		for j := 0; j < element.NumField(); j++ {
			configsField := element.Field(j)
			switch configsField.Kind() {
			case reflect.Map:
				// Loop over the mapping
				for _, e := range configsField.MapKeys() {
					config := configsField.MapIndex(e)
					configType, _ := GetConfigType(config.Interface())
					b, _ := json.Marshal(config.Interface())
					configs = append(configs, &Config{
						Name:       e.String(),
						TenantID:   tenantID,
						ConfigType: configType,
						Config:     b,
					})
				}
			default:
				return nil, errors.FromError(fmt.Errorf("unknown %d config type ", configsField.Kind())).ExtendComponent(component)
			}
		}
	}

	return configs, nil
}

// GetConfigStruct maps a config type to a traefik dynamic config
func GetConfigStruct(configType ConfigType) (interface{}, error) {
	switch configType {
	case HTTPROUTER:
		return &dynamic.Router{}, nil
	case HTTPMIDDLEWARE:
		return &dynamic.Middleware{}, nil
	case HTTPSERVICE:
		return &dynamic.Service{}, nil
	case TCPROUTER:
		return &dynamic.TCPRouter{}, nil
	case TCPSERVICE:
		return &dynamic.TCPService{}, nil
	case TLSCERTIFICATE:
		return &tls.CertAndStores{}, errors.FromError(fmt.Errorf("%v config type are not allowed", configType)).ExtendComponent(component)
	case TLSOPTIONS:
		return &tls.Options{}, nil
	case TLSSTORES:
		return &tls.Store{}, nil
	default:
		return nil, errors.FromError(fmt.Errorf("unknown %d config type ", configType)).ExtendComponent(component)
	}
}

// GetConfigType maps a traefik dynamic config to a config type
func GetConfigType(configStruct interface{}) (ConfigType, error) {
	switch configStruct.(type) {
	case *dynamic.Router:
		return HTTPROUTER, nil
	case *dynamic.Middleware:
		return HTTPMIDDLEWARE, nil
	case *dynamic.Service:
		return HTTPSERVICE, nil
	case *dynamic.TCPRouter:
		return TCPROUTER, nil
	case *dynamic.TCPService:
		return TCPSERVICE, nil
	case *tls.CertAndStores:
		return TLSCERTIFICATE, errors.FromError(fmt.Errorf("%v config type are not allowed", configStruct)).ExtendComponent(component)
	case *tls.Options:
		return TLSOPTIONS, nil
	case *tls.Store:
		return TLSSTORES, nil
	default:
		return UNKNOWN, errors.FromError(fmt.Errorf("unknown %v config type ", configStruct)).ExtendComponent(component)
	}
}
