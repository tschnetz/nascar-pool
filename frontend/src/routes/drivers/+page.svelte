<script lang="ts">
	import { onMount } from 'svelte';
	import { api, type Driver } from '$lib/api';

	let drivers = $state<Driver[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let manufacturerFilter = $state<string>('all');
	let teamFilter = $state<string>('all');
	let charteredFilter = $state<string>('chartered');

	onMount(async () => {
		try {
			drivers = await api.getDrivers();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load drivers';
		} finally {
			loading = false;
		}
	});

	// Get unique teams
	let teams = $derived(
		[...new Set(drivers.map(d => d.team_name).filter(Boolean))].sort() as string[]
	);

	// Get unique manufacturers
	let manufacturers = $derived(
		[...new Set(drivers.map(d => d.manufacturer).filter(Boolean))].sort() as string[]
	);

	function getManufacturerStyle(manufacturer?: string) {
		switch (manufacturer) {
			case 'Chevrolet': return 'bg-yellow-100 text-yellow-800';
			case 'Ford': return 'bg-blue-100 text-blue-800';
			case 'Toyota': return 'bg-red-100 text-red-800';
			default: return 'bg-gray-100 text-gray-600';
		}
	}

	function getManufacturerActiveStyle(manufacturer: string) {
		switch (manufacturer) {
			case 'Chevrolet': return 'bg-yellow-500 text-white';
			case 'Ford': return 'bg-blue-600 text-white';
			case 'Toyota': return 'bg-red-600 text-white';
			default: return 'bg-nascar-blue text-white';
		}
	}

	function getManufacturerCount(manufacturer: string) {
		return drivers.filter(d => d.manufacturer === manufacturer).length;
	}

	function getTeamCount(team: string) {
		return drivers.filter(d => d.team_name === team).length;
	}

	function getCharteredCount(chartered: boolean) {
		return drivers.filter(d => d.is_chartered === chartered).length;
	}

	// Filter drivers by manufacturer, team, and chartered status, sorted by car number
	let filteredDrivers = $derived(
		drivers
			.filter(d => manufacturerFilter === 'all' || d.manufacturer === manufacturerFilter)
			.filter(d => teamFilter === 'all' || d.team_name === teamFilter)
			.filter(d => charteredFilter === 'all' ||
				(charteredFilter === 'chartered' && d.is_chartered) ||
				(charteredFilter === 'non-chartered' && !d.is_chartered))
			.sort((a, b) => a.car_number - b.car_number)
	);
</script>

<svelte:head>
	<title>Drivers - NASCAR Pool</title>
</svelte:head>

<div>
	<h1 class="text-2xl font-bold text-gray-900 mb-4">2026 Drivers</h1>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-pulse text-gray-500">Loading drivers...</div>
		</div>
	{:else if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{error}
		</div>
	{:else}
		<!-- Chartered filter (top level) -->
		<div class="mb-3">
			<div class="text-xs text-gray-500 mb-1">Status</div>
			<div class="flex gap-2 overflow-x-auto pb-2">
				<button
					onclick={() => charteredFilter = 'all'}
					class="px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors
						   {charteredFilter === 'all' ? 'bg-nascar-blue text-white' : 'bg-gray-100 text-gray-600'}"
				>
					All ({drivers.length})
				</button>
				<button
					onclick={() => charteredFilter = 'chartered'}
					class="px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors
						   {charteredFilter === 'chartered' ? 'bg-green-600 text-white' : 'bg-green-100 text-green-800'}"
				>
					Chartered ({getCharteredCount(true)})
				</button>
				<button
					onclick={() => charteredFilter = 'non-chartered'}
					class="px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors
						   {charteredFilter === 'non-chartered' ? 'bg-orange-600 text-white' : 'bg-orange-100 text-orange-800'}"
				>
					Non-Chartered ({getCharteredCount(false)})
				</button>
			</div>
		</div>

		<!-- Manufacturer filter -->
		<div class="mb-3">
			<div class="text-xs text-gray-500 mb-1">Manufacturer</div>
			<div class="flex gap-2 overflow-x-auto pb-2">
				<button
					onclick={() => manufacturerFilter = 'all'}
					class="px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors
						   {manufacturerFilter === 'all' ? 'bg-nascar-blue text-white' : 'bg-gray-100 text-gray-600'}"
				>
					All ({drivers.length})
				</button>
				{#each manufacturers as mfr}
					<button
						onclick={() => manufacturerFilter = mfr}
						class="px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors
							   {manufacturerFilter === mfr ? getManufacturerActiveStyle(mfr) : getManufacturerStyle(mfr)}"
					>
						{mfr} ({getManufacturerCount(mfr)})
					</button>
				{/each}
			</div>
		</div>

		<!-- Team filter -->
		<div class="mb-4">
			<div class="text-xs text-gray-500 mb-1">Team</div>
			<div class="flex gap-2 overflow-x-auto pb-2">
				<button
					onclick={() => teamFilter = 'all'}
					class="px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors
						   {teamFilter === 'all' ? 'bg-nascar-blue text-white' : 'bg-gray-100 text-gray-600'}"
				>
					All
				</button>
				{#each teams as team}
					<button
						onclick={() => teamFilter = team}
						class="px-3 py-1.5 rounded-full text-sm font-medium whitespace-nowrap transition-colors
							   {teamFilter === team ? 'bg-gray-700 text-white' : 'bg-gray-100 text-gray-600'}"
					>
						{team} ({getTeamCount(team)})
					</button>
				{/each}
			</div>
		</div>

		<!-- Drivers list (numerical order) -->
		<div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
			<div class="divide-y divide-gray-100">
				{#each filteredDrivers as driver}
					<div class="p-4 flex items-center gap-4">
						<!-- Car number -->
						<div class="w-12 h-12 rounded-lg bg-nascar-blue text-white flex items-center justify-center font-bold text-lg flex-shrink-0">
							{driver.car_number}
						</div>
						<!-- Driver info -->
						<div class="flex-1 min-w-0">
							<div class="font-semibold text-gray-900">{driver.name}</div>
							<div class="text-sm text-gray-500">{driver.team_name || 'Unknown'}</div>
						</div>
						<!-- Manufacturer badge -->
						{#if driver.manufacturer}
							<span class="text-xs px-2 py-0.5 rounded flex-shrink-0 {getManufacturerStyle(driver.manufacturer)}">
								{driver.manufacturer}
							</span>
						{/if}
					</div>
				{/each}
			</div>
		</div>

		<!-- Summary -->
		<div class="mt-6 p-4 bg-gray-50 rounded-lg">
			<div class="text-sm text-gray-600 text-center">
				Showing {filteredDrivers.length} of {drivers.length} drivers
			</div>
		</div>
	{/if}
</div>
