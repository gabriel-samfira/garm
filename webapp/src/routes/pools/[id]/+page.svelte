<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { garmApi } from '$lib/api/client.js';
	import type { Pool, UpdatePoolParams } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import UpdatePoolModal from '$lib/components/UpdatePoolModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';
	import InstancesSection from '$lib/components/InstancesSection.svelte';
	import DetailHeader from '$lib/components/DetailHeader.svelte';
	import { websocketStore, type WebSocketEvent } from '$lib/stores/websocket.js';
	import type { Instance } from '$lib/api/types.js';
	import { toastStore } from '$lib/stores/toast.js';

	let pool: Pool | null = null;
	let loading = true;
	let error = '';
	let showUpdateModal = false;
	let showDeleteModal = false;
	let showDeleteInstanceModal = false;
	let selectedInstance: Instance | null = null;
	let unsubscribeWebsocket: (() => void) | null = null;

	$: poolId = $page.params.id;

	async function loadPool() {
		if (!poolId) return;
		
		try {
			loading = true;
			error = '';
			pool = await garmApi.getPool(poolId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load pool';
		} finally {
			loading = false;
		}
	}

	async function handleUpdate(params: UpdatePoolParams) {
		if (!pool) return;
		try {
			// Update pool and get the updated object from API response
			const updatedPool = await garmApi.updatePool(pool.id, params);
			// Update the local state directly instead of re-rendering
			pool = updatedPool;
			showUpdateModal = false;
			toastStore.success(
				'Pool Updated',
				`Pool ${pool.id} has been updated successfully.`
			);
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to update pool';
			toastStore.error(
				'Update Failed',
				errorMessage
			);
		}
	}

	async function handleDelete() {
		if (!pool) return;
		try {
			await garmApi.deletePool(pool.id);
			goto(`${base}/pools`);
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to delete pool';
			toastStore.error(
				'Delete Failed',
				errorMessage
			);
		}
		showDeleteModal = false;
	}

	async function handleDeleteInstance() {
		if (!selectedInstance) return;
		try {
			await garmApi.deleteInstance(selectedInstance.name);
			toastStore.success(
				'Instance Deleted',
				`Instance ${selectedInstance.name} has been deleted successfully.`
			);
			showDeleteInstanceModal = false;
			selectedInstance = null;
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to delete instance';
			toastStore.error(
				'Delete Failed',
				errorMessage
			);
		}
		showDeleteInstanceModal = false;
		selectedInstance = null;
	}

	function openDeleteInstanceModal(instance: Instance) {
		selectedInstance = instance;
		showDeleteInstanceModal = true;
	}

	function getEntityName(pool: Pool): string {
		if (pool.repo_name) return pool.repo_name;
		if (pool.org_name) return pool.org_name;
		if (pool.enterprise_name) return pool.enterprise_name;
		return 'Unknown Entity';
	}

	function getEntityType(pool: Pool): string {
		if (pool.repo_id) return 'repository';
		if (pool.org_id) return 'organization';
		if (pool.enterprise_id) return 'enterprise';
		return 'unknown';
	}

	function getEntityUrl(pool: Pool): string {
		if (pool.repo_id) return `${base}/repositories/${pool.repo_id}`;
		if (pool.org_id) return `${base}/organizations/${pool.org_id}`;
		if (pool.enterprise_id) return `${base}/enterprises/${pool.enterprise_id}`;
		return '#';
	}

	function getForgeIcon(endpointType?: string) {
		if (endpointType === 'gitea') {
			return `<svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else {
			return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}

	function formatExtraSpecs(extraSpecs: any): string {
		if (!extraSpecs) return '{}';
		try {
			if (typeof extraSpecs === 'string') {
				const parsed = JSON.parse(extraSpecs);
				return JSON.stringify(parsed, null, 2);
			}
			return JSON.stringify(extraSpecs, null, 2);
		} catch (e) {
			return extraSpecs.toString();
		}
	}

	function handlePoolEvent(event: WebSocketEvent) {
		
		if (event.operation === 'update') {
			const updatedPool = event.payload as Pool;
			// Only update if this is the pool we're viewing
			if (pool && updatedPool.id === pool.id) {
				pool = updatedPool;
			}
		} else if (event.operation === 'delete') {
			const deletedPoolId = event.payload.id || event.payload;
			// If this pool was deleted, redirect to pools list
			if (pool && pool.id === deletedPoolId) {
				goto(`${base}/pools`);
			}
		}
	}

	function handleInstanceEvent(event: WebSocketEvent) {
		
		if (!pool || !pool.instances) return;
		
		const instance = event.payload;
		// Only handle instances that belong to this pool
		if (instance.pool_id !== pool.id) return;

		if (event.operation === 'create') {
			// Add new instance to the list
			pool.instances = [...pool.instances, instance];
		} else if (event.operation === 'update') {
			// Update existing instance
			pool.instances = pool.instances.map(inst => 
				inst.id === instance.id ? instance : inst
			);
		} else if (event.operation === 'delete') {
			// Remove deleted instance
			const instanceId = instance.id || instance;
			pool.instances = pool.instances.filter(inst => inst.id !== instanceId);
		}
		
		// Force reactivity
		pool = pool;
	}

	onMount(() => {
		loadPool();
		
		// Subscribe to pool events
		const unsubscribePool = websocketStore.subscribeToEntity(
			'pool',
			['update', 'delete'],
			handlePoolEvent
		);

		// Subscribe to instance events
		const unsubscribeInstance = websocketStore.subscribeToEntity(
			'instance',
			['create', 'update', 'delete'],
			handleInstanceEvent
		);

		// Combine unsubscribe functions
		unsubscribeWebsocket = () => {
			unsubscribePool();
			unsubscribeInstance();
		};
	});

	onDestroy(() => {
		if (unsubscribeWebsocket) {
			unsubscribeWebsocket();
			unsubscribeWebsocket = null;
		}
	});
</script>

<svelte:head>
	<title>{pool ? `Pool ${pool.id} - Pool Details` : 'Pool Details'} - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Breadcrumbs -->
	<nav class="flex" aria-label="Breadcrumb">
		<ol class="inline-flex items-center space-x-1 md:space-x-3">
			<li class="inline-flex items-center">
				<a href={`${base}/pools`} class="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-600 dark:text-gray-400 dark:hover:text-white">
					<svg class="w-3 h-3 mr-2.5" fill="currentColor" viewBox="0 0 20 20">
						<path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
					</svg>
					Pools
				</a>
			</li>
			<li>
				<div class="flex items-center">
					<svg class="w-3 h-3 text-gray-400 mx-1" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
					</svg>
					<span class="ml-1 text-sm font-medium text-gray-500 md:ml-2 dark:text-gray-400">
						{pool ? pool.id : 'Loading...'}
					</span>
				</div>
			</li>
		</ol>
	</nav>

	{#if loading}
		<div class="p-6 text-center">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
			<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Loading pool...</p>
		</div>
	{:else if error}
		<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
			<p class="text-sm font-medium text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if pool}
		<!-- Header -->
		<DetailHeader
			title={pool.id}
			subtitle="Pool for {getEntityName(pool)} ({getEntityType(pool)})"
			forgeIcon={getForgeIcon(pool.endpoint?.endpoint_type)}
			onEdit={() => showUpdateModal = true}
			onDelete={() => showDeleteModal = true}
		/>

		<!-- Pool Details -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Basic Information -->
			<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Basic Information</h2>
					<dl class="space-y-4">
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Pool ID</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white font-mono">{pool.id}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Provider</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.provider_name}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Image</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">
								<code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-xs">{pool.image}</code>
							</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Flavor</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.flavor}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Status</dt>
							<dd class="mt-1">
								<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {pool.enabled ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'}">
									{pool.enabled ? 'Enabled' : 'Disabled'}
								</span>
							</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Entity</dt>
							<dd class="mt-1">
								<div class="flex items-center space-x-2">
									<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
										{getEntityType(pool)}
									</span>
									<a href={getEntityUrl(pool)} class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
										{getEntityName(pool)}
									</a>
								</div>
							</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Created At</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{formatDate(pool.created_at)}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Updated At</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{formatDate(pool.updated_at)}</dd>
						</div>
					</dl>
				</div>
			</div>

			<!-- Configuration -->
			<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Configuration</h2>
					<dl class="space-y-4">
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Max Runners</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.max_runners}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Min Idle Runners</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.min_idle_runners}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Bootstrap Timeout</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.runner_bootstrap_timeout} minutes</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Priority</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.priority}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Runner Prefix</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.runner_prefix || 'garm'}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">OS Type / Architecture</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.os_type} / {pool.os_arch}</dd>
						</div>
						{#if pool.github_runner_group}
							<div>
								<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">GitHub Runner Group</dt>
								<dd class="mt-1 text-sm text-gray-900 dark:text-white">{pool.github_runner_group}</dd>
							</div>
						{/if}
						{#if pool.tags && pool.tags.length > 0}
							<div>
								<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Tags</dt>
								<dd class="mt-1">
									<div class="flex flex-wrap gap-2">
										{#each pool.tags as tag}
											<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200">
												{typeof tag === 'string' ? tag : tag.name}
											</span>
										{/each}
									</div>
								</dd>
							</div>
						{/if}
					</dl>
				</div>
			</div>
		</div>

		<!-- Extra Specs -->
		{#if pool.extra_specs}
			<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Extra Specifications</h2>
					<pre class="bg-gray-100 dark:bg-gray-700 p-4 rounded-md overflow-x-auto text-sm text-gray-900 dark:text-white font-mono">{formatExtraSpecs(pool.extra_specs)}</pre>
				</div>
			</div>
		{/if}

		<!-- Instances -->
		{#if pool.instances}
			<InstancesSection instances={pool.instances} entityType="pool" onDeleteInstance={openDeleteInstanceModal} />
		{/if}

	{/if}
</div>

<!-- Modals -->
{#if showUpdateModal && pool}
	<UpdatePoolModal
		{pool}
		on:close={() => showUpdateModal = false}
		on:submit={(e) => handleUpdate(e.detail)}
	/>
{/if}

{#if showDeleteModal && pool}
	<DeleteModal
		title="Delete Pool"
		message="Are you sure you want to delete this pool? This action cannot be undone and will remove all associated runners."
		itemName={`Pool ${pool.id} (${getEntityName(pool)})`}
		on:close={() => showDeleteModal = false}
		on:confirm={handleDelete}
	/>
{/if}

{#if showDeleteInstanceModal && selectedInstance}
	<DeleteModal
		title="Delete Instance"
		message="Are you sure you want to delete this instance? This action cannot be undone."
		itemName={selectedInstance.name}
		on:close={() => { showDeleteInstanceModal = false; selectedInstance = null; }}
		on:confirm={handleDeleteInstance}
	/>
{/if}