<script lang="ts">
	import type { EntityEvent } from '$lib/api/types.js';

	export let events: EntityEvent[] | undefined;
	export let eventsContainer: HTMLElement | undefined = undefined;

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleString();
	}
</script>

{#if events && events.length > 0}
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
		<div class="px-4 py-5 sm:p-6">
			<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Events</h2>
			<div bind:this={eventsContainer} class="space-y-3 max-h-96 overflow-y-auto scroll-smooth">
				{#each events as event}
					<div class="bg-gray-50 dark:bg-gray-700 p-4 rounded-lg">
						<div class="flex justify-between items-start">
							<p class="text-sm text-gray-900 dark:text-white flex-1 mr-4">{event.message}</p>
							<div class="flex items-center space-x-2 flex-shrink-0">
								<span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full
									{event.event_level === 'error' ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200' :
									event.event_level === 'warning' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200' :
									'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'}">
									{event.event_level}
								</span>
								<span class="text-xs text-gray-500 dark:text-gray-400">{formatDate(event.created_at)}</span>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
{:else}
	<div class="bg-white dark:bg-gray-800 shadow rounded-lg">
		<div class="px-4 py-5 sm:p-6">
			<h2 class="text-lg font-medium text-gray-900 dark:text-white mb-4">Events</h2>
			<div class="text-center py-8">
				<svg class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
				</svg>
				<p class="text-sm text-gray-500 dark:text-gray-400">No events available</p>
			</div>
		</div>
	</div>
{/if}