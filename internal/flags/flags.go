package flags

import (
	"BankingLedgerSystem/constants"
	"flag"
	"os"
	"strconv"
)

type DockerBasedDeployment struct{}

type DeploymentMode interface {
	Env() string
	Port() int
	BaseConfigPath() string
	GetConfigPath() string
}

var (
	env            = flag.String(constants.EnvKey, constants.EnvDefaultValue, constants.EnvUsage)
	port           = flag.Int(constants.PortKey, constants.PortDefaultValue, constants.PortUsage)
	baseConfigPath = flag.String(constants.BaseConfigPathKey, constants.BaseConfigPathDefaultValue, constants.BaseConfigPathUsage)
)

func init() {
	flag.Parse()
}

func (vm DockerBasedDeployment) Env() string {
	return *env
}

func (vm DockerBasedDeployment) BaseConfigPath() string {
	return *baseConfigPath
}

func (vm DockerBasedDeployment) GetConfigPath() string {
	return *baseConfigPath + "/" + *env
}

func (vm DockerBasedDeployment) Port() int {
	port := os.Getenv(constants.PortKey)
	if port == "" {
		return constants.PortDefaultValue
	}
	portInt, _ := strconv.Atoi(port)
	return portInt
}
func (vm DockerBasedDeployment) AppVersion() string {
	appVer := os.Getenv(constants.AppVersion)
	if appVer == "" {
		return constants.AppVersionDefaultValue
	}
	return appVer
}

func ReadDeploymentModeKey() string {
	deploymentModeKey := os.Getenv(constants.DeploymentModeKey)
	if deploymentModeKey != "" {
		return deploymentModeKey
	}
	return constants.DefaultDeploymentModeKey
}

func GetDeploymentMode() DeploymentMode {
	return DockerBasedDeployment{}
}
