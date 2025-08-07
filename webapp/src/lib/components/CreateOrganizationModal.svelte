<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { garmApi } from '$lib/api/client.js';
	import type { CreateOrgParams, ForgeCredentials, Endpoint, PoolBalancerType } from '$lib/api/types.js';
	import Modal from './Modal.svelte';
	import ForgeTypeSelector from './ForgeTypeSelector.svelte';

	const dispatch = createEventDispatcher<{
		close: void;
		submit: CreateOrgParams;
	}>();

	let loading = false;
	let error = '';
	let credentials: ForgeCredentials[] = [];
	let selectedForgeType: 'github' | 'gitea' | '' = 'github';

	// Form data
	let formData: CreateOrgParams = {
		name: '',
		credentials_name: '',
		webhook_secret: '',
		pool_balancer_type: 'roundrobin'
	};

	let installWebhook = true;
	let generateWebhookSecret = true;

	// Filtered credentials based on selected forge type
	$: filteredCredentials = credentials.filter(cred => {
		if (!selectedForgeType) return true;
		return cred.forge_type === selectedForgeType;
	});

	async function loadAllCredentials() {
		try {
			loading = true;
			error = '';
			credentials = await garmApi.listAllCredentials();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load credentials';
		} finally {
			loading = false;
		}
	}

	function handleForgeTypeSelect(event: CustomEvent<'github' | 'gitea'>) {
		selectedForgeType = event.detail;
		// Reset credential selection when forge type changes
		formData.credentials_name = '';
	}

	function handleCredentialChange() {
		// Auto-detect forge type when credential is selected
		if (formData.credentials_name) {
			const credential = credentials.find(c => c.name === formData.credentials_name);
			if (credential && credential.forge_type) {
				selectedForgeType = credential.forge_type;
			}
		}
	}

	// Generate secure random webhook secret
	function generateSecureWebhookSecret(): string {
		const array = new Uint8Array(32);
		crypto.getRandomValues(array);
		return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('');
	}

	// Auto-generate webhook secret when checkbox is checked
	$: if (generateWebhookSecret) {
		formData.webhook_secret = generateSecureWebhookSecret();
	} else if (!generateWebhookSecret) {
		// Clear the secret if user unchecks auto-generate
		formData.webhook_secret = '';
	}

	// Check if all mandatory fields are filled
	$: isFormValid = formData.name.trim() !== '' && 
					 formData.credentials_name !== '' &&
					 (generateWebhookSecret || formData.webhook_secret.trim() !== '');

	async function handleSubmit() {
		if (!formData.name.trim()) {
			error = 'Organization name is required';
			return;
		}

		if (!formData.credentials_name) {
			error = 'Please select credentials';
			return;
		}

		try {
			loading = true;
			error = '';

			const submitData: CreateOrgParams = {
				...formData
			};

			dispatch('submit', submitData);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create organization';
			loading = false;
		}
	}

	function getForgeIcon(forgeType: 'github' | 'gitea') {
		if (forgeType === 'gitea') {
			return `<svg class="w-8 h-8" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else {
			return `<div class="inline-flex w-8 h-8"><svg class="w-8 h-8 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-8 h-8 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
		}
	}

	onMount(() => {
		loadAllCredentials();
	});
</script>

<Modal on:close={() => dispatch('close')}>
	<div class="w-full p-6" style="min-width: 600px;">
		<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Create Organization</h3>

		{#if error}
			<div class="mb-4 rounded-md bg-red-50 dark:bg-red-900 p-4">
				<p class="text-sm font-medium text-red-800 dark:text-red-200">{error}</p>
			</div>
		{/if}

		{#if loading}
			<div class="text-center py-4">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
				<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Loading...</p>
			</div>
		{:else}
			<form on:submit|preventDefault={handleSubmit} class="space-y-4">
				<!-- Forge Type Selection -->
				<ForgeTypeSelector 
					bind:selectedForgeType 
					on:select={handleForgeTypeSelect}
				/>

				<!-- Organization Name -->
				<div>
					<label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Organization Name
					</label>
					<input
						id="name"
						type="text"
						bind:value={formData.name}
						required
						class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white sm:text-sm"
						placeholder="Enter organization name"
					/>
				</div>

				<!-- Credentials -->
				<div>
					<label for="credentials" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
						Credentials
					</label>
					<select
						id="credentials"
						bind:value={formData.credentials_name}
						on:change={handleCredentialChange}
						required
						class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white sm:text-sm"
					>
						<option value="">Select credentials...</option>
						{#each filteredCredentials as credential}
							<option value={credential.name}>
								{credential.name} ({credential.endpoint?.name || 'Unknown endpoint'})
							</option>
						{/each}
					</select>
				</div>

				<!-- Pool Balancer Type -->
				<div>
					<div class="flex items-center mb-1">
						<label for="pool_balancer_type" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
							Pool Balancer Type
						</label>
						<div class="ml-2 relative group">
							<svg class="w-4 h-4 text-gray-400 cursor-help" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
							</svg>
							<div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 w-80 p-3 bg-gray-900 text-white text-xs rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
								<div class="mb-2">
									<strong>Round Robin:</strong> Cycles through pools in turn. Job 1 → Pool 1, Job 2 → Pool 2, etc.
								</div>
								<div>
									<strong>Pack:</strong> Uses first available pool until full, then moves to next pool.
								</div>
								<div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
							</div>
						</div>
					</div>
					<select
						id="pool_balancer_type"
						bind:value={formData.pool_balancer_type}
						class="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white sm:text-sm"
					>
						<option value="roundrobin">Round Robin</option>
						<option value="pack">Pack</option>
					</select>
				</div>

				<!-- Webhook Configuration -->
				<div>
					<div class="flex items-center mb-3">
						<input
							id="install-webhook"
							type="checkbox"
							bind:checked={installWebhook}
							class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 dark:border-gray-600 rounded"
						/>
						<label for="install-webhook" class="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300">
							Install Webhook
						</label>
					</div>
					
					<div class="space-y-3">
						<div class="flex items-center">
							<input
								id="generate-webhook-secret"
								type="checkbox"
								bind:checked={generateWebhookSecret}
								class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 dark:border-gray-600 rounded"
							/>
							<label for="generate-webhook-secret" class="ml-2 text-sm text-gray-700 dark:text-gray-300">
								Auto-generate webhook secret
							</label>
						</div>
						
						{#if !generateWebhookSecret}
							<input
								type="password"
								bind:value={formData.webhook_secret}
								class="block w-full px-3 py-2 mt-3 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-white sm:text-sm"
								placeholder="Enter webhook secret"
							/>
						{:else}
							<p class="text-sm text-gray-500 dark:text-gray-400">
								Webhook secret will be automatically generated
							</p>
						{/if}
					</div>
				</div>

				<!-- Actions -->
				<div class="flex justify-end space-x-3 pt-4">
					<button
						type="button"
						on:click={() => dispatch('close')}
						class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:focus:ring-offset-gray-900"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={loading || !isFormValid}
						class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:focus:ring-offset-gray-900 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{loading ? 'Creating...' : 'Create Organization'}
					</button>
				</div>
			</form>
		{/if}
	</div>
</Modal>