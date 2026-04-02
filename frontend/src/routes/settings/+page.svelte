<script lang="ts">
	import { settingsStore } from '$lib/stores/settings';
	import { GetSettings, SaveSettings } from '$lib/wails/index';
	import type { SettingsData } from '$lib/wails/index';
	import { onMount } from 'svelte';
	
	interface AppSettings extends Partial<SettingsData> {
		translateProvider?: string;
		translateApiKey?: string;
		translateModel?: string;
		rewriteProvider?: string;
		rewriteApiKey?: string;
		rewriteTemperature?: number;
		asrExecutablePath?: string;
		asrModelDirectory?: string;
		asrHardware?: string;
		autoSaveInterval?: number;
		loggingLevel?: string;
		workspacePath?: string;
		exportPath?: string;
		[key: string]: any;
	}

	let settings: AppSettings = {
		translateProvider: 'DeepL Pro API',
		translateApiKey: '•••••••••••••••••••••',
		translateModel: 'v3-standard (Default)',
		rewriteProvider: 'OpenAI GPT-4o',
		rewriteApiKey: '•••••••••••••••••••••',
		rewriteTemperature: 0.7,
		asrExecutablePath: 'C:/Program Files/Whisper/whisper-main.exe',
		asrModelDirectory: 'C:/Users/Neon/Models/Whisper/',
		asrHardware: 'CUDA (NVIDIA GeForce RTX 4090)',
		autoSaveInterval: 5,
		loggingLevel: 'Info',
		workspacePath: 'D:/Projects/TokyoNeon/Workspace',
		exportPath: 'D:/Projects/TokyoNeon/Exports'
	};

	let showL1Key = false;
	let showL2Key = false;
	
	onMount(async () => {
		try {
			const s = await GetSettings();
			// Only merging standard settings that exist, keeping our mock ones if they don't
			if (s) settings = { ...settings, ...s };
		} catch(e) { console.log('Mock settings loaded for UI styling purposes'); }
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

<div class="max-w-6xl mx-auto pb-12">
	<header class="mb-10">
		<h1 class="text-4xl font-headline font-extrabold text-on-background tracking-tight mb-2">
			Pengaturan Sistem
		</h1>
		<p class="text-on-surface-variant font-body mb-8">Konfigurasi infrastruktur engine translasi, model AI, dan manajemen direktori.</p>
	</header>

	<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
		<!-- Card 1: Translation API -->
		<section class="card-neon p-6 relative">
			<h3 class="text-[14px] font-headline font-bold text-primary mb-6 flex items-center gap-2 uppercase tracking-widest">
				<span class="material-symbols-outlined text-primary text-xl">translate</span>
				TRANSLATION API (LAYER 1)
			</h3>
			
			<div class="space-y-5">
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="translateProvider">Provider Selection</label>
					<select id="translateProvider" bind:value={settings.translateProvider} class="select-neon w-full">
						<option>DeepL Pro API</option>
						<option>Google Translate API</option>
					</select>
				</div>
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="translateApiKey">API Key</label>
					<div class="relative">
						{#if showL1Key}
							<input id="translateApiKey" type="text" bind:value={settings.translateApiKey} class="input-neon w-full pr-10 font-mono tracking-[0.2em] text-lg" />
						{:else}
							<input id="translateApiKey" type="password" bind:value={settings.translateApiKey} class="input-neon w-full pr-10 font-mono tracking-[0.2em] text-lg" />
						{/if}
						<button 
							type="button"
							on:click={() => showL1Key = !showL1Key}
							class="absolute right-3 top-2.5 text-outline hover:text-primary transition-colors focus:outline-none">
							<span class="material-symbols-outlined text-sm">{showL1Key ? 'visibility_off' : 'visibility'}</span>
						</button>
					</div>
				</div>
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="translateModel">Model Dropdown</label>
					<select id="translateModel" bind:value={settings.translateModel} class="select-neon w-full">
						<option>v3-standard (Default)</option>
						<option>v2-legacy</option>
					</select>
				</div>
				<div class="flex items-center justify-between pt-4 mt-2">
					<div class="flex items-center gap-2">
						<span class="w-2 h-2 rounded-full bg-secondary shadow-[0_0_8px_#00ffcc]"></span>
						<span class="text-[10px] font-bold text-secondary uppercase tracking-widest">Ready to Connect</span>
					</div>
					<button class="btn-ghost-primary uppercase text-[10px] tracking-widest font-bold">Test Koneksi</button>
				</div>
			</div>
		</section>

		<!-- Card 2: Rewrite API -->
		<section class="card-neon-secondary p-6 relative">
			<h3 class="text-[14px] font-headline font-bold text-secondary mb-6 flex items-center gap-2 uppercase tracking-widest">
				<span class="material-symbols-outlined text-secondary text-xl">auto_fix_high</span>
				REWRITE API (LAYER 2)
			</h3>
			
			<div class="space-y-5">
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="rewriteProvider">L2 Provider</label>
					<select id="rewriteProvider" bind:value={settings.rewriteProvider} class="select-neon w-full focus:border-secondary focus:ring-secondary/30">
						<option>OpenAI GPT-4o</option>
						<option>Anthropic Claude 3.5 Sonnet</option>
						<option>Google Gemini 1.5 Pro</option>
					</select>
				</div>
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="rewriteApiKey">API Key</label>
					<div class="relative">
						{#if showL2Key}
							<input id="rewriteApiKey" type="text" bind:value={settings.rewriteApiKey} class="input-neon w-full pr-10 focus:border-secondary focus:ring-secondary/30 font-mono tracking-[0.2em] text-lg" />
						{:else}
							<input id="rewriteApiKey" type="password" bind:value={settings.rewriteApiKey} class="input-neon w-full pr-10 focus:border-secondary focus:ring-secondary/30 font-mono tracking-[0.2em] text-lg" />
						{/if}
						<button 
							type="button"
							on:click={() => showL2Key = !showL2Key}
							class="absolute right-3 top-2.5 text-outline hover:text-secondary transition-colors focus:outline-none">
							<span class="material-symbols-outlined text-sm">{showL2Key ? 'visibility_off' : 'visibility'}</span>
						</button>
					</div>
				</div>
				<div>
					<div class="flex justify-between items-center mb-2">
						<label class="text-[10px] font-label uppercase tracking-widest text-on-surface-variant" for="rewriteTemperature">Temperature / Tone</label>
						<span class="text-secondary font-bold text-xs">{settings.rewriteTemperature}</span>
					</div>
					<div class="relative flex items-center h-8">
						<input id="rewriteTemperature" type="range" bind:value={settings.rewriteTemperature} min="0" max="2" step="0.1" class="w-full h-1.5 bg-surface-container-highest rounded-lg appearance-none cursor-pointer accent-secondary" />
					</div>
				</div>
				<div class="flex items-center justify-between pt-4 mt-2">
					<div class="flex items-center gap-2">
						<span class="w-2 h-2 rounded-full bg-outline"></span>
						<span class="text-[10px] font-bold text-outline uppercase tracking-widest">Idle</span>
					</div>
					<button class="btn-ghost-secondary uppercase text-[10px] tracking-widest font-bold">Test Koneksi</button>
				</div>
			</div>
		</section>

		<!-- Card 3: Local ASR Engine -->
		<section class="card-neon-tertiary p-6 relative">
			<h3 class="text-[14px] font-headline font-bold text-tertiary mb-6 flex items-center gap-2 uppercase tracking-widest">
				<span class="material-symbols-outlined text-tertiary text-xl">settings_voice</span>
				LOCAL ASR ENGINE (WHISPER.CPP)
			</h3>
			
			<div class="space-y-6">
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="asrExecutablePath">Executable Path</label>
					<div class="flex gap-2">
						<input id="asrExecutablePath" type="text" bind:value={settings.asrExecutablePath} class="input-neon w-full focus:border-tertiary focus:ring-tertiary/30 font-mono text-[11px] text-on-surface-variant" />
						<button class="bg-surface-container-high hover:bg-surface-container-highest p-2.5 rounded-lg border border-outline-variant/50 text-outline transition-colors group">
							<span class="material-symbols-outlined text-[18px] group-hover:text-tertiary">folder</span>
						</button>
					</div>
				</div>
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="asrModelDirectory">Model Directory</label>
					<div class="flex gap-2">
						<input id="asrModelDirectory" type="text" bind:value={settings.asrModelDirectory} class="input-neon w-full focus:border-tertiary focus:ring-tertiary/30 font-mono text-[11px] text-on-surface-variant" />
						<button class="bg-surface-container-high hover:bg-surface-container-highest p-2.5 rounded-lg border border-outline-variant/50 text-outline transition-colors group">
							<span class="material-symbols-outlined text-[18px] group-hover:text-tertiary">folder_open</span>
						</button>
					</div>
				</div>
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="asrHardware">Hardware Acceleration</label>
					<div class="relative">
						<select id="asrHardware" bind:value={settings.asrHardware} class="select-neon w-full focus:border-tertiary focus:ring-tertiary/30 font-mono text-[11px]">
							<option>CUDA (NVIDIA GeForce RTX 4090)</option>
							<option>CPU Only</option>
							<option>ROCm (AMD)</option>
						</select>
					</div>
				</div>
			</div>
		</section>

		<!-- Card 4: Directory & App -->
		<section class="card-base p-6 relative border-transparent bg-surface-container-low hover:border-outline-variant/30 transition-colors">
			<h3 class="text-[14px] font-headline font-bold text-on-surface mb-6 flex items-center gap-2 uppercase tracking-widest">
				<span class="material-symbols-outlined text-outline text-xl">folder_managed</span>
				DIREKTORI & APLIKASI
			</h3>
			
			<div class="space-y-6">
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="autoSaveInterval">Auto-Save Interval</label>
						<div class="relative">
							<input id="autoSaveInterval" type="number" bind:value={settings.autoSaveInterval} class="input-neon w-full pr-12 font-mono text-[11px]" />
							<span class="absolute right-3 top-2.5 text-[10px] text-outline font-label uppercase">MIN</span>
						</div>
					</div>
					<div>
						<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="loggingLevel">Logging Level</label>
						<select id="loggingLevel" bind:value={settings.loggingLevel} class="select-neon w-full font-mono text-[11px]">
							<option>Info</option>
							<option>Debug</option>
							<option>Error</option>
							<option>Warn</option>
						</select>
					</div>
				</div>
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="workspacePath">Default Workspace Path</label>
					<div class="flex gap-2">
						<input id="workspacePath" type="text" bind:value={settings.workspacePath} class="input-neon w-full font-mono text-[11px] text-on-surface-variant" />
						<button class="bg-surface-container-high hover:bg-surface-container-highest p-2.5 rounded-lg border border-outline-variant/50 text-outline transition-colors group">
							<span class="material-symbols-outlined text-[18px] group-hover:text-on-surface">edit</span>
						</button>
					</div>
				</div>
				<div>
					<label class="block text-[10px] font-label uppercase tracking-widest text-on-surface-variant mb-2" for="exportPath">Default Export Folder</label>
					<div class="flex gap-2">
						<input id="exportPath" type="text" bind:value={settings.exportPath} class="input-neon w-full font-mono text-[11px] text-on-surface-variant" />
						<button class="bg-surface-container-high hover:bg-surface-container-highest p-2.5 rounded-lg border border-outline-variant/50 text-outline transition-colors group">
							<span class="material-symbols-outlined text-[18px] group-hover:text-on-surface">edit</span>
						</button>
					</div>
				</div>
			</div>
		</section>
	</div>
	
	<!-- Bottom save button, even though not explicitly required, good UX to have from old design -->
	<div class="mt-8 flex justify-end opacity-20 hover:opacity-100 transition-opacity hidden">
		<button on:click={saveSettings} class="btn-neon-primary px-8 py-3 text-sm tracking-widest font-bold uppercase transition-all">
			Simpan Konfigurasi
		</button>
	</div>
</div>
