<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { Button, Select, Toggle, Label, Badge, Alert } from 'flowbite-svelte';
	import { 
		ExportSubtitle, 
		SelectDirectory,
		type ExportOptions 
	} from '$lib/wails';
	import { editor } from '$lib/stores/editor';

	export let projectId: string;

	const dispatch = createEventDispatcher<{
		export: { path: string };
		close: void;
	}>();

	type ExportFormat = 'srt' | 'vtt' | 'ass' | 'txt';
	type ExportLayer = 'source' | 'l1' | 'l2';

	let format: ExportFormat = 'srt';
	let layer: ExportLayer = 'l2';
	let dualSubtitle = false;
	let outputDir = '';
	let isExporting = false;
	let error: string | null = null;
	let success: string | null = null;

	const formatOptions = [
		{ value: 'srt', name: 'SRT - SubRip (.srt)' },
		{ value: 'vtt', name: 'VTT - WebVTT (.vtt)' },
		{ value: 'ass', name: 'ASS - Advanced SSA (.ass)' },
		{ value: 'txt', name: 'TXT - Plain Text (.txt)' }
	];

	const layerOptions = [
		{ value: 'source', name: 'Source (Original)' },
		{ value: 'l1', name: 'Layer 1 (Terjemahan Mentah)' },
		{ value: 'l2', name: 'Layer 2 (Hasil Akhir)' }
	];

	$: qaStats = getQAStats($editor.segments);

	function getQAStats(segments: typeof $editor.segments) {
		let pass = 0, warn = 0, err = 0, pending = 0;
		for (const seg of segments) {
			switch (seg.qaStatus) {
				case 'pass': pass++; break;
				case 'warn': warn++; break;
				case 'error': err++; break;
				case 'pending': pending++; break;
			}
		}
		return { pass, warn, error: err, pending, total: segments.length };
	}

	async function selectOutputDir() {
		try {
			const dir = await SelectDirectory('Pilih Folder Output');
			if (dir) {
				outputDir = dir;
				error = null;
			}
		} catch (e) {
			error = `Gagal memilih folder: ${e}`;
		}
	}

	async function doExport() {
		if (!outputDir) {
			error = 'Pilih folder output terlebih dahulu';
			return;
		}

		if (qaStats.error > 0) {
			const confirm = window.confirm(
				`Ada ${qaStats.error} segmen dengan QA error. Lanjutkan ekspor?`
			);
			if (!confirm) return;
		}

		isExporting = true;
		error = null;
		success = null;

		try {
			const options: ExportOptions = {
				projectId,
				format,
				outputDir,
				layer: dualSubtitle ? 'l2' : layer,
				dualSubtitle
			};

			const outputPath = await ExportSubtitle(options);
			success = `Berhasil diekspor ke: ${outputPath}`;
			dispatch('export', { path: outputPath });
		} catch (e) {
			error = `Gagal ekspor: ${e}`;
		} finally {
			isExporting = false;
		}
	}
</script>

