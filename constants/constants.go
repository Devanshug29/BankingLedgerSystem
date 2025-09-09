package constants

const (
	EnvKey                     = "env"
	EnvDefaultValue            = "dev"
	EnvUsage                   = "environment of the application, can be dev, uat, test, pre-prod, prod"
	PortKey                    = "port"
	PortDefaultValue           = 8080
	PortUsage                  = "application.yml port"
	BaseConfigPathKey          = "base-config-path"
	BaseConfigPathDefaultValue = "resources"
	BaseConfigPathUsage        = "path to folder that stores your configurations"
	AppVersion                 = "app-version"
	AppVersionDefaultValue     = "0.0.0"
	DeploymentModeKey          = "DEPLOYMENT_MODE"
	DefaultDeploymentModeKey   = "vm"
)
