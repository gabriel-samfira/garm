<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { garmApi } from '$lib/api/client.js';
	import type { Enterprise, CreateEnterpriseParams, UpdateEntityParams, Endpoint, Credential } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import CreateEnterpriseModal from '$lib/components/CreateEnterpriseModal.svelte';
	import UpdateEntityModal from '$lib/components/UpdateEntityModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';
	import { websocketStore, type WebSocketEvent } from '$lib/stores/websocket.js';

	let enterprises: Enterprise[] = [];
	let loading = true;
	let error = '';
	let searchTerm = '';
	let currentPage = 1;
	let itemsPerPage = 25;
	let showCreateModal = false;
	let showUpdateModal = false;
	let showDeleteModal = false;
	let selectedEnterprise: Enterprise | null = null;
	let unsubscribeWebsocket: (() => void) | null = null;


	// Filtered and paginated data
	$: filteredEnterprises = enterprises.filter((ent) =>
		ent.name.toLowerCase().includes(searchTerm.toLowerCase())
	);
	$: totalPages = Math.ceil(filteredEnterprises.length / itemsPerPage);
	$: paginatedEnterprises = filteredEnterprises.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	async function loadEnterprises() {
		try {
			loading = true;
			error = '';
			enterprises = await garmApi.listEnterprises();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load enterprises';
		} finally {
			loading = false;
		}
	}

	function handleEnterpriseEvent(event: WebSocketEvent) {
		
		if (event.operation === 'create') {
			const newEnterprise = event.payload as Enterprise;
			enterprises = [...enterprises, newEnterprise];
		} else if (event.operation === 'update') {
			const updatedEnterprise = event.payload as Enterprise;
			enterprises = enterprises.map(enterprise => 
				enterprise.id === updatedEnterprise.id ? updatedEnterprise : enterprise
			);
		} else if (event.operation === 'delete') {
			const enterpriseId = event.payload.id || event.payload;
			enterprises = enterprises.filter(enterprise => enterprise.id !== enterpriseId);
		}
	}

	async function handleCreateEnterprise(params: CreateEnterpriseParams) {
		try {
			await garmApi.createEnterprise(params);
			// No need to reload - websocket will handle the update
			showCreateModal = false;
		} catch (err) {
			throw err; // Let the modal handle the error
		}
	}

	async function handleUpdateEnterprise(params: UpdateEntityParams) {
		if (!selectedEnterprise) return;
		try {
			await garmApi.updateEnterprise(selectedEnterprise.id, params);
			// No need to reload - websocket will handle the update
			showUpdateModal = false;
			selectedEnterprise = null;
		} catch (err) {
			throw err; // Let the modal handle the error
		}
	}

	async function handleDeleteEnterprise() {
		if (!selectedEnterprise) return;
		try {
			await garmApi.deleteEnterprise(selectedEnterprise.id);
			// No need to reload - websocket will handle the update
			showDeleteModal = false;
			selectedEnterprise = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete enterprise';
		}
	}

	function openCreateModal() {
		showCreateModal = true;
	}

	function openUpdateModal(enterprise: Enterprise) {
		selectedEnterprise = enterprise;
		showUpdateModal = true;
	}

	function openDeleteModal(enterprise: Enterprise) {
		selectedEnterprise = enterprise;
		showDeleteModal = true;
	}

	function getForgeIcon() {
		// Enterprises are GitHub only
		return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
	}

	onMount(() => {
		loadEnterprises();
		
		// Subscribe to real-time enterprise events
		unsubscribeWebsocket = websocketStore.subscribeToEntity(
			'enterprise',
			['create', 'update', 'delete'],
			handleEnterpriseEvent
		);
	});

	onDestroy(() => {
		if (unsubscribeWebsocket) {
			unsubscribeWebsocket();
			unsubscribeWebsocket = null;
		}
	});
</script>

