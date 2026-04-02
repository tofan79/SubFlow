<script lang="ts">
	import { page } from '$app/stores';
	import '../app.css';

	const mainNavItems = [
		{ name: 'Home', href: '/', icon: 'home' },
		{ name: 'Editor', href: '/editor', icon: 'edit_note' },
		{ name: 'Glossary', href: '/glossary', icon: 'menu_book' },
		{ name: 'Projects', href: '/projects', icon: 'folder_open' },
		{ name: 'Settings', href: '/settings', icon: 'settings' }
	];

    let isSidebarExpanded = true;
</script>

<svelte:head>
	<link href="https://fonts.googleapis.com/css2?family=Sora:wght@400;600;700;800&family=Inter:wght@400;500;600&family=Space+Grotesk:wght@500;700&display=swap" rel="stylesheet" />
	<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet" />
</svelte:head>

<div class="w-full min-h-screen relative flex flex-col bg-background text-on-background font-body selection:bg-primary/30 scanline-bg">
	<!-- TopNavBar -->
	<header class="fixed top-0 left-0 right-0 z-50 bg-[#0a0a12]/90 backdrop-blur-md border-b border-primary/30 shadow-[0_0_12px_rgba(255,45,120,0.1)] flex justify-between items-center w-full px-6 py-3 h-16">
		<div class="flex items-center gap-6">
			<!-- Dedicated Sidebar Toggle Button -->
			<button on:click={() => isSidebarExpanded = !isSidebarExpanded} class="text-slate-300 hover:text-primary transition-colors p-1 rounded-md hover:bg-primary/10 flex items-center justify-center" title="Toggle Sidebar">
				<span class="material-symbols-outlined text-[26px]">menu</span>
			</button>
			
			<div class="text-xl font-black text-primary drop-shadow-[0_0_8px_rgba(255,45,120,0.8)] font-headline tracking-tight">
				SubFlow AI
			</div>
			
			<!-- Search Bar -->
			<div class="hidden lg:flex items-center bg-surface-container-highest/50 px-4 py-1.5 rounded-full border border-outline/20 ml-2">
				<span class="material-symbols-outlined text-slate-400 text-sm mr-2">search</span>
				<input class="bg-transparent border-none focus:ring-0 text-sm text-on-surface w-64 placeholder:text-slate-500" placeholder="Cari proyek..." type="text"/>
			</div>
		</div>
		<div class="flex items-center gap-6">
			<!-- Only Notification bell left -->
			<button class="text-slate-400 hover:text-secondary hover:drop-shadow-[0_0_8px_#00ffcc] transition-all duration-300 active:scale-95 relative">
				<span class="material-symbols-outlined">notifications</span>
				<span class="absolute top-0 right-0 w-2 h-2 rounded-full bg-primary animate-pulse shadow-[0_0_8px_#ff2d78]"></span>
			</button>
		</div>
	</header>

	<!-- SideNavBar -->
	<nav class="fixed left-0 top-0 h-full bg-background border-r border-secondary/20 shadow-[inset_0_0_15px_rgba(0,255,204,0.05)] flex-col pt-20 pb-6 hidden md:flex z-40 transition-all duration-300 {isSidebarExpanded ? 'w-64' : 'w-[5.5rem]'}">
		<div class="flex-1 space-y-2 overflow-hidden px-2">
			{#each mainNavItems as item}
				<a
					href={item.href}
					class="flex items-center gap-3 py-3 rounded-lg font-label text-sm font-medium transition-all overflow-hidden {$page.url.pathname === item.href ? 'text-primary bg-primary/10 drop-shadow-[0_0_8px_rgba(255,45,120,0.4)] shadow-[inset_0_0_12px_rgba(255,45,120,0.2)] border border-primary/20' : 'text-slate-400 hover:bg-secondary/10 hover:text-secondary border border-transparent'} {isSidebarExpanded ? 'px-4' : 'justify-center px-0'}"
					title={item.name}
				>
					<span
						class="material-symbols-outlined shrink-0"
						style={$page.url.pathname === item.href ? "font-variation-settings: 'FILL' 1;" : ""}
					>
						{item.icon}
					</span>
					
					<span class="whitespace-nowrap transition-all duration-300 {isSidebarExpanded ? 'opacity-100 max-w-full' : 'opacity-0 max-w-0 hidden'}">
						{item.name}
					</span>
				</a>
			{/each}
		</div>
		<div class="mt-auto border-t border-outline-variant/30 pt-4 px-2 overflow-hidden flex flex-col gap-2">
			<a class="flex items-center gap-3 py-3 rounded-lg text-slate-400 hover:bg-secondary/10 hover:text-secondary transition-colors font-label text-sm font-medium {isSidebarExpanded ? 'px-4' : 'justify-center px-0'}" href="/help" title="Help">
				<span class="material-symbols-outlined shrink-0">help</span>
				<span class="whitespace-nowrap transition-all duration-300 {isSidebarExpanded ? 'opacity-100 max-w-full' : 'opacity-0 max-w-0 hidden'}">Help</span>
			</a>
            <!-- Removed Logout Button per instruction -->
		</div>
	</nav>

	<!-- Main Content Canvas -->
	<main class="transition-all duration-300 {isSidebarExpanded ? 'md:ml-64' : 'md:ml-[5.5rem]'} pt-20 p-6 flex-1 w-full bg-transparent z-10 overflow-x-hidden min-h-screen">
		<slot />
	</main>

	<!-- BottomNavBar (Mobile Only) -->
	<nav class="fixed bottom-0 left-0 right-0 z-50 bg-[#0a0a12]/95 backdrop-blur-lg border-t border-primary/20 md:hidden flex justify-around items-center h-16 px-2">
		{#each mainNavItems as item}
			<a href={item.href} class="flex flex-col items-center justify-center {$page.url.pathname === item.href ? 'text-primary drop-shadow-[0_0_5px_currentColor]' : 'text-slate-400'}">
				<span class="material-symbols-outlined" style={$page.url.pathname === item.href ? "font-variation-settings: 'FILL' 1;" : ""}>{item.icon}</span>
				<span class="text-[10px] font-label mt-1">{item.name}</span>
			</a>
		{/each}
	</nav>
</div>
