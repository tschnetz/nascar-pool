<script lang="ts">
	import '../app.css';
	import { page } from '$app/stores';

	interface Props {
		children: import('svelte').Snippet;
	}

	let { children }: Props = $props();

	// Admin mode check
	let isAdmin = $derived($page.url.searchParams.get('admin') === 'true');

	// Navigation items
	const navItems = [
		{ href: '/', label: 'Standings', icon: 'trophy' },
		{ href: '/summary', label: 'Summary', icon: 'chart' },
		{ href: '/races', label: 'Races', icon: 'flag' },
		{ href: '/schedule', label: 'Schedule', icon: 'calendar' },
		{ href: '/drivers', label: 'Drivers', icon: 'users' }
	];

	// Check if route is active
	function isActive(href: string): boolean {
		if (href === '/') {
			return $page.url.pathname === '/';
		}
		return $page.url.pathname.startsWith(href);
	}

	// Get href with admin param preserved
	function getHref(href: string): string {
		return isAdmin ? `${href}?admin=true` : href;
	}
</script>

<div class="min-h-screen bg-gray-100 pb-20 md:pb-0">
	<!-- Top header (simplified on mobile) -->
	<header class="bg-nascar-blue text-white sticky top-0 z-40 safe-top">
		<div class="px-4 h-14 flex items-center justify-between">
			<h1 class="text-lg font-bold">NASCAR Pool</h1>
			{#if isAdmin}
				<span class="text-xs bg-nascar-red px-2 py-1 rounded font-medium">ADMIN</span>
			{/if}
		</div>
	</header>

	<!-- Main content area with safe padding -->
	<main class="px-4 py-4 max-w-3xl mx-auto">
		{@render children()}
	</main>

	<!-- Bottom navigation (mobile-first, hidden on larger screens) -->
	<nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-50 safe-bottom md:hidden">
		<div class="flex justify-around items-center h-16">
			{#each navItems as item}
				<a
					href={getHref(item.href)}
					class="flex flex-col items-center justify-center w-full h-full min-h-[44px] transition-colors
						   {isActive(item.href) ? 'text-nascar-blue' : 'text-gray-500'}"
				>
					<!-- Simple SVG icons -->
					{#if item.icon === 'trophy'}
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
						</svg>
					{:else if item.icon === 'chart'}
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
						</svg>
					{:else if item.icon === 'flag'}
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9" />
						</svg>
					{:else if item.icon === 'calendar'}
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
						</svg>
					{:else if item.icon === 'users'}
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
								d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
						</svg>
					{/if}
					<span class="text-xs mt-1 font-medium">{item.label}</span>
				</a>
			{/each}
		</div>
	</nav>

	<!-- Desktop top navigation (hidden on mobile) -->
	<nav class="hidden md:block fixed top-14 left-0 right-0 bg-nascar-blue/95 backdrop-blur z-30">
		<div class="max-w-3xl mx-auto px-4">
			<div class="flex gap-1">
				{#each navItems as item}
					<a
						href={getHref(item.href)}
						class="px-4 py-3 text-sm font-medium transition-colors
							   {isActive(item.href)
								 ? 'text-white bg-white/20'
								 : 'text-white/70 hover:text-white hover:bg-white/10'}"
					>
						{item.label}
					</a>
				{/each}
			</div>
		</div>
	</nav>
</div>

<style>
	/* Safe area insets for notched phones */
	.safe-top {
		padding-top: env(safe-area-inset-top);
	}
	.safe-bottom {
		padding-bottom: env(safe-area-inset-bottom);
	}
</style>
