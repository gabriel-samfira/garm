<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { garmApi } from '$lib/api/client.js';
	import type { ScaleSet, CreateScaleSetParams } from '$lib/api/types.js';
	import { base } from '$app/paths';
	import UpdateScaleSetModal from '$lib/components/UpdateScaleSetModal.svelte';
	import DeleteModal from '$lib/components/DeleteModal.svelte';

	let scaleSet: ScaleSet | null = null;
	let loading = true;
	let error = '';
	let showUpdateModal = false;
	let showDeleteModal = false;

	$: scaleSetId = parseInt($page.params.id);

	async function loadScaleSet() {
		if (!scaleSetId || isNaN(scaleSetId)) return;
		
		try {
			loading = true;
			error = '';
			scaleSet = await garmApi.getScaleSet(scaleSetId);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load scale set';
		} finally {
			loading = false;
		}
	}

	async function handleUpdate(params: Partial<CreateScaleSetParams>) {
		if (!scaleSet) return;
		try {
			await garmApi.updateScaleSet(scaleSet.id, params);
			await loadScaleSet();
			showUpdateModal = false;
		} catch (err) {
			throw err; // Let the modal handle the error
		}
	}

	async function handleDelete() {
		if (!scaleSet) return;
		try {
			await garmApi.deleteScaleSet(scaleSet.id);
			goto(`${base}/scalesets`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to delete scale set';
		}
		showDeleteModal = false;
	}

	function getEntityName(scaleSet: ScaleSet): string {
		if (scaleSet.repo_name) return scaleSet.repo_name;
		if (scaleSet.org_name) return scaleSet.org_name;
		if (scaleSet.enterprise_name) return scaleSet.enterprise_name;
		return 'Unknown Entity';
	}

	function getEntityType(scaleSet: ScaleSet): string {
		if (scaleSet.repo_id) return 'repository';
		if (scaleSet.org_id) return 'organization';  
		if (scaleSet.enterprise_id) return 'enterprise';
		return 'unknown';
	}

	function getEntityUrl(scaleSet: ScaleSet): string {
		if (scaleSet.repo_id) return `${base}/repositories/${scaleSet.repo_id}`;
		if (scaleSet.org_id) return `${base}/organizations/${scaleSet.org_id}`;
		if (scaleSet.enterprise_id) return `${base}/enterprises/${scaleSet.enterprise_id}`;
		return '#';
	}

	// Scale sets are GitHub only
	function getForgeIcon() {
		return `<div class="inline-flex w-4 h-4"><svg class="w-4 h-4 dark:hidden" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#24292f"/></svg><svg class="w-4 h-4 hidden dark:block" width="98" height="96" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 98 96"><path fill-rule="evenodd" clip-rule="evenodd" d="M48.854 0C21.839 0 0 22 0 49.217c0 21.756 13.993 40.172 33.405 46.69 2.427.49 3.316-1.059 3.316-2.362 0-1.141-.08-5.052-.08-9.127-13.59 2.934-16.42-5.867-16.42-5.867-2.184-5.704-5.42-7.17-5.42-7.17-4.448-3.015.324-3.015.324-3.015 4.934.326 7.523 5.052 7.523 5.052 4.367 7.496 11.404 5.378 14.235 4.074.404-3.178 1.699-5.378 3.074-6.6-10.839-1.141-22.243-5.378-22.243-24.283 0-5.378 1.94-9.778 5.014-13.2-.485-1.222-2.184-6.275.486-13.038 0 0 4.125-1.304 13.426 5.052a46.97 46.97 0 0 1 12.214-1.63c4.125 0 8.33.571 12.213 1.63 9.302-6.356 13.427-5.052 13.427-5.052 2.67 6.763.97 11.816.485 13.038 3.155 3.422 5.015 7.822 5.015 13.2 0 18.905-11.404 23.06-22.324 24.283 1.78 1.548 3.316 4.481 3.316 9.126 0 6.6-.08 11.897-.08 13.526 0 1.304.89 2.853 3.316 2.364 19.412-6.52 33.405-24.935 33.405-46.691C97.707 22 75.788 0 48.854 0z" fill="#fff"/></svg></div>`;
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

	onMount(() => {
		loadScaleSet();
	});
</script>

<svelte:head>
	<title>{scaleSet ? `${scaleSet.name} - Scale Set Details` : 'Scale Set Details'} - GARM</title>
</svelte:head>

<div class="space-y-6">
	<!-- Breadcrumbs -->
	<nav class="flex" aria-label="Breadcrumb">
		<ol class="inline-flex items-center space-x-1 md:space-x-3">
			<li class="inline-flex items-center">
				<a href={`${base}/scalesets`} class="inline-flex items-center text-sm font-medium text-gray-700 hover:text-blue-600 dark:text-gray-400 dark:hover:text-white">
					<svg class="w-3 h-3 mr-2.5" fill="currentColor" viewBox="0 0 20 20">
						<path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
					</svg>
					Scale Sets
				</a>
			</li>
			<li>
				<div class="flex items-center">
					<svg class="w-3 h-3 text-gray-400 mx-1" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
					</svg>
					<span class="ml-1 text-sm font-medium text-gray-500 md:ml-2 dark:text-gray-400">
						{scaleSet ? scaleSet.name : 'Loading...'}
					</span>
				</div>
			</li>
		</ol>
	</nav>

	{#if loading}
		<div class="p-6 text-center">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
			<p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Loading scale set...</p>
		</div>
	{:else if error}
		<div class="rounded-md bg-red-50 dark:bg-red-900 p-4">
			<p class="text-sm font-medium text-red-800 dark:text-red-200">{error}</p>
		</div>
	{:else if scaleSet}
		<!-- Header -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="px-4 py-5 sm:p-6">
				<div class="sm:flex sm:items-center sm:justify-between">
					<div class="flex items-center space-x-3">
						<div class="flex-shrink-0">
							{@html getForgeIcon()}
						</div>
						<div>
							<h1 class="text-2xl font-bold text-gray-900 dark:text-white">{scaleSet.name}</h1>
							<p class="text-sm text-gray-500 dark:text-gray-400">
								Scale set for {getEntityName(scaleSet)} ({getEntityType(scaleSet)}) • GitHub Runner Scale Set
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

		<!-- Scale Set Details -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Basic Information -->
			<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Basic Information</h2>
					<dl class="space-y-4">
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Scale Set ID</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.id}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Name</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white font-medium">{scaleSet.name}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Provider</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.provider_name}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Image</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">
								<code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-xs">{scaleSet.image}</code>
							</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Flavor</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.flavor}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Status</dt>
							<dd class="mt-1">
								<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {scaleSet.enabled ? 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200' : 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200'}">
									{scaleSet.enabled ? 'Enabled' : 'Disabled'}
								</span>
							</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Entity</dt>
							<dd class="mt-1">
								<div class="flex items-center space-x-2">
									<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200">
										{getEntityType(scaleSet)}
									</span>
									<a href={getEntityUrl(scaleSet)} class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">
										{getEntityName(scaleSet)}
									</a>
								</div>
							</dd>
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
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.max_runners}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Min Idle Runners</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.min_idle_runners}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Bootstrap Timeout</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.runner_bootstrap_timeout} minutes</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Priority</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.priority}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Runner Prefix</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.runner_prefix || 'garm'}</dd>
						</div>
						<div>
							<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">OS Type / Architecture</dt>
							<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.os_type} / {scaleSet.os_arch}</dd>
						</div>
						{#if scaleSet.github_runner_group}
							<div>
								<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">GitHub Runner Group</dt>
								<dd class="mt-1 text-sm text-gray-900 dark:text-white">{scaleSet.github_runner_group}</dd>
							</div>
						{/if}
					</dl>
				</div>
			</div>
		</div>

		<!-- Tags -->
		{#if scaleSet.tags && scaleSet.tags.length > 0}
			<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Tags</h2>
					<div class="flex flex-wrap gap-2">
						{#each scaleSet.tags as tag}
							<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200">
								{typeof tag === 'string' ? tag : tag.name}
							</span>
						{/each}
					</div>
				</div>
			</div>
		{/if}

		<!-- Extra Specs -->
		{#if scaleSet.extra_specs}
			<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
				<div class="px-4 py-5 sm:p-6">
					<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Extra Specifications</h2>
					<pre class="bg-gray-100 dark:bg-gray-700 p-4 rounded-md overflow-x-auto text-sm text-gray-900 dark:text-white font-mono">{formatExtraSpecs(scaleSet.extra_specs)}</pre>
				</div>
			</div>
		{/if}

		<!-- Timestamps -->
		<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
			<div class="px-4 py-5 sm:p-6">
				<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Timestamps</h2>
				<dl class="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Created At</dt>
						<dd class="mt-1 text-sm text-gray-900 dark:text-white">{formatDate(scaleSet.created_at)}</dd>
					</div>
					<div>
						<dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Updated At</dt>
						<dd class="mt-1 text-sm text-gray-900 dark:text-white">{formatDate(scaleSet.updated_at)}</dd>
					</div>
				</dl>
			</div>
		</div>
	{/if}
</div>

<!-- Modals -->
{#if showUpdateModal && scaleSet}
	<UpdateScaleSetModal
		{scaleSet}
		on:close={() => showUpdateModal = false}
		on:submit={(e) => handleUpdate(e.detail)}
	/>
{/if}

{#if showDeleteModal && scaleSet}
	<DeleteModal
		title="Delete Scale Set"
		message="Are you sure you want to delete this scale set? This action cannot be undone and will remove all associated runners."
		itemName={`Scale Set ${scaleSet.name} (${getEntityName(scaleSet)})`}
		on:close={() => showDeleteModal = false}
		on:confirm={handleDelete}
	/>
{/if}