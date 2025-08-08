<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { garmApi } from '$lib/api/client.js';
	import type { Endpoint } from '$lib/api/types.js';
	import ForgeTypeSelector from '$lib/components/ForgeTypeSelector.svelte';
	import { eagerCache, eagerCacheManager } from '$lib/stores/eager-cache.js';
	import { toastStore } from '$lib/stores/toast.js';

	let loading = true;
	let endpoints: Endpoint[] = [];
	let error = '';

	// Subscribe to eager cache for endpoints
	$: endpoints = $eagerCache.endpoints;
	$: loading = $eagerCache.loading.endpoints;
	$: cacheError = $eagerCache.errorMessages.endpoints;
	let showCreateModal = false;
	let showEditModal = false;
	let showDeleteModal = false;
	let selectedForgeType: 'github' | 'gitea' | '' = 'github';
	let editingEndpoint: Endpoint | null = null;
	let deletingEndpoint: Endpoint | null = null;
	// Form state
	let formData = {
		name: '',
		description: '',
		endpoint_type: '',
		base_url: '',
		api_base_url: '',
		upload_base_url: '',
		ca_cert_bundle: ''
	};
	// Track original values for comparison during updates
	let originalFormData: typeof formData = { ...formData };


	onMount(async () => {
		// Load endpoints through eager cache (priority load + background load others)
		try {
			await eagerCacheManager.getEndpoints();
		} catch (err) {
			// Cache error is already handled by the eager cache system
			// We don't need to set error here anymore since it's in the cache state
			console.error('Failed to load endpoints:', err);
		}
	});

	async function retryLoadEndpoints() {
		try {
			await eagerCacheManager.retryResource('endpoints');
		} catch (err) {
			console.error('Retry failed:', err);
		}
	}

	// Endpoints are now handled by eager cache with websocket subscriptions


	function getForgeIcon(forgeType: string) {
		if (forgeType === 'gitea') {
			return `<svg class="w-8 h-8" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else {
			return `<div class="inline-flex w-8 h-8"><svg class="w-8 h-8 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-8 h-8 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
		}
	}

	function getForgeIconTable(endpointType: string) {
		if (endpointType === 'gitea') {
			return `<svg class="w-5 h-5" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else {
			return `<svg class="w-5 h-5 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
				<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/>
			</svg>
			<svg class="w-5 h-5 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
				<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.80-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#ffffff"/>
			</svg>`;
		}
	}

	function showCreateEndpointModal() {
		selectedForgeType = 'github';
		resetForm();
		showCreateModal = true;
	}

	function handleForgeTypeSelect(event: CustomEvent<'github' | 'gitea'>) {
		selectedForgeType = event.detail;
		formData.endpoint_type = event.detail;
	}

	function showEditEndpointModal(endpoint: Endpoint) {
		editingEndpoint = endpoint;
		formData = {
			name: endpoint.name || '',
			description: endpoint.description || '',
			endpoint_type: endpoint.endpoint_type || '',
			base_url: endpoint.base_url || '',
			api_base_url: endpoint.api_base_url || '',
			upload_base_url: endpoint.upload_base_url || '',
			ca_cert_bundle: endpoint.ca_cert_bundle || ''
		};
		// Store original values for comparison
		originalFormData = { ...formData };
		showEditModal = true;
	}

	function showDeleteEndpointModal(endpoint: Endpoint) {
		deletingEndpoint = endpoint;
		showDeleteModal = true;
	}

	function resetForm() {
		formData = {
			name: '',
			description: '',
			endpoint_type: '',
			base_url: '',
			api_base_url: '',
			upload_base_url: '',
			ca_cert_bundle: ''
		};
		originalFormData = { ...formData };
	}

	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		selectedForgeType = 'github';
		editingEndpoint = null;
		deletingEndpoint = null;
		resetForm();
	}

	function buildUpdateParams() {
		const updateParams: any = {};
		
		// Only include fields that have changed from original values
		if (formData.description !== originalFormData.description) {
			// Only set if not empty or if it was intentionally cleared
			if (formData.description.trim() !== '' || originalFormData.description !== '') {
				updateParams.description = formData.description.trim();
			}
		}
		
		if (formData.base_url !== originalFormData.base_url) {
			if (formData.base_url.trim() !== '') {
				updateParams.base_url = formData.base_url.trim();
			}
		}
		
		if (formData.api_base_url !== originalFormData.api_base_url) {
			// For Gitea, api_base_url is optional, so allow empty
			// For GitHub, it's required so only set if not empty
			if (formData.api_base_url.trim() !== '' || originalFormData.api_base_url !== '') {
				updateParams.api_base_url = formData.api_base_url.trim();
			}
		}
		
		// GitHub-only field
		if (editingEndpoint?.endpoint_type === 'github' && formData.upload_base_url !== originalFormData.upload_base_url) {
			if (formData.upload_base_url.trim() !== '' || originalFormData.upload_base_url !== '') {
				updateParams.upload_base_url = formData.upload_base_url.trim();
			}
		}
		
		if (formData.ca_cert_bundle !== originalFormData.ca_cert_bundle) {
			// CA cert can be cleared by setting to empty
			if (formData.ca_cert_bundle !== '') {
				// Convert base64 string to byte array for API
				try {
					const bytes = atob(formData.ca_cert_bundle);
					updateParams.ca_cert_bundle = Array.from(bytes, char => char.charCodeAt(0));
				} catch (e) {
					// If not valid base64, treat as empty
					if (originalFormData.ca_cert_bundle !== '') {
						updateParams.ca_cert_bundle = [];
					}
				}
			} else if (originalFormData.ca_cert_bundle !== '') {
				// User intentionally cleared the CA cert
				updateParams.ca_cert_bundle = [];
			}
		}
		
		return updateParams;
	}

	async function handleCreateEndpoint() {
		try {
			if (formData.endpoint_type === 'github') {
				await garmApi.createGithubEndpoint(formData);
			} else {
				await garmApi.createGiteaEndpoint(formData);
			}
			// No need to reload - websocket will handle the update
			toastStore.success(
				'Endpoint Created',
				`Endpoint ${formData.name} has been created successfully.`
			);
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create endpoint';
		}
	}

	async function handleUpdateEndpoint() {
		if (!editingEndpoint) return;
		
		try {
			const updateParams = buildUpdateParams();
			
			// Only proceed if there are changes to apply
			if (Object.keys(updateParams).length === 0) {
				toastStore.info(
					'No Changes',
					'No fields were modified.'
				);
				closeModals();
				return;
			}
			
			if (editingEndpoint.endpoint_type === 'github') {
				await garmApi.updateGithubEndpoint(editingEndpoint.name, updateParams);
			} else {
				await garmApi.updateGiteaEndpoint(editingEndpoint.name, updateParams);
			}
			// No need to reload - websocket will handle the update
			toastStore.success(
				'Endpoint Updated',
				`Endpoint ${editingEndpoint.name} has been updated successfully.`
			);
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update endpoint';
		}
	}

	async function handleDeleteEndpoint() {
		if (!deletingEndpoint) return;
		
		try {
			if (deletingEndpoint.endpoint_type === 'github') {
				await garmApi.deleteGithubEndpoint(deletingEndpoint.name);
			} else {
				await garmApi.deleteGiteaEndpoint(deletingEndpoint.name);
			}
			// No need to reload - websocket will handle the update
			toastStore.success(
				'Endpoint Deleted',
				`Endpoint ${deletingEndpoint.name} has been deleted successfully.`
			);
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete endpoint';
		}
	}

	function handleFileUpload(event: Event) {
		const target = event.target as HTMLInputElement;
		const file = target.files?.[0];
		
		if (!file) {
			formData.ca_cert_bundle = '';
			return;
		}

		const reader = new FileReader();
		reader.onload = (e) => {
			const content = e.target?.result as string;
			formData.ca_cert_bundle = btoa(content);
		};
		reader.readAsText(file);
	}

	function isFormValid() {
		if (!formData.name || !formData.description || !formData.base_url) return false;
		if (formData.endpoint_type === 'github' && !formData.api_base_url) return false;
		return true;
	}
</script>

<svelte:head>
	<title>Endpoints - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="sm:flex sm:items-center">
		<div class="sm:flex-auto">
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Endpoints</h1>
			<p class="mt-2 text-sm text-gray-700 dark:text-gray-300">
				Manage your GitHub and Gitea endpoints for runner management.
			</p>
		</div>
		<div class="mt-4 sm:ml-16 sm:mt-0 sm:flex-none flex items-center space-x-4">
			<button
				type="button"
				on:click={showCreateEndpointModal}
				class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors duration-200 flex items-center space-x-2"
			>
				<span>Add Endpoint</span>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
				</svg>
			</button>
		</div>
	</div>

	{#if error || cacheError}
		<!-- Error state -->
		<div class="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
			<div class="flex">
				<div class="flex-shrink-0">
					<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
					</svg>
				</div>
				<div class="ml-3 flex-1">
					<h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error loading endpoints</h3>
					<p class="mt-2 text-sm text-red-700 dark:text-red-300">{cacheError || error}</p>
					{#if cacheError}
						<div class="mt-3">
							<button
								on:click={retryLoadEndpoints}
								class="inline-flex items-center px-3 py-1 border border-transparent text-sm leading-5 font-medium rounded text-red-700 dark:text-red-200 bg-red-100 dark:bg-red-800 hover:bg-red-200 dark:hover:bg-red-700 focus:outline-none focus:bg-red-200 dark:focus:bg-red-700 transition duration-150 ease-in-out"
							>
								<svg class="-ml-1 mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
								</svg>
								Retry
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}

	<!-- Table -->
	<div class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
		{#if loading}
			<div class="px-4 py-8 text-center">
				<div class="animate-spin mx-auto h-8 w-8 border-b-2 border-blue-600 rounded-full"></div>
				<p class="mt-4 text-sm text-gray-500 dark:text-gray-400">Loading endpoints...</p>
			</div>
		{:else if endpoints.length === 0}
			<div class="px-4 py-8 text-center">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No endpoints</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Get started by creating a new endpoint.</p>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Description</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">API URL</th>
							<th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Forge Type</th>
							<th scope="col" class="relative px-6 py-3"><span class="sr-only">Actions</span></th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each endpoints as endpoint}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
									{endpoint.name}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
									{endpoint.description}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
									{endpoint.api_base_url || endpoint.base_url}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">
									<div class="flex items-center">
										{@html getForgeIconTable(endpoint.endpoint_type || '')}
										<span class="ml-2 capitalize">{endpoint.endpoint_type}</span>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex justify-end space-x-2">
										<button
											on:click={() => showEditEndpointModal(endpoint)}
											class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-900 dark:hover:text-indigo-300 cursor-pointer"
											title="Edit endpoint"
										>
											<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
											</svg>
										</button>
										<button
											on:click={() => showDeleteEndpointModal(endpoint)}
											class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300 cursor-pointer"
											title="Delete endpoint"
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
	<div class="fixed inset-0 bg-black/30 dark:bg-black/50 flex items-center justify-center z-50" on:click={closeModals}>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-screen overflow-y-auto" on:click|stopPropagation>
			<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
				<div>
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white">
						Add Endpoint
					</h3>
					<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
						Connect to GitHub or Gitea for runner management
					</p>
				</div>
				<button on:click={closeModals} class="text-gray-400 hover:text-gray-600 dark:text-gray-300 dark:hover:text-gray-100">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			
			<form on:submit|preventDefault={handleCreateEndpoint} class="p-6 space-y-4">
				<!-- Forge Type Selection -->
				<ForgeTypeSelector 
					bind:selectedForgeType 
					on:select={handleForgeTypeSelect}
				/>
				<div>
					<label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Endpoint Name <span class="text-red-500">*</span>
					</label>
					<input
						type="text"
						id="name"
						bind:value={formData.name}
						required
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						placeholder={selectedForgeType === 'github' ? 'e.g., github-enterprise or github-com' : 'e.g., gitea-main or my-gitea'}
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
						placeholder="Brief description of this endpoint"
					></textarea>
				</div>

				<div>
					<label for="base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Base URL <span class="text-red-500">*</span>
					</label>
					<input
						type="url"
						id="base_url"
						bind:value={formData.base_url}
						required
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						placeholder={selectedForgeType === 'github' ? 'https://github.com or https://github.example.com' : 'https://gitea.example.com'}
					/>
				</div>

				{#if selectedForgeType === 'github'}
					<div>
						<label for="api_base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							API Base URL <span class="text-red-500">*</span>
						</label>
						<input
							type="url"
							id="api_base_url"
							bind:value={formData.api_base_url}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="https://api.github.com or https://github.example.com/api/v3"
						/>
					</div>

					<div>
						<label for="upload_base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Upload Base URL
						</label>
						<input
							type="url"
							id="upload_base_url"
							bind:value={formData.upload_base_url}
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="https://uploads.github.com"
						/>
					</div>
				{:else}
					<div>
						<label for="api_base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							API Base URL <span class="text-xs text-gray-500">(optional)</span>
						</label>
						<input
							type="url"
							id="api_base_url"
							bind:value={formData.api_base_url}
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
							placeholder="https://gitea.example.com/api/v1 (leave empty to use Base URL)"
						/>
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">If empty, Base URL will be used as API Base URL</p>
					</div>
				{/if}

				<!-- CA Certificate Upload -->
				<div class="space-y-3 border-t border-gray-200 dark:border-gray-700 pt-4">
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">CA Certificate Bundle (Optional)</label>
					<div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 text-center hover:border-blue-400 dark:hover:border-blue-400 transition-colors">
						<input
							type="file"
							id="ca_cert_file"
							accept=".pem,.crt,.cer,.cert"
							on:change={handleFileUpload}
							class="hidden"
						/>
						<div class="space-y-2">
							<svg class="mx-auto h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
							</svg>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								<button type="button" on:click={() => document.getElementById('ca_cert_file')?.click()} class="text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 hover:underline">
									Choose a file
								</button>
								or drag and drop
							</p>
							<p class="text-xs text-gray-500 dark:text-gray-400">PEM, CRT, CER, CERT files only</p>
						</div>
					</div>
				</div>

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
						Create Endpoint
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Edit Modal -->
{#if showEditModal && editingEndpoint}
	<div class="fixed inset-0 bg-black/30 dark:bg-black/50 flex items-center justify-center z-50" on:click={closeModals}>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-screen overflow-y-auto" on:click|stopPropagation>
			<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
				<div>
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white">
						Edit {editingEndpoint.endpoint_type === 'github' ? 'GitHub' : 'Gitea'} Endpoint
					</h3>
					<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
						Update endpoint configuration
					</p>
				</div>
				<button on:click={closeModals} class="text-gray-400 hover:text-gray-600 dark:text-gray-300 dark:hover:text-gray-100">
					<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
					</svg>
				</button>
			</div>
			
			<form on:submit|preventDefault={handleUpdateEndpoint} class="p-6 space-y-4">
				<div>
					<label for="edit_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Endpoint Name <span class="text-red-500">*</span>
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
					<label for="edit_base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Base URL <span class="text-red-500">*</span>
					</label>
					<input
						type="url"
						id="edit_base_url"
						bind:value={formData.base_url}
						required
						class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
					/>
				</div>

				{#if editingEndpoint.endpoint_type === 'github'}
					<div>
						<label for="edit_api_base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							API Base URL <span class="text-red-500">*</span>
						</label>
						<input
							type="url"
							id="edit_api_base_url"
							bind:value={formData.api_base_url}
							required
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						/>
					</div>

					<div>
						<label for="edit_upload_base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							Upload Base URL
						</label>
						<input
							type="url"
							id="edit_upload_base_url"
							bind:value={formData.upload_base_url}
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						/>
					</div>
				{:else}
					<div>
						<label for="edit_api_base_url" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
							API Base URL <span class="text-xs text-gray-500">(optional)</span>
						</label>
						<input
							type="url"
							id="edit_api_base_url"
							bind:value={formData.api_base_url}
							class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
						/>
						<p class="text-xs text-gray-500 dark:text-gray-400 mt-1">If empty, Base URL will be used as API Base URL</p>
					</div>
				{/if}

				<!-- CA Certificate Upload -->
				<div class="space-y-3 border-t border-gray-200 dark:border-gray-700 pt-4">
					<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">CA Certificate Bundle (Optional)</label>
					<div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 text-center hover:border-blue-400 dark:hover:border-blue-400 transition-colors">
						<input
							type="file"
							id="edit_ca_cert_file"
							accept=".pem,.crt,.cer,.cert"
							on:change={handleFileUpload}
							class="hidden"
						/>
						<div class="space-y-2">
							<svg class="mx-auto h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
							</svg>
							<p class="text-sm text-gray-600 dark:text-gray-400">
								<button type="button" on:click={() => document.getElementById('edit_ca_cert_file')?.click()} class="text-gray-900 dark:text-white hover:text-gray-700 dark:hover:text-gray-300 hover:underline">
									Choose a file
								</button>
								or drag and drop
							</p>
							<p class="text-xs text-gray-500 dark:text-gray-400">PEM, CRT, CER, CERT files only</p>
						</div>
					</div>
				</div>

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
						Update Endpoint
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Modal -->
{#if showDeleteModal && deletingEndpoint}
	<div class="fixed inset-0 bg-black/30 dark:bg-black/50 flex items-center justify-center z-50" on:click={closeModals}>
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full mx-4" on:click|stopPropagation>
			<div class="px-6 py-4">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-lg font-medium text-gray-900 dark:text-white">Delete Endpoint</h3>
						<p class="mt-2 text-sm text-gray-500 dark:text-gray-300">
							Are you sure you want to delete the endpoint "{deletingEndpoint.name}"? This action cannot be undone.
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
					on:click={handleDeleteEndpoint}
					class="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
				>
					Delete
				</button>
			</div>
		</div>
	</div>
{/if}