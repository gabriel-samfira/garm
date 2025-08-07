<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { garmApi } from '$lib/api/client.js';
	import type { ForgeCredentials, Endpoint } from '$lib/api/types.js';
	import { AuthType } from '$lib/api/types.js';
	import ForgeTypeSelector from '$lib/components/ForgeTypeSelector.svelte';
	import { websocketStore, type WebSocketEvent } from '$lib/stores/websocket.js';

	let loading = true;
	let credentials: ForgeCredentials[] = [];
	let endpoints: Endpoint[] = []; // Only used for modal dropdowns
	let error = '';
	let showCreateModal = false;
	let showEditModal = false;
	let showDeleteModal = false;
	let selectedAuthType = AuthType.PAT;
	let editingCredential: ForgeCredentials | null = null;
	let deletingCredential: ForgeCredentials | null = null;
	let unsubscribeGithubWebsocket: (() => void) | null = null;
	let unsubscribeGiteaWebsocket: (() => void) | null = null;


	// Form state
	let formData = {
		name: '',
		description: '',
		endpoint: '',
		auth_type: AuthType.PAT,
		pat_token: '',
		app_id: '',
		app_installation_id: '',
		private_key_bytes: ''
	};

	function handleCredentialEvent(event: WebSocketEvent) {
		console.log('[Credentials] Received websocket event:', event);
		
		if (event.operation === 'create') {
			const newCredential = event.payload as ForgeCredentials;
			credentials = [...credentials, newCredential];
		} else if (event.operation === 'update') {
			const updatedCredential = event.payload as ForgeCredentials;
			credentials = credentials.map(credential => 
				credential.name === updatedCredential.name ? updatedCredential : credential
			);
		} else if (event.operation === 'delete') {
			const credentialName = event.payload.name || event.payload;
			credentials = credentials.filter(credential => credential.name !== credentialName);
		}
	}

	onMount(async () => {
		// Only load credentials - they contain endpoint information
		await loadCredentials();
		
		// Subscribe to both GitHub and Gitea credential events
		unsubscribeGithubWebsocket = websocketStore.subscribeToEntity(
			'github_credentials',
			['create', 'update', 'delete'],
			handleCredentialEvent
		);
		
		unsubscribeGiteaWebsocket = websocketStore.subscribeToEntity(
			'gitea_credentials',
			['create', 'update', 'delete'],
			handleCredentialEvent
		);
	});

	onDestroy(() => {
		if (unsubscribeGithubWebsocket) {
			unsubscribeGithubWebsocket();
			unsubscribeGithubWebsocket = null;
		}
		if (unsubscribeGiteaWebsocket) {
			unsubscribeGiteaWebsocket();
			unsubscribeGiteaWebsocket = null;
		}
	});

	async function loadCredentials() {
		try {
			loading = true;
			error = '';
			credentials = await garmApi.listAllCredentials();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load credentials';
			console.error('Credentials error:', err);
		} finally {
			loading = false;
		}
	}

	async function loadEndpoints() {
		try {
			endpoints = await garmApi.listAllEndpoints();
		} catch (err) {
			console.error('Failed to load endpoints:', err);
			endpoints = [];
		}
	}

	async function showCreateCredentialsModal() {
		resetForm();
		await loadEndpoints();
		showCreateModal = true;
		selectedForgeType = 'github'; // Default to github
	}

	// Add forge type selection state
	let selectedForgeType: 'github' | 'gitea' | '' = '';

	function handleForgeTypeSelect(event: CustomEvent<'github' | 'gitea'>) {
		selectedForgeType = event.detail;
		// Reset form when forge type changes
		formData.auth_type = AuthType.PAT;
		selectedAuthType = AuthType.PAT;
	}

	async function showEditCredentialsModal(credential: ForgeCredentials) {
		editingCredential = credential;
		formData = {
			name: credential.name || '',
			description: credential.description || '',
			endpoint: credential.endpoint.name || '',
			auth_type: credential['auth-type'] || AuthType.PAT,
			pat_token: '',
			app_id: '',
			app_installation_id: '',
			private_key_bytes: ''
		};
		selectedAuthType = credential['auth-type'] || AuthType.PAT;
		await loadEndpoints();
		showEditModal = true;
	}

	function showDeleteCredentialsModal(credential: ForgeCredentials) {
		deletingCredential = credential;
		showDeleteModal = true;
	}

	function resetForm() {
		formData = {
			name: '',
			description: '',
			endpoint: '',
			auth_type: AuthType.PAT,
			pat_token: '',
			app_id: '',
			app_installation_id: '',
			private_key_bytes: ''
		};
		selectedAuthType = AuthType.PAT;
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		editingCredential = null;
		deletingCredential = null;
		selectedForgeType = '';
		resetForm();
	}

	function handleAuthTypeChange(authType: AuthType) {
		selectedAuthType = authType;
		formData.auth_type = authType;
	}

	async function handleCreateCredentials() {
		try {
			// Use selected forge type to determine which API to call
			if (selectedForgeType === 'github') {
				await garmApi.createGithubCredentials(formData);
			} else if (selectedForgeType === 'gitea') {
				await garmApi.createGiteaCredentials(formData);
			} else {
				throw new Error('Please select a forge type');
			}
			// No need to reload - websocket will handle the update
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create credentials';
		}
	}

	async function handleUpdateCredentials() {
		if (!editingCredential) return;
		
		try {
			const endpointType = editingCredential.forge_type;
			if (endpointType === 'github') {
				await garmApi.updateGithubCredentials(editingCredential.name, formData);
			} else {
				await garmApi.updateGiteaCredentials(editingCredential.name, formData);
			}
			// No need to reload - websocket will handle the update
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update credentials';
		}
	}

	async function handleDeleteCredentials() {
		if (!deletingCredential) return;
		
		try {
			const endpointType = deletingCredential.forge_type;
			if (endpointType === 'github') {
				await garmApi.deleteGithubCredentials(deletingCredential.name);
			} else {
				await garmApi.deleteGiteaCredentials(deletingCredential.name);
			}
			// No need to reload - websocket will handle the update
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete credentials';
		}
	}

	function handlePrivateKeyUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		
		if (!file) {
			formData.private_key_bytes = '';
			return;
		}

		const reader = new FileReader();
		reader.onload = (e) => {
			const content = e.target?.result as string;
			formData.private_key_bytes = btoa(content);
		};
		reader.readAsText(file);
	}

	function isFormValid() {
		if (!formData.name || !formData.description || !formData.endpoint) return false;
		
		if (formData.auth_type === AuthType.PAT) {
			return !!formData.pat_token;
		} else {
			return !!formData.app_id && !!formData.app_installation_id && !!formData.private_key_bytes;
		}
	}

	function getEndpointForgeType(endpointName: string): string {
		const endpoint = endpoints.find(e => e.name === endpointName);
		return endpoint?.endpoint_type || '';
	}

	function isGiteaEndpoint(endpointName: string): boolean {
		return getEndpointForgeType(endpointName) === 'gitea';
	}

	// Get filtered endpoints based on selected forge type
	$: filteredEndpoints = selectedForgeType ? endpoints.filter(e => e.endpoint_type === selectedForgeType) : endpoints;

	function getForgeIcon(credential: ForgeCredentials): string;
	function getForgeIcon(forgeType: 'github' | 'gitea'): string;
	function getForgeIcon(input: ForgeCredentials | 'github' | 'gitea'): string {
		const forgeType = typeof input === 'string' ? input : input.forge_type;
		
		if (forgeType === 'gitea') {
			return `<svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else if (forgeType === 'github') {
			return `<svg class="w-5 h-5 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
				<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/>
			</svg>
			<svg class="w-5 h-5 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
				<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.80 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#ffffff"/>
			</svg>`;
		} else {
			// Return a generic placeholder icon if endpoint type is unknown
			return `<svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
			</svg>`;
		}
	}

	function getForgeIconModal(forgeType: 'github' | 'gitea') {
		if (forgeType === 'gitea') {
			return `<svg class="w-8 h-8" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else {
			return `<div class="inline-flex w-8 h-8"><svg class="w-8 h-8 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-8 h-8 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.80 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
		}
	}
</script>

<svelte:head>
	<title>Credentials - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="sm:flex sm:items-center">
		<div class="sm:flex-auto">
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Credentials</h1>
			<p class="mt-2 text-sm text-gray-700 dark:text-gray-300">
				Manage authentication credentials for your GitHub and Gitea endpoints.
			</p>
		</div>
		<div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none flex items-center space-x-4">
			<button
				type="button"
				on:click={showCreateCredentialsModal}
				class="block rounded-md bg-blue-600 px-3 py-2 text-center text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
			>
				Add Credentials
			</button>
		</div>
	</div>

	{#if error}
		<!-- Error state -->
		<div class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
					</svg>
				</div>
				<div class="ml-3">
					<h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error</h3>
					<p class="mt-2 text-sm text-red-700 dark:text-red-300">{error}</p>
				</div>
			</div>
		</div>
	{/if}

	<!-- Table -->
	<div class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
		{#if loading}
			<div class="px-4 py-8 text-center">
				<div class="animate-spin mx-auto h-8 w-8 border-b-2 border-blue-600 rounded-full"></div>
				<p class="mt-4 text-sm text-gray-500 dark:text-gray-400">Loading credentials...</p>
			</div>
		{:else if credentials.length === 0}
			<div class="px-4 py-8 text-center">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No credentials</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Get started by creating new credentials.</p>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Description</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Endpoint</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Auth Type</th>
							<th scope="col" class="relative px-6 py-3"><span class="sr-only">Actions</span></th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each credentials as credential}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
									{credential.name}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
									{credential.description}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
									<div class="flex items-center">
										{@html getForgeIcon(credential)}
										<span class="ml-2">{credential.endpoint.name}</span>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
									<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full
										{credential['auth-type'] === AuthType.PAT ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'}">
										{credential['auth-type'] === AuthType.PAT ? 'PAT' : 'App'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex justify-end space-x-2">
										<button
											on:click={() => showEditCredentialsModal(credential)}
											class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-900 dark:hover:text-indigo-300"
											title="Edit credentials"
										>
											<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
											</svg>
										</button>
										<button
											on:click={() => showDeleteCredentialsModal(credential)}
											class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300"
											title="Delete credentials"
										>
											<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
											</svg>
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>

<!-- Create Modal -->
{#if showCreateModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 dark:bg-gray-900 dark:bg-opacity-75 flex items-center justify-center z-50" on:click={closeModals}>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-screen overflow-y-auto" on:click|stopPropagation>
			<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
				<div>
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white">
						Add Credentials
					</h3>
					<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
						Create new authentication credentials
					</p>
				</div>
				<button on:click={closeModals} class="text-gray-400 hover:text-gray-600 dark:text-gray-300 dark:hover:text-gray-100">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			
				<form on:submit|preventDefault={handleCreateCredentials} class="p-6 space-y-4">
					<!-- Forge Type Selection -->
					<ForgeTypeSelector 
						bind:selectedForgeType 
						on:select={handleForgeTypeSelect}
					/>

					<div>
						<label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Credentials Name <span class="text-red-500">*</span>
						</label>
						<input
							type="text"
							id="name"
							bind:value={formData.name}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="e.g., my-github-credentials"
						/>
					</div>

				<div>
					<label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Description <span class="text-red-500">*</span>
					</label>
					<textarea
						id="description"
						bind:value={formData.description}
						rows="2"
						required
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						placeholder="Brief description of these credentials"
					></textarea>
				</div>

				<div>
					<label for="endpoint" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Endpoint <span class="text-red-500">*</span>
					</label>
					<select
						id="endpoint"
						bind:value={formData.endpoint}
						required
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
					>
						<option value="">Select an endpoint</option>
						{#each filteredEndpoints as endpoint}
							<option value={endpoint.name}>
								{endpoint.name} ({endpoint.endpoint_type})
							</option>
						{/each}
					</select>
					{#if selectedForgeType}
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
							Showing only {selectedForgeType} endpoints
						</p>
					{/if}
				</div>

				<!-- Authentication Type Selection -->
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
						Authentication Type <span class="text-red-500">*</span>
					</label>
					<div class="flex space-x-4">
						<button
							type="button"
							on:click={() => handleAuthTypeChange(AuthType.PAT)}
							class="flex-1 py-2 px-4 text-sm font-medium rounded-md border focus:outline-none focus:ring-2 focus:ring-blue-500
								{selectedAuthType === AuthType.PAT 
									? 'bg-blue-600 text-white border-blue-600' 
									: 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-600'}
								{formData.endpoint && isGiteaEndpoint(formData.endpoint) ? '' : ''}"
						>
							PAT
						</button>
						<button
							type="button"
							on:click={() => handleAuthTypeChange(AuthType.APP)}
							disabled={selectedForgeType === 'gitea'}
							class="flex-1 py-2 px-4 text-sm font-medium rounded-md border focus:outline-none focus:ring-2 focus:ring-blue-500
								{selectedAuthType === AuthType.APP 
									? 'bg-blue-600 text-white border-blue-600' 
									: 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-600'}
								{selectedForgeType === 'gitea' ? 'opacity-50 cursor-not-allowed' : ''}"
						>
							App
						</button>
					</div>
					{#if selectedForgeType === 'gitea'}
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Gitea only supports PAT authentication</p>
					{/if}
				</div>

				<!-- PAT Fields -->
				{#if selectedAuthType === AuthType.PAT}
					<div>
						<label for="pat_token" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Personal Access Token <span class="text-red-500">*</span>
						</label>
						<input
							type="password"
							id="pat_token"
							bind:value={formData.pat_token}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="ghp_xxxxxxxxxxxxxxxxxxxx"
						/>
					</div>
				{/if}

				<!-- App Fields -->
				{#if selectedAuthType === AuthType.APP}
					<div>
						<label for="app_id" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							App ID <span class="text-red-500">*</span>
						</label>
						<input
							type="text"
							id="app_id"
							bind:value={formData.app_id}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="123456"
						/>
					</div>

					<div>
						<label for="app_installation_id" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							App Installation ID <span class="text-red-500">*</span>
						</label>
						<input
							type="text"
							id="app_installation_id"
							bind:value={formData.app_installation_id}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="12345678"
						/>
					</div>

					<div>
						<label for="private_key" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Private Key <span class="text-red-500">*</span>
						</label>
						<div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 text-center hover:border-blue-400 dark:hover:border-blue-400 transition-colors">
							<input
								type="file"
								id="private_key"
								accept=".pem,.key"
								on:change={handlePrivateKeyUpload}
								class="hidden"
							/>
							<div class="space-y-2">
								<svg class="mx-auto h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
								</svg>
								<p class="text-sm text-gray-600 dark:text-gray-400">
									<button type="button" on:click={() => document.getElementById('private_key')?.click()} class="text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 hover:underline">
										Choose a file
									</button>
									or drag and drop
								</p>
								<p class="text-xs text-gray-500 dark:text-gray-400">PEM, KEY files only</p>
							</div>
						</div>
					</div>
				{/if}

				<div class="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-700">
					<button
						type="button"
						on:click={closeModals}
						class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-600 hover:bg-gray-200 dark:hover:bg-gray-500 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={!isFormValid()}
						class="px-4 py-2 text-sm font-medium text-white rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors
							{isFormValid() ? 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500' : 'bg-gray-400 cursor-not-allowed'}"
					>
						Create Credentials
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal && editingCredential}
	<div class="fixed inset-0 bg-black bg-opacity-50 dark:bg-gray-900 dark:bg-opacity-75 flex items-center justify-center z-50" on:click={closeModals}>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-screen overflow-y-auto" on:click|stopPropagation>
			<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
				<div>
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white">
						Edit Credentials
					</h3>
					<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
						Update credentials for {editingCredential.name}
					</p>
				</div>
				<button on:click={closeModals} class="text-gray-400 hover:text-gray-600 dark:text-gray-300 dark:hover:text-gray-100">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			
			<form on:submit|preventDefault={handleUpdateCredentials} class="p-6 space-y-4">
				<div>
					<label for="edit_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Credentials Name <span class="text-red-500">*</span>
					</label>
					<input
						type="text"
						id="edit_name"
						bind:value={formData.name}
						required
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
					/>
				</div>

				<div>
					<label for="edit_description" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Description <span class="text-red-500">*</span>
					</label>
					<textarea
						id="edit_description"
						bind:value={formData.description}
						rows="2"
						required
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
					></textarea>
				</div>

				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Endpoint
					</label>
					<input
						type="text"
						value={formData.endpoint}
						disabled
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-100 dark:bg-gray-600 text-gray-500 dark:text-gray-400 cursor-not-allowed"
					/>
					<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Endpoint cannot be changed after creation</p>
				</div>

				<!-- Authentication Type Display -->
				<div>
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Authentication Type
					</label>
					<div class="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-100 dark:bg-gray-600">
						<span class="text-sm font-medium text-gray-700 dark:text-gray-300">
							{editingCredential['auth-type'] === AuthType.PAT ? 'Personal Access Token (PAT)' : 'GitHub App'}
						</span>
					</div>
					<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">Authentication type cannot be changed after creation</p>
				</div>

				<!-- PAT Fields -->
				{#if editingCredential['auth-type'] === AuthType.PAT}
					<div>
						<label for="edit_pat_token" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Personal Access Token <span class="text-red-500">*</span>
						</label>
						<input
							type="password"
							id="edit_pat_token"
							bind:value={formData.pat_token}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="Enter new token or leave empty to keep current"
						/>
					</div>
				{/if}

				<!-- App Fields -->
				{#if editingCredential['auth-type'] === AuthType.APP}
					<div>
						<label for="edit_app_id" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							App ID <span class="text-red-500">*</span>
						</label>
						<input
							type="text"
							id="edit_app_id"
							bind:value={formData.app_id}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						/>
					</div>

					<div>
						<label for="edit_app_installation_id" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							App Installation ID <span class="text-red-500">*</span>
						</label>
						<input
							type="text"
							id="edit_app_installation_id"
							bind:value={formData.app_installation_id}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						/>
					</div>

					<div>
						<label for="edit_private_key" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Private Key
						</label>
						<div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 text-center hover:border-blue-400 dark:hover:border-blue-400 transition-colors">
							<input
								type="file"
								id="edit_private_key"
								accept=".pem,.key"
								on:change={handlePrivateKeyUpload}
								class="hidden"
							/>
							<div class="space-y-2">
								<svg class="mx-auto h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
								</svg>
								<p class="text-sm text-gray-600 dark:text-gray-400">
									<button type="button" on:click={() => document.getElementById('edit_private_key')?.click()} class="text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 hover:underline">
										Choose a file
									</button>
									or drag and drop
								</p>
								<p class="text-xs text-gray-500 dark:text-gray-400">PEM, KEY files only. Leave empty to keep current key.</p>
							</div>
						</div>
					</div>
				{/if}

				<div class="flex justify-end space-x-3 pt-4 border-t border-gray-200 dark:border-gray-700">
					<button
						type="button"
						on:click={closeModals}
						class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-600 hover:bg-gray-200 dark:hover:bg-gray-500 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
					>
						Update Credentials
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Modal -->
{#if showDeleteModal && deletingCredential}
	<div class="fixed inset-0 bg-black bg-opacity-50 dark:bg-gray-900 dark:bg-opacity-75 flex items-center justify-center z-50" on:click={closeModals}>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full mx-4" on:click|stopPropagation>
			<div class="px-6 py-4">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-lg font-medium text-gray-900 dark:text-white">Delete Credentials</h3>
						<p class="mt-2 text-sm text-gray-500 dark:text-gray-300">
							Are you sure you want to delete the credentials "{deletingCredential.name}"? This action cannot be undone.
						</p>
					</div>
				</div>
			</div>
			<div class="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end space-x-3">
				<button
					type="button"
					on:click={closeModals}
					class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-600 hover:bg-gray-200 dark:hover:bg-gray-500 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
				>
					Cancel
				</button>
				<button
					type="button"
					on:click={handleDeleteCredentials}
					class="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
				>
					Delete
				</button>
			</div>
		</div>
	</div>
{/if}