<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { garmApi } from '$lib/api/client.js';
	import type { Instance } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import DeleteModal from '$lib/components/DeleteModal.svelte';
	import { websocketStore, type WebSocketEvent } from '$lib/stores/websocket.js';

	let instances: Instance[] = [];
	let loading = true;
	let error = '';
	let statusFilter = '';
	let unsubscribeWebsocket: (() => void) | null = null;


	// Pagination
	let currentPage = 1;
	let itemsPerPage = 25;
	let searchTerm = '';

	// Modal state
	let showDeleteModal = false;
	let instanceToDelete: Instance | null = null;

	$: filteredInstances = instances.filter(instance => {
		const matchesSearch = searchTerm === '' || 
			instance.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
			instance.provider_id.toLowerCase().includes(searchTerm.toLowerCase());
		const matchesStatus = statusFilter === '' || instance.status === statusFilter || instance.runner_status === statusFilter;
		return matchesSearch && matchesStatus;
	});

	$: totalPages = Math.ceil(filteredInstances.length / itemsPerPage);
	$: paginatedInstances = filteredInstances.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	async function loadInstances() {
		try {
			loading = true;
			error = '';
			instances = await garmApi.listInstances();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load instances';
		} finally {
			loading = false;
		}
	}

	function handleDelete(instance: Instance) {
		instanceToDelete = instance;
		showDeleteModal = true;
	}

	async function confirmDelete() {
		if (!instanceToDelete) return;
		
		try {
			await garmApi.deleteInstance(instanceToDelete.name);
			// No need to reload - websocket will handle the update
			showDeleteModal = false;
			instanceToDelete = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete instance';
		}
	}

	function handleStatusFilterChange(event: Event) {
		statusFilter = (event.target as HTMLSelectElement).value;
		currentPage = 1;
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}

	function getStatusBadgeClass(status: string): string {
		switch (status) {
			case 'running':
				return 'bg-green-50 text-green-700 ring-green-600/20 dark:bg-green-500/10 dark:text-green-400 dark:ring-green-500/20';
			case 'idle':
				return 'bg-blue-50 text-blue-700 ring-blue-600/20 dark:bg-blue-500/10 dark:text-blue-400 dark:ring-blue-500/20';
			case 'active':
				return 'bg-yellow-50 text-yellow-700 ring-yellow-600/20 dark:bg-yellow-500/10 dark:text-yellow-400 dark:ring-yellow-500/20';
			case 'pending':
			case 'pending_create':
			case 'creating':
			case 'installing':
				return 'bg-gray-50 text-gray-700 ring-gray-600/20 dark:bg-gray-500/10 dark:text-gray-400 dark:ring-gray-500/20 animate-pulse';
			case 'pending_delete':
			case 'deleting':
			case 'terminating':
				return 'bg-orange-50 text-orange-700 ring-orange-600/20 dark:bg-orange-500/10 dark:text-orange-400 dark:ring-orange-500/20 animate-pulse';
			case 'failed':
			case 'terminated':
				return 'bg-red-50 text-red-700 ring-red-600/20 dark:bg-red-500/10 dark:text-red-400 dark:ring-red-500/20';
			default:
				return 'bg-gray-50 text-gray-700 ring-gray-600/20 dark:bg-gray-500/10 dark:text-gray-400 dark:ring-gray-500/20';
		}
	}

	function handleInstanceEvent(event: WebSocketEvent) {
		
		if (event.operation === 'create') {
			// Add new instance
			const newInstance = event.payload as Instance;
			instances = [...instances, newInstance];
		} else if (event.operation === 'update') {
			// Update existing instance
			const updatedInstance = event.payload as Instance;
			instances = instances.map(instance => 
				instance.id === updatedInstance.id ? updatedInstance : instance
			);
		} else if (event.operation === 'delete') {
			// Remove instance - payload might only contain ID
			const instanceId = event.payload.id || event.payload;
			instances = instances.filter(instance => instance.id !== instanceId);
		}
	}

	onMount(() => {
		// Initial load
		loadInstances();
		
		// Subscribe to real-time instance events
		unsubscribeWebsocket = websocketStore.subscribeToEntity(
			'instance',
			['create', 'update', 'delete'],
			handleInstanceEvent
		);
	});

	onDestroy(() => {
		// Clean up websocket subscription
		if (unsubscribeWebsocket) {
			unsubscribeWebsocket();
			unsubscribeWebsocket = null;
		}
	});
</script>

