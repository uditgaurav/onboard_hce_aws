package register

type InfraDetails struct {
	Name      string
	Namespace string
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

type InfraParameters struct {
	ApiKey        string
	AccountId     string
	Organisation  string
	Project       string
	Infra         InfraDetails
	InfraScope    string
	InfraNsExists bool
}
