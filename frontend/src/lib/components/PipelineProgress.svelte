<script lang="ts">
	import { pipeline } from '$lib/stores/pipeline';

	const steps = [
		{ id: 'import', label: 'Impor' },
		{ id: 'asr', label: 'ASR' },
		{ id: 'correct', label: 'Koreksi' },
		{ id: 'context', label: 'Konteks' },
		{ id: 'translate', label: 'Terjemahan' },
		{ id: 'rewrite', label: 'Tulis Ulang' },
		{ id: 'qa', label: 'QA' }
	];

	// Mocking active step if store is empty
	$: activeIndex = steps.findIndex(s => s.id === $pipeline.currentStep) >= 0 
		? steps.findIndex(s => s.id === $pipeline.currentStep) 
		: 1; 
</script>

<div class="relative w-full">
	<div class="flex justify-between mb-2 relative z-10">
		{#each steps as step, i}
			<div class="flex flex-col items-center gap-2 w-full text-center">
				<div class="w-8 h-8 rounded-full flex items-center justify-center font-bold text-sm transition-all duration-300
					{i < activeIndex ? 'bg-[#00ffcc] text-black shadow-[0_0_10px_#00ffcc]' : 
					 i === activeIndex ? 'bg-[#ff2d78] text-white shadow-[0_0_15px_#ff2d78] scale-110' : 
					 'bg-[#141422] text-gray-500'}">
					{i < activeIndex ? '✓' : i + 1}
				</div>
				<span class="text-xs {i === activeIndex ? 'text-[#ff2d78] font-bold' : 'text-gray-400'}">
					{step.label}
				</span>
			</div>
		{/each}
	</div>
	
	<!-- Progress Line Background -->
	<div class="absolute top-4 left-[5%] right-[5%] h-1 bg-[#141422] -z-0">
		<div class="h-full bg-[#00ffcc] transition-all duration-500 shadow-[0_0_10px_#00ffcc]" 
			 style="width: {activeIndex > 0 ? (activeIndex / (steps.length - 1)) * 100 : 0}%">
		</div>
	</div>
</div>
