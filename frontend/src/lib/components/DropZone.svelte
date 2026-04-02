<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { EventsOn, EventsOff, WailsEvents } from '$lib/wails/index';

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

	onDestroy(() => {
		try {
			EventsOff(WailsEvents.FILE_DROPPED);
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

	function handleFileInput(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			const file = target.files[0] as File & { path?: string };
			dispatch('file', { file, path: file.path || file.name });
		}
	}
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="relative h-full bg-surface-container/60 border-2 border-dashed {isDragging ? 'border-primary bg-primary/10 shadow-[0_0_20px_rgba(255,45,120,0.2)]' : 'border-primary/40'} rounded-xl p-10 flex flex-col items-center justify-center text-center transition-all hover:border-primary hover:bg-surface-container/80 cursor-pointer"
	on:dragover={handleDragOver}
	on:dragleave={handleDragLeave}
	on:drop={handleDrop}
	on:click={() => document.getElementById('fileUpload')?.click()}
	on:keydown={(e) => e.key === 'Enter' && document.getElementById('fileUpload')?.click()}
>
	<div class="w-16 h-16 rounded-full flex items-center justify-center mb-6 {isDragging ? 'bg-primary/30 neon-border-primary animate-neon-pulse' : 'bg-primary/20 neon-border-primary'}">
		<span class="material-symbols-outlined text-primary text-3xl font-bold {isDragging ? 'neon-text-primary' : ''}">upload_file</span>
	</div>
	<h2 class="text-2xl font-headline font-bold text-on-surface mb-2">Mulai Proyek Baru</h2>
	<p class="text-slate-400 max-w-md mb-8">Tarik dan lepas file video atau subtitle di sini untuk memulai proses automasi AI secara instan.</p>
	<button class="bg-primary hover:bg-primary-container text-on-primary px-8 py-3 rounded-lg font-bold transition-all shadow-[0_0_16px_rgba(255,45,120,0.4)] active:scale-95 pointer-events-none">
		Pilih File
	</button>
	<input id="fileUpload" type="file" class="hidden" accept="video/*,.srt,.vtt,.ass" on:change={handleFileInput} />
</div>
