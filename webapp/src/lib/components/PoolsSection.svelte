<script lang="ts">
	import type { Pool } from '$lib/api/types.js';
	import { base } from '$app/paths';

	export let pools: Pool[];
	export let entityType: 'repository' | 'organization' | 'enterprise';
</script>

<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
	<div class="px-4 py-5 sm:p-6">
		<div class="flex items-center justify-between mb-4">
			<h2 class="text-lg font-medium text-gray-900 dark:text-white">Pools ({pools.length})</h2>
			<a href={`${base}/pools`} class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">View all pools</a>
		</div>
		{#if pools.length === 0}
			<p class="text-sm text-gray-500 dark:text-gray-400">No pools configured for this {entityType}.</p>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider w-auto min-w-0">ID</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Image</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Provider</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider w-20">Status</th>
						</tr>
					</thead>
					<tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
						{#each pools as pool}
							<tr>
								<td class="px-6 py-4 w-auto min-w-0">
									<a href={`${base}/pools/${pool.id}`} class="text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300 font-mono truncate block">
										{pool.id}
									</a>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
									<code class="bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded text-xs">{pool.image}</code>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">{pool.provider_name}</td>
								<td class="px-6 py-4 whitespace-nowrap w-20">
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