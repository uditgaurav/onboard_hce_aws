package types

type InfraDetails struct {
	Name             string
	Namespace        string
	InfraDescription string
	PlatformName     string
	ServiceAccount   string
	InfraSaExists    bool
	InstallationType string
	InfraScope       string
	InfraNsExists    bool
	SkipSsl          bool
}

type EnvironmentDetails struct {
	EnvironmentName        string
	EnvironmentDescription string
	EnvironmentType        string
}

type Identifiers struct {
	OrgIdentifier     string `json:"orgIdentifier"`
	AccountIdentifier string `json:"accountIdentifier"`
	ProjectIdentifier string `json:"projectIdentifier"`
}

type Request struct {
	Name             string `json:"name"`
	EnvironmentID    string `json:"environmentID"`
	Description      string `json:"description"`
	PlatformName     string `json:"platformName"`
	InfraNamespace   string `json:"infraNamespace"`
	ServiceAccount   string `json:"serviceAccount"`
	InfraScope       string `json:"infraScope"`
	InfraNsExists    bool   `json:"infraNsExists"`
	InfraSaExists    bool   `json:"infraSaExists"`
	InstallationType string `json:"installationType"`
	SkipSsl          bool   `json:"skipSsl"`
}

type Payload struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables"`
}

type Variables struct {
	Identifiers Identifiers `json:"identifiers"`
	Request     Request     `json:"request"`
}

type OnboardingParameters struct {
	ApiKey                       string
	AccountId                    string
	Organisation                 string
	Project                      string
	Infra                        InfraDetails
	Environment                  EnvironmentDetails
	Timeout                      int
	Delay                        int
	ProviderUrl                  string
	ProviderARN                  string
	RoleName                     string
	RoleARN                      string
	Resources                    string
	Region                       string
	ExperimentServiceAccountName string
	KubeConfigPath               string
	Actions                      string
	AWSCredentialFile            string
	AWSProfile                   string
}

type Response struct {
	Data struct {
		RegisterInfra struct {
			Token    string `json:"token"`
			InfraID  string `json:"infraID"`
			Name     string `json:"name"`
			Manifest string `json:"manifest"`
		} `json:"registerInfra"`
	} `json:"data"`
}

type HarnessEnvironment struct {
	OrgIdentifier     string            `json:"orgIdentifier"`
	ProjectIdentifier string            `json:"projectIdentifier"`
	Identifier        string            `json:"identifier"`
	Tags              map[string]string `json:"tags"`
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	Color             string            `json:"color"`
	Type              string            `json:"type"`
	Yaml              string            `json:"yaml"`
}
