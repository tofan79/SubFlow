<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Modal, Button, Alert } from 'flowbite-svelte';
	import { 
		GetProjects, 
		DeleteProject,
		GetStats,
		type ProjectInfo as WailsProjectInfo,
		type AppStats
	} from '$lib/wails';
	import { projectStore, type ProjectInfo } from '$lib/stores/project';

	let projects: WailsProjectInfo[] = [];
	let stats: AppStats | null = null;
	let searchQuery = '';
	let statusFilter = 'all';
	let sortBy = 'newest';
	let isLoading = true;
	let error: string | null = null;
	let success: string | null = null;

	let showDeleteConfirm = false;
	let deletingProject: WailsProjectInfo | null = null;

	$: filteredProjects = projects
		.filter(p => {
			const matchesSearch = p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
			                      p.sourcePath.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || p.state === statusFilter;
			return matchesSearch && matchesStatus;
		})
		.sort((a, b) => {
			if (sortBy === 'newest') return b.updatedAt - a.updatedAt;
			if (sortBy === 'oldest') return a.updatedAt - b.updatedAt;
			if (sortBy === 'name') return a.name.localeCompare(b.name);
			return 0;
		});

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		isLoading = true;
		error = null;
		try {
			const [projectsData, statsData] = await Promise.all([
				GetProjects(),
				GetStats()
			]);
			projects = projectsData;
			stats = statsData;
			// Map wails type to store type (state is string at runtime, matches union values)
			projectStore.setProjects(projectsData as unknown as ProjectInfo[]);
		} catch (e) {
			error = `Gagal memuat data: ${e}`;
			projects = [];
		} finally {
			isLoading = false;
		}
	}

	function getStatusConfig(state: string) {
		switch (state) {
			case 'completed':
				return { label: 'Selesai', color: 'text-green-400 border-green-400/30 bg-green-400/10', icon: 'check_circle' };
			case 'running':
				return { label: 'Memproses', color: 'text-secondary border-secondary/30 bg-secondary/10', icon: 'sync', animate: true };
			case 'paused':
				return { label: 'Dijeda', color: 'text-tertiary border-tertiary/30 bg-tertiary/10', icon: 'pause_circle' };
			case 'error':
				return { label: 'Gagal', color: 'text-red-400 border-red-400/30 bg-red-400/10', icon: 'error' };
			default:
				return { label: 'Idle', color: 'text-slate-400 border-slate-400/30 bg-slate-400/10', icon: 'schedule' };
		}
	}

	function formatDate(timestamp: number): string {
		if (!timestamp) return '-';
		const date = new Date(timestamp * 1000);
		return date.toLocaleDateString('id-ID', { 
			day: 'numeric', 
			month: 'short', 
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatDuration(minutes: number): string {
		if (minutes < 60) return `${minutes.toFixed(1)} mnt`;
		const hours = Math.floor(minutes / 60);
		const mins = minutes % 60;
		return `${hours}j ${mins.toFixed(0)}m`;
	}

	function openProject(project: WailsProjectInfo) {
		projectStore.selectProject(project.id);
		goto(`/editor?project=${project.id}`);
	}

	function confirmDelete(project: WailsProjectInfo) {
		deletingProject = project;
		showDeleteConfirm = true;
	}

	async function deleteProject() {
		if (!deletingProject) return;

		try {
			await DeleteProject(deletingProject.id);
			projects = projects.filter(p => p.id !== deletingProject!.id);
			projectStore.removeProject(deletingProject.id);
			showDeleteConfirm = false;
			deletingProject = null;
			success = 'Proyek berhasil dihapus';
			setTimeout(() => success = null, 3000);
		} catch (e) {
			error = `Gagal menghapus proyek: ${e}`;
		}
	}
</script>

<div class="max-w-7xl mx-auto pb-12">
	<header class="flex flex-col md:flex-row justify-between md:items-end mb-10 gap-6">
		<div class="max-w-xl">
			<h1 class="text-4xl font-headline font-extrabold text-on-background tracking-tight mb-3">
				Proyek <span class="text-primary drop-shadow-[0_0_8px_rgba(255,45,120,0.6)] neon-glow-primary">Historis</span>
			</h1>
			<p class="text-sm font-body text-slate-400 leading-relaxed">Kelola dan tinjau semua transkripsi video masa lalu Anda dalam satu hub terpusat.</p>
		</div>
		<div class="flex items-center gap-3">
			<a 
				href="/"
				class="flex items-center gap-2 px-6 py-2.5 bg-primary/10 border border-primary text-primary rounded-md font-headline font-bold text-[10px] tracking-widest uppercase transition-all shadow-[0_0_12px_rgba(255,45,120,0.2)] hover:shadow-[0_0_20px_rgba(255,45,120,0.4)] hover:bg-primary/20"
			>
				<span class="material-symbols-outlined text-sm">add_circle</span>
				BARU
			</a>
			<button 
				on:click={loadData}
				class="flex items-center gap-2 px-6 py-2.5 bg-transparent border border-outline-variant text-on-surface hover:bg-surface-variant rounded-md font-headline font-bold text-[10px] tracking-widest uppercase transition-all"
			>
				<span class="material-symbols-outlined text-sm">refresh</span>
				REFRESH
			</button>
		</div>
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
		<div class="bg-surface-container-low border border-secondary/30 rounded-xl p-6 hover:border-secondary/50 transition-colors shadow-[inset_0_0_20px_rgba(0,255,204,0.02)]">
			<div class="flex justify-between items-start mb-6">
				<p class="text-[10px] font-label uppercase tracking-widest text-[#00ffcc] font-bold">TOTAL PROYEK</p>
				<span class="material-symbols-outlined text-[#00ffcc] text-xl drop-shadow-[0_0_8px_rgba(0,255,204,0.3)]">folder</span>
			</div>
			<div class="flex items-end gap-3">
				<p class="text-4xl font-headline font-bold text-on-surface leading-none">
					{stats?.totalProjects.toLocaleString() ?? '0'}
				</p>
			</div>
		</div>
		
		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-6 hover:border-primary/50 transition-colors">
			<div class="flex justify-between items-start mb-6">
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400">TOTAL SEGMEN</p>
				<span class="material-symbols-outlined text-primary text-xl drop-shadow-[0_0_8px_rgba(255,45,120,0.3)]">subtitles</span>
			</div>
			<div class="flex items-end gap-1">
				<p class="text-4xl font-headline font-bold text-on-surface leading-none">
					{stats?.totalSegments.toLocaleString() ?? '0'}
				</p>
			</div>
		</div>

		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-6 hover:border-tertiary/50 transition-colors">
			<div class="flex justify-between items-start mb-6">
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400">MENIT ASR</p>
				<span class="material-symbols-outlined text-tertiary text-xl drop-shadow-[0_0_8px_rgba(255,224,74,0.3)]">mic</span>
			</div>
			<div class="flex items-end gap-1">
				<p class="text-4xl font-headline font-bold text-on-surface leading-none">
					{formatDuration(stats?.totalMinutesAsr ?? 0)}
				</p>
			</div>
		</div>
	</div>

	<div class="bg-surface-container-lowest border border-outline-variant/40 rounded-xl overflow-hidden shadow-2xl mb-8">
		<div class="p-5 border-b border-outline-variant/30 flex flex-col md:flex-row gap-4 justify-between items-center bg-surface-container-low/50">
			<div class="relative w-full md:w-96 flex-shrink-0">
				<span class="material-symbols-outlined absolute left-4 top-2.5 text-slate-500 text-sm">search</span>
				<input 
					type="text" 
					bind:value={searchQuery}
					placeholder="Cari nama proyek..." 
					class="w-full bg-surface-container-highest border-none rounded-md py-2.5 pl-10 pr-4 text-sm text-on-surface placeholder:text-slate-500 focus:ring-1 focus:ring-secondary/30 outline-none transition-all" 
				/>
			</div>
			<div class="flex gap-2 w-full md:w-auto">
				<select 
					bind:value={statusFilter}
					class="px-4 py-2 bg-surface-container border border-outline-variant/30 rounded-md text-sm text-slate-300"
				>
					<option value="all">Semua Status</option>
					<option value="completed">Selesai</option>
					<option value="running">Memproses</option>
					<option value="paused">Dijeda</option>
					<option value="error">Gagal</option>
					<option value="idle">Idle</option>
				</select>
				<select 
					bind:value={sortBy}
					class="px-4 py-2 bg-surface-container border border-outline-variant/30 rounded-md text-sm text-slate-300"
				>
					<option value="newest">Terbaru</option>
					<option value="oldest">Terlama</option>
					<option value="name">Nama A-Z</option>
				</select>
			</div>
		</div>

		<div class="overflow-x-auto">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr class="border-b border-outline-variant/30 text-[9px] font-label font-bold text-slate-500 tracking-[0.2em] uppercase bg-surface-container-low/20">
						<th class="py-4 px-6 font-medium">NAMA PROYEK</th>
						<th class="py-4 px-6 font-medium">STATUS</th>
						<th class="py-4 px-6 font-medium">SEGMEN</th>
						<th class="py-4 px-6 font-medium">TERAKHIR DIUBAH</th>
						<th class="py-4 px-6 font-medium text-right">AKSI</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-outline-variant/20">
					{#if isLoading}
						<tr>
							<td colspan="5" class="py-12 text-center text-slate-400">
								<span class="material-symbols-outlined animate-spin text-3xl mb-2">progress_activity</span>
								<p>Memuat proyek...</p>
							</td>
						</tr>
					{:else if filteredProjects.length === 0}
						<tr>
							<td colspan="5" class="py-12 text-center text-slate-400">
								<span class="material-symbols-outlined text-3xl mb-2">folder_off</span>
								<p>{searchQuery ? 'Tidak ada hasil pencarian' : 'Belum ada proyek'}</p>
							</td>
						</tr>
					{:else}
						{#each filteredProjects as project (project.id)}
							{@const status = getStatusConfig(project.state)}
							<tr class="hover:bg-surface-container-low/40 transition-colors group">
								<td class="py-4 px-6">
									<div class="flex items-center gap-4">
										<div class="w-8 h-8 rounded-md bg-primary/10 text-primary flex items-center justify-center shrink-0 border border-current/20">
											<span class="material-symbols-outlined text-[16px]">movie</span>
										</div>
										<div>
											<p class="font-bold text-on-surface mb-0.5 text-sm truncate max-w-xs">{project.name}</p>
											<p class="text-[9px] font-mono text-slate-500 truncate max-w-xs">{project.sourcePath}</p>
										</div>
									</div>
								</td>
								<td class="py-4 px-6">
									<span class="inline-flex items-center py-1 px-3 rounded-full text-[8px] font-bold uppercase tracking-widest border {status.color} shadow-[0_0_8px_currentColor] opacity-90">
										<span class="material-symbols-outlined text-[10px] mr-1 {status.animate ? 'animate-spin' : ''}">{status.icon}</span>
										{status.label}
									</span>
								</td>
								<td class="py-4 px-6">
									<span class="text-sm font-mono text-slate-400">{project.segmentCount}</span>
								</td>
								<td class="py-4 px-6">
									<p class="text-xs font-mono text-slate-400">{formatDate(project.updatedAt)}</p>
								</td>
								<td class="py-4 px-6 text-right">
									<div class="flex items-center justify-end gap-2 opacity-60 group-hover:opacity-100 transition-opacity">
										<button 
											on:click={() => openProject(project)}
											class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-secondary/20 text-slate-400 hover:text-secondary transition-colors" 
											title="Buka Editor"
										>
											<span class="material-symbols-outlined text-[18px]">edit</span>
										</button>
										<button 
											on:click={() => confirmDelete(project)}
											class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-red-900/30 text-slate-400 hover:text-red-400 transition-colors" 
											title="Hapus"
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

		<div class="p-5 border-t border-outline-variant/30 flex justify-between items-center text-xs text-slate-500 font-label">
			<p>Menampilkan {filteredProjects.length} dari {projects.length} proyek</p>
		</div>
	</div>
</div>

<Modal bind:open={showDeleteConfirm} size="sm" class="bg-dark-800">
	<svelte:fragment slot="header">
		<h3 class="text-lg font-bold text-red-400 flex items-center gap-2">
			<span class="material-symbols-outlined">warning</span>
			Konfirmasi Hapus
		</h3>
	</svelte:fragment>
	
	{#if deletingProject}
		<p class="text-gray-300">
			Apakah Anda yakin ingin menghapus proyek <strong class="text-white">"{deletingProject.name}"</strong>?
		</p>
		<p class="text-sm text-gray-500 mt-2">Semua data termasuk segmen dan hasil terjemahan akan dihapus. Tindakan ini tidak dapat dibatalkan.</p>
	{/if}

	<svelte:fragment slot="footer">
		<Button color="alternative" on:click={() => { showDeleteConfirm = false; deletingProject = null; }}>Batal</Button>
		<Button color="red" on:click={deleteProject}>
			<span class="material-symbols-outlined text-sm mr-1">delete</span>
			Hapus Proyek
		</Button>
	</svelte:fragment>
</Modal>
