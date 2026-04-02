<script lang="ts">
	import { onMount } from 'svelte';
	import { Modal, Button, Alert, Toggle } from 'flowbite-svelte';
	import { 
		GetGlossary, 
		AddGlossaryTerm, 
		UpdateGlossaryTerm, 
		DeleteGlossaryTerm,
		ImportGlossary,
		ExportGlossary,
		SelectFile,
		SelectDirectory,
		type GlossaryTerm
	} from '$lib/wails';

	let terms: GlossaryTerm[] = [];
	let filteredTerms: GlossaryTerm[] = [];
	let searchQuery = '';
	let isLoading = true;
	let error: string | null = null;
	let success: string | null = null;

	let showAddModal = false;
	let showEditModal = false;
	let showDeleteConfirm = false;
	let editingTerm: GlossaryTerm | null = null;
	let deletingTerm: GlossaryTerm | null = null;

	let newTerm: Partial<GlossaryTerm> = {
		sourceTerm: '',
		targetTerm: '',
		caseSensitive: false,
		notes: ''
	};

	$: filteredTerms = terms.filter(t => {
		const q = searchQuery.toLowerCase();
		return t.sourceTerm.toLowerCase().includes(q) || 
		       t.targetTerm.toLowerCase().includes(q) ||
		       (t.notes?.toLowerCase().includes(q) ?? false);
	});

	$: stats = {
		total: terms.length,
		caseSensitive: terms.filter(t => t.caseSensitive).length
	};

	onMount(async () => {
		await loadTerms();
	});

	async function loadTerms() {
		isLoading = true;
		error = null;
		try {
			terms = await GetGlossary();
		} catch (e) {
			error = `Gagal memuat glosarium: ${e}`;
			terms = [];
		} finally {
			isLoading = false;
		}
	}

	function openAddModal() {
		newTerm = { sourceTerm: '', targetTerm: '', caseSensitive: false, notes: '' };
		showAddModal = true;
	}

	function openEditModal(term: GlossaryTerm) {
		editingTerm = { ...term };
		showEditModal = true;
	}

	function confirmDelete(term: GlossaryTerm) {
		deletingTerm = term;
		showDeleteConfirm = true;
	}

	async function addTerm() {
		if (!newTerm.sourceTerm || !newTerm.targetTerm) {
			error = 'Istilah asal dan target wajib diisi';
			return;
		}

		try {
			const added = await AddGlossaryTerm({
				id: '',
				sourceTerm: newTerm.sourceTerm!,
				targetTerm: newTerm.targetTerm!,
				caseSensitive: newTerm.caseSensitive ?? false,
				notes: newTerm.notes ?? ''
			});
			terms = [added, ...terms];
			showAddModal = false;
			success = 'Istilah berhasil ditambahkan';
			setTimeout(() => success = null, 3000);
		} catch (e) {
			error = `Gagal menambahkan istilah: ${e}`;
		}
	}

	async function updateTerm() {
		if (!editingTerm) return;

		try {
			await UpdateGlossaryTerm(editingTerm);
			terms = terms.map(t => t.id === editingTerm!.id ? editingTerm! : t);
			showEditModal = false;
			editingTerm = null;
			success = 'Istilah berhasil diperbarui';
			setTimeout(() => success = null, 3000);
		} catch (e) {
			error = `Gagal memperbarui istilah: ${e}`;
		}
	}

	async function deleteTerm() {
		if (!deletingTerm) return;

		try {
			await DeleteGlossaryTerm(deletingTerm.id);
			terms = terms.filter(t => t.id !== deletingTerm!.id);
			showDeleteConfirm = false;
			deletingTerm = null;
			success = 'Istilah berhasil dihapus';
			setTimeout(() => success = null, 3000);
		} catch (e) {
			error = `Gagal menghapus istilah: ${e}`;
		}
	}

	async function handleImport() {
		try {
			const filePath = await SelectFile('Pilih File JSON Glosarium', ['*.json']);
			if (!filePath) return;

			const count = await ImportGlossary(filePath);
			await loadTerms();
			success = `Berhasil mengimpor ${count} istilah`;
			setTimeout(() => success = null, 3000);
		} catch (e) {
			error = `Gagal mengimpor: ${e}`;
		}
	}

	async function handleExport() {
		try {
			const dir = await SelectDirectory('Pilih Folder Ekspor');
			if (!dir) return;

			const filePath = `${dir}/glossary_export.json`;
			await ExportGlossary(filePath);
			success = `Berhasil mengekspor ke: ${filePath}`;
			setTimeout(() => success = null, 5000);
		} catch (e) {
			error = `Gagal mengekspor: ${e}`;
		}
	}
