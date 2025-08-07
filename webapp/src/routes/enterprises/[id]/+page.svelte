<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { garmApi } from '$lib/api/client.js';
	import type { Enterprise, Pool, Instance } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import UpdateEntityModal from '$lib/components/UpdateEntityModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';
	import EntityInformation from '$lib/components/EntityInformation.svelte';
	import PoolsSection from '$lib/components/PoolsSection.svelte';
	import InstancesSection from '$lib/components/InstancesSection.svelte';
	import EventsSection from '$lib/components/EventsSection.svelte';
	import { websocketStore, type WebSocketEvent } from '$lib/stores/websocket.js';
	import { toastStore } from '$lib/stores/toast.js';

	let enterprise: Enterprise | null = null;
	let pools: Pool[] = [];
	let instances: Instance[] = [];
	let loading = true;
	let error = '';
	let showUpdateModal = false;
	let showDeleteModal = false;
	let showDeleteInstanceModal = false;
	let selectedInstance: Instance | null = null;
	let unsubscribeWebsocket: (() => void) | null = null;
	let eventsContainer: HTMLElement;

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
			toastStore.success(
				'Enterprise Updated',
				`Enterprise ${enterprise.name} has been updated successfully.`
			);
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

	async function handleDeleteInstance() {
		if (!selectedInstance) return;
		try {
			await garmApi.deleteInstance(selectedInstance.name);
			toastStore.success(
				'Instance Deleted',
				`Instance ${selectedInstance.name} has been deleted successfully.`
			);
			// No need to reload - websocket events will update the UI automatically
			showDeleteInstanceModal = false;
			selectedInstance = null;
		} catch (err) {
			const errorMessage = err instanceof Error ? err.message : 'Failed to delete instance';
			toastStore.error(
				'Delete Failed',
				errorMessage
			);
			showDeleteInstanceModal = false;
			selectedInstance = null;
		}
	}

	function openDeleteInstanceModal(instance: Instance) {
		selectedInstance = instance;
		showDeleteInstanceModal = true;
	}

	// Enterprises are GitHub only
	function getForgeIcon() {
		return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
	}


	function scrollToBottomEvents() {
		if (eventsContainer) {
			eventsContainer.scrollTop = eventsContainer.scrollHeight;
		}
	}

	function handleEnterpriseEvent(event: WebSocketEvent) {
		console.log('[Enterprise Detail] Received websocket event:', event);
		
		if (event.operation === 'update') {
			const updatedEnterprise = event.payload as Enterprise;
			// Only update if this is the enterprise we're viewing
			if (enterprise && updatedEnterprise.id === enterprise.id) {
				// Check if events have been updated
				const oldEventCount = enterprise.events?.length || 0;
				const newEventCount = updatedEnterprise.events?.length || 0;
				
				// Update enterprise
				enterprise = updatedEnterprise;
				
				// Auto-scroll if new events were added
				if (newEventCount > oldEventCount) {
					// Use setTimeout to ensure the DOM has updated
					setTimeout(() => {
						scrollToBottomEvents();
					}, 100);
				}
			}
		} else if (event.operation === 'delete') {
			const deletedEnterpriseId = event.payload.id || event.payload;
			// If this enterprise was deleted, redirect to enterprises list
			if (enterprise && enterprise.id === deletedEnterpriseId) {
				goto(`${base}/enterprises`);
			}
		}
	}

	function handlePoolEvent(event: WebSocketEvent) {
		console.log('[Enterprise Detail] Received pool websocket event:', event);
		
		if (!enterprise) return;
		
		const pool = event.payload;
		// Only handle pools that belong to this enterprise
		if (pool.enterprise_id !== enterprise.id) return;

		if (event.operation === 'create') {
			// Add new pool to the list
			pools = [...pools, pool];
		} else if (event.operation === 'update') {
			// Update existing pool
			pools = pools.map(p => 
				p.id === pool.id ? pool : p
			);
		} else if (event.operation === 'delete') {
			// Remove deleted pool
			const poolId = pool.id || pool;
			pools = pools.filter(p => p.id !== poolId);
		}
	}

	function handleInstanceEvent(event: WebSocketEvent) {
		console.log('[Enterprise Detail] Received instance websocket event:', event);
		
		if (!enterprise || !pools) return;
		
		const instance = event.payload;
		// Check if instance belongs to any pool that belongs to this enterprise
		const belongsToEnterprise = pools.some(pool => pool.id === instance.pool_id);
		if (!belongsToEnterprise) return;

		if (event.operation === 'create') {
			// Add new instance to the list
			instances = [...instances, instance];
		} else if (event.operation === 'update') {
			// Update existing instance
			instances = instances.map(inst => 
				inst.id === instance.id ? instance : inst
			);
		} else if (event.operation === 'delete') {
			// Remove deleted instance
			const instanceId = instance.id || instance;
			instances = instances.filter(inst => inst.id !== instanceId);
		}
	}

	onMount(() => {
		loadEnterprise().then(() => {
			// Scroll to bottom on initial load if there are events
			if (enterprise?.events?.length) {
				setTimeout(() => {
					scrollToBottomEvents();
				}, 100);
			}
		});
		
		// Subscribe to enterprise events
		const unsubscribeEnt = websocketStore.subscribeToEntity(
			'enterprise',
			['update', 'delete'],
			handleEnterpriseEvent
		);

		// Subscribe to pool events
		const unsubscribePool = websocketStore.subscribeToEntity(
			'pool',
			['create', 'update', 'delete'],
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
			unsubscribeEnt();
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
		<EntityInformation entity={enterprise} entityType="enterprise" />

		<!-- Pools -->
		<PoolsSection {pools} entityType="enterprise" />

		<!-- Instances -->
		<InstancesSection {instances} entityType="enterprise" onDeleteInstance={openDeleteInstanceModal} />

		<!-- Events -->
		<EventsSection events={enterprise?.events} bind:eventsContainer />
	{/if}
</div>

<!-- Modals -->
{#if showUpdateModal && enterprise}
	<UpdateEntityModal
		entity={enterprise}
		entityType="enterprise"
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

{#if showDeleteInstanceModal && selectedInstance}
	<DeleteModal
		title="Delete Instance"
		message="Are you sure you want to delete this instance? This action cannot be undone."
		itemName={selectedInstance.name}
		on:close={() => { showDeleteInstanceModal = false; selectedInstance = null; }}
		on:confirm={handleDeleteInstance}
	/>
{/if}