<svelte:head>
	<title>Instances - GARM</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Runner Instances</h1>
			<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
				Monitor your running instances
			</p>
		</div>
		<div class="flex items-center space-x-2">
		</div>
	</div>

	{#if error}
		<div class="bg-red-50 dark:bg-red-900/50 border border-red-200 dark:border-red-800 rounded-md p-4">
			<div class="flex">
				<div class="ml-3">
					<h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error</h3>
					<div class="mt-2 text-sm text-red-700 dark:text-red-300">{error}</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Search and Filter Controls -->
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
		<div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
			<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-3 sm:space-y-0">
				<div class="flex-1 max-w-lg">
					<label for="search" class="sr-only">Search instances</label>
					<div class="relative">
						<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
							<svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
							</svg>
						</div>
						<input
							id="search"
							bind:value={searchTerm}
							on:input={() => currentPage = 1}
							type="text"
							class="block w-full pl-10 pr-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md leading-5 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
							placeholder="Search instances..."
						/>
					</div>
				</div>
				<div class="flex items-center space-x-3">
					<select
						bind:value={statusFilter}
						on:change={handleStatusFilterChange}
						class="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
					>
						<option value="">All statuses</option>
						<option value="idle">Idle</option>
						<option value="active">Active</option>
						<option value="running">Running</option>
						<option value="pending">Pending</option>
						<option value="installing">Installing</option>
						<option value="failed">Failed</option>
						<option value="terminated">Terminated</option>
						<option value="offline">Offline</option>
						<option value="online">Online</option>
						<option value="unknown">Unknown</option>
					</select>
				</div>
			</div>
		</div>

		{#if loading}
			<div class="px-6 py-4 text-center">
				<div class="inline-flex items-center">
					<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-gray-500" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Loading instances...
				</div>
			</div>
		{:else if filteredInstances.length === 0}
			<div class="px-6 py-8 text-center text-gray-500 dark:text-gray-400">
				{searchTerm || statusFilter ? 'No instances match your search criteria.' : 'No instances found.'}
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider w-1/4">Name</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Pool/Scale Set</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider w-32">Created</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider w-24">Status</th>
							<th class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider w-28">Runner Status</th>
							<th class="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider w-20">Actions</th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each paginatedInstances as instance}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-4 py-4 whitespace-nowrap">
									<div class="text-sm font-medium text-gray-900 dark:text-white">
										<a href="{base}/instances/{encodeURIComponent(instance.name)}" class="hover:text-blue-600 dark:hover:text-blue-400">
											{instance.name}
										</a>
									</div>
									<div class="text-sm text-gray-500 dark:text-gray-400">
										{instance.provider_id}
									</div>
								</td>
								<td class="px-4 py-4 text-sm text-gray-900 dark:text-white">
									{#if instance.pool_id}
										<a href="{base}/pools/{instance.pool_id}" class="text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300 break-all">
											Pool: {instance.pool_id}
										</a>
									{:else if instance.scale_set_id}
										<a href="{base}/scalesets/{instance.scale_set_id}" class="text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300 break-all">
											Scale Set: {instance.scale_set_id}
										</a>
									{:else}
										<span class="text-gray-400 dark:text-gray-500">-</span>
									{/if}
								</td>
								<td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
									{formatDate(instance.created_at)}
								</td>
								<td class="px-4 py-4 whitespace-nowrap">
									<span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset {getStatusBadgeClass(instance.status)}">
										{instance.status}
									</span>
								</td>
								<td class="px-4 py-4 whitespace-nowrap">
									<span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset {getStatusBadgeClass(instance.runner_status)}">
										{instance.runner_status}
									</span>
								</td>
								<td class="px-4 py-4 whitespace-nowrap text-right text-sm font-medium">
									<button
										on:click={() => handleDelete(instance)}
										class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300"
										title="Delete instance"
									>
										<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
										</svg>
									</button>
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
						<div class="flex items-center space-x-2">
							<p class="text-sm text-gray-700 dark:text-gray-300">
								Show
							</p>
							<select
								bind:value={itemsPerPage}
								on:change={() => currentPage = 1}
								class="px-2 py-1 border border-gray-300 dark:border-gray-600 rounded text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
							>
								<option value={25}>25</option>
								<option value={50}>50</option>
								<option value={100}>100</option>
							</select>
							<p class="text-sm text-gray-700 dark:text-gray-300">
								per page. Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(currentPage * itemsPerPage, filteredInstances.length)} of {filteredInstances.length} results.
							</p>
						</div>
						<div>
							<nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
								<button
									on:click={() => currentPage = Math.max(1, currentPage - 1)}
									disabled={currentPage === 1}
									class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									Previous
								</button>
								{#each Array.from({length: Math.min(5, totalPages)}, (_, i) => i + Math.max(1, Math.min(currentPage - 2, totalPages - 4))) as page}
									<button
										on:click={() => currentPage = page}
										class="relative inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium {page === currentPage ? 'z-10 bg-blue-50 dark:bg-blue-900/50 border-blue-500 dark:border-blue-400 text-blue-600 dark:text-blue-400' : 'bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600'}"
									>
										{page}
									</button>
								{/each}
								<button
									on:click={() => currentPage = Math.min(totalPages, currentPage + 1)}
									disabled={currentPage === totalPages}
									class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-sm font-medium text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									Next
								</button>
							</nav>
						</div>
					</div>
				</div>
			{/if}
		{/if}
	</div>
</div>

<!-- Delete Modal -->
{#if showDeleteModal && instanceToDelete}
	<DeleteModal
		title="Delete Instance"
		message="Are you sure you want to delete this instance? This action cannot be undone."
		itemName={instanceToDelete.name}
		on:close={() => {
			showDeleteModal = false;
			instanceToDelete = null;
		}}
		on:confirm={confirmDelete}
	/>
{/if}