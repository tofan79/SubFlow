<script lang="ts">
	export let durationMinutes: number = 0;
	export const targetLanguage: string = 'id';
	
	// Mock calc
	$: estimatedCost = (durationMinutes * 0.02).toFixed(2);
</script>

<div class="p-4 bg-[#141422] rounded-lg border border-[#00ffcc]/20 flex justify-between items-center">
	<div>
		<div class="text-sm text-gray-400">Estimasi Biaya API</div>
		<div class="text-xs text-gray-500">Durasi: {durationMinutes} m</div>
	</div>
	<div class="text-xl font-bold text-[#00ffcc] drop-shadow-[0_0_5px_#00ffcc]">
		${estimatedCost}
	</div>
</div>