<div class="export-panel bg-dark-800 rounded-xl p-6 border border-dark-600">
	<h2 class="text-xl font-semibold text-white mb-6 flex items-center gap-2">
		<span class="material-symbols-outlined text-neon-cyan">download</span>
		Ekspor Subtitle
	</h2>

	<!-- QA Status Summary -->
	<div class="qa-summary mb-6 p-4 rounded-lg bg-dark-700/50">
		<h3 class="text-sm font-medium text-gray-400 mb-3">Status QA</h3>
		<div class="flex flex-wrap gap-2">
			<Badge color="green" class="flex items-center gap-1">
				<span class="material-symbols-outlined text-sm">check_circle</span>
				{qaStats.pass} Pass
			</Badge>
			{#if qaStats.warn > 0}
				<Badge color="yellow" class="flex items-center gap-1">
					<span class="material-symbols-outlined text-sm">warning</span>
					{qaStats.warn} Warning
				</Badge>
			{/if}
			{#if qaStats.error > 0}
				<Badge color="red" class="flex items-center gap-1">
					<span class="material-symbols-outlined text-sm">error</span>
					{qaStats.error} Error
				</Badge>
			{/if}
			{#if qaStats.pending > 0}
				<Badge color="dark" class="flex items-center gap-1">
					<span class="material-symbols-outlined text-sm">pending</span>
					{qaStats.pending} Pending
				</Badge>
			{/if}
		</div>
		{#if qaStats.error > 0}
			<p class="mt-2 text-sm text-amber-400">
				⚠️ Disarankan untuk memperbaiki error sebelum ekspor
			</p>
		{/if}
	</div>

	<!-- Format Selection -->
	<div class="mb-5">
		<Label for="format" class="text-gray-300 mb-2 block">Format Output</Label>
		<Select 
			id="format" 
			bind:value={format}
			class="bg-dark-700 border-dark-600 text-white"
		>
			{#each formatOptions as opt}
				<option value={opt.value}>{opt.name}</option>
			{/each}
		</Select>
	</div>

	<!-- Layer Selection -->
	<div class="mb-5">
		<Label for="layer" class="text-gray-300 mb-2 block">Layer Teks</Label>
		<Select 
			id="layer" 
			bind:value={layer}
			disabled={dualSubtitle}
			class="bg-dark-700 border-dark-600 text-white disabled:opacity-50"
		>
			{#each layerOptions as opt}
				<option value={opt.value}>{opt.name}</option>
			{/each}
		</Select>
	</div>

	<!-- Dual Subtitle Toggle -->
	<div class="mb-5 flex items-center justify-between p-3 rounded-lg bg-dark-700/50">
		<div>
			<Label class="text-gray-300">Dual Subtitle</Label>
			<p class="text-xs text-gray-500 mt-1">
				Tampilkan source + terjemahan bersamaan
			</p>
		</div>
		<Toggle bind:checked={dualSubtitle} color="cyan" />
	</div>

	<!-- Output Directory -->
	<div class="mb-6">
		<Label class="text-gray-300 mb-2 block">Folder Output</Label>
		<div class="flex gap-2">
			<input
				type="text"
				readonly
				value={outputDir}
				placeholder="Pilih folder..."
				class="flex-1 px-3 py-2 rounded-lg bg-dark-700 border border-dark-600 text-white 
				       placeholder-gray-500 focus:border-neon-cyan focus:ring-1 focus:ring-neon-cyan"
			/>
			<Button color="alternative" on:click={selectOutputDir} class="whitespace-nowrap">
				<span class="material-symbols-outlined text-sm mr-1">folder_open</span>
				Pilih
			</Button>
		</div>
	</div>

	<!-- Error/Success Messages -->
	{#if error}
		<Alert color="red" class="mb-4" dismissable on:close={() => error = null}>
			{error}
		</Alert>
	{/if}
	{#if success}
		<Alert color="green" class="mb-4" dismissable on:close={() => success = null}>
			{success}
		</Alert>
	{/if}

	<!-- Action Buttons -->
	<div class="flex gap-3 justify-end">
		<Button color="alternative" on:click={() => dispatch('close')}>
			Batal
		</Button>
		<Button 
			color="primary" 
			on:click={doExport}
			disabled={isExporting || !outputDir}
			class="!bg-neon-cyan !text-dark-900 hover:!bg-neon-cyan/80 disabled:opacity-50"
		>
			{#if isExporting}
				<span class="material-symbols-outlined animate-spin text-sm mr-2">progress_activity</span>
				Mengekspor...
			{:else}
				<span class="material-symbols-outlined text-sm mr-2">download</span>
				Ekspor
			{/if}
		</Button>
	</div>
</div>

<style>
	.export-panel {
		min-width: 400px;
		max-width: 500px;
	}
</style>
