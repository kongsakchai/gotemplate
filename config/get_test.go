package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrefix(t *testing.T) {
	type testcase struct {
		title    string
		env      string
		key      string
		expected string
	}

	testcases := []testcase{
		{
			title:    "should return local key when environment is local",
			env:      Local,
			key:      "TESTKEY",
			expected: "LOCAL_TESTKEY",
		},
		{
			title:    "should return dev key when environment is dev",
			env:      Dev,
			key:      "TESTKEY",
			expected: "DEV_TESTKEY",
		},
		{
			title:    "should return key when environment is empty",
			env:      "",
			key:      "TESTKEY",
			expected: "TESTKEY",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			Env = tc.env

			result := prefix(tc.key)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetDuration(t *testing.T) {
	Env = "" // Reset environment variable for testing

	type testcase struct {
		title          string
		envKey         string
		envValue       string
		defaultValue   time.Duration
		expectedResult time.Duration
	}

	testcases := []testcase{
		{
			title:          "should return default value when environment variable is not set",
			envKey:         "TEST_DURATION",
			envValue:       "",
			defaultValue:   5 * time.Second,
			expectedResult: 5 * time.Second,
		},
		{
			title:          "should return parsed duration when environment variable is set",
			envKey:         "TEST_DURATION",
			envValue:       "10s",
			defaultValue:   5 * time.Second,
			expectedResult: 10 * time.Second,
		},
		{
			title:          "should return default value when environment variable is invalid",
			envKey:         "TEST_DURATION",
			envValue:       "invalid",
			defaultValue:   5 * time.Second,
			expectedResult: 5 * time.Second,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			t.Setenv(tc.envKey, tc.envValue)

			result := getDuration(tc.envKey, tc.defaultValue)

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetSecret(t *testing.T) {
	Env = "" // Reset environment variable for testing

	type testcase struct {
		title       string
		secretKey   string
		secretValue string
		searchKey   string
		expected    string
	}

	testcases := []testcase{
		{
			title:       "should return secret value when environment variable is set",
			secretKey:   "TEST_SECRET",
			secretValue: "secretValue",
			searchKey:   "$TEST_SECRET",
			expected:    "secretValue",
		},
		{
			title:       "should return key when environment variable is not set",
			secretKey:   "TEST_SECRET",
			secretValue: "",
			searchKey:   "$TEST_SECRET",
			expected:    "$TEST_SECRET",
		},
		{
			title:       "should return key when not prefixed with $",
			secretKey:   "TEST_SECRET",
			secretValue: "secretValue",
			searchKey:   "TEST_SECRET",
			expected:    "TEST_SECRET",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			t.Setenv(tc.secretKey, tc.secretValue)

			result := getSecret(tc.searchKey)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetString(t *testing.T) {
	Env = "" // Reset environment variable for testing

	type testcase struct {
		title          string
		envKey         string
		envValue       string
		defaultValue   string
		expectedResult string
	}

	testcases := []testcase{
		{
			title:          "should return default value when environment variable is not set",
			envKey:         "TEST_STRING",
			envValue:       "",
			defaultValue:   "default",
			expectedResult: "default",
		},
		{
			title:          "should return environment variable value when set",
			envKey:         "TEST_STRING",
			envValue:       "value",
			defaultValue:   "default",
			expectedResult: "value",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			t.Setenv(tc.envKey, tc.envValue)

			result := getString(tc.envKey, tc.defaultValue)

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetInt(t *testing.T) {
	Env = "" // Reset environment variable for testing

	type testcase struct {
		title          string
		envKey         string
		envValue       string
		defaultValue   int
		expectedResult int
	}

	testcases := []testcase{
		{
			title:          "should return default value when environment variable is not set",
			envKey:         "TEST_INT",
			envValue:       "",
			defaultValue:   42,
			expectedResult: 42,
		},
		{
			title:          "should return parsed integer when environment variable is set",
			envKey:         "TEST_INT",
			envValue:       "100",
			defaultValue:   42,
			expectedResult: 100,
		},
		{
			title:          "should return default value when environment variable is invalid",
			envKey:         "TEST_INT",
			envValue:       "invalid",
			defaultValue:   42,
			expectedResult: 42,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			t.Setenv(tc.envKey, tc.envValue)

			result := getInt(tc.envKey, tc.defaultValue)

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetConfig(t *testing.T) {
	Env = "" // Reset environment variable for testing
	t.Run("should return correct App configuration", func(t *testing.T) {
		t.Setenv("APP_NAME", "TestApp")
		t.Setenv("APP_PORT", "8080")
		t.Setenv("APP_VERSION", "1.0.0")

		expectedApp := App{
			Name:    "TestApp",
			Port:    "8080",
			Version: "1.0.0",
		}

		appConfig := getAppConfig()

		assert.Equal(t, expectedApp, appConfig)
	})

	t.Run("should return correct Header configuration", func(t *testing.T) {
		t.Setenv("HEADER_REF_ID_KEY", "X-Ref-ID")

		expectedHeader := Header{
			RefIDKey: "X-Ref-ID",
		}

		headerConfig := getHeaderConfig()

		assert.Equal(t, expectedHeader, headerConfig)
	})

	t.Run("should return correct Migration configuration", func(t *testing.T) {
		t.Setenv("MIGRATION_DIRECTORY", "migrations")
		t.Setenv("MIGRATION_VERSION", "1001")

		expectedMigration := Migration{
			Directory: "migrations",
			Version:   "1001",
		}

		migrationConfig := getMigrationConfig()

		assert.Equal(t, expectedMigration, migrationConfig)
	})

	t.Run("should return correct Database configuration", func(t *testing.T) {
		t.Setenv("DATABASE_URL", "user:password@tcp(localhost:5432)/testdb")

		expectedDatabase := Database{
			URL: "user:password@tcp(localhost:5432)/testdb",
		}

		dbConfig := getDatabaseConfig()

		assert.Equal(t, expectedDatabase, dbConfig)
	})

	t.Run("should return correct Redis configuration", func(t *testing.T) {
		t.Setenv("REDIS_HOST", "localhost")
		t.Setenv("REDIS_PORT", "6379")
		t.Setenv("REDIS_USERNAME", "redisuser")
		t.Setenv("REDIS_PASSWORD", "redispassword")
		t.Setenv("REDIS_DB", "1")
		t.Setenv("REDIS_TIMEOUT", "5s")

		expectedRedis := Redis{
			Host:     "localhost",
			Port:     "6379",
			Username: "redisuser",
			Password: "redispassword",
			DB:       1,
			Timeout:  5 * time.Second,
		}

		redisConfig := getRedisConfig()

		assert.Equal(t, expectedRedis, redisConfig)
	})
}
