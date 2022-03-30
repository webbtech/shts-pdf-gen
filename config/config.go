package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	config
	DefaultsFilePath string
}

// StageEnvironment string
type StageEnvironment string

// DB type constants
const (
	DevEnv   StageEnvironment = "dev"
	StageEnv StageEnvironment = "stage"
	TestEnv  StageEnvironment = "test"
	ProdEnv  StageEnvironment = "prod"
)

const defaultFileName = "defaults.yml"

var defs = &defaults{}

/**
Steps to initializing:
1. Unmarshal yaml defaults into config struct
2. Fetch KMS parameters and overwrite
3. Fetch any environment vars and overwrite above
The order of above is important
*/

// ========================== Public Methods =============================== //

// Init method
func (c *Config) Init() (err error) {

	if err = c.setDefaults(); err != nil {
		return err
	}

	if err = c.setEnvVars(); err != nil {
		return err
	}

	if err = c.setSSMParams(); err != nil {
		return err
	}

	return err
}

// GetStageEnv method
func (c *Config) GetStageEnv() StageEnvironment {
	return c.Stage
}

// SetStageEnv method
func (c *Config) SetStageEnv(env string) (err error) {
	defs.Stage = env
	return c.validateStage()
}

// GetMongoConnectURL method
/* func (c *Config) GetMongoConnectURL() string {
	return c.DBConnectURL
} */

// ========================== Private Methods =============================== //

func (c *Config) setDefaults() (err error) {

	if c.DefaultsFilePath == "" {
		dir, _ := os.Getwd()
		c.DefaultsFilePath = path.Join(dir, defaultFileName)
	}

	file, err := ioutil.ReadFile(c.DefaultsFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(file), &defs)
	if err != nil {
		return err
	}

	err = c.validateStage()
	if err != nil {
		return err
	}

	return err
}

// validateStage method to validate Stage value
func (c *Config) validateStage() (err error) {

	validEnv := true

	switch defs.Stage {
	case "dev":
		c.Stage = DevEnv
	case "stage":
		c.Stage = StageEnv
	case "test":
		c.Stage = TestEnv
	case "prod":
		c.Stage = ProdEnv
	case "production":
		c.Stage = ProdEnv
	default:
		validEnv = false
	}

	if !validEnv {
		return errors.New("Invalid StageEnvironment requested")
	}

	return err
}

// setEnvVars sets any environment variables that match the default struct fields
func (c *Config) setEnvVars() (err error) {

	vals := reflect.Indirect(reflect.ValueOf(defs))
	for i := 0; i < vals.NumField(); i++ {
		nm := vals.Type().Field(i).Name
		if e := os.Getenv(nm); e != "" {
			vals.Field(i).SetString(e)
		}
		// If field is Stage, validate and return error if required
		if nm == "Stage" {
			err = c.validateStage()
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (c *Config) setSSMParams() (err error) {

	s := []string{"", string(c.GetStageEnv()), defs.SsmPath}
	paramPath := aws.String(strings.Join(s, "/"))

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(defs.AwsRegion),
	})
	if err != nil {
		return err
	}

	svc := ssm.New(sess)
	res, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           paramPath,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	paramLen := len(res.Parameters)
	if paramLen == 0 {
		err = fmt.Errorf("Error fetching ssm params, total number found: %d", paramLen)
	}

	// Get struct keys so we can test before attempting to set
	t := reflect.ValueOf(defs).Elem()
	for _, r := range res.Parameters {
		paramName := strings.Split(*r.Name, "/")[3]
		structKey := t.FieldByName(paramName)
		if structKey.IsValid() {
			structKey.Set(reflect.ValueOf(*r.Value))
		}
	}

	return err
}

// setDBConnectString Build a connection string to MongoDB Atlas
// should look like: mongodb+srv://<defs.DbUser>:<defs.DbPassword>@<defs.DbCluster>/<defs.DbName>?retryWrites=true&w=majority
func (c *Config) setDBConnectString() {
	c.DbConnectString = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", defs.DbUser, defs.DbPassword, defs.DbCluster, defs.DbName)
}

// Copies required fields from the defaults to the Config struct
func (c *Config) setFinal() (err error) {
	c.AwsRegion = defs.AwsRegion
	return err
}
