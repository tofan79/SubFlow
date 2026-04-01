<script lang="ts">
	export let status: 'success' | 'warning' | 'error' | 'info' = 'info';
	export let text: string = '';

	const styles = {
		success: 'bg-[#00ffcc]/10 text-[#00ffcc] border-[#00ffcc]/30',
		warning: 'bg-[#ffe04a]/10 text-[#ffe04a] border-[#ffe04a]/30',
		error: 'bg-[#ff2d78]/10 text-[#ff2d78] border-[#ff2d78]/30',
		info: 'bg-blue-500/10 text-blue-400 border-blue-500/30'
	};
</script>

<span class="px-2 py-1 rounded text-xs font-semibold border {styles[status]}">
	{text}
</span>
