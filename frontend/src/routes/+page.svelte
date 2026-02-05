<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';

	let standings = $state([]);
	let loading = $state(true);
	let error = $state(null);

	onMount(async () => {
		try {
			standings = await api.getStandings();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load standings';
		} finally {
			loading = false;
		}
	});

	// Get rank styling
	function getRankStyle(rank) {
		if (rank === 1) return 'bg-yellow-400 text-yellow-900';
		if (rank === 2) return 'bg-gray-300 text-gray-700';
		if (rank === 3) return 'bg-amber-600 text-amber-100';
		return 'bg-gray-100 text-gray-600';
	}
</script>

<svelte:head>
	<title>Standings - NASCAR Pool</title>
</svelte:head>

<div>
	<h1 class="text-2xl font-bold text-gray-900 mb-4">Season Standings</h1>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-pulse text-gray-500">Loading standings...</div>
		</div>
	{:else if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{error}
		</div>
	{:else if !standings || standings.length === 0}
		<div class="bg-amber-50 border border-amber-200 text-amber-700 px-4 py-3 rounded-lg">
			No standings yet. Create a race and enter results to see standings.
		</div>
	{:else}
		<!-- Mobile card layout -->
		<div class="space-y-3">
			{#each standings as standing}
				<div
					class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 flex items-center gap-4
						   {standing.rank === 1 ? 'ring-2 ring-yellow-400' : ''}"
				>
					<!-- Rank badge -->
					<div
						class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-lg flex-shrink-0
							   {getRankStyle(standing.rank)}"
					>
						{standing.rank}
					</div>

					<!-- Name and details -->
					<div class="flex-1 min-w-0">
						<div class="font-semibold text-gray-900 truncate">
							{standing.participant_name}
						</div>
						<div class="text-sm text-gray-500">
							{standing.races_completed} {standing.races_completed === 1 ? 'race' : 'races'}
						</div>
					</div>

					<!-- Points -->
					<div class="text-right flex-shrink-0">
						<div class="text-2xl font-bold text-nascar-blue">
							{standing.total_points}
						</div>
						<div class="text-xs text-gray-400 uppercase">pts</div>
					</div>
				</div>
			{/each}
		</div>

		<!-- Points breakdown legend -->
		<div class="mt-6 p-4 bg-gray-50 rounded-lg">
			<h3 class="text-sm font-medium text-gray-700 mb-2">Points per race</h3>
			<div class="grid grid-cols-2 gap-2 text-sm text-gray-600">
				<div>1st place: <span class="font-medium">135 pts</span></div>
				<div>2nd place: <span class="font-medium">25 pts</span></div>
				<div>Last place: <span class="font-medium">15 pts</span></div>
				<div>Stage wins: <span class="font-medium">25 pts each</span></div>
			</div>
			<div class="mt-2 text-xs text-gray-500">
				Special races (Daytona 500, etc.) = 2x points
			</div>
		</div>
	{/if}
</div>
