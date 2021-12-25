package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/metrics/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/database"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability"
	server "github.com/prixfixeco/api_server/internal/server"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
)

func TestServerConfig_EncodeToFile(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &InstanceConfig{
			Server: server.Config{
				HTTPPort:        1234,
				Debug:           false,
				StartupDeadline: time.Minute,
			},
			Meta: MetaSettings{
				RunMode: DevelopmentRunMode,
			},
			Encoding: encoding.Config{
				ContentType: "application/json",
			},
			Observability: observability.Config{
				Metrics: config.Config{
					Provider: "",
				},
			},
			Services: ServicesConfigurations{
				Auth: authservice.Config{
					Cookies: authservice.CookieConfig{
						Name:     "prixfixecookie",
						Domain:   "https://verygoodsoftwarenotvirus.ru",
						Lifetime: time.Second,
					},
					MinimumUsernameLength: 4,
					MinimumPasswordLength: 8,
					EnableUserSignup:      true,
				},
			},
			Database: dbconfig.Config{
				Debug:             true,
				RunMigrations:     true,
				ConnectionDetails: database.ConnectionDetails("postgres://username:password@host/table"),
			},
		}

		f, err := ioutil.TempFile("", "")
		require.NoError(t, err)

		assert.NoError(t, cfg.EncodeToFile(f.Name(), json.Marshal))
	})

	T.Run("with error marshaling", func(t *testing.T) {
		t.Parallel()

		cfg := &InstanceConfig{}

		f, err := ioutil.TempFile("", "")
		require.NoError(t, err)

		assert.Error(t, cfg.EncodeToFile(f.Name(), func(interface{}) ([]byte, error) {
			return nil, errors.New("blah")
		}))
	})
}
