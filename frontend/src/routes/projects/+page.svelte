<script lang="ts">
	const projects = [
		{ name: 'Vlog_Tokyo_Neon_Final.mp4', size: '124.5 MB', duration: '14:20', status: 'Selesai', color: 'text-green-400 border-green-400/30 bg-green-400/10', lang: 'JP', date: '24 Okt 2023, 14:20', icon: 'description', isProcessing: false, isFailed: false },
		{ name: 'Wawancara_Startup_Tech_V2.mp4', size: '85.2 MB', duration: '05:45', status: 'Memproses', color: 'text-secondary border-secondary/30 bg-secondary/10', lang: 'ID', date: 'Sedang berjalan...', icon: 'settings', isProcessing: true, isFailed: false },
		{ name: 'Tutorial_Gaming_Epic_Fail.mkv', size: '2.4 GB', duration: '45:10', status: 'Gagal', color: 'text-error border-error/30 bg-error/10', lang: 'EN', date: '22 Okt 2023, 09:12', icon: 'error', isProcessing: false, isFailed: true },
		{ name: 'Podcast_Malam_Jumat_Eps10.mp3', size: '42.8 MB', duration: '56:30', status: 'Selesai', color: 'text-green-400 border-green-400/30 bg-green-400/10', lang: 'ID', date: '21 Okt 2023, 23:45', icon: 'description', isProcessing: false, isFailed: false }
	];
</script>

