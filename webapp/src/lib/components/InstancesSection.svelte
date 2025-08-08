<script lang="ts">
	import type { Instance } from '$lib/api/types.js';
	import { base } from '$app/paths';

	export let instances: Instance[];
	export let entityType: 'repository' | 'organization' | 'enterprise';
	export let onDeleteInstance: (instance: Instance) => void;

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
</script>

<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
	<div class="px-4 py-5 sm:p-6">
		<div class="flex items-center justify-between mb-4">
			<h2 class="text-lg font-medium text-gray-900 dark:text-white">Instances ({instances.length})</h2>
			<a href={`${base}/instances`} class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300">View all instances</a>
		</div>
		{#if instances.length === 0}
			<p class="text-sm text-gray-500 dark:text-gray-400">No instances running for this {entityType}.</p>
		{:else}
			<div class="overflow-x-auto">
				<table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
					<thead class="bg-gray-50 dark:bg-gray-700">
						<tr>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Name</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Runner Status</th>
							<th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Created</th>
							<th class="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Actions</th>
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
									<span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset {getStatusBadgeClass(instance.status)}">
										{instance.status}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap">
									<span class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset {getStatusBadgeClass(instance.runner_status)}">
										{instance.runner_status}
									</span>
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
									{formatDate(instance.created_at)}
								</td>
								<td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
									<button
										on:click={() => onDeleteInstance(instance)}
										class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300"
										title="Delete instance"
									>
										<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
										</svg>
									</button>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>