<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api, type Race, type Driver } from '$lib/api';

	// Admin mode check
	let isAdmin = $derived($page.url.searchParams.get('admin') === 'true');

	let races = $state<Race[]>([]);
	let drivers = $state<Driver[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	// New race form
	let showForm = $state(false);
	let newRace = $state({
		name: '',
		race_number: 1,
		date: '',
		is_special_race: false
	});
	let submitting = $state(false);

	onMount(async () => {
		try {
			[races, drivers] = await Promise.all([api.getRaces(), api.getDrivers()]);
			// Set default race number to next in sequence
			newRace.race_number = races.length + 1;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load data';
		} finally {
			loading = false;
		}
	});

	async function createRace() {
		submitting = true;
		error = null;
		try {
			await api.createRace(newRace);
			races = await api.getRaces();
			showForm = false;
			newRace = { name: '', race_number: races.length + 1, date: '', is_special_race: false };
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create race';
		} finally {
			submitting = false;
		}
	}

	function getStatusStyle(status: string) {
		switch (status) {
			case 'completed':
				return 'bg-green-100 text-green-700';
			case 'in_progress':
				return 'bg-amber-100 text-amber-700';
			default:
				return 'bg-gray-100 text-gray-600';
		}
	}

	function getStatusLabel(status: string) {
		switch (status) {
			case 'completed': return 'Complete';
			case 'in_progress': return 'In Progress';
			default: return 'Upcoming';
		}
	}
</script>

<svelte:head>
	<title>Races - NASCAR Pool</title>
</svelte:head>

<div>
	<div class="flex items-center justify-between mb-4">
		<h1 class="text-2xl font-bold text-gray-900">Races</h1>
		{#if isAdmin}
			<button
				onclick={() => (showForm = !showForm)}
				class="min-h-[44px] px-4 py-2 bg-nascar-red text-white font-medium rounded-lg
					   active:bg-red-700 transition-colors"
			>
				{showForm ? 'Cancel' : '+ Add Race'}
			</button>
		{/if}
	</div>

	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4 flex items-center justify-between">
			<span>{error}</span>
			<button onclick={() => (error = null)} class="p-1 font-bold">&times;</button>
		</div>
	{/if}

	<!-- New Race Form (admin only) -->
	{#if isAdmin}
	{#if showForm}
		<div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 mb-4">
			<h2 class="text-lg font-bold text-gray-900 mb-4">New Race</h2>
			<form onsubmit={(e) => { e.preventDefault(); createRace(); }} class="space-y-4">
				<div>
					<label for="name" class="block text-sm font-medium text-gray-700 mb-1">
						Race Name
					</label>
					<input
						type="text"
						id="name"
						bind:value={newRace.name}
						required
						placeholder="e.g., Daytona 500"
						class="w-full min-h-[44px] px-3 py-2 border border-gray-300 rounded-lg
							   focus:ring-2 focus:ring-nascar-blue focus:border-transparent"
					/>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="race_number" class="block text-sm font-medium text-gray-700 mb-1">
							Race #
						</label>
						<input
							type="number"
							id="race_number"
							bind:value={newRace.race_number}
							min="1"
							required
							class="w-full min-h-[44px] px-3 py-2 border border-gray-300 rounded-lg
								   focus:ring-2 focus:ring-nascar-blue focus:border-transparent"
						/>
					</div>
					<div>
						<label for="date" class="block text-sm font-medium text-gray-700 mb-1">
							Date
						</label>
						<input
							type="date"
							id="date"
							bind:value={newRace.date}
							class="w-full min-h-[44px] px-3 py-2 border border-gray-300 rounded-lg
								   focus:ring-2 focus:ring-nascar-blue focus:border-transparent"
						/>
					</div>
				</div>

				<label class="flex items-center gap-3 min-h-[44px] cursor-pointer">
					<input
						type="checkbox"
						bind:checked={newRace.is_special_race}
						class="w-5 h-5 text-nascar-red rounded focus:ring-nascar-red"
					/>
					<span class="text-gray-700">
						Special Race <span class="text-gray-500">(2x points)</span>
					</span>
				</label>

				<button
					type="submit"
					disabled={submitting}
					class="w-full min-h-[48px] bg-nascar-blue text-white font-medium rounded-lg
						   active:bg-blue-800 transition-colors disabled:opacity-50"
				>
					{submitting ? 'Creating...' : 'Create Race'}
				</button>
			</form>
		</div>
	{/if}
	{/if}

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-pulse text-gray-500">Loading races...</div>
		</div>
	{:else if !races || races.length === 0}
		<div class="bg-amber-50 border border-amber-200 text-amber-700 px-4 py-3 rounded-lg">
			No races yet. Tap "+ Add Race" to create your first race.
		</div>
	{:else}
		<div class="space-y-3">
			{#each races as race}
				<a
					href="/races/{race.id}"
					class="block bg-white rounded-xl shadow-sm border border-gray-100 p-4
						   active:bg-gray-50 transition-colors"
				>
					<div class="flex items-start justify-between gap-3">
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2 flex-wrap">
								<span class="text-xs font-medium text-gray-400">#{race.race_number}</span>
								<h2 class="font-semibold text-gray-900">{race.name}</h2>
								{#if race.is_special_race}
									<span class="text-xs font-bold bg-yellow-400 text-yellow-900 px-2 py-0.5 rounded">
										2X
									</span>
								{/if}
							</div>
							{#if race.date}
								<p class="text-sm text-gray-500 mt-1">{race.date}</p>
							{/if}
						</div>
						<span class="text-xs font-medium px-2 py-1 rounded-full flex-shrink-0 {getStatusStyle(race.status)}">
							{getStatusLabel(race.status)}
						</span>
					</div>

					{#if race.status === 'completed' && race.results}
						{@const winner = race.results.find(r => r.is_first_place)}
						{#if winner}
							<div class="mt-3 pt-3 border-t border-gray-100">
								<div class="flex items-center gap-2 text-sm">
									<span class="text-gray-500">Winner:</span>
									<span class="font-medium text-gray-900">#{winner.car_number} {winner.driver_name}</span>
								</div>
							</div>
						{/if}
					{/if}
				</a>
			{/each}
		</div>
	{/if}
</div>
