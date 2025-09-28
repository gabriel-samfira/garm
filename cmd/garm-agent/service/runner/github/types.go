package github

type RunnerConfig struct {
	AgentID     uint   `json:"agentId"`
	AgentName   string `json:"agentName"`
	PoolID      uint   `json:"poolId"`
	PoolName    string `json:"poolName"`
	Ephemeral   bool   `json:"ephemeral"`
	ServerURL   string `json:"serverUrl"`
	GitHubURL   string `json:"gitHubUrl"`
	WorkFolder  string `json:"workFolder"`
	UseV2Flow   bool   `json:"useV2Flow"`
	ServerURLV2 string `json:"serverUrlV2"`
}
