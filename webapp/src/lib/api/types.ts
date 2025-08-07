// API Types generated from GARM Swagger spec

export interface Repository {
	id: string;
	name: string;
	owner: string;
	credentials_name: string;
	endpoint: Endpoint;
	pool_manager_status: PoolManagerStatus;
	pool_balancer_type: string;
	webhook_secret?: string;
	events?: EntityEvent[];
}

export interface Organization {
	id: string;
	name: string;
	credentials_name: string;
	endpoint: Endpoint;
	pool_manager_status: PoolManagerStatus;
	pool_balancer_type: string;
	webhook_secret?: string;
	events?: EntityEvent[];
}

export interface Enterprise {
	id: string;
	name: string;
	credentials_name: string;
	endpoint: Endpoint;
	pool_manager_status: PoolManagerStatus;
	pool_balancer_type: string;
	webhook_secret?: string;
	events?: EntityEvent[];
}

export interface Endpoint {
	name: string;
	description: string;
	api_base_url: string;
	base_url: string;
	endpoint_type: EndpointType;
	upload_base_url?: string;
	ca_cert_bundle?: string;
}

export interface Pool {
	id: string;
	provider_name: string;
	max_runners: number;
	min_idle_runners: number;
	image: string;
	flavor: string;
	os_type: string;
	os_arch: string;
	tags: string[];
	enabled: boolean;
	runner_bootstrap_timeout: number;
	extra_specs?: string;
	github_runner_group?: string;
	repo_id?: string;
	org_id?: string;
	enterprise_id?: string;
	priority?: number;
	runner_prefix: string;
	instances?: Instance[];
	endpoint?: Endpoint;
}

export interface ScaleSet {
	id: number;
	scale_set_id: number;
	name: string;
	provider_name: string;
	max_runners: number;
	min_idle_runners: number;
	image: string;
	flavor: string;
	os_type: string;
	os_arch: string;
	enabled: boolean;
	runner_bootstrap_timeout: number;
	extra_specs?: string;
	github_runner_group?: string;
	repo_id?: string;
	org_id?: string;
	enterprise_id?: string;
	runner_prefix: string;
	instances?: Instance[];
	endpoint?: Endpoint;
}

export interface Instance {
	id: string;
	name: string;
	os_type: string;
	os_arch: string;
	status: InstanceStatus;
	runner_status: RunnerStatus;
	pool_id?: string;
	scale_set_id?: number;
	provider_id: string;
	created_at: string;
	updated_at: string;
	addresses?: InstanceAddress[];
	status_messages?: InstanceStatusMessage[];
}

export interface InstanceAddress {
	address: string;
	type: AddressType;
}

export interface InstanceStatusMessage {
	created_at: string;
	message: string;
}

export interface EntityEvent {
	id: string;
	created_at: string;
	message: string;
	event_level: EventLevel;
	event_type: string;
}

export interface PoolManagerStatus {
	running: boolean;
	failure_reason?: string;
}

export interface ForgeCredentials {
	id: number;
	name: string;
	description: string;
	api_base_url: string;
	upload_base_url?: string;
	base_url: string;
	ca_bundle?: Uint8Array;
	'auth-type': AuthType;
	forge_type: 'github' | 'gitea';
	repositories?: Repository[];
	organizations?: Organization[];
	enterprises?: Enterprise[];
	endpoint: Endpoint;
	created_at: string;
	updated_at: string;
	rate_limit?: any; // GithubRateLimit type - not needed for UI
}

// Alias for backward compatibility
export type Credential = ForgeCredentials;

export interface Provider {
	name: string;
	description: string;
	provider_type: string;
}

export interface ControllerInfo {
	controller_id: string;
	hostname: string;
	metadata_url: string;
	callback_url: string;
	webhook_url: string;
	websocket_url: string;
	controller_webhook_url: string;
	version: string;
}

// Enums
export enum EndpointType {
	GITHUB = "github",
	GITEA = "gitea"
}

export enum InstanceStatus {
	PENDING_CREATE = "pending_create",
	CREATING = "creating", 
	RUNNING = "running",
	PENDING_DELETE = "pending_delete",
	DELETING = "deleting",
	ERROR = "error",
	STOPPED = "stopped"
}

export enum RunnerStatus {
	PENDING = "pending",
	IDLE = "idle", 
	ACTIVE = "active",
	ERROR = "error"
}

export enum EventLevel {
	INFO = "info",
	WARNING = "warning", 
	ERROR = "error"
}

export enum AuthType {
	PAT = "pat",
	APP = "app"
}

export enum AddressType {
	PUBLIC = "public",
	PRIVATE = "private"
}

export enum PoolBalancerType {
	ROUNDROBIN = "roundrobin",
	PACK = "pack",
	NONE = ""
}

// Request/Response types
export interface CreateRepoParams {
	name: string;
	owner: string;
	credentials_name: string;
	webhook_secret?: string;
	pool_balancer_type?: string;
}

export interface CreateOrgParams {
	name: string;
	credentials_name: string;
	webhook_secret?: string;
	pool_balancer_type?: string;
}

export interface CreateEnterpriseParams {
	name: string;
	credentials_name: string;
	webhook_secret?: string;
	pool_balancer_type?: string;
}

export interface CreateEndpointParams {
	name: string;
	description: string;
	endpoint_type: string;
	base_url: string;
	api_base_url?: string;
	upload_base_url?: string;
	ca_cert_bundle?: string;
}

export interface UpdateEndpointParams {
	name?: string;
	description?: string;
	base_url?: string;
	api_base_url?: string;
	upload_base_url?: string;
	ca_cert_bundle?: string;
}

export interface CreateCredentialsParams {
	name: string;
	description: string;
	endpoint: string;
	auth_type: AuthType;
	pat_token?: string;
	app_id?: string;
	app_installation_id?: string;
	private_key_bytes?: string;
}

export interface UpdateCredentialsParams {
	name?: string;
	description?: string;
	pat_token?: string;
	app_id?: string;
	app_installation_id?: string;
	private_key_bytes?: string;
}

export interface CreatePoolParams {
	provider_name: string;
	max_runners: number;
	min_idle_runners: number;
	image: string;
	flavor: string;
	os_type: string;
	os_arch: string;
	tags: string[];
	enabled: boolean;
	runner_bootstrap_timeout: number;
	extra_specs?: string;
	github_runner_group?: string;
	runner_prefix: string;
	priority?: number;
}

export interface CreateScaleSetParams {
	name: string;
	provider_name: string;
	max_runners: number;
	min_idle_runners: number;
	image: string;
	flavor: string;
	os_type: string;
	os_arch: string;
	enabled: boolean;
	runner_bootstrap_timeout: number;
	extra_specs?: string;
	github_runner_group?: string;
	runner_prefix: string;
}

export interface UpdateEntityParams {
	credentials_name?: string;
	webhook_secret?: string;
	pool_balancer_type?: string;
}

export interface UpdatePoolParams {
	max_runners?: number;
	min_idle_runners?: number;
	enabled?: boolean;
	runner_bootstrap_timeout?: number;
	extra_specs?: string;
	github_runner_group?: string;
	tags?: string[];
	priority?: number;
}

export interface LoginRequest {
	username: string;
	password: string;
}

export interface LoginResponse {
	token: string;
	expires_in: number;
}

export interface APIError {
	error: string;
	details?: string;
}