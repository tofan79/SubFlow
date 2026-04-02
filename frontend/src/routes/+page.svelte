<script lang="ts">
	import DropZone from '$lib/components/DropZone.svelte';
	import PipelineProgress from '$lib/components/PipelineProgress.svelte';
	import { onMount } from 'svelte';
	import { pipeline, isRunning } from '$lib/stores/pipeline';
	import { recentProjects } from '$lib/stores/project';
	import { GetStats } from '$lib/wails/index';
	import type { AppStats } from '$lib/wails/index';

	let stats: AppStats = { totalProjects: 128, totalSegments: 450, totalCharsProcessed: 12500, totalMinutesAsr: 42.5 };

	onMount(async () => {
		try {
			stats = await GetStats();
		} catch (e) {
			console.log('Using mock stats');
		}
	});

	const quickActions = [
		{ label: 'Proyek Baru', desc: 'Mulai dari nol', icon: 'add_box', color: 'text-primary', bgColor: 'bg-primary/10', hoverColor: 'group-hover:text-primary' },
		{ label: 'Buka Editor', desc: 'Lanjutkan pengeditan', icon: 'video_settings', color: 'text-secondary', bgColor: 'bg-secondary/10', hoverColor: 'group-hover:text-secondary' },
		{ label: 'Lihat Glosarium', desc: 'Atur istilah khusus', icon: 'translate', color: 'text-tertiary', bgColor: 'bg-tertiary/10', hoverColor: 'group-hover:text-tertiary' }
	];

	const mockProjects = [
		{ name: 'Iklan_Produk_Neon.mp4', time: 'Diedit 2 jam yang lalu', duration: '04:32', state: 'exported', img: 'https://lh3.googleusercontent.com/aida-public/AB6AXuDVQuEzUWt4nWriYCKOUlX2Yf1Ctguw-eR0v2so78vewqYx4iwlz-daOv-UWIBRduL2ZIfEDmG-cnX7Dk4Q7TZq_9K69ScUSilLtgU-rW6TbdchlRaUUdj2qz2-Hi6oZGPC8dRN8u9Ab5yelUIcCDMa01JIYGYcKdmUNb-43VEZ86_rale9oKnYQFO_KcYwh1C65u7rVBr3gps9uToNDSO3KEZe3U1LkaPlyiHEfRojeUA5y-X5afyfZPYrDTdFq0Qrrm8WibsK' },
		{ name: 'Tutorial_Blender_NextGen.mov', time: 'Diedit kemarin', duration: '16:07', state: 'qa_done', img: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAyRlXOIb5EC8rCjenmItJHG_143cKihuJARJsHOssfoqHA2w8it953pflPPoo-16876n_YGhK0tLwZmHlGMOzb6hlDm4OokLkKlYrwW60ZM1az__e46_TYO_goiZI87DtWJ0EWFqVYDvIa1IGslUaikr_20Y_zqtHOZIi-j-nELSzMTMrw6ZiEWiAr7cPZCBEQQ0qm0h07ANmM1uPHT-uiZG4luzTtmCWy85mdbjDKiaotOTF1YnDZzjLo3lRjfF7bV6dQT2Ze' },
		{ name: 'Wawancara_Eksklusif_01.mp4', time: 'Diedit 3 hari yang lalu', duration: '08:48', state: 'translated', img: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAjy-s1YH0Ca9V-rFh5wIDsLd7NF3i1VsyhjYxqrBBPCWgndlkF93KBXqEVdUajTcQV4PUrZTM3nbEHKxsXpJcSDTycRypzUQ6tyE3fPb1HCrV78jmpD7xtGZd4YsexrAKhW1ckUpnoDl3P2KOIaCs9S-gXrvHJpUYZluMkxvHgQTt4pCHcOi9Tw5MWrMnYS8nYYCxWFWT4izHFSNelIycqxMekKK0N78aLNSJ3-cIrZUzh6hiaBc9zttdprsvLmQMzdtfJ8Zpe' },
		{ name: 'Short_Movie_Trailer.mp4', time: 'Diedit minggu lalu', duration: '12:55', state: 'exported', img: 'https://lh3.googleusercontent.com/aida-public/AB6AXuA9hGoBJ7WvNP_afR0qboOUY8HTccqtLFQyfR2yrPlnvohnTR0QgdAlIHH5iwbBtgSF5eczZUtYyN_5L0_si3uM8rAPWICziG32bBakvzPEBxfmXLTJusqeE8Y5V9nTrNBfhjMu_7rr9em2KURnio0asEqw1LJrCB_hpnNhYPulTeGs-dxbVpZwhd26TvwSJxq6qfet2JPd2JBkCJVtnFJgC4Hr85AkisVIUMwg3tKet6qe4BPgGdjIFPnFCJBjgg8YzQ4ME7PZ' }
	];
</script>

<div class="max-w-7xl mx-auto space-y-8 z-10 relative">
	<!-- Hero Header & Drop Zone -->
	<section class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<div class="lg:col-span-2 relative group h-full block">
			<div class="absolute inset-0 bg-gradient-to-br from-primary/10 to-secondary/10 rounded-xl blur-xl transition-all group-hover:blur-2xl opacity-50 z-0"></div>
			<div class="h-full relative z-10">
				<DropZone />
			</div>
		</div>

		<!-- Quick Actions Grid -->
		<div class="grid grid-cols-1 gap-4">
			{#each quickActions as action}
				<div class="bg-surface-container p-5 rounded-xl border border-outline-variant/30 flex items-center gap-4 hover:bg-surface-container-high transition-colors cursor-pointer group hover:border-[#ff2d78]/30">
					<div class="w-12 h-12 {action.bgColor} rounded-lg flex items-center justify-center {action.color} group-hover:scale-110 transition-transform">
						<span class="material-symbols-outlined">{action.icon}</span>
					</div>
					<div>
						<p class="font-headline font-bold text-on-surface">{action.label}</p>
						<p class="text-xs text-slate-500 font-label">{action.desc}</p>
					</div>
				</div>
			{/each}
		</div>
	</section>

	<!-- Middle Section: Pipeline & Stats -->
	<section class="grid grid-cols-1 xl:grid-cols-12 gap-6">
		<!-- Active Pipeline (Bento Style) -->
		<div class="xl:col-span-8 space-y-4">
			<div class="flex items-center justify-between px-2">
				<h3 class="font-headline font-bold text-lg flex items-center gap-2">
					<span class="material-symbols-outlined text-secondary">analytics</span>
					Pipa Aktif
				</h3>
				<span class="text-xs font-label uppercase tracking-widest flex items-center gap-1 {$isRunning ? 'text-secondary neon-glow-secondary' : 'text-slate-500'}">
					{#if $isRunning}
						<span class="w-2 h-2 rounded-full bg-secondary animate-pulse ml-1 inline-block"></span> BERJALAN
					{:else}
						MENUNGGU
					{/if}
				</span>
			</div>
			<PipelineProgress />
		</div>

		<!-- System Stats Tiles -->
		<div class="xl:col-span-4 grid grid-cols-2 gap-4">
			<div class="bg-surface-container-low p-5 rounded-xl border border-outline-variant/10 flex flex-col justify-between h-32 hover:border-secondary/30 transition-all">
				<span class="material-symbols-outlined text-secondary/50 text-xl">folder_managed</span>
				<div>
					<p class="text-2xl font-headline font-extrabold text-on-surface">{stats.totalProjects}</p>
					<p class="text-[10px] font-label text-slate-500 uppercase">Total Proyek</p>
				</div>
			</div>
			<div class="bg-surface-container-low p-5 rounded-xl border border-outline-variant/10 flex flex-col justify-between h-32 hover:border-primary/30 transition-all cursor-pointer group">
				<span class="material-symbols-outlined text-primary/50 text-xl group-hover:scale-110 transition-transform">schedule</span>
				<div>
					<p class="text-2xl font-headline font-extrabold text-on-surface">{stats.totalMinutesAsr}</p>
					<p class="text-[10px] font-label text-slate-500 uppercase">Jam Transkripsi</p>
				</div>
			</div>
			
			<div class="col-span-2 bg-surface-container-low p-5 rounded-xl border border-outline-variant/10 hover:border-tertiary/30 transition-all">
				<div class="flex justify-between items-center mb-4">
					<p class="text-[10px] font-label text-slate-500 uppercase">Beban Sistem (GPU)</p>
					<span class="w-2 h-2 rounded-full bg-secondary animate-pulse shadow-[0_0_8px_#00ffcc]"></span>
				</div>
				<div class="flex items-center gap-4">
					<div class="flex-1 space-y-2">
						<div class="flex justify-between text-[10px]">
							<span class="text-slate-400">RTX 4090 v2</span>
							<span class="text-tertiary">68%</span>
						</div>
						<div class="w-full bg-surface-container-highest rounded-full h-1">
							<div class="h-full bg-tertiary shadow-[0_0_5px_#ffe04a]" style="width: 68%"></div>
						</div>
					</div>
					<div class="flex-1 space-y-2">
						<div class="flex justify-between text-[10px]">
							<span class="text-slate-400">VRAM</span>
							<span class="text-secondary">42%</span>
						</div>
						<div class="w-full bg-surface-container-highest rounded-full h-1">
							<div class="h-full bg-secondary shadow-[0_0_5px_#00ffcc]" style="width: 42%"></div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Bottom Section: Recent Projects -->
	<section class="space-y-4 pb-12 md:pb-0">
		<div class="flex items-center justify-between px-2">
			<h3 class="font-headline font-bold text-lg">Proyek Terbaru</h3>
			<a href="/projects" class="text-xs font-label text-primary hover:text-white transition-all neon-glow-primary">Lihat Semua</a>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			{#each ($recentProjects.length > 0 ? $recentProjects : mockProjects) as proj}
				<div class="bg-surface-container-low p-4 rounded-xl border border-outline-variant/20 hover:bg-surface-container hover:border-outline transition-all group cursor-pointer">
					<div class="aspect-video mb-3 rounded-lg overflow-hidden bg-surface-container-highest relative">
						<img alt="Project thumbnail" class="w-full h-full object-cover opacity-60 group-hover:opacity-100 transition-opacity duration-300 group-hover:scale-105" src={'img' in proj ? proj.img : 'https://lh3.googleusercontent.com/aida-public/AB6AXuDVQuEzUWt4nWriYCKOUlX2Yf1Ctguw-eR0v2so78vewqYx4iwlz-daOv-UWIBRduL2ZIfEDmG-cnX7Dk4Q7TZq_9K69ScUSilLtgU-rW6TbdchlRaUUdj2qz2-Hi6oZGPC8dRN8u9Ab5yelUIcCDMa01JIYGYcKdmUNb-43VEZ86_rale9oKnYQFO_KcYwh1C65u7rVBr3gps9uToNDSO3KEZe3U1LkaPlyiHEfRojeUA5y-X5afyfZPYrDTdFq0Qrrm8WibsK'}/>
						<div class="absolute inset-0 bg-gradient-to-t from-background to-transparent opacity-60"></div>
						{#if 'duration' in proj}
							<span class="absolute bottom-2 right-2 bg-black/70 text-[10px] px-2 py-0.5 rounded text-white font-label backdrop-blur-sm border border-white/10">{proj.duration}</span>
						{/if}
					</div>
					<h4 class="font-semibold text-on-surface text-sm mb-1 truncate group-hover:text-primary transition-colors">{proj.name}</h4>
					<p class="text-[10px] text-slate-500 font-label">{'time' in proj ? proj.time : proj.state}</p>
				</div>
			{/each}
		</div>
	</section>
</div>
