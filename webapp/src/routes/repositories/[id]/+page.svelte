<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { garmApi } from '$lib/api/client.js';
	import type { Repository, Pool, Instance } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import UpdateEntityModal from '$lib/components/UpdateEntityModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';
	import EntityInformation from '$lib/components/EntityInformation.svelte';
	import DetailHeader from '$lib/components/DetailHeader.svelte';
	import PoolsSection from '$lib/components/PoolsSection.svelte';
	import InstancesSection from '$lib/components/InstancesSection.svelte';
	import EventsSection from '$lib/components/EventsSection.svelte';
	import { websocketStore, type WebSocketEvent } from '$lib/stores/websocket.js';
	import { toastStore } from '$lib/stores/toast.js';

	let repository: Repository | null = null;
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

	$: repositoryId = $page.params.id;

	async function loadRepository() {
		if (!repositoryId) return;
		
		try {
			loading = true;
			error = '';
			
			const [repo, repoPools, repoInstances] = await Promise.all([
				garmApi.getRepository(repositoryId),
				garmApi.listRepositoryPools(repositoryId).catch(() => []),
				garmApi.listRepositoryInstances(repositoryId).catch(() => [])
			]);
			
			repository = repo;
			pools = repoPools;
			instances = repoInstances;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load repository';
		} finally {
			loading = false;
		}
	}

	function updateEntityFields(currentEntity: any, updatedFields: any): any {
		// Preserve only fields that are definitely not in the API response
		const { events: originalEvents } = currentEntity;
		
		// Use the API response as the primary source, add back preserved fields
		const result = {
			...updatedFields,
			events: originalEvents // Always preserve events since they're managed by websockets
		};
		
		return result;
	}

	async function handleUpdate(params: any) {
		if (!repository) return;
		try {
			// Update repository
			await garmApi.updateRepository(repository.id, params);
			
			// Reload fresh data to ensure UI is up to date
			await loadRepository();
			
			toastStore.success(
				'Repository Updated',
				`Repository ${repository.owner}/${repository.name} has been updated successfully.`
			);
			showUpdateModal = false;
		} catch (err) {
			throw err; // Let the modal handle the error
		}
	}

	async function handleDelete() {
		if (!repository) return;
		try {
			await garmApi.deleteRepository(repository.id);
			goto(`${base}/repositories`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete repository';
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

	function getForgeIcon(endpointType?: string) {
		if (endpointType === 'gitea') {
			return `<svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" xml:space="preserve" viewBox="0 0 640 640"><path d="m395.9 484.2-126.9-61c-12.5-6-17.9-21.2-11.8-33.8l61-126.9c6-12.5 21.2-17.9 33.8-11.8 17.2 8.3 27.1 13 27.1 13l-.1-109.2 16.7-.1.1 117.1s57.4 24.2 83.1 40.1c3.7 2.3 10.2 6.8 12.9 14.4 2.1 6.1 2 13.1-1 19.3l-61 126.9c-6.2 12.7-21.4 18.1-33.9 12" style="fill:#fff"/><path d="M622.7 149.8c-4.1-4.1-9.6-4-9.6-4s-117.2 6.6-177.9 8c-13.3.3-26.5.6-39.6.7v117.2c-5.5-2.6-11.1-5.3-16.6-7.9 0-36.4-.1-109.2-.1-109.2-29 .4-89.2-2.2-89.2-2.2s-141.4-7.1-156.8-8.5c-9.8-.6-22.5-2.1-39 1.5-8.7 1.8-33.5 7.4-53.8 26.9C-4.9 212.4 6.6 276.2 8 285.8c1.7 11.7 6.9 44.2 31.7 72.5 45.8 56.1 144.4 54.8 144.4 54.8s12.1 28.9 30.6 55.5c25 33.1 50.7 58.9 75.7 62 63 0 188.9-.1 188.9-.1s12 .1 28.3-10.3c14-8.5 26.5-23.4 26.5-23.4S547 483 565 451.5c5.5-9.7 10.1-19.1 14.1-28 0 0 55.2-117.1 55.2-231.1-1.1-34.5-9.6-40.6-11.6-42.6M125.6 353.9c-25.9-8.5-36.9-18.7-36.9-18.7S69.6 321.8 60 295.4c-16.5-44.2-1.4-71.2-1.4-71.2s8.4-22.5 38.5-30c13.8-3.7 31-3.1 31-3.1s7.1 59.4 15.7 94.2c7.2 29.2 24.8 77.7 24.8 77.7s-26.1-3.1-43-9.1m300.3 107.6s-6.1 14.5-19.6 15.4c-5.8.4-10.3-1.2-10.3-1.2s-.3-.1-5.3-2.1l-112.9-55s-10.9-5.7-12.8-15.6c-2.2-8.1 2.7-18.1 2.7-18.1L322 273s4.8-9.7 12.2-13c.6-.3 2.3-1 4.5-1.5 8.1-2.1 18 2.8 18 2.8L467.4 315s12.6 5.7 15.3 16.2c1.9 7.4-.5 14-1.8 17.2-6.3 15.4-55 113.1-55 113.1" style="fill:#609926"/><path d="M326.8 380.1c-8.2.1-15.4 5.8-17.3 13.8s2 16.3 9.1 20c7.7 4 17.5 1.8 22.7-5.4 5.1-7.1 4.3-16.9-1.8-23.1l24-49.1c1.5.1 3.7.2 6.2-.5 4.1-.9 7.1-3.6 7.1-3.6 4.2 1.8 8.6 3.8 13.2 6.1 4.8 2.4 9.3 4.9 13.4 7.3.9.5 1.8 1.1 2.8 1.9 1.6 1.3 3.4 3.1 4.7 5.5 1.9 5.5-1.9 14.9-1.9 14.9-2.3 7.6-18.4 40.6-18.4 40.6-8.1-.2-15.3 5-17.7 12.5-2.6 8.1 1.1 17.3 8.9 21.3s17.4 1.7 22.5-5.3c5-6.8 4.6-16.3-1.1-22.6 1.9-3.7 3.7-7.4 5.6-11.3 5-10.4 13.5-30.4 13.5-30.4.9-1.7 5.7-10.3 2.7-21.3-2.5-11.4-12.6-16.7-12.6-16.7-12.2-7.9-29.2-15.2-29.2-15.2s0-4.1-1.1-7.1c-1.1-3.1-2.8-5.1-3.9-6.3 4.7-9.7 9.4-19.3 14.1-29-4.1-2-8.1-4-12.2-6.1-4.8 9.8-9.7 19.7-14.5 29.5-6.7-.1-12.9 3.5-16.1 9.4-3.4 6.3-2.7 14.1 1.9 19.8z" style="fill:#609926"/></svg>`;
		} else {
			return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
		}
	}


	function scrollToBottomEvents() {
		if (eventsContainer) {
			eventsContainer.scrollTop = eventsContainer.scrollHeight;
		}
	}

	function handleRepositoryEvent(event: WebSocketEvent) {
		
		if (event.operation === 'update') {
			const updatedRepository = event.payload as Repository;
			// Only update if this is the repository we're viewing
			if (repository && updatedRepository.id === repository.id) {
				// Check if events have been updated
				const oldEventCount = repository.events?.length || 0;
				const newEventCount = updatedRepository.events?.length || 0;
				
				// Update repository using selective field updates
				repository = updateEntityFields(repository, updatedRepository);
				
				// Auto-scroll if new events were added
				if (newEventCount > oldEventCount) {
					// Use setTimeout to ensure the DOM has updated
					setTimeout(() => {
						scrollToBottomEvents();
					}, 100);
				}
			}
		} else if (event.operation === 'delete') {
			const deletedRepositoryId = event.payload.id || event.payload;
			// If this repository was deleted, redirect to repositories list
			if (repository && repository.id === deletedRepositoryId) {
				goto(`${base}/repositories`);
			}
		}
	}

	function handlePoolEvent(event: WebSocketEvent) {
		
		if (!repository) return;
		
		const pool = event.payload;
		// Only handle pools that belong to this repository
		if (pool.repo_id !== repository.id) return;

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
		
		if (!repository || !pools) return;
		
		const instance = event.payload;
		// Check if instance belongs to any pool that belongs to this repository
		const belongsToRepository = pools.some(pool => pool.id === instance.pool_id);
		if (!belongsToRepository) return;

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
		loadRepository().then(() => {
			// Scroll to bottom on initial load if there are events
			if (repository?.events?.length) {
				setTimeout(() => {
					scrollToBottomEvents();
				}, 100);
			}
		});
		
		// Subscribe to repository events
		const unsubscribeRepo = websocketStore.subscribeToEntity(
			'repository',
			['update', 'delete'],
			handleRepositoryEvent
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
			unsubscribeRepo();
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
	<title>{repository ? `${repository.name} - Repository Details` : 'Repository Details'} - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Breadcrumbs -->
	<nav class="flex" aria-label="Breadcrumb">
		<ol class="inline-flex items-center space-x-1 md:space-x-3">
			<li class="inline-flex items-center">
				<a href={`${base}/repositories`} class="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-600 dark:text-gray-400 dark:hover:text-white">
					<svg class="w-3 h-3 mr-2.5" fill="currentColor" viewBox="0 0 20 20">
						<path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
					</svg>
					Repositories
				</a>
			</li>
			<li>
				<div class="flex items-center">
					<svg class="w-3 h-3 text-gray-400 mx-1" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
					</svg>
					<span class="ml-1 text-sm font-medium text-gray-500 md:ml-2 dark:text-gray-400">
						{repository ? repository.name : 'Loading...'}
					</span>
				</div>
			</li>
		</ol>
	</nav>

	{#if loading}
		<div class="p-6 text-center">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
			<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Loading repository...</p>
		</div>
	{:else if error}
		<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
			<p class="text-sm font-medium text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if repository}
		<!-- Header -->
		<DetailHeader
			title={repository.name}
			subtitle="Owner: {repository.owner} • Endpoint: {repository.endpoint?.name}"
			forgeIcon={getForgeIcon(repository.endpoint?.endpoint_type)}
			onEdit={() => showUpdateModal = true}
			onDelete={() => showDeleteModal = true}
		/>

		<!-- Repository Details -->
		<EntityInformation entity={repository} entityType="repository" />

		<!-- Pools -->
		<PoolsSection {pools} entityType="repository" />

		<!-- Instances -->
		<InstancesSection {instances} entityType="repository" onDeleteInstance={openDeleteInstanceModal} />

		<!-- Events -->
		<EventsSection events={repository?.events} bind:eventsContainer />
	{/if}
</div>

<!-- Modals -->
{#if showUpdateModal && repository}
	<UpdateEntityModal
		entity={repository}
		entityType="repository"
		on:close={() => showUpdateModal = false}
		on:submit={(e) => handleUpdate(e.detail)}
	/>
{/if}

{#if showDeleteModal && repository}
	<DeleteModal
		title="Delete Repository"
		message="Are you sure you want to delete this repository? This action cannot be undone and will remove all associated pools and instances."
		itemName={`${repository.owner}/${repository.name}`}
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