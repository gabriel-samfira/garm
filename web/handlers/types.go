package handlers

import (
	"html/template"
	
	commonParams "github.com/cloudbase/garm-provider-common/params"
)

// CSS class constants for event levels
const (
	EventLevelInfoCSS    = "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200"
	EventLevelWarningCSS = "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
	EventLevelErrorCSS   = "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
	EventLevelDefaultCSS = "bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200"
)

// CSS class constants for instance status
const (
	InstanceStatusRunningCSS  = "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200"
	InstanceStatusStoppedCSS  = "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200"
	InstanceStatusPendingCSS  = "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200"
	InstanceStatusStoppingCSS = "bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200"
	InstanceStatusDefaultCSS  = "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200"
)

// GetInstanceStatusClass returns the appropriate CSS class for an instance status
func GetInstanceStatusClass(status commonParams.InstanceStatus) string {
	switch status {
	case commonParams.InstanceRunning:
		return InstanceStatusRunningCSS
	case commonParams.InstancePendingCreate:
		return InstanceStatusPendingCSS
	case commonParams.InstanceStopped, commonParams.InstanceError, 
		 commonParams.InstanceDeleted, commonParams.InstancePendingDelete,
		 commonParams.InstancePendingForceDelete:
		return InstanceStatusStoppedCSS
	case commonParams.InstanceDeleting:
		return InstanceStatusStoppingCSS
	default:
		return InstanceStatusDefaultCSS
	}
}

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

// OrganizationPoolRowData represents data for an organization pool table row
type OrganizationPoolRowData struct {
	ID           string
	EntityID     string
	Image        string
	ProviderName string
	Status       string
	StatusClass  string
}

// OrganizationInstanceRowData represents data for an organization instance table row
type OrganizationInstanceRowData struct {
	Name         string
	EntityID     string
	CreatedTime  string
	Status       string
	StatusClass  string
	RunnerStatus string
}

// EmptyEventsData represents data for empty events section
type EmptyEventsData struct {
	Message    string
	SubMessage string
}

// EnterpriseRowData represents data for an enterprise table row
type EnterpriseRowData struct {
	ID              string
	Name            string
	ForgeType       string
	ForgeIcon       template.HTML
	EndpointName    string
	CredentialsName string
	Status          string
	StatusClass     string
}

// EnterprisePoolRowData represents data for an enterprise pool table row
type EnterprisePoolRowData struct {
	ID           string
	EntityID     string
	Image        string
	ProviderName string
	Status       string
	StatusClass  string
}

// EnterpriseInstanceRowData represents data for an enterprise instance table row
type EnterpriseInstanceRowData struct {
	Name         string
	EntityID     string
	CreatedTime  string
	Status       string
	StatusClass  string
	RunnerStatus string
}

// RepositoryPoolRowData represents data for a repository pool table row
type RepositoryPoolRowData struct {
	ID           string
	EntityID     string
	Image        string
	ProviderName string
	Status       string
	StatusClass  string
}

// RepositoryInstanceRowData represents data for a repository instance table row
type RepositoryInstanceRowData struct {
	Name         string
	EntityID     string
	CreatedTime  string
	Status       string
	StatusClass  string
	RunnerStatus string
}