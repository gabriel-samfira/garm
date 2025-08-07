<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { garmApi } from '$lib/api/client.js';
	import type { Enterprise, Pool, Instance } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import UpdateEnterpriseModal from '$lib/components/UpdateEnterpriseModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';

	let enterprise: Enterprise | null = null;
	let pools: Pool[] = [];
	let instances: Instance[] = [];
	let loading = true;
	let error = '';
	let showUpdateModal = false;
	let showDeleteModal = false;

	$: enterpriseId = $page.params.id;

	async function loadEnterprise() {
		if (!enterpriseId) return;
		
		try {
			loading = true;
			error = '';
			
			const [ent, entPools, entInstances] = await Promise.all([
				garmApi.getEnterprise(enterpriseId),
				garmApi.listEnterprisePools(enterpriseId).catch(() => []),
				garmApi.listEnterpriseInstances(enterpriseId).catch(() => [])
			]);
			
			enterprise = ent;
			pools = entPools;
			instances = entInstances;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load enterprise';
		} finally {
			loading = false;
		}
	}

	async function handleUpdate(params: any) {
		if (!enterprise) return;
		try {
			await garmApi.updateEnterprise(enterprise.id, params);
			await loadEnterprise();
			showUpdateModal = false;
		} catch (err) {
			throw err; // Let the modal handle the error
		}
	}

	async function handleDelete() {
		if (!enterprise) return;
		try {
			await garmApi.deleteEnterprise(enterprise.id);
			goto(`${base}/enterprises`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete enterprise';
		}
		showDeleteModal = false;
	}

	// Enterprises are GitHub only
	function getForgeIcon() {
		return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}

	onMount(() => {
		loadEnterprise();
	});
</script>

<svelte:head>
	<title>{enterprise ? `${enterprise.name} - Enterprise Details` : 'Enterprise Details'} - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Breadcrumbs -->
	<nav class="flex" aria-label="Breadcrumb">
		<ol class="inline-flex items-center space-x-1 md:space-x-3">
			<li class="inline-flex items-center">
				<a href={`${base}/enterprises`} class="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-600 dark:text-gray-400 dark:hover:text-white">
					<svg class="w-3 h-3 mr-2.5" fill="currentColor" viewBox="0 0 20 20">
						<path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
					</svg>
					Enterprises
				</a>
			</li>
			<li>
				<div class="flex items-center">
					<svg class="w-3 h-3 text-gray-400 mx-1" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
					</svg>
					<span class="ml-1 text-sm font-medium text-gray-500 md:ml-2 dark:text-gray-400">
						{enterprise ? enterprise.name : 'Loading...'}
					</span>
				</div>
			</li>
		</ol>
	</nav>

	{#if loading}
		<div class="p-6 text-center">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
			<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Loading enterprise...</p>
		</div>
	{:else if error}
		<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
			<p class="text-sm font-medium text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if enterprise}
		<!-- Header -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="px-4 py-5 sm:p-6">
				<div class="sm:flex sm:items-center sm:justify-between">
					<div class="flex items-center space-x-3">
						<div class="flex-shrink-0">
							{@html getForgeIcon()}
						</div>
						<div>
							<h1 class="text-2xl font-bold text-gray-900 dark:text-white">{enterprise.name}</h1>
							<p class="text-sm text-gray-500 dark:text-gray-400">
								Endpoint: {enterprise.endpoint?.name} • GitHub Enterprise
							</p>
						</div>
					</div>
					<div class="mt-4 sm:mt-0 flex space-x-3">
						<button
							on:click={() => showUpdateModal = true}
							class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
						>
							<svg class="-ml-1 mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
							</svg>
							Edit
						</button>
						<button
							on:click={() => showDeleteModal = true}
							class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
						>
							<svg class="-ml-1 mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
							</svg>
							Delete
						</button>
					</div>
				</div>
			</div>
		</div>

		<!-- Enterprise Details -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="px-4 py-5 sm:p-6">
				<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Enterprise Information</h2>
				<dl class="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">ID</dt>
						<dd class="mt-1 text-sm text-gray-900 dark:text-white font-mono">{enterprise.id}</dd>
					</div>
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Created At</dt>
						<dd class="mt-1 text-sm text-gray-900 dark:text-white">{formatDate(enterprise.created_at)}</dd>
					</div>
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Updated At</dt>
						<dd class="mt-1 text-sm text-gray-900 dark:text-white">{formatDate(enterprise.updated_at)}</dd>
					</div>
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Status</dt>
						<dd class="mt-1">
							<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {enterprise.pool_manager_status?.is_running ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'}">
								{enterprise.pool_manager_status?.is_running ? 'Running' : 'Stopped'}
							</span>
						</dd>
					</div>
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Pool Balancer Type</dt>
						<dd class="mt-1 text-sm text-gray-900 dark:text-white">{enterprise.pool_balancer_type || 'None'}</dd>
					</div>
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Enterprise URL</dt>
						<dd class="mt-1 text-sm">
							<a href="{enterprise.endpoint?.base_url}/enterprises/{enterprise.name}" target="_blank" rel="noopener noreferrer" class="text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
								{enterprise.endpoint?.base_url}/enterprises/{enterprise.name}
								<svg class="inline w-3 h-3 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
								</svg>
							</a>
						</dd>
					</div>
				</dl>
			</div>
		</div>

		<!-- Pools -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="px-4 py-5 sm:p-6">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white">Pools ({pools.length})</h2>
					<a href={`${base}/pools`} class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">View all pools</a>
				</div>
				{#if pools.length === 0}
					<p class="text-sm text-gray-500 dark:text-gray-400">No pools configured for this enterprise.</p>
				{:else}
					<div class="overflow-x-auto">
						<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
							<thead class="bg-gray-50 dark:bg-gray-700">
								<tr>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">ID</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Provider</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Image</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
								</tr>
							</thead>
							<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
								{#each pools as pool}
									<tr>
										<td class="px-6 py-4 whitespace-nowrap">
											<a href={`${base}/pools/${pool.id}`} class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300 font-mono">
												{pool.id.slice(0, 8)}...
											</a>
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">{pool.provider_name}</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
											<code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-xs">{pool.image}</code>
										</td>
										<td class="px-6 py-4 whitespace-nowrap">
											<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {pool.enabled ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'}">
												{pool.enabled ? 'Enabled' : 'Disabled'}
											</span>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</div>
		</div>

		<!-- Instances -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="px-4 py-5 sm:p-6">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white">Instances ({instances.length})</h2>
					<a href={`${base}/instances`} class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">View all instances</a>
				</div>
				{#if instances.length === 0}
					<p class="text-sm text-gray-500 dark:text-gray-400">No instances running for this enterprise.</p>
				{:else}
					<div class="overflow-x-auto">
						<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
							<thead class="bg-gray-50 dark:bg-gray-700">
								<tr>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
									<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Created</th>
								</tr>
							</thead>
							<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
								{#each instances as instance}
									<tr>
										<td class="px-6 py-4 whitespace-nowrap">
											<a href={`${base}/instances/${instance.name}`} class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
												{instance.name}
											</a>
										</td>
										<td class="px-6 py-4 whitespace-nowrap">
											<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200">
												{instance.status}
											</span>
										</td>
										<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
											{formatDate(instance.created_at)}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<!-- Modals -->
{#if showUpdateModal && enterprise}
	<UpdateEnterpriseModal
		enterprise={enterprise}
		on:close={() => showUpdateModal = false}
		on:submit={(e) => handleUpdate(e.detail)}
	/>
{/if}

{#if showDeleteModal && enterprise}
	<DeleteModal
		title="Delete Enterprise"
		message="Are you sure you want to delete this enterprise? This action cannot be undone and will remove all associated pools and instances."
		itemName={enterprise.name}
		on:close={() => showDeleteModal = false}
		on:confirm={handleDelete}
	/>
{/if}