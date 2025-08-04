package handlers

import "html/template"

// EmptyTableRowData represents data for an empty table row with icon
type EmptyTableRowData struct {
	ColSpan    int
	Icon       template.HTML
	Message    string
	SubMessage string
}

// SimpleEmptyTableRowData represents data for a simple empty table row
type SimpleEmptyTableRowData struct {
	ColSpan    int
	ItemType   string
	SearchTerm string
}

// EmptyContentData represents data for empty content sections
type EmptyContentData struct {
	Message string
}

// CredentialRowData represents data for a credential table row
type CredentialRowData struct {
	ID            uint
	Name          string
	Description   string
	ForgeIcon     template.HTML
	EndpointName  string
	AuthTypeBadge string
	AuthTypeText  string
}

// EndpointRowData represents data for an endpoint table row
type EndpointRowData struct {
	Name        string
	Description string
	APIURL      string
	ForgeIcon   template.HTML
	ForgeType   string
	BaseURL     string
}

// PoolRowData represents data for a pool table row
type PoolRowData struct {
	ID           string
	Image        string
	EntityType   string
	EntityName   string
	EndpointName string
	ForgeIcon    template.HTML
	Status       string
	StatusClass  string
}

// RepositoryRowData represents data for a repository table row
type RepositoryRowData struct {
	ID              string
	DisplayName     string
	ForgeType       string
	ForgeIcon       template.HTML
	EndpointName    string
	CredentialsName string
	Status          string
	StatusClass     string
}

// EventRowData represents data for an event row
type EventRowData struct {
	Level      string
	LevelClass string
	Message    string
	Timestamp  string
}

// ScaleSetRowData represents data for a scale set table row
type ScaleSetRowData struct {
	ID             uint
	Name           string
	Image          string
	ImageTitle     string
	EntityType     string
	EntityName     string
	ProviderName   string
	InstanceCount  int
	Status         string
	StatusClass    string
}

// InstanceRowData represents data for an instance table row
type InstanceRowData struct {
	Name              string
	ID                string
	PoolID            string
	CreatedTime       string
	Status            string
	StatusClass       string
	RunnerStatus      string
	RunnerStatusClass string
}

// OrganizationRowData represents data for an organization table row
type OrganizationRowData struct {
	ID              string
	Name            string
	ForgeType       string
	ForgeIcon       template.HTML
	EndpointName    string
	CredentialsName string
	Status          string
	StatusClass     string
}

// EmptyContentData represents data for empty content sections
type EmptyContentSectionData struct {
	Message string
}