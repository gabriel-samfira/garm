<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { garmApi } from '$lib/api/client.js';
	import type { Repository, ForgeCredentials, CreateRepoParams, UpdateEntityParams } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import CreateRepositoryModal from '$lib/components/CreateRepositoryModal.svelte';
	import { websocketStore, type WebSocketEvent } from '$lib/stores/websocket.js';

	let repositories: Repository[] = [];
	let loading = true;
	let error = '';
	let searchTerm = '';
	let credentials: ForgeCredentials[] = [];
	let credentialsLoading = false;

	// Modal states
	let showCreateModal = false;
	let showEditModal = false;
	let showDeleteModal = false;
	let editingRepository: Repository | null = null;
	let deletingRepository: Repository | null = null;
	let unsubscribeWebsocket: (() => void) | null = null;


	let updateFormData: UpdateEntityParams & { change_webhook_secret?: boolean } = {
		credentials_name: '',
		webhook_secret: '',
		pool_balancer_type: '',
		change_webhook_secret: false
	};

	// Pagination
	let currentPage = 1;
	let perPage = 25;
	let totalPages = 1;

	$: filteredRepositories = repositories.filter(repo =>
		repo.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
		repo.owner.toLowerCase().includes(searchTerm.toLowerCase())
	);

	$: {
		totalPages = Math.ceil(filteredRepositories.length / perPage);
		if (currentPage > totalPages && totalPages > 0) {
			currentPage = totalPages;
		}
	}

	$: paginatedRepositories = filteredRepositories.slice(
		(currentPage - 1) * perPage,
		currentPage * perPage
	);

	function handleRepositoryEvent(event: WebSocketEvent) {
		console.log('[Repositories] Received websocket event:', event);
		
		if (event.operation === 'create') {
			// Add new repository
			const newRepository = event.payload as Repository;
			repositories = [...repositories, newRepository];
		} else if (event.operation === 'update') {
			// Update existing repository
			const updatedRepository = event.payload as Repository;
			repositories = repositories.map(repo => 
				repo.id === updatedRepository.id ? updatedRepository : repo
			);
		} else if (event.operation === 'delete') {
			// Remove repository - payload might only contain ID
			const repositoryId = event.payload.id || event.payload;
			repositories = repositories.filter(repo => repo.id !== repositoryId);
		}
	}

	onMount(async () => {
		// Initial load
		await loadRepositories();
		
		// Subscribe to real-time repository events
		unsubscribeWebsocket = websocketStore.subscribeToEntity(
			'repository',
			['create', 'update', 'delete'],
			handleRepositoryEvent
		);
	});

	onDestroy(() => {
		// Clean up websocket subscription
		if (unsubscribeWebsocket) {
			unsubscribeWebsocket();
			unsubscribeWebsocket = null;
		}
	});

	async function loadRepositories() {
		try {
			loading = true;
			error = '';
			repositories = await garmApi.listRepositories();
			console.log('Loaded repositories:', repositories);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load repositories';
			console.error('Error loading repositories:', err);
		} finally {
			loading = false;
		}
	}

	async function showEditRepositoryModal(repository: Repository) {
		editingRepository = repository;
		updateFormData = {
			credentials_name: repository.credentials_name || '',
			webhook_secret: '',
			pool_balancer_type: repository.pool_balancing_type || 'roundrobin',
			change_webhook_secret: false
		};
		
		// Load credentials for this repository's forge type
		const forgeType = repository.endpoint.endpoint_type;
		try {
			credentialsLoading = true;
			if (forgeType === 'github') {
				credentials = await garmApi.listGithubCredentials();
			} else {
				credentials = await garmApi.listGiteaCredentials();
			}
		} catch (err) {
			console.error('Error loading credentials:', err);
			credentials = [];
		} finally {
			credentialsLoading = false;
		}
		
		showEditModal = true;
	}

	function showDeleteRepositoryModal(repository: Repository) {
		deletingRepository = repository;
		showDeleteModal = true;
	}


	function closeModals() {
		showCreateModal = false;
		showEditModal = false;
		showDeleteModal = false;
		editingRepository = null;
		deletingRepository = null;
		error = '';
	}

	async function handleCreateRepository(event: CustomEvent<CreateRepoParams & { install_webhook?: boolean; auto_generate_secret?: boolean }>) {
		try {
			error = '';
			
			const data = event.detail;
			const repoParams: CreateRepoParams = {
				name: data.name,
				owner: data.owner,
				credentials_name: data.credentials_name,
				webhook_secret: data.webhook_secret
			};

			const createdRepo = await garmApi.createRepository(repoParams);
			
			// If install_webhook is checked, install the webhook
			if (data.install_webhook) {
				try {
					await garmApi.installRepoWebhook(createdRepo.id);
				} catch (webhookError) {
					console.warn('Repository created but webhook installation failed:', webhookError);
				}
			}

			// No need to reload - websocket will handle the update
			showCreateModal = false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create repository';
			throw err; // Let the modal handle the error display
		}
	}

	async function handleUpdateRepository() {
		if (!editingRepository) return;
		
		try {
			error = '';
			
			const updateParams: UpdateEntityParams = {};
			
			if (updateFormData.credentials_name !== editingRepository.credentials_name) {
				updateParams.credentials_name = updateFormData.credentials_name;
			}
			
			if (updateFormData.pool_balancer_type !== editingRepository.pool_balancing_type) {
				updateParams.pool_balancer_type = updateFormData.pool_balancer_type;
			}'Screenshot from 2025-08-07 02-30-06-1.png'
			
			if (updateFormData.change_webhook_secret && updateFormData.webhook_secret) {
				updateParams.webhook_secret = updateFormData.webhook_secret;
			}

			await garmApi.updateRepository(editingRepository.id, updateParams);
			// No need to reload - websocket will handle the update
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update repository';
		}
	}

	async function handleDeleteRepository() {
		if (!deletingRepository) return;
		
		try {
			error = '';
			await garmApi.deleteRepository(deletingRepository.id);
			// No need to reload - websocket will handle the update
			closeModals();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete repository';
		}
	}

	function getStatusBadge(repo: Repository) {
		if (repo.pool_manager_status.running) {
			return {
				text: 'Running',
				class: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
			};
		} else {
			return {
				text: 'Stopped',
				class: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
			};
		}
	}

	function getForgeIcon(endpointType: string) {
		if (endpointType === 'gitea') {
			return `<svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640">
				<path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/>
				<path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/>
			</svg>`;
		} else {
			return `<div class="inline-flex w-4 h-4">
				<svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
					<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/>
				</svg>
				<svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96">
					<path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/>
				</svg>
			</div>`;
		}
	}

	function changePage(page: number) {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
		}
	}

	function changePerPage(newPerPage: number) {
		perPage = newPerPage;
		currentPage = 1;
	}