<div class="max-w-7xl mx-auto pb-12">
	<!-- Header -->
	<header class="flex flex-col md:flex-row justify-between md:items-end mb-10 gap-6">
		<div class="max-w-xl">
			<h1 class="text-4xl font-headline font-extrabold text-on-background tracking-tight mb-3">
				Proyek <span class="text-primary drop-shadow-[0_0_8px_rgba(255,45,120,0.6)] neon-glow-primary">Historis</span>
			</h1>
			<p class="text-sm font-body text-slate-400 leading-relaxed">Kelola dan tinjau semua transkripsi video masa lalu Anda dalam satu hub terpusat.</p>
		</div>
		<div class="flex items-center gap-3">
			<button class="flex items-center gap-2 px-6 py-2.5 bg-primary/10 border border-primary text-primary rounded-md font-headline font-bold text-[10px] tracking-widest uppercase transition-all shadow-[0_0_12px_rgba(255,45,120,0.2)] hover:shadow-[0_0_20px_rgba(255,45,120,0.4)] hover:bg-primary/20">
				<span class="material-symbols-outlined text-sm">add_circle</span>
				BARU
			</button>
			<button class="flex items-center gap-2 px-6 py-2.5 bg-transparent border border-outline-variant text-on-surface hover:bg-surface-variant rounded-md font-headline font-bold text-[10px] tracking-widest uppercase transition-all">
				<span class="material-symbols-outlined text-sm">download</span>
				EXPORT
			</button>
		</div>
	</header>

	<!-- Stat Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<!-- Total Proyek -->
		<div class="bg-surface-container-low border border-secondary/30 rounded-xl p-6 hover:border-secondary/50 transition-colors shadow-[inset_0_0_20px_rgba(0,255,204,0.02)]">
			<div class="flex justify-between items-start mb-6">
				<p class="text-[10px] font-label uppercase tracking-widest text-[#00ffcc] font-bold">TOTAL PROYEK</p>
				<span class="material-symbols-outlined text-[#00ffcc] text-xl drop-shadow-[0_0_8px_rgba(0,255,204,0.3)]">folder</span>
			</div>
			<div class="flex items-end gap-3">
				<p class="text-4xl font-headline font-bold text-on-surface leading-none">1,284</p>
			</div>
			<p class="text-[9px] font-label text-slate-500 mt-2">
				<span class="text-green-400 font-bold">+12%</span> dari bulan lalu
			</p>
		</div>
		
		<!-- Jam Transkripsi -->
		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-6 hover:border-primary/50 transition-colors">
			<div class="flex justify-between items-start mb-6">
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400">JAM TRANSKRIPSI</p>
				<span class="material-symbols-outlined text-primary text-xl drop-shadow-[0_0_8px_rgba(255,45,120,0.3)]">schedule</span>
			</div>
			<div class="flex items-end gap-1">
				<p class="text-4xl font-headline font-bold text-on-surface leading-none">45.2</p>
				<p class="text-sm font-headline text-slate-400 pb-1">jam</p>
			</div>
			<p class="text-[9px] font-label text-slate-500 mt-2">
				<span class="text-primary font-bold">+5.4j</span> minggu ini
			</p>
		</div>

		<!-- Penyimpanan -->
		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-6 hover:border-tertiary/50 transition-colors">
			<div class="flex justify-between items-start mb-6">
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400">PENYIMPANAN</p>
				<span class="material-symbols-outlined text-tertiary text-xl drop-shadow-[0_0_8px_rgba(255,224,74,0.3)]">storage</span>
			</div>
			<div class="flex items-end gap-3 mb-3">
				<p class="text-4xl font-headline font-bold text-on-surface leading-none">82%</p>
			</div>
			<div class="w-full bg-surface-container-highest rounded-full h-1.5 mb-2 overflow-hidden">
				<div class="h-full rounded-full w-[82%]" style="background: linear-gradient(90deg, #ffe04a, #ccb33b); box-shadow: 0 0 10px rgba(255,224,74,0.5);"></div>
			</div>
			<p class="text-[9px] font-label text-slate-500 uppercase">20.5 GB / 25 GB digunakan</p>
		</div>
	</div>

	<!-- Main List Section -->
	<div class="bg-surface-container-lowest border border-outline-variant/40 rounded-xl overflow-hidden shadow-2xl mb-8">
		<!-- Toolbar -->
		<div class="p-5 border-b border-outline-variant/30 flex justify-between items-center bg-surface-container-low/50">
			<div class="relative w-full md:w-96 flex-shrink-0">
				<span class="material-symbols-outlined absolute left-4 top-2.5 text-slate-500 text-sm">search</span>
				<input type="text" placeholder="Cari nama proyek atau ID..." class="w-full bg-surface-container-highest border-none rounded-md py-2.5 pl-10 pr-4 text-[11px] font-label uppercase tracking-widest text-on-surface placeholder:text-slate-500 focus:ring-1 focus:ring-secondary/30 outline-none transition-all" />
			</div>
			<div class="flex gap-2 w-full md:w-auto overflow-x-auto justify-end">
				<!-- Custom Select / Dropdowns mimicking the design -->
				<button class="hidden md:flex items-center justify-between gap-6 px-4 py-2 bg-surface-container border border-outline-variant/30 rounded-md text-[10px] font-label font-bold text-slate-300 uppercase tracking-widest min-w-[160px]">
					SEMUA STATUS <span class="material-symbols-outlined text-sm">expand_more</span>
				</button>
				<button class="hidden md:flex items-center justify-between gap-6 px-4 py-2 bg-surface-container border border-outline-variant/30 rounded-md text-[10px] font-label font-bold text-slate-300 uppercase tracking-widest min-w-[160px]">
					URUTKAN: TERBARU <span class="material-symbols-outlined text-sm">expand_more</span>
				</button>
				<button class="flex items-center justify-center p-2 bg-surface-container border border-outline-variant/30 rounded-md text-slate-300">
					<span class="material-symbols-outlined text-sm">filter_list</span>
				</button>
			</div>
		</div>

		<!-- Table -->
		<div class="overflow-x-auto">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr class="border-b border-outline-variant/30 text-[9px] font-label font-bold text-slate-500 tracking-[0.2em] uppercase bg-surface-container-low/20">
						<th class="py-4 px-6 font-medium">NAMA PROYEK</th>
						<th class="py-4 px-6 font-medium">STATUS</th>
						<th class="py-4 px-6 font-medium">BAHASA</th>
						<th class="py-4 px-6 font-medium">TANGGAL SELESAI</th>
						<th class="py-4 px-6 font-medium text-right">AKSI</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-outline-variant/20">
					{#each projects as proj}
						<tr class="hover:bg-surface-container-low/40 transition-colors group">
							<td class="py-4 px-6">
								<div class="flex items-center gap-4">
									<div class="w-8 h-8 rounded-md {proj.isFailed ? 'bg-error/10 text-error' : proj.isProcessing ? 'bg-secondary/10 text-secondary' : 'bg-primary/10 text-primary'} flex items-center justify-center shrink-0 border border-current/20">
										<span class="material-symbols-outlined text-[16px]">{proj.icon}</span>
									</div>
									<div>
										<p class="font-bold text-on-surface mb-0.5 text-sm truncate max-w-xs">{proj.name}</p>
										<div class="flex gap-2 items-center">
											<p class="text-[9px] font-mono text-slate-500">{proj.size}</p>
											<span class="w-1 h-1 rounded-full bg-slate-700"></span>
											<p class="text-[9px] font-mono text-slate-500">{proj.duration}</p>
										</div>
									</div>
								</div>
							</td>
							<td class="py-4 px-6">
								<span class="inline-flex py-1 px-3 rounded-full text-[8px] font-bold uppercase tracking-widest border {proj.color} shadow-[0_0_8px_currentColor] opacity-90 drop-shadow-[0_0_5px_currentColor]">
									<span class="mr-1.5 w-1 h-1 rounded-full bg-current inline-block {proj.isProcessing ? 'animate-pulse' : ''}"></span>
									{proj.status}
								</span>
							</td>
							<td class="py-4 px-6">
								<span class="text-xs font-mono text-slate-400 uppercase">{proj.lang}</span>
							</td>
							<td class="py-4 px-6">
								<p class="text-xs font-mono text-slate-400 {proj.isProcessing ? 'animate-pulse italic' : ''}">{proj.date}</p>
							</td>
							<td class="py-4 px-6 text-right">
								<div class="flex items-center justify-end gap-3 opacity-60 group-hover:opacity-100 transition-opacity">
									{#if !proj.isFailed}
									<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-slate-800 text-slate-400 hover:text-white transition-colors" title="Buka Proyek">
										<span class="material-symbols-outlined text-[18px]">open_in_new</span>
									</button>
									<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-slate-800 text-slate-400 hover:text-white transition-colors" title="Unduh">
										<span class="material-symbols-outlined text-[18px]">download</span>
									</button>
									{/if}
									{#if proj.isFailed}
										<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-primary/20 text-slate-400 hover:text-primary transition-colors" title="Coba Lagi">
											<span class="material-symbols-outlined text-[18px]">refresh</span>
										</button>
										<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-error/20 text-slate-400 hover:text-error transition-colors" title="Hapus">
											<span class="material-symbols-outlined text-[18px]">delete</span>
										</button>
									{:else if proj.isProcessing}
										<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-slate-800 text-slate-400 hover:text-white transition-colors" title="Batal">
											<span class="material-symbols-outlined text-[18px]">cancel</span>
										</button>
									{:else}
										<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-error/20 text-slate-400 hover:text-error transition-colors" title="Hapus">
											<span class="material-symbols-outlined text-[18px]">delete</span>
										</button>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Footer Pagination -->
		<div class="p-5 border-t border-outline-variant/30 flex flex-col md:flex-row justify-between items-center gap-4 text-[9px] text-slate-500 font-label tracking-widest uppercase">
			<p>Menampilkan 1-4 dari 1,284 proyek</p>
			
			<div class="flex items-center gap-1.5">
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant/50 hover:bg-surface-variant transition-colors group">
					<span class="material-symbols-outlined text-[14px] group-hover:text-white">chevron_left</span>
				</button>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-secondary text-secondary bg-secondary/10 font-bold drop-shadow-[0_0_5px_currentColor]">1</button>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant/50 hover:bg-surface-variant transition-colors hover:text-white">2</button>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant/50 hover:bg-surface-variant transition-colors hover:text-white">3</button>
				<span class="w-8 h-8 flex items-center justify-center text-outline">...</span>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant/50 hover:bg-surface-variant transition-colors hover:text-white">321</button>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant/50 hover:bg-surface-variant transition-colors group">
					<span class="material-symbols-outlined text-[14px] group-hover:text-white">chevron_right</span>
				</button>
			</div>
		</div>
	</div>

	<!-- Promo Banners -->
	<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
		<!-- Pink Promo -->
		<div class="rounded-xl overflow-hidden relative border border-primary/20 bg-surface-container-low hover:border-primary/50 transition-all p-8 flex flex-col items-start min-h-[200px] group cursor-pointer">
			<span class="material-symbols-outlined absolute right-[-20px] bottom-[-20px] text-[150px] text-primary/10 group-hover:scale-110 transition-transform duration-500 transform rotate-[-15deg]">rocket_launch</span>
			<div class="relative z-10">
				<h3 class="text-xl font-headline font-bold text-white mb-3">Butuh Kecepatan<br />Lebih?</h3>
				<p class="text-[11px] font-body text-slate-400 mb-6 max-w-[200px] leading-relaxed">Upgrade ke paket Pro untuk transkripsi 3x lebih cepat dan penyimpanan tanpa batas.</p>
				<button class="px-6 py-2.5 bg-primary hover:bg-[#ff4088] text-white rounded font-headline font-bold text-[10px] tracking-widest uppercase transition-all shadow-[0_0_15px_rgba(255,45,120,0.5)] active:scale-95">
					Upgrade Pro
				</button>
			</div>
		</div>

		<!-- Cyan Promo -->
		<div class="rounded-xl overflow-hidden relative border border-secondary/20 bg-surface-container-low hover:border-secondary/50 transition-all p-8 flex flex-col items-start min-h-[200px] group cursor-pointer">
			<span class="material-symbols-outlined absolute right-[-10px] bottom-[-20px] text-[150px] text-secondary/10 group-hover:scale-110 transition-transform duration-500">integration_instructions</span>
			<div class="relative z-10">
				<h3 class="text-xl font-headline font-bold text-white mb-3">Automasi API</h3>
				<p class="text-[11px] font-body text-slate-400 mb-6 max-w-[220px] leading-relaxed">Integrasikan SubFlow AI langsung ke workflow Anda menggunakan API Developer kami.</p>
				<button class="px-6 py-2.5 bg-secondary hover:bg-[#33ffe6] text-[#001a1a] rounded font-headline font-bold text-[10px] tracking-widest uppercase transition-all shadow-[0_0_15px_rgba(0,255,204,0.5)] active:scale-95">
					Dokumentasi
				</button>
			</div>
		</div>
	</div>
</div>
