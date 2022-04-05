package config

type defaults struct {
	AwsRegion  string `yaml:"AwsRegion"`
	CoDomain   string `yaml:"CoDomain"`
	CoEmail    string `yaml:"CoEmail"`
	CoPhone    string `yaml:"CoPhone"`
	DbCluster  string `yaml:"DbCluster"`
	DbName     string `yaml:"DbName"`
	DbPassword string `yaml:"DbPassword"`
	HST        string `yaml:"HST"`
	LogoURI    string `yaml:"LogoURI"`
	DbUser     string `yaml:"DbUser"`
	SsmPath    string `yaml:"SsmPath"`
	Stage      string `yaml:"Stage"`
}

type config struct {
	AwsRegion       string
	DbConnectString string
	DbName          string
	Stage           StageEnvironment
}

type companyInfo struct {
	Domain  string
	Email   string
	HST     string
	LogoURI string
	Phone   string
}