</script>

<svelte:head>
	<title>Repositories - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex justify-between items-center">
		<div>
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Repositories</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Manage your GitHub repositories and their runners
			</p>
		</div>
		<div class="flex items-center space-x-4">
			<button 
				id="add-repo-button"
				class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors duration-200 flex items-center space-x-2"
				on:click={() => { showCreateModal = true; }}
			>
				<span>Add Repository</span>
				<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
				</svg>
			</button>
		</div>
	</div>

	<!-- Search and filters -->
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg p-4">
		<div class="sm:flex sm:items-center sm:justify-between">
			<div class="flex-1 min-w-0">
				<div class="max-w-md">
					<label for="search" class="sr-only">Search repositories</label>
					<div class="relative">
						<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
							<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
							</svg>
						</div>
						<input
							id="search"
							bind:value={searchTerm}
							on:input={() => currentPage = 1}
							type="text"
							class="block w-full pl-10 pr-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md leading-5 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
							placeholder="Search repositories by name or owner..."
						/>
					</div>
				</div>
			</div>
			<div class="mt-4 sm:mt-0 sm:ml-4 flex items-center space-x-4">
				<div class="flex items-center space-x-2">
					<label for="per-page" class="text-sm text-gray-700 dark:text-gray-300">Show:</label>
					<select
						id="per-page"
						bind:value={perPage}
						on:change={() => changePerPage(perPage)}
						class="block w-20 pl-2 pr-8 py-1 text-sm border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white rounded-md focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
					>
						<option value={25}>25</option>
						<option value={50}>50</option>
						<option value={100}>100</option>
					</select>
				</div>
			</div>
		</div>
	</div>

	<!-- Repositories table -->
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
		{#if loading}
			<div class="p-6 text-center">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
				<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">Loading repositories...</p>
			</div>
		{:else if error}
			<div class="p-6">
				<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
					<div class="flex">
						<div class="flex-shrink-0">
							<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
							</svg>
						</div>
						<div class="ml-3">
							<h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error loading repositories</h3>
							<p class="mt-2 text-sm text-red-700 dark:text-red-300">{error}</p>
						</div>
					</div>
				</div>
			</div>
		{:else if paginatedRepositories.length === 0}
			<div class="p-6 text-center">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
				<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
					{searchTerm ? `No repositories found matching "${searchTerm}"` : 'No repositories found'}
				</p>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
								Repository
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
								Endpoint
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
								Credentials
							</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
								Status
							</th>
							<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
								Actions
							</th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each paginatedRepositories as repo}
							{@const status = getStatusBadge(repo)}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="ml-2">
											<div class="text-sm font-medium text-gray-900 dark:text-white">
												<a href="{base}/repositories/{repo.id}" class="hover:text-blue-600 dark:hover:text-blue-400">
													{repo.owner}/{repo.name}
												</a>
											</div>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="mr-2">
											{@html getForgeIcon(repo.endpoint.endpoint_type)}
										</div>
										<div class="text-sm text-gray-900 dark:text-white">{repo.endpoint.name}</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
									{repo.credentials_name}
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full {status.class}">
										{status.text}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex justify-end space-x-2">
										<button
											on:click={() => showEditRepositoryModal(repo)}
											class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-900 dark:hover:text-indigo-300"
											title="Edit repository"
										>
											<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
											</svg>
										</button>
										<button
											on:click={() => showDeleteRepositoryModal(repo)}
											class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300"
											title="Delete repository"
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

			<!-- Pagination -->
			{#if totalPages > 1}
				<div class="bg-white dark:bg-gray-800 px-4 py-3 flex items-center justify-between border-t border-gray-200 dark:border-gray-700 sm:px-6">
					<div class="flex-1 flex justify-between sm:hidden">
						<button
							on:click={() => changePage(currentPage - 1)}
							disabled={currentPage === 1}
							class="relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							Previous
						</button>
						<button
							on:click={() => changePage(currentPage + 1)}
							disabled={currentPage === totalPages}
							class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							Next
						</button>
					</div>
					<div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
						<div>
							<p class="text-sm text-gray-700 dark:text-gray-300">
								Showing <span class="font-medium">{(currentPage - 1) * perPage + 1}</span>
								to <span class="font-medium">{Math.min(currentPage * perPage, filteredRepositories.length)}</span>
								of <span class="font-medium">{filteredRepositories.length}</span> results
							</p>
						</div>
						<div>
							<nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
								<button
									on:click={() => changePage(currentPage - 1)}
									disabled={currentPage === 1}
									class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm font-medium text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
									</svg>
								</button>
								
								{#each Array(totalPages) as _, i}
									{@const page = i + 1}
									<button
										on:click={() => changePage(page)}
										class="relative inline-flex items-center px-4 py-2 border text-sm font-medium
											{page === currentPage 
												? 'z-10 bg-blue-50 dark:bg-blue-900 border-blue-500 text-blue-600 dark:text-blue-200' 
												: 'bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-600'}"
									>
										{page}
									</button>
								{/each}

								<button
									on:click={() => changePage(currentPage + 1)}
									disabled={currentPage === totalPages}
									class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm font-medium text-gray-500 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
									</svg>
								</button>
							</nav>
						</div>
					</div>
				</div>
			{/if}
		{/if}
	</div>
</div>

<!-- Create Repository Modal -->
{#if showCreateModal}
	<CreateRepositoryModal
		on:close={() => showCreateModal = false}
		on:submit={handleCreateRepository}
	/>
{/if}

<!-- Edit Repository Modal -->
{#if showEditModal && editingRepository}
	<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" on:click={closeModals}>
		<div class="relative top-20 mx-auto p-5 border w-11/12 max-w-md shadow-lg rounded-md bg-white dark:bg-gray-800" on:click|stopPropagation>
			<div class="mt-3">
				<h3 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Edit Repository</h3>
				
				<form on:submit|preventDefault={handleUpdateRepository} class="space-y-4">
					{#if error}
						<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
							<div class="flex">
								<div class="flex-shrink-0">
									<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
										<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
									</svg>
								</div>
								<div class="ml-3">
									<p class="text-sm text-red-800 dark:text-red-200">{error}</p>
								</div>
							</div>
						</div>
					{/if}

					<div>
						<label class="block text-sm font-medium text-gray-700 dark:text-gray-300">Repository</label>
						<div class="mt-1 text-sm text-gray-900 dark:text-white">
							{editingRepository.owner}/{editingRepository.name}
						</div>
					</div>

					<div>
						<label for="edit-credentials" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Credentials</label>
						<select
							id="edit-credentials"
							bind:value={updateFormData.credentials_name}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:text-white sm:text-sm"
						>
							{#each credentials as credential}
								<option value={credential.name}>
									{credential.name} ({credential.endpoint?.name || credential.endpoint_name || 'Unknown'})
								</option>
							{/each}
						</select>
					</div>

					<div>
						<label for="edit-pool-balancer" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Pool Balancer Type</label>
						<select
							id="edit-pool-balancer"
							bind:value={updateFormData.pool_balancer_type}
							class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:text-white sm:text-sm"
						>
							<option value="roundrobin">Round Robin</option>
							<option value="pack">Pack</option>
						</select>
					</div>

					<div class="space-y-3">
						<div class="flex items-center">
							<input
								id="change-webhook-secret"
								type="checkbox"
								bind:checked={updateFormData.change_webhook_secret}
								class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 dark:border-gray-600 rounded dark:bg-gray-700"
							/>
							<label for="change-webhook-secret" class="ml-2 block text-sm text-gray-900 dark:text-white">
								I want to change the webhook secret
							</label>
						</div>

						{#if updateFormData.change_webhook_secret}
							<div>
								<label for="edit-webhook-secret" class="block text-sm font-medium text-gray-700 dark:text-gray-300">New Webhook Secret</label>
								<input
									id="edit-webhook-secret"
									type="password"
									bind:value={updateFormData.webhook_secret}
									required={updateFormData.change_webhook_secret}
									class="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white sm:text-sm"
									placeholder="Enter new webhook secret"
								/>
							</div>
						{/if}
					</div>

					<div class="flex justify-end space-x-3 pt-4">
						<button
							type="button"
							class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
							on:click={closeModals}
						>
							Cancel
						</button>
						<button
							type="submit"
							disabled={updateFormData.change_webhook_secret && !updateFormData.webhook_secret}
							class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							Update Repository
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Delete Repository Modal -->
{#if showDeleteModal && deletingRepository}
	<div class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" on:click={closeModals}>
		<div class="relative top-20 mx-auto p-5 border w-11/12 max-w-md shadow-lg rounded-md bg-white dark:bg-gray-800" on:click|stopPropagation>
			<div class="mt-3">
				<div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 dark:bg-red-900">
					<svg class="h-6 w-6 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
					</svg>
				</div>
				<div class="mt-3 text-center sm:mt-5">
					<h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">Delete Repository</h3>
					<div class="mt-2">
						<p class="text-sm text-gray-500 dark:text-gray-400">
							Are you sure you want to delete <strong>{deletingRepository.owner}/{deletingRepository.name}</strong>? 
							This action cannot be undone and will remove all associated pools and runners.
						</p>
					</div>
				</div>
			</div>

			{#if error}
				<div class="mt-4 rounded-md bg-red-50 dark:bg-red-900 p-4">
					<div class="flex">
						<div class="flex-shrink-0">
							<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
							</svg>
						</div>
						<div class="ml-3">
							<p class="text-sm text-red-800 dark:text-red-200">{error}</p>
						</div>
					</div>
				</div>
			{/if}

			<div class="mt-5 sm:mt-6 flex justify-end space-x-3">
				<button
					type="button"
					class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
					on:click={closeModals}
				>
					Cancel
				</button>
				<button
					type="button"
					class="px-4 py-2 text-sm font-medium text-white bg-red-600 border border-transparent rounded-md shadow-sm hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
					on:click={handleDeleteRepository}
				>
					Delete Repository
				</button>
			</div>
		</div>
	</div>
{/if}
