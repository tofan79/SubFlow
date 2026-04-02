<script lang="ts">
	import { pipeline, isRunning } from '$lib/stores/pipeline';

	const steps = [
		{ id: 'import', label: 'Impor', icon: 'file_download' },
		{ id: 'asr', label: 'ASR', icon: 'graphic_eq' },
		{ id: 'correct', label: 'Koreksi', icon: 'spellcheck' },
		{ id: 'context', label: 'Konteks', icon: 'category' },
		{ id: 'translate', label: 'Trans', icon: 'translate' },
		{ id: 'rewrite', label: 'Rewrite', icon: 'edit_note' },
		{ id: 'qa', label: 'QA', icon: 'fact_check' }
	];

	// Mocking active step if store is empty
	$: activeIndex = steps.findIndex(s => s.id === $pipeline.currentStep) >= 0 
		? steps.findIndex(s => s.id === $pipeline.currentStep) 
		: 4; // Mock translate
	
	$: progressPercent = Math.round((activeIndex / (steps.length - 1)) * 100);
</script>

<div class="bg-surface-container rounded-xl overflow-hidden border border-outline-variant/20 shadow-xl">
	<!-- Pipeline Header -->
	<div class="p-4 border-b border-outline-variant/30 flex items-center justify-between bg-surface-container-high/30">
		<div class="flex items-center gap-3">
			<span class="material-symbols-outlined {$isRunning ? 'text-secondary animate-spin' : 'text-slate-500'} text-sm" style="animation-duration: 3s;">cycle</span>
			<span class="font-label text-sm text-on-surface">Dokumentasi_Teknis_V2.mp4</span>
		</div>
		<span class="text-[10px] font-label px-2 py-0.5 rounded bg-secondary/20 text-secondary border border-secondary/30 uppercase">
			{steps[activeIndex]?.id || 'IDLE'}
		</span>
	</div>

	<!-- Progress Body -->
	<div class="p-6">
		<div class="flex justify-between items-end mb-3">
			<div class="space-y-1">
				<p class="text-xs text-slate-400">Tahap Saat Ini:</p>
				<p class="text-sm font-semibold text-on-surface">{steps[activeIndex]?.label || 'Menunggu'}</p>
			</div>
			<span class="text-xl font-headline font-bold text-secondary neon-glow-secondary">{progressPercent}%</span>
		</div>
		
		<!-- Linear Progress Bar -->
		<div class="w-full bg-surface-container-highest rounded-full h-2 relative overflow-hidden mb-8">
			<div class="absolute inset-0 bg-secondary/10 w-full"></div>
			<div class="h-full bg-secondary shadow-[0_0_10px_#00ffcc] transition-all duration-700" style="width: {progressPercent}%"></div>
		</div>
		
		<!-- Steps Grid / Line -->
		<div class="relative w-full">
			<div class="absolute top-3 left-[7%] right-[7%] h-0.5 bg-surface-container-highest z-0">
				<div class="h-full bg-secondary transition-all duration-500 shadow-[0_0_10px_#00ffcc]" 
					style="width: {activeIndex > 0 ? (activeIndex / (steps.length - 1)) * 100 : 0}%">
				</div>
			</div>
			<div class="flex justify-between relative z-10 w-full">
				{#each steps as step, i}
					<div class="flex flex-col items-center gap-2 group cursor-default">
						<div class="w-6 h-6 rounded-full flex items-center justify-center font-bold text-xs transition-all duration-300 border-2
							{i < activeIndex ? 'bg-secondary border-secondary text-black shadow-[0_0_10px_#00ffcc]' : 
							 i === activeIndex ? 'bg-surface-container border-secondary text-secondary shadow-[0_0_15px_#00ffcc] scale-125' : 
							 'bg-surface border-outline-variant text-outline'}">
							{#if i < activeIndex}
								<span class="material-symbols-outlined text-[14px]">check</span>
							{:else}
								<span>{i + 1}</span>
							{/if}
						</div>
						<span class="text-[10px] font-label hidden md:block transition-colors {i === activeIndex ? 'text-secondary font-bold' : i < activeIndex ? 'text-slate-300' : 'text-slate-500'}">
							{step.label}
						</span>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>
