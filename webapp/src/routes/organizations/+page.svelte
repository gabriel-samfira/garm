<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { garmApi } from '$lib/api/client.js';
	import type { Organization, CreateOrgParams, UpdateEntityParams, Endpoint, Credential } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import CreateOrganizationModal from '$lib/components/CreateOrganizationModal.svelte';
	import UpdateEntityModal from '$lib/components/UpdateEntityModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';
	import { eagerCache, eagerCacheManager } from '$lib/stores/eager-cache.js';
	import { toastStore } from '$lib/stores/toast.js';

	let organizations: Organization[] = [];
	let loading = true;
	let error = '';

	// Subscribe to eager cache for organizations
	$: organizations = $eagerCache.organizations;
	$: loading = $eagerCache.loading.organizations;
	$: cacheError = $eagerCache.errorMessages.organizations;
	let searchTerm = '';
	let currentPage = 1;
	let itemsPerPage = 25;
	let showCreateModal = false;
	let showUpdateModal = false;
	let showDeleteModal = false;
	let selectedOrganization: Organization | null = null;
	// Filtered and paginated data
	$: filteredOrganizations = organizations.filter((org) =>
		org.name.toLowerCase().includes(searchTerm.toLowerCase())
	);
	$: totalPages = Math.ceil(filteredOrganizations.length / itemsPerPage);
	$: paginatedOrganizations = filteredOrganizations.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);


	async function handleCreateOrganization(params: CreateOrgParams) {
		try {
			error = '';
			await garmApi.createOrganization(params);
			// No need to reload - eager cache websocket will handle the update
			toastStore.success(
				'Organization Created',
				`Organization ${params.name} has been created successfully.`
			);
			showCreateModal = false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create organization';
			throw err; // Let the modal handle the error
		}
	}

	async function handleUpdateOrganization(params: UpdateEntityParams) {
		if (!selectedOrganization) return;
		try {
			await garmApi.updateOrganization(selectedOrganization.id, params);
			// No need to reload - eager cache websocket will handle the update
			toastStore.success(
				'Organization Updated',
				`Organization ${selectedOrganization.name} has been updated successfully.`
			);
			showUpdateModal = false;
			selectedOrganization = null;
		} catch (err) {
			throw err; // Let the modal handle the error
		}
	}

	async function handleDeleteOrganization() {
		if (!selectedOrganization) return;
		try {
			error = '';
			await garmApi.deleteOrganization(selectedOrganization.id);
			// No need to reload - eager cache websocket will handle the update
			toastStore.success(
				'Organization Deleted',
				`Organization ${selectedOrganization.name} has been deleted successfully.`
			);
			showDeleteModal = false;
			selectedOrganization = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete organization';
		}
	}

	function openCreateModal() {
		showCreateModal = true;
	}

	function openUpdateModal(organization: Organization) {
		selectedOrganization = organization;
		showUpdateModal = true;
	}

	function openDeleteModal(organization: Organization) {
		selectedOrganization = organization;
		showDeleteModal = true;
	}

	function getForgeIcon(endpointType: string) {
		if (endpointType === 'gitea') {
			return `<svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else {
			return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
		}
	}


	onMount(async () => {
		// Load organizations through eager cache (priority load + background load others)
		try {
			await eagerCacheManager.getOrganizations();
		} catch (err) {
			// Cache error is already handled by the eager cache system
			// We don't need to set error here anymore since it's in the cache state
			console.error('Failed to load organizations:', err);
		}
	});

	async function retryLoadOrganizations() {
		try {
			await eagerCacheManager.retryResource('organizations');
		} catch (err) {
			console.error('Retry failed:', err);
		}
	}

	// Organizations are now handled by eager cache with websocket subscriptions
</script>

<svelte:head>
	<title>Organizations - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="sm:flex sm:items-center sm:justify-between">
		<div>
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white">Organizations</h1>
			<p class="mt-2 text-sm text-gray-700 dark:text-gray-300">Manage GitHub and Gitea organizations</p>
		</div>
		<div class="mt-4 sm:mt-0 flex items-center space-x-4">
			<button
				on:click={openCreateModal}
				class="inline-flex items-center justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
			>
				<svg class="-ml-1 mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
				</svg>
				Add Organization
			</button>
		</div>
	</div>

	<!-- Search and filters -->
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg p-4">
		<div class="sm:flex sm:items-center sm:justify-between">
			<div class="flex-1 min-w-0">
				<div class="max-w-md">
					<label for="search" class="sr-only">Search organizations</label>
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
							placeholder="Search organizations..."
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

	<!-- Organizations table -->
	<div class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
		{#if loading}
			<div class="p-6 text-center">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
				<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Loading organizations...</p>
			</div>
		{:else if error || cacheError}
			<div class="p-6">
				<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
					<div class="flex">
						<div class="flex-shrink-0">
							<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
							</svg>
						</div>
						<div class="ml-3 flex-1">
							<h3 class="text-sm font-medium text-red-800 dark:text-red-200">Error loading organizations</h3>
							<p class="mt-2 text-sm text-red-700 dark:text-red-300">{cacheError || error}</p>
							{#if cacheError}
								<div class="mt-3">
									<button
										on:click={retryLoadOrganizations}
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
			</div>
		{:else if paginatedOrganizations.length === 0}
			<div class="p-6 text-center">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">No organizations found</h3>
				<p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
					{#if searchTerm}
						No organizations match your search criteria.
					{:else}
						Get started by creating your first organization.
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
						{#each paginatedOrganizations as organization}
							<tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div>
											<a href={`${base}/organizations/${organization.id}`} class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
												{organization.name}
											</a>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="flex items-center">
										<div class="flex-shrink-0">
											{@html getForgeIcon(organization.endpoint.endpoint_type)}
										</div>
										<div class="ml-2">
											<div class="text-sm text-gray-900 dark:text-white">{organization.endpoint.name}</div>
										</div>
									</div>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
									{organization.credentials_name}
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {organization.pool_manager_status?.running ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'}">
										{organization.pool_manager_status?.running ? 'Running' : 'Stopped'}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<div class="flex justify-end space-x-2">
										<button
											on:click={() => openUpdateModal(organization)}
											class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-900 dark:hover:text-indigo-300 cursor-pointer"
											title="Edit organization"
										>
											<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
											</svg>
										</button>
										<button
											on:click={() => openDeleteModal(organization)}
											class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300 cursor-pointer"
											title="Delete organization"
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
								to <span class="font-medium">{Math.min(currentPage * itemsPerPage, filteredOrganizations.length)}</span>
								of <span class="font-medium">{filteredOrganizations.length}</span> results
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
	<CreateOrganizationModal
		on:close={() => showCreateModal = false}
		on:submit={(e) => handleCreateOrganization(e.detail)}
	/>
{/if}

{#if showUpdateModal && selectedOrganization}
	<UpdateEntityModal
		entity={selectedOrganization}
		entityType="organization"
		on:close={() => { showUpdateModal = false; selectedOrganization = null; }}
		on:submit={(e) => handleUpdateOrganization(e.detail)}
	/>
{/if}

{#if showDeleteModal && selectedOrganization}
	<DeleteModal
		title="Delete Organization"
		message="Are you sure you want to delete this organization? This action cannot be undone."
		itemName={selectedOrganization.name}
		on:close={() => { showDeleteModal = false; selectedOrganization = null; }}
		on:confirm={handleDeleteOrganization}
	/>
{/if}