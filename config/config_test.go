package config

import (
	"fmt"
	"os"
	"path"
	"testing"
)

var cfg *Config

func TestInitConfig(t *testing.T) {
	t.Run("Successful Init", func(t *testing.T) {
		cfg = &Config{}
		err := cfg.Init()
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}

		// if cfg.Stage != ProdEnv {
		// 	t.Fatalf("Expected Stage value should be: %s, have: %s", ProdEnv, cfg.Stage)
		// }
	})
}

func TestGetStageEnv(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()

	t.Run("successfully get stage value", func(t *testing.T) {
		stg := cfg.GetStageEnv()
		if stg != ProdEnv {
			t.Fatalf("Stage should be %s, have: %s", ProdEnv, stg)
		}
	})

	t.Run("successfully set then get stage value", func(t *testing.T) {
		cfg.SetStageEnv("test")
		stg := cfg.GetStageEnv()
		if stg != TestEnv {
			t.Fatalf("Stage should be %s, have: %s", TestEnv, stg)
		}
	})
}

func TestSetDefaults(t *testing.T) {
	t.Run("test setting DefaultsFilePath", func(t *testing.T) {

		cfg = &Config{}
		cfg.setDefaults()
		dir, _ := os.Getwd()
		expectedFilePath := path.Join(dir, defaultFileName)
		if expectedFilePath != cfg.DefaultsFilePath {
			t.Fatalf("DefaultsFilePath should be %s, have: %s", expectedFilePath, cfg.DefaultsFilePath)
		}

		fp := path.Join("/tmp", defaultFileName)
		cfg = &Config{DefaultsFilePath: fp}
		err := cfg.setDefaults()
		expectedFilePath = fp

		if expectedFilePath != cfg.DefaultsFilePath {
			t.Fatalf("DefaultsFilePath should be %s, have: %s", expectedFilePath, cfg.DefaultsFilePath)
		}
		if err == nil {
			t.Fatalf("setDefaults should return error")
		}
	})
}

// TestValidateStage tests the validateStage method
// validateStage is called at various times including in setEnvVars
func TestValidateStage(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()

	t.Run("stage set from defaults file", func(t *testing.T) {
		if cfg.Stage != ProdEnv {
			t.Fatalf("Stage value should be: %s, have: %s", ProdEnv, cfg.Stage)
		}
	})

	t.Run("stage set from environment", func(t *testing.T) {
		os.Setenv("Stage", "test")
		cfg.setEnvVars() // calls validateStage
		if cfg.Stage != TestEnv {
			t.Fatalf("Stage value should be: %s, have: %s", TestEnv, cfg.Stage)
		}
	})

	t.Run("stage set from invalid environment variable", func(t *testing.T) {
		os.Setenv("Stage", "testit")
		err := cfg.setEnvVars()
		if err == nil {
			t.Fatalf("Expected validateStage to return error")
		}
	})

	t.Run("stage set with SetStageEnv method", func(t *testing.T) {
		err := cfg.SetStageEnv("stage")
		if err != nil {
			t.Fatalf("Expected null error received: %s", err)
		}
	})

	t.Run("invalid stage set with SetStageEnv method", func(t *testing.T) {
		err := cfg.SetStageEnv("stageit")
		if err == nil {
			t.Fatalf("Expected validateStage error")
		}
	})
}

func TestSetSSMParams(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()

	t.Run("DBName is accurate", func(t *testing.T) {
		err := cfg.setSSMParams()
		if err != nil {
			t.Fatalf("Expected null error, received: %s", err)
		}

		if defs.DbName == "" {
			t.Fatalf("Expected defs.DbName to have value")
		}
		if defs.DbPassword == "" {
			t.Fatalf("Expected defs.DbPassword to have value")
		}

	})
}

func TestSetDbConnectString(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()
	cfg.setSSMParams()

	t.Run("DbConnectionString is properly set", func(t *testing.T) {
		expectString := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", defs.DbUser, defs.DbPassword, defs.DbCluster, defs.DbName)
		cfg.setDBConnectString()
		if expectString != cfg.DbConnectString {
			t.Fatalf("DbConnectString should be: %s, have: %s", expectString, cfg.DbConnectString)
		}
	})
}

func TestSetFinal(t *testing.T) {
	cfg = &Config{}
	cfg.setDefaults()
	cfg.setSSMParams()

	t.Run("complete config struct population", func(t *testing.T) {
		cfg.setFinal()
		expectAwsRegion := defs.AwsRegion
		if expectAwsRegion != cfg.config.AwsRegion {
			t.Fatalf("config.AwsRegion should be %s, have: %s", expectAwsRegion, cfg.config.AwsRegion)
		}
	})
}