package config

type defaults struct {
	AwsRegion  string `yaml:"AwsRegion"`
	CoAddress1 string `yaml:"CoAddress1"`
	CoAddress2 string `yaml:"CoAddress2"`
	CoDomain   string `yaml:"CoDomain"`
	CoEmail    string `yaml:"CoEmail"`
	CoName     string `yaml:"CoName"`
	CoPhone    string `yaml:"CoPhone"`
	DbCluster  string `yaml:"DbCluster"`
	DbName     string `yaml:"DbName"`
	HST        string `yaml:"HST"`
	LogoURI    string `yaml:"LogoURI"`
	S3Bucket   string `yaml:"S3Bucket"`
	SsmPath    string `yaml:"SsmPath"`
	Stage      string `yaml:"Stage"`
}

type config struct {
	AwsRegion       string
	DbConnectString string
	DbName          string
	S3Bucket        string
	Stage           StageEnvironment
}

type companyInfo struct {
	Address1 string
	Address2 string
	Domain   string
	Email    string
	HST      string
	LogoURI  string
	Name     string
	Phone    string
}
