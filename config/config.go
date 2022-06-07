package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	IsDefaultsLocal bool
	companyInfo     *companyInfo
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

const (
	defaultFileName    = "pdf-defaults.yml"
	defaultsRemotePath = "https://shts-pdf.s3.ca-central-1.amazonaws.com/public/pdf-defaults.yml"
)

var (
	defs             = &defaults{}
	defaultsFilePath string
)

/**
Steps to initializing:
1. Unmarshal yaml defaults into config struct
2. Fetch SSM parameters and overwrite
3. Fetch any environment vars and overwrite above
4.
The order of above is important
*/

// ========================== Public Methods =============================== //

// Init method
func (c *Config) Init() (err error) {

	if err = c.setDefaults(); err != nil {
		return err
	}

	if err = c.setSSMParams(); err != nil {
		return err
	}

	if err = c.setEnvVars(); err != nil {
		return err
	}

	// TODO: complete connection with AWS
	// if c.Stage == ProdEnv {
	// c.setAWSConnectString()
	// } else {
	c.setDBConnectString()
	// }

	c.setCompanyInfo()
	c.setFinal()

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
func (c *Config) GetMongoConnectString() string {
	return c.config.DbConnectString
}

// GetDbName method
func (c *Config) GetDbName() string {
	return c.config.DbName
}

// GetCompanyInfo method
func (c *Config) GetCompanyInfo() *companyInfo {
	return c.companyInfo
}

// ========================== Private Methods =============================== //

func (c *Config) setDefaults() (err error) {

	var file []byte
	if c.IsDefaultsLocal == true { // DefaultsRemote is explicitly set to true

		dir, _ := os.Getwd()
		defaultsFilePath = path.Join(dir, defaultFileName)
		if _, err = os.Stat(defaultsFilePath); os.IsNotExist(err) {
			return err
		}

		file, err = ioutil.ReadFile(defaultsFilePath)
		if err != nil {
			return err
		}

	} else { // using remote file path
		res, err := http.Get(defaultsRemotePath)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		file, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
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
	case "development":
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
// Result should look like: mongodb+srv://<defs.DbUser>:<defs.DbPassword>@<defs.DbCluster>/<defs.DbName>?retryWrites=true&w=majority
func (c *Config) setDBConnectString() {
	c.DbConnectString = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", defs.DbUser, defs.DbPassword, defs.DbCluster, defs.DbName)
}

// Result should look like: mongodb+srv://<AWS access key>:<AWS secret key>@<defs.DbCluster>/<defs.DbName>?authSource=%24external&authMechanism=MONGODB-AWS&retryWrites=true&w=majority&authMechanismProperties=AWS_SESSION_TOKEN:<AWS_SESSION_TOKEN>
func (c *Config) setAWSConnectString() {

	AWS_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID")
	// AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY")
	AWS_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	AWS_SESSION_TOKEN := os.Getenv("AWS_SESSION_TOKEN")
	// fmt.Printf("AWS_ACCESS_KEY: %+v\n", AWS_ACCESS_KEY)
	// fmt.Printf("AWS_SECRET_ACCESS_KEY: %+v\n", AWS_SECRET_ACCESS_KEY)

	// mongodb+srv://<AWS access key>:<AWS secret key>@cluster0.q0zlz.mongodb.net/myFirstDatabase?authSource=%24external&authMechanism=MONGODB-AWS&retryWrites=true&w=majority&authMechanismProperties=AWS_SESSION_TOKEN:<session token (for AWS IAM Roles)>
	c.DbConnectString = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?authSource=%sexternal&authMechanism=MONGODB-AWS&retryWrites=true&w=majority&authMechanismProperties=AWS_SESSION_TOKEN:%s", AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, defs.DbCluster, defs.DbName, "%24", AWS_SESSION_TOKEN)
	// fmt.Printf("c.DbConnectString =: %+v\n", c.DbConnectString)

	// mongodb+srv://ASIAV5T2ZJR5XVVNKQNM:kArqB36wJtTlNqA+kAPgnHS96OKsOcGoB1J4PeAV@cluster0.q0zlz.mongodb.net/shts-test?authSource=%24external&authMechanism=MONGODB-AWS&retryWrites=true&w=majority&authMechanismProperties=AWS_SESSION_TOKEN:
}

func (c *Config) setCompanyInfo() {
	c.companyInfo = &companyInfo{
		Address1: defs.CoAddress1,
		Address2: defs.CoAddress2,
		Domain:   defs.CoDomain,
		Email:    defs.CoEmail,
		HST:      defs.HST,
		LogoURI:  defs.LogoURI,
		Name:     defs.CoName,
		Phone:    defs.CoPhone,
	}
}

// Copies required fields from the defaults to the Config struct
func (c *Config) setFinal() {
	c.AwsRegion = defs.AwsRegion
	c.DbName = defs.DbName
	c.S3Bucket = defs.S3Bucket
}