</script>

<div class="max-w-7xl mx-auto pb-12">
	<header class="flex justify-between items-end mb-10">
		<div>
			<div class="flex items-center gap-3 mb-2">
				<div class="w-8 h-0.5 bg-primary shadow-[0_0_8px_#ff2d78]"></div>
				<span class="text-[10px] font-label font-bold text-primary uppercase tracking-[0.2em] neon-glow-primary">Pusat Linguistik</span>
			</div>
			<h1 class="text-4xl font-headline font-extrabold text-on-background tracking-tight">
				Glosarium <span class="text-secondary drop-shadow-[0_0_8px_rgba(0,255,204,0.6)]">Istilah</span>
			</h1>
		</div>
		<button 
			on:click={openAddModal}
			class="flex items-center gap-2 px-6 py-3 border border-primary text-primary hover:bg-primary/10 rounded-lg font-headline font-bold text-sm tracking-widest uppercase transition-all shadow-[0_0_12px_rgba(255,45,120,0.2)] hover:shadow-[0_0_20px_rgba(255,45,120,0.4)]"
		>
			<span class="material-symbols-outlined text-lg">add</span>
			Tambah Istilah Baru
		</button>
	</header>

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

	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-5 flex items-center gap-5 hover:border-primary/50 transition-colors">
			<div class="w-14 h-14 rounded-lg bg-primary/10 flex items-center justify-center border border-primary/20">
				<span class="material-symbols-outlined text-primary text-2xl drop-shadow-[0_0_8px_rgba(255,45,120,0.5)]">translate</span>
			</div>
			<div>
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400 mb-1">Total Istilah</p>
				<p class="text-3xl font-headline font-bold text-on-surface">{stats.total.toLocaleString()}</p>
			</div>
		</div>
		
		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-5 flex items-center gap-5 hover:border-secondary/50 transition-colors">
			<div class="w-14 h-14 rounded-lg bg-secondary/10 flex items-center justify-center border border-secondary/20">
				<span class="material-symbols-outlined text-secondary text-2xl drop-shadow-[0_0_8px_rgba(0,255,204,0.5)]">match_case</span>
			</div>
			<div>
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400 mb-1">Case Sensitive</p>
				<p class="text-3xl font-headline font-bold text-on-surface">{stats.caseSensitive}</p>
			</div>
		</div>

		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-5 flex items-center gap-5 hover:border-tertiary/50 transition-colors">
			<div class="w-14 h-14 rounded-lg bg-tertiary/10 flex items-center justify-center border border-tertiary/20">
				<span class="material-symbols-outlined text-tertiary text-2xl drop-shadow-[0_0_8px_rgba(255,224,74,0.5)]">filter_alt</span>
			</div>
			<div>
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400 mb-1">Hasil Filter</p>
				<p class="text-3xl font-headline font-bold text-on-surface">{filteredTerms.length}</p>
			</div>
		</div>
	</div>

	<div class="bg-surface-container-lowest border border-outline-variant/40 rounded-xl overflow-hidden shadow-2xl">
		<div class="p-5 border-b border-outline-variant/30 flex flex-col md:flex-row gap-4 justify-between items-center bg-surface-container-low/50">
			<div class="relative w-full md:w-96">
				<span class="material-symbols-outlined absolute left-4 top-2.5 text-slate-500 text-sm">search</span>
				<input 
					type="text" 
					bind:value={searchQuery}
					placeholder="Cari istilah atau terjemahan..." 
					class="w-full bg-background border border-outline-variant/50 rounded-lg py-2.5 pl-10 pr-4 text-sm font-body text-on-surface placeholder:text-slate-500 focus:border-primary focus:ring-1 focus:ring-primary/30 outline-none transition-all" 
				/>
			</div>
			<div class="flex gap-3 w-full md:w-auto">
				<button 
					on:click={handleImport}
					class="flex items-center gap-2 px-4 py-2 border border-outline-variant hover:border-slate-400 bg-background rounded-lg text-sm text-on-surface transition-colors flex-1 md:flex-none justify-center font-medium"
				>
					<span class="material-symbols-outlined text-sm">upload</span>
					Impor JSON
				</button>
				<button 
					on:click={handleExport}
					class="flex items-center gap-2 px-4 py-2 border border-outline-variant hover:border-slate-400 bg-background rounded-lg text-sm text-on-surface transition-colors flex-1 md:flex-none justify-center font-medium"
				>
					<span class="material-symbols-outlined text-sm">download</span>
					Ekspor JSON
				</button>
			</div>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr class="border-b border-outline-variant/30 text-[10px] font-label font-bold text-primary tracking-widest uppercase">
						<th class="py-4 px-6">Istilah Asal</th>
						<th class="py-4 px-6 text-secondary">Terjemahan Target</th>
						<th class="py-4 px-6 text-slate-400">Catatan</th>
						<th class="py-4 px-6 text-slate-400 text-center">Case</th>
						<th class="py-4 px-6 text-slate-400 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-outline-variant/20">
					{#if isLoading}
						<tr>
							<td colspan="5" class="py-12 text-center text-slate-400">
								<span class="material-symbols-outlined animate-spin text-3xl mb-2">progress_activity</span>
								<p>Memuat glosarium...</p>
							</td>
						</tr>
					{:else if filteredTerms.length === 0}
						<tr>
							<td colspan="5" class="py-12 text-center text-slate-400">
								<span class="material-symbols-outlined text-3xl mb-2">search_off</span>
								<p>{searchQuery ? 'Tidak ada hasil pencarian' : 'Belum ada istilah'}</p>
							</td>
						</tr>
					{:else}
						{#each filteredTerms as term (term.id)}
							<tr class="hover:bg-surface-container/50 transition-colors group">
								<td class="py-5 px-6">
									<p class="font-bold text-on-surface text-sm">{term.sourceTerm}</p>
								</td>
								<td class="py-5 px-6">
									<p class="font-bold text-secondary text-sm drop-shadow-[0_0_5px_rgba(0,255,204,0.4)] neon-text-secondary">{term.targetTerm}</p>
								</td>
								<td class="py-5 px-6 text-xs text-slate-400 max-w-sm truncate">
									{term.notes || '-'}
								</td>
								<td class="py-5 px-6 text-center">
									{#if term.caseSensitive}
										<span class="inline-flex items-center px-2 py-1 rounded-full bg-tertiary/20 text-tertiary text-[10px] font-bold">
											Aa
										</span>
									{:else}
										<span class="text-slate-600">—</span>
									{/if}
								</td>
								<td class="py-5 px-6 text-right">
									<div class="flex items-center justify-end gap-3 opacity-50 group-hover:opacity-100 transition-opacity">
										<button 
											on:click={() => openEditModal(term)}
											class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-slate-800 text-slate-400 hover:text-white transition-colors"
										>
											<span class="material-symbols-outlined text-[18px]">edit</span>
										</button>
										<button 
											on:click={() => confirmDelete(term)}
											class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-red-900/30 text-slate-400 hover:text-red-400 transition-colors"
										>
											<span class="material-symbols-outlined text-[18px]">delete</span>
										</button>
									</div>
								</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>

		<div class="p-5 border-t border-outline-variant/30 flex justify-between items-center bg-surface-container-low/30 text-xs text-slate-400 font-label">
			<p>Menampilkan {filteredTerms.length} dari {stats.total} istilah</p>
		</div>
	</div>
</div>

<Modal bind:open={showAddModal} size="md" class="bg-dark-800">
	<svelte:fragment slot="header">
		<h3 class="text-lg font-bold text-white flex items-center gap-2">
			<span class="material-symbols-outlined text-neon-cyan">add_circle</span>
			Tambah Istilah Baru
		</h3>
	</svelte:fragment>
	
	<div class="space-y-4">
		<div>
			<label class="block text-sm font-medium text-gray-300 mb-1">Istilah Asal *</label>
			<input 
				type="text" 
				bind:value={newTerm.sourceTerm}
				placeholder="Contoh: Machine Learning"
				class="w-full px-3 py-2 rounded-lg bg-dark-700 border border-dark-600 text-white placeholder-gray-500 focus:border-neon-cyan"
			/>
		</div>
		<div>
			<label class="block text-sm font-medium text-gray-300 mb-1">Terjemahan Target *</label>
			<input 
				type="text" 
				bind:value={newTerm.targetTerm}
				placeholder="Contoh: Pembelajaran Mesin"
				class="w-full px-3 py-2 rounded-lg bg-dark-700 border border-dark-600 text-white placeholder-gray-500 focus:border-neon-cyan"
			/>
		</div>
		<div>
			<label class="block text-sm font-medium text-gray-300 mb-1">Catatan</label>
			<textarea 
				bind:value={newTerm.notes}
				rows="2"
				placeholder="Catatan tambahan (opsional)"
				class="w-full px-3 py-2 rounded-lg bg-dark-700 border border-dark-600 text-white placeholder-gray-500 focus:border-neon-cyan resize-none"
			></textarea>
		</div>
		<div class="flex items-center justify-between p-3 rounded-lg bg-dark-700/50">
			<div>
				<p class="text-sm text-gray-300">Case Sensitive</p>
				<p class="text-xs text-gray-500">Perhatikan huruf besar/kecil</p>
			</div>
			<Toggle bind:checked={newTerm.caseSensitive} color="cyan" />
		</div>
	</div>

	<svelte:fragment slot="footer">
		<Button color="alternative" on:click={() => showAddModal = false}>Batal</Button>
		<Button color="primary" on:click={addTerm} class="!bg-neon-cyan !text-dark-900">
			<span class="material-symbols-outlined text-sm mr-1">add</span>
			Tambah
		</Button>
	</svelte:fragment>
</Modal>

<Modal bind:open={showEditModal} size="md" class="bg-dark-800">
	<svelte:fragment slot="header">
		<h3 class="text-lg font-bold text-white flex items-center gap-2">
			<span class="material-symbols-outlined text-neon-cyan">edit</span>
			Edit Istilah
		</h3>
	</svelte:fragment>
	
	{#if editingTerm}
		<div class="space-y-4">
			<div>
				<label class="block text-sm font-medium text-gray-300 mb-1">Istilah Asal *</label>
				<input 
					type="text" 
					bind:value={editingTerm.sourceTerm}
					class="w-full px-3 py-2 rounded-lg bg-dark-700 border border-dark-600 text-white focus:border-neon-cyan"
				/>
			</div>
			<div>
				<label class="block text-sm font-medium text-gray-300 mb-1">Terjemahan Target *</label>
				<input 
					type="text" 
					bind:value={editingTerm.targetTerm}
					class="w-full px-3 py-2 rounded-lg bg-dark-700 border border-dark-600 text-white focus:border-neon-cyan"
				/>
			</div>
			<div>
				<label class="block text-sm font-medium text-gray-300 mb-1">Catatan</label>
				<textarea 
					bind:value={editingTerm.notes}
					rows="2"
					class="w-full px-3 py-2 rounded-lg bg-dark-700 border border-dark-600 text-white focus:border-neon-cyan resize-none"
				></textarea>
			</div>
			<div class="flex items-center justify-between p-3 rounded-lg bg-dark-700/50">
				<div>
					<p class="text-sm text-gray-300">Case Sensitive</p>
					<p class="text-xs text-gray-500">Perhatikan huruf besar/kecil</p>
				</div>
				<Toggle bind:checked={editingTerm.caseSensitive} color="cyan" />
			</div>
		</div>
	{/if}

	<svelte:fragment slot="footer">
		<Button color="alternative" on:click={() => { showEditModal = false; editingTerm = null; }}>Batal</Button>
		<Button color="primary" on:click={updateTerm} class="!bg-neon-cyan !text-dark-900">
			<span class="material-symbols-outlined text-sm mr-1">save</span>
			Simpan
		</Button>
	</svelte:fragment>
</Modal>

<Modal bind:open={showDeleteConfirm} size="sm" class="bg-dark-800">
	<svelte:fragment slot="header">
		<h3 class="text-lg font-bold text-red-400 flex items-center gap-2">
			<span class="material-symbols-outlined">warning</span>
			Konfirmasi Hapus
		</h3>
	</svelte:fragment>
	
	{#if deletingTerm}
		<p class="text-gray-300">
			Apakah Anda yakin ingin menghapus istilah <strong class="text-white">"{deletingTerm.sourceTerm}"</strong>?
		</p>
		<p class="text-sm text-gray-500 mt-2">Tindakan ini tidak dapat dibatalkan.</p>
	{/if}

	<svelte:fragment slot="footer">
		<Button color="alternative" on:click={() => { showDeleteConfirm = false; deletingTerm = null; }}>Batal</Button>
		<Button color="red" on:click={deleteTerm}>
			<span class="material-symbols-outlined text-sm mr-1">delete</span>
			Hapus
		</Button>
	</svelte:fragment>
</Modal>
