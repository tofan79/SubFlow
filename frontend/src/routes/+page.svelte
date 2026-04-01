<script lang="ts">
	import DropZone from '$lib/components/DropZone.svelte';
	import PipelineProgress from '$lib/components/PipelineProgress.svelte';
	import { onMount } from 'svelte';
	import { pipeline, isRunning } from '$lib/stores/pipeline';
	import { recentProjects } from '$lib/stores/project';
	import { GetStats } from '$lib/wails/index';
	import type { AppStats } from '$lib/wails/index';

	let stats: AppStats = { totalProjects: 12, totalSegments: 450, totalCharsProcessed: 12500, totalMinutesAsr: 45 };

	onMount(async () => {
		try {
			stats = await GetStats();
		} catch (e) {
			console.log('Using mock stats');
		}
	});
</script>

<div class="max-w-6xl mx-auto space-y-8">
	<div class="flex justify-between items-end">
		<div>
			<h2 class="text-3xl font-bold neon-text-primary">Beranda</h2>
			<p class="text-gray-400 mt-1">Selamat datang di SubFlow, mulai terjemahkan video Anda.</p>
		</div>
	</div>

	<!-- Stats Grid -->
	<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
		{#each [
			{ label: 'Total Proyek', val: stats.totalProjects, color: 'text-[#ff2d78]' },
			{ label: 'Total Segmen', val: stats.totalSegments, color: 'text-[#00ffcc]' },
			{ label: 'Karakter Diproses', val: stats.totalCharsProcessed, color: 'text-[#ffe04a]' },
			{ label: 'Menit ASR', val: stats.totalMinutesAsr, color: 'text-white' }
		] as stat}
			<div class="bg-[#0f0f1a] p-4 rounded-xl border border-[#141422] hover:border-[#ff2d78]/50 transition-colors card-neon">
				<div class="text-gray-400 text-sm mb-1">{stat.label}</div>
				<div class="text-2xl font-bold {stat.color} drop-shadow-[0_0_8px_currentColor]">{stat.val}</div>
			</div>
		{/each}
	</div>

	<!-- Drop Zone -->
	<DropZone />

	<!-- Active Pipeline -->
	{#if $isRunning || true}
		<div class="bg-[#0f0f1a] p-6 rounded-xl border border-[#00ffcc]/30 card-neon">
			<h3 class="text-xl font-semibold text-[#00ffcc] mb-4">Proses Saat Ini</h3>
			<PipelineProgress />
		</div>
	{/if}

	<!-- Recent Projects -->
	<div>
		<h3 class="text-xl font-semibold mb-4 text-[#ffe04a]">Proyek Terbaru</h3>
		<div class="bg-[#0f0f1a] rounded-xl border border-[#141422] overflow-hidden">
			{#if $recentProjects.length > 0}
				<ul class="divide-y divide-[#141422]">
					{#each $recentProjects as proj}
						<li class="p-4 hover:bg-[#141422] transition-colors flex justify-between items-center">
							<span>{proj.name}</span>
							<span class="text-sm text-gray-400">{proj.state}</span>
						</li>
					{/each}
				</ul>
			{:else}
				<div class="p-8 text-center text-gray-500">
					Belum ada proyek terbaru. Tarik file ke atas untuk memulai.
				</div>
			{/if}
		</div>
	</div>
</div>