<svelte:head>
	<title>Enterprises - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="sm:flex sm:items-center sm:justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Enterprises</h1>
			<p class="mt-2 text-sm text-gray-700 dark:text-gray-300">Manage GitHub enterprises</p>
		</div>
		<div class="mt-4 sm:mt-0 flex items-center space-x-4">
			<button
				on:click={openCreateModal}
				class="inline-flex items-center justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
			>
				<svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
				</svg>
				Add Enterprise
			</button>
		</div>
	</div>

	<!-- Search and filters -->
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg p-4">
		<div class="sm:flex sm:items-center sm:justify-between">
			<div class="flex-1 min-w-0">
				<div class="max-w-md">
					<label for="search" class="sr-only">Search enterprises</label>
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
							placeholder="Search enterprises..."
						/>
					</div>
				</div>
			</div>
			<div class="mt-4 sm:mt-0 sm:ml-4 flex items-center space-x-4">
				<div class="flex items-center space-x-2">
					<label for="per-page" class="text-sm text-gray-700 dark:text-gray-300">Show:</label>
					<select
						id="per-page"
						bind:value={itemsPerPage}
						on:change={() => currentPage = 1}
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

	<!-- Enterprises table -->
	<div class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
		{#if loading}
			<div class="p-6 text-center">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
				<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Loading enterprises...</p>
			</div>
		{:else if error}
			<div class="p-6 text-center">
				<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
					<p class="text-sm font-medium text-red-800 dark:text-red-200">{error}</p>
				</div>
			</div>
		{:else if paginatedEnterprises.length === 0}
			<div class="p-6 text-center">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No enterprises found</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
					{#if searchTerm}
						No enterprises match your search criteria.
					{:else}
						Get started by creating your first enterprise.
					{/if}
				</p>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Endpoint</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Credential</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
							<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Actions</th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each paginatedEnterprises as enterprise}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div>
											<a href={`${base}/enterprises/${enterprise.id}`} class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
												{enterprise.name}
											</a>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="flex-shrink-0">
											{@html getForgeIcon()}
										</div>
										<div class="ml-2">
											<div class="text-sm text-gray-900 dark:text-white">{enterprise.endpoint.name}</div>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
									{enterprise.credentials_name}
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {enterprise.pool_manager_status?.running ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'}">
										{enterprise.pool_manager_status?.running ? 'Running' : 'Stopped'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex justify-end space-x-2">
										<button
											on:click={() => openUpdateModal(enterprise)}
											class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-900 dark:hover:text-indigo-300"
											title="Edit enterprise"
										>
											<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
											</svg>
										</button>
										<button
											on:click={() => openDeleteModal(enterprise)}
											class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300"
											title="Delete enterprise"
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
							on:click={() => currentPage = Math.max(1, currentPage - 1)}
							disabled={currentPage === 1}
							class="relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							Previous
						</button>
						<button
							on:click={() => currentPage = Math.min(totalPages, currentPage + 1)}
							disabled={currentPage === totalPages}
							class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							Next
						</button>
					</div>
					<div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
						<div>
							<p class="text-sm text-gray-700 dark:text-gray-300">
								Showing <span class="font-medium">{(currentPage - 1) * itemsPerPage + 1}</span>
								to <span class="font-medium">{Math.min(currentPage * itemsPerPage, filteredEnterprises.length)}</span>
								of <span class="font-medium">{filteredEnterprises.length}</span> results
							</p>
						</div>
						<div>
							<nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
								<button
									on:click={() => currentPage = Math.max(1, currentPage - 1)}
									disabled={currentPage === 1}
									class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
									</svg>
								</button>
								
								{#each Array.from({length: Math.min(5, totalPages)}, (_, i) => {
									const start = Math.max(1, currentPage - 2);
									const end = Math.min(totalPages, start + 4);
									return start + i;
								}).filter(p => p <= totalPages) as page}
									<button
										on:click={() => currentPage = page}
										class="relative inline-flex items-center px-4 py-2 border text-sm font-medium {currentPage === page ? 'z-10 bg-blue-50 dark:bg-blue-900 border-blue-500 text-blue-600 dark:text-blue-200' : 'bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600'}"
									>
										{page}
									</button>
								{/each}
								
								<button
									on:click={() => currentPage = Math.min(totalPages, currentPage + 1)}
									disabled={currentPage === totalPages}
									class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
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

<!-- Modals -->
{#if showCreateModal}
	<CreateEnterpriseModal
		on:close={() => showCreateModal = false}
		on:submit={(e) => handleCreateEnterprise(e.detail)}
	/>
{/if}

{#if showUpdateModal && selectedEnterprise}
	<UpdateEntityModal
		entity={selectedEnterprise}
		entityType="enterprise"
		on:close={() => { showUpdateModal = false; selectedEnterprise = null; }}
		on:submit={(e) => handleUpdateEnterprise(e.detail)}
	/>
{/if}

{#if showDeleteModal && selectedEnterprise}
	<DeleteModal
		title="Delete Enterprise"
		message="Are you sure you want to delete this enterprise? This action cannot be undone."
		itemName={selectedEnterprise.name}
		on:close={() => { showDeleteModal = false; selectedEnterprise = null; }}
		on:confirm={handleDeleteEnterprise}
	/>
{/if}