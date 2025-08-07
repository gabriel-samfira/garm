<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher<{
		close: void;
	}>();

	function handleBackdropClick() {
		dispatch('close');
	}

	function handleContentClick(event: Event) {
		event.stopPropagation();
	}
</script>

<div 
	class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center p-4" 
	on:click={handleBackdropClick}
	role="dialog"
	aria-modal="true"
>
	<div 
		class="relative w-full mx-auto" 
		on:click={handleContentClick}
		role="document"
	>
		<slot />
	</div>
</div>