<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { EventsOn, WailsEvents } from '$lib/wails/index';

	const dispatch = createEventDispatcher();
	let isDragging = false;

	onMount(() => {
		try {
			EventsOn(WailsEvents.FILE_DROPPED, (path: unknown) => {
				if (typeof path === 'string') {
					dispatch('file', { path });
				}
			});
		} catch(e) {}
	});

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		if (e.dataTransfer?.files.length) {
			const file = e.dataTransfer.files[0] as File & { path?: string };
			dispatch('file', { file, path: file.path || file.name });
		}
	}
</script>

<div
	role="button"
	tabindex="0"
	class="w-full border-2 border-dashed rounded-xl p-12 text-center transition-all duration-300
		{isDragging ? 'border-[#ff2d78] bg-[#ff2d78]/10 shadow-[0_0_20px_#ff2d78] scale-[1.02]' : 'border-[#141422] bg-[#0f0f1a] hover:border-[#00ffcc]/50'}"
	on:dragover={handleDragOver}
	on:dragleave={handleDragLeave}
	on:drop={handleDrop}
>
	<span class="material-symbols-outlined text-5xl mb-4 {isDragging ? 'text-[#ff2d78]' : 'text-gray-500'}">
		upload_file
	</span>
	<h3 class="text-xl font-bold {isDragging ? 'text-white' : 'text-gray-300'} mb-2">
		Tarik & Lepas File Video/Audio di Sini
	</h3>
	<p class="text-gray-500 text-sm">Atau klik untuk memilih file dari komputer Anda</p>
</div>
