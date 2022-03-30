package config

type defaults struct {
	AwsRegion  string `yaml:"AwsRegion"`
	DbCluster  string `yaml:"DbCluster"`
	DbName     string `yaml:"DbName"`
	DbPassword string `yaml:"DbPassword"`
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
