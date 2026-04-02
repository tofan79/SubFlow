<script lang="ts">
	const terms = [
		{ source: 'Synthesize', sourceLang: 'EN-US', target: 'Mensintesis', targetLang: 'ID-ID', note: 'Gunakan dalam konteks penggabungan data t...' },
		{ source: 'Neural Pathway', sourceLang: 'EN-US', target: 'Jalur Saraf Tiruan', targetLang: 'ID-ID', note: 'Istilah khusus untuk arsitektur deep learning.' },
		{ source: 'Cybernetic Core', sourceLang: 'EN-US', target: 'Inti Sibernetik', targetLang: 'ID-ID', note: "Jangan diterjemahkan sebagai 'Pusat Siber'." },
		{ source: 'Batch Processing', sourceLang: 'EN-US', target: 'Pemrosesan Tumpak', targetLang: 'ID-ID', note: 'Sesuai dengan standar teknis Bahasa Indonesia.' }
	];
</script>

<div class="max-w-7xl mx-auto pb-12">
	<!-- Header -->
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
		<button class="flex items-center gap-2 px-6 py-3 border border-primary text-primary hover:bg-primary/10 rounded-lg font-headline font-bold text-sm tracking-widest uppercase transition-all shadow-[0_0_12px_rgba(255,45,120,0.2)] hover:shadow-[0_0_20px_rgba(255,45,120,0.4)]">
			<span class="material-symbols-outlined text-lg">add</span>
			Tambah Istilah Baru
		</button>
	</header>

	<!-- Stat Cards -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<!-- Total Istilah -->
		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-5 flex items-center gap-5 hover:border-primary/50 transition-colors">
			<div class="w-14 h-14 rounded-lg bg-primary/10 flex items-center justify-center border border-primary/20">
				<span class="material-symbols-outlined text-primary text-2xl drop-shadow-[0_0_8px_rgba(255,45,120,0.5)]">translate</span>
			</div>
			<div>
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400 mb-1">Total Istilah</p>
				<p class="text-3xl font-headline font-bold text-on-surface">1,284</p>
			</div>
		</div>
		
		<!-- Akurasi AI -->
		<div class="bg-surface-container-low border border-outline-variant/30 rounded-xl p-5 flex items-center gap-5 hover:border-secondary/50 transition-colors">
			<div class="w-14 h-14 rounded-lg bg-secondary/10 flex items-center justify-center border border-secondary/20">
				<span class="material-symbols-outlined text-secondary text-2xl drop-shadow-[0_0_8px_rgba(0,255,204,0.5)]">bolt</span>
			</div>
			<div>
				<p class="text-[10px] font-label uppercase tracking-widest text-slate-400 mb-1">Akurasi AI</p>
				<p class="text-3xl font-headline font-bold text-on-surface">99.2%</p>
			</div>
		</div>

		<!-- Sync Status -->
		<div class="bg-surface-container-low border border-tertiary/40 rounded-xl p-5 flex items-center gap-5 relative overflow-hidden shadow-[inset_0_0_20px_rgba(255,224,74,0.05)]">
			<div class="w-14 h-14 rounded-lg bg-tertiary/10 flex items-center justify-center border border-tertiary/30">
				<span class="material-symbols-outlined text-tertiary text-2xl drop-shadow-[0_0_8px_rgba(255,224,74,0.6)] animate-pulse">sync</span>
			</div>
			<div class="flex-1">
				<div class="flex justify-between items-start w-full">
					<p class="text-[10px] font-label uppercase tracking-widest text-slate-400 mb-1">Sinkronisasi Terakhir</p>
					<span class="text-[8px] font-bold text-tertiary border border-tertiary/30 bg-tertiary/10 px-2 py-0.5 rounded-full uppercase tracking-widest">2 Menit Lalu</span>
				</div>
				<p class="text-sm font-headline font-bold text-on-surface mt-1">Sinkronisasi Cloud Aktif</p>
			</div>
		</div>
	</div>

	<!-- Main Table Section -->
	<div class="bg-surface-container-lowest border border-outline-variant/40 rounded-xl overflow-hidden shadow-2xl">
		<!-- Toolbar -->
		<div class="p-5 border-b border-outline-variant/30 flex flex-col md:flex-row gap-4 justify-between items-center bg-surface-container-low/50">
			<div class="relative w-full md:w-96">
				<span class="material-symbols-outlined absolute left-4 top-2.5 text-slate-500 text-sm">search</span>
				<input type="text" placeholder="Cari istilah atau terjemahan..." class="w-full bg-background border border-outline-variant/50 rounded-lg py-2.5 pl-10 pr-4 text-sm font-body text-on-surface placeholder:text-slate-500 focus:border-primary focus:ring-1 focus:ring-primary/30 outline-none transition-all" />
			</div>
			<div class="flex gap-3 w-full md:w-auto">
				<button class="flex items-center gap-2 px-4 py-2 border border-outline-variant hover:border-slate-400 bg-background rounded-lg text-sm text-on-surface transition-colors flex-1 md:flex-none justify-center font-medium">
					<span class="material-symbols-outlined text-sm">filter_list</span>
					Filter
				</button>
				<button class="flex items-center gap-2 px-4 py-2 border border-outline-variant hover:border-slate-400 bg-background rounded-lg text-sm text-on-surface transition-colors flex-1 md:flex-none justify-center font-medium">
					<span class="material-symbols-outlined text-sm">download</span>
					Ekspor
				</button>
			</div>
		</div>

		<!-- Table -->
		<div class="overflow-x-auto">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr class="border-b border-outline-variant/30 text-[10px] font-label font-bold text-primary tracking-widest uppercase">
						<th class="py-4 px-6">Istilah Asal</th>
						<th class="py-4 px-6 text-secondary">Terjemahan Target</th>
						<th class="py-4 px-6 text-slate-400">Catatan</th>
						<th class="py-4 px-6 text-slate-400 text-right">Aksi</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-outline-variant/20">
					{#each terms as term}
						<tr class="hover:bg-surface-container/50 transition-colors group">
							<td class="py-5 px-6">
								<p class="font-bold text-on-surface mb-1 text-sm">{term.source}</p>
								<p class="text-[9px] font-mono text-slate-500 uppercase">{term.sourceLang}</p>
							</td>
							<td class="py-5 px-6">
								<p class="font-bold text-secondary text-sm drop-shadow-[0_0_5px_rgba(0,255,204,0.4)] mb-1 neon-text-secondary">{term.target}</p>
								<p class="text-[9px] font-mono text-slate-500 uppercase">{term.targetLang}</p>
							</td>
							<td class="py-5 px-6 text-xs text-slate-400 max-w-sm truncate whitespace-normal leading-relaxed">
								{term.note}
							</td>
							<td class="py-5 px-6 text-right">
								<div class="flex items-center justify-end gap-3 opacity-50 group-hover:opacity-100 transition-opacity">
									<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-slate-800 text-slate-400 hover:text-white transition-colors">
										<span class="material-symbols-outlined text-[18px]">edit</span>
									</button>
									<button class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-red-900/30 text-slate-400 hover:text-red-400 transition-colors">
										<span class="material-symbols-outlined text-[18px]">delete</span>
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Footer Pagination -->
		<div class="p-5 border-t border-outline-variant/30 flex flex-col md:flex-row justify-between items-center gap-4 bg-surface-container-low/30 text-xs text-slate-400 font-label">
			<p>Menampilkan 1-4 dari 1,284 istilah</p>
			
			<div class="flex items-center gap-1">
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant hover:bg-surface-variant transition-colors">
					<span class="material-symbols-outlined text-[14px]">chevron_left</span>
				</button>
				<button class="w-8 h-8 flex items-center justify-center rounded-md bg-primary text-white font-bold drop-shadow-[0_0_8px_rgba(255,45,120,0.5)]">1</button>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant hover:bg-surface-variant transition-colors hover:text-white">2</button>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant hover:bg-surface-variant transition-colors hover:text-white">3</button>
				<span class="w-8 h-8 flex items-center justify-center">...</span>
				<button class="w-8 h-8 flex items-center justify-center rounded-md border border-outline-variant hover:bg-surface-variant transition-colors">
					<span class="material-symbols-outlined text-[14px]">chevron_right</span>
				</button>
			</div>
		</div>
	</div>

	<!-- Promo Banner -->
	<div class="mt-8 rounded-xl overflow-hidden relative border border-secondary/20 h-40 group cursor-pointer">
		<div class="absolute inset-0 bg-surface-container mix-blend-multiply"></div>
		<!-- Mocked monitors background using a subtle gradient representation to match the dark look -->
		<div class="absolute inset-0 bg-[radial-gradient(ellipse_at_center,_var(--tw-gradient-stops))] from-surface-container-highest via-background to-background opacity-80"></div>
		<div class="absolute inset-x-0 bottom-0 h-1/2 bg-gradient-to-t from-background to-transparent z-10"></div>
		
		<div class="relative z-20 h-full p-8 flex flex-col justify-center">
			<h3 class="text-3xl font-headline font-extrabold text-white mb-2">Optimalkan Alur Kerja Anda</h3>
			<p class="text-[10px] font-label font-bold uppercase tracking-[0.2em] text-secondary drop-shadow-[0_0_8px_rgba(0,255,204,0.6)] animate-pulse">
				Sistem AI Memantau 24/7
			</p>
		</div>
	</div>
</div>
