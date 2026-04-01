<script lang="ts">
	import { page } from '$app/stores';
	import { ui, isSidebarVisible } from '$lib/stores/ui';
	import '../app.css';

	const navItems = [
		{ name: 'Beranda', href: '/', icon: 'home' },
		{ name: 'Editor', href: '/editor', icon: 'edit_document' },
		{ name: 'Glosarium', href: '/glossary', icon: 'dictionary' },
		{ name: 'Proyek', href: '/projects', icon: 'folder_open' },
		{ name: 'Pengaturan', href: '/settings', icon: 'settings' }
	];
</script>

<div class="flex h-screen w-full bg-[#0a0a12] overflow-hidden text-white font-['Geist']">
	<!-- Sidebar -->
	{#if $isSidebarVisible}
		<aside class="w-64 flex flex-col bg-[#0f0f1a] border-r border-[#ff2d78]/20">
			<div class="h-16 flex items-center px-6 border-b border-[#00ffcc]/20">
				<h1 class="text-xl font-bold neon-text-primary uppercase tracking-widest">SubFlow</h1>
			</div>
			<nav class="flex-1 p-4 space-y-2">
				{#each navItems as item}
					<a
						href={item.href}
						class="flex items-center gap-3 px-4 py-3 rounded-lg transition-all duration-300 group hover:bg-[#141422] 
							{$page.url.pathname === item.href ? 'bg-[#141422] neon-border-primary' : 'border border-transparent'}"
					>
						<span class="material-symbols-outlined { $page.url.pathname === item.href ? 'text-[#ff2d78]' : 'text-gray-400 group-hover:text-[#00ffcc]' }">
							{item.icon}
						</span>
						<span class="{ $page.url.pathname === item.href ? 'text-white' : 'text-gray-400 group-hover:text-white' }">
							{item.name}
						</span>
					</a>
				{/each}
			</nav>
		</aside>
	{/if}

	<!-- Main Content -->
	<main class="flex-1 flex flex-col min-w-0">
		<header class="h-16 flex items-center justify-between px-6 bg-[#0f0f1a] border-b border-[#ffe04a]/20">
			<button class="text-gray-400 hover:text-[#00ffcc]" on:click={() => ui.toggleSidebar()}>
				<span class="material-symbols-outlined">menu</span>
			</button>
			<div class="flex items-center gap-4">
				<span class="text-sm text-gray-400">Status: <span class="text-[#00ffcc]">Online</span></span>
			</div>
		</header>
		<div class="flex-1 overflow-auto p-6 bg-[#0a0a12]">
			<slot />
		</div>
	</main>
</div>
