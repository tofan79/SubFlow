<script lang="ts">
	import { settingsStore, hardwareDisplay } from '$lib/stores/settings';
	import { GetSettings, SaveSettings } from '$lib/wails/index';
	import type { SettingsData } from '$lib/wails/index';
	import { onMount } from 'svelte';
	
	let settings: Partial<SettingsData> = {
		deeplApiKey: '', openaiApiKey: '', anthropicApiKey: '', geminiApiKey: '',
		groqApiKey: '', deepgramApiKey: '', xaiApiKey: '', qwenApiKey: '', ollamaEndpoint: 'http://localhost:11434',
		asrBackend: 'auto', whisperModel: 'base', preferredAsr: 'local'
	};
	
	onMount(async () => {
		try {
			const s = await GetSettings();
			if (s) settings = { ...settings, ...s };
		} catch(e) { console.log('Mock settings loaded'); }
	});

	async function saveSettings() {
		try {
			await SaveSettings(settings as SettingsData);
			alert('Pengaturan disimpan!');
		} catch (e) {
			console.error('Save error or mock', e);
		}
	}
</script>

<div class="max-w-4xl mx-auto space-y-8">
	<h2 class="text-3xl font-bold neon-text-primary">Pengaturan</h2>

	<div class="bg-[#0f0f1a] p-6 rounded-xl border border-[#141422] card-neon">
		<h3 class="text-xl font-semibold text-[#00ffcc] mb-4 flex items-center gap-2">
			<span class="material-symbols-outlined">key</span> Kunci API
		</h3>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="deepl">DeepL API Key</label>
				<input id="deepl" type="password" bind:value={settings.deeplApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="openai">OpenAI API Key</label>
				<input id="openai" type="password" bind:value={settings.openaiApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="anthropic">Anthropic API Key</label>
				<input id="anthropic" type="password" bind:value={settings.anthropicApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="gemini">Gemini API Key</label>
				<input id="gemini" type="password" bind:value={settings.geminiApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="groq">Groq API Key</label>
				<input id="groq" type="password" bind:value={settings.groqApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="deepgram">Deepgram API Key</label>
				<input id="deepgram" type="password" bind:value={settings.deepgramApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="xai">xAI API Key</label>
				<input id="xai" type="password" bind:value={settings.xaiApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div>
				<label class="block text-sm text-gray-400 mb-1" for="qwen">Qwen API Key</label>
				<input id="qwen" type="password" bind:value={settings.qwenApiKey} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
			<div class="md:col-span-2">
				<label class="block text-sm text-gray-400 mb-1" for="ollama">Ollama Endpoint</label>
				<input id="ollama" type="text" bind:value={settings.ollamaEndpoint} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-[#ff2d78] focus:ring-1 focus:ring-[#ff2d78] input-neon transition-all" />
			</div>
		</div>
	</div>

	<div class="bg-[#0f0f1a] p-6 rounded-xl border border-[#141422] card-neon">
		<h3 class="text-xl font-semibold text-[#ffe04a] mb-4 flex items-center gap-2">
			<span class="material-symbols-outlined">memory</span> Pengaturan ASR (Hardware & Model)
		</h3>
		<div class="mb-4 p-3 bg-[#141422] rounded-lg border border-[#00ffcc]/30">
			<span class="text-sm text-gray-400">Hardware Terdeteksi: </span>
			<span class="text-[#00ffcc] font-mono">{$hardwareDisplay || 'CPU / Unidentified'}</span>
		</div>
		
		<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
			<div>
				<label for="asr-backend" class="block text-sm text-gray-400 mb-1">Backend ASR</label>
				<select id="asr-backend" bind:value={settings.asrBackend} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:border-[#ffe04a]">
					{#each ['auto', 'cpu', 'cuda', 'rocm', 'coreml', 'openvino'] as b}
						<option value={b}>{b}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="whisper-model" class="block text-sm text-gray-400 mb-1">Model ASR Lokal</label>
				<select id="whisper-model" bind:value={settings.whisperModel} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:border-[#ffe04a]">
					{#each ['tiny', 'base', 'small', 'medium', 'large-v3'] as m}
						<option value={m}>{m}</option>
					{/each}
				</select>
			</div>
			<div>
				<label for="preferred-asr" class="block text-sm text-gray-400 mb-1">Layanan Utama ASR</label>
				<select id="preferred-asr" bind:value={settings.preferredAsr} class="w-full bg-[#141422] border border-gray-700 rounded-lg px-3 py-2 text-white focus:border-[#ffe04a]">
					{#each ['local', 'groq', 'deepgram'] as p}
						<option value={p}>{p}</option>
					{/each}
				</select>
			</div>
		</div>
	</div>

	<div class="flex justify-end">
		<button 
			on:click={saveSettings}
			class="px-6 py-2 bg-[#ff2d78] text-white font-bold rounded-lg hover:bg-[#ff2d78]/80 hover:shadow-[0_0_15px_#ff2d78] transition-all btn-neon-primary"
		>
			Simpan Pengaturan
		</button>
	</div>
</div>
