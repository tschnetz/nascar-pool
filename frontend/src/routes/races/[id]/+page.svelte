<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api.js';

	// Admin mode check
	let isAdmin = $derived($page.url.searchParams.get('admin') === 'true');

	let race = $state(null);
	let teams = $state([]);
	let drivers = $state([]);
	let loading = $state(true);
	let error = $state(null);
	let actionLoading = $state(false);

	// Edit form
	let showEditForm = $state(false);
	let editData = $state({
		name: '',
		date: '',
		is_special_race: false,
		extra_drivers: ''
	});

	// Results form - uses car numbers
	let showResultsForm = $state(false);
	let results = $state({
		first_place_car_number: 0,
		second_place_car_number: 0,
		last_place_car_number: 0,
		stage1_winner_car_number: 0,
		stage2_winner_car_number: 0
	});

	$effect(() => {
		loadData();
	});

	async function loadData() {
		const id = Number($page.params.id);
		if (isNaN(id)) {
			error = 'Invalid race ID';
			loading = false;
			return;
		}

		try {
			const [raceData, driversData] = await Promise.all([api.getRace(id), api.getDrivers()]);
			race = raceData;
			drivers = driversData;

			if (race.status !== 'upcoming') {
				teams = await api.getRaceTeams(id);
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load race data';
		} finally {
			loading = false;
		}
	}

	async function generateTeams() {
		if (!race) return;
		actionLoading = true;
		error = null;
		try {
			await api.generateTeams(race.id);
			await loadData();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to generate teams';
		} finally {
			actionLoading = false;
		}
	}

	async function submitResults() {
		if (!race) return;
		actionLoading = true;
		error = null;
		try {
			await api.enterRaceResults(race.id, results);
			showResultsForm = false;
			await loadData();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to enter results';
		} finally {
			actionLoading = false;
		}
	}

	function openEditForm() {
		if (!race) return;
		editData = {
			name: race.name,
			date: race.date || '',
			is_special_race: race.is_special_race,
			extra_drivers: race.extra_drivers || ''
		};
		showEditForm = true;
	}

	async function submitEdit() {
		if (!race) return;
		actionLoading = true;
		error = null;
		try {
			await api.updateRace(race.id, {
				name: editData.name,
				race_number: race.race_number,
				date: editData.date,
				is_special_race: editData.is_special_race,
				extra_drivers: editData.extra_drivers
			});
			showEditForm = false;
			await loadData();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to update race';
		} finally {
			actionLoading = false;
		}
	}

	function getStatusStyle(status) {
		switch (status) {
			case 'completed': return 'bg-green-100 text-green-700';
			case 'in_progress': return 'bg-amber-100 text-amber-700';
			default: return 'bg-gray-100 text-gray-600';
		}
	}

	function getStatusLabel(status) {
		switch (status) {
			case 'completed': return 'Complete';
			case 'in_progress': return 'In Progress';
			default: return 'Upcoming';
		}
	}

	function getResultByType(type) {
		if (!race?.results) return undefined;
		return race.results.find(r => {
			switch (type) {
				case 'first': return r.is_first_place;
				case 'second': return r.is_second_place;
				case 'last': return r.is_last_place;
				case 'stage1': return r.is_stage1_winner;
				case 'stage2': return r.is_stage2_winner;
			}
		});
	}
</script>

<svelte:head>
	<title>{race?.name || 'Race'} - NASCAR Pool</title>
</svelte:head>

<div>
	<!-- Back button -->
	<a
		href="/races"
		class="inline-flex items-center gap-1 text-gray-600 mb-4 min-h-[44px]"
	>
		<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
		</svg>
		Back
	</a>

	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4 flex items-center justify-between">
			<span>{error}</span>
			<button onclick={() => (error = null)} class="p-1 font-bold">&times;</button>
		</div>
	{/if}

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-pulse text-gray-500">Loading race...</div>
		</div>
	{:else if race}
		{@const hasRollover = race.rollover_first > 0 || race.rollover_second > 0 || race.rollover_last > 0 || race.rollover_stage1 > 0 || race.rollover_stage2 > 0}
		{@const mult = race.is_special_race ? 2 : 1}

		<!-- Rollover alert -->
		{#if hasRollover && race.status !== 'completed'}
			<div class="bg-amber-50 border border-amber-200 text-amber-800 px-4 py-3 rounded-lg mb-4">
				<div class="font-medium mb-1">Rollover Points Available!</div>
				<div class="text-sm space-x-3">
					{#if race.rollover_first > 0}<span>1st: +{race.rollover_first * mult}</span>{/if}
					{#if race.rollover_second > 0}<span>2nd: +{race.rollover_second * mult}</span>{/if}
					{#if race.rollover_last > 0}<span>Last: +{race.rollover_last * mult}</span>{/if}
					{#if race.rollover_stage1 > 0}<span>S1: +{race.rollover_stage1 * mult}</span>{/if}
					{#if race.rollover_stage2 > 0}<span>S2: +{race.rollover_stage2 * mult}</span>{/if}
				</div>
			</div>
		{/if}

		<!-- Race header card -->
		<div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 mb-4">
			<div class="flex items-start justify-between gap-3 mb-3">
				<div>
					<div class="flex items-center gap-2 flex-wrap">
						<span class="text-sm text-gray-400">#{race.race_number}</span>
						{#if race.is_special_race}
							<span class="text-xs font-bold bg-yellow-400 text-yellow-900 px-2 py-0.5 rounded">
								2X POINTS
							</span>
						{/if}
					</div>
					<h1 class="text-xl font-bold text-gray-900 mt-1">{race.name}</h1>
					{#if race.date}
						<p class="text-sm text-gray-500 mt-1">{race.date}</p>
					{/if}
				</div>
				<div class="flex items-center gap-2">
					{#if isAdmin}
						<button
							onclick={openEditForm}
							class="p-2 text-gray-400 hover:text-gray-600"
							title="Edit race"
						>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
									d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
							</svg>
						</button>
					{/if}
					<span class="text-xs font-medium px-2 py-1 rounded-full {getStatusStyle(race.status)}">
						{getStatusLabel(race.status)}
					</span>
				</div>
			</div>

			<!-- Action buttons (admin only) -->
			{#if isAdmin}
				{#if race.status === 'upcoming'}
					<button
						onclick={generateTeams}
						disabled={actionLoading}
						class="w-full min-h-[48px] bg-nascar-red text-white font-medium rounded-lg
							   active:bg-red-700 transition-colors disabled:opacity-50"
					>
						{actionLoading ? 'Generating...' : 'Generate Teams'}
					</button>
				{:else if race.status === 'in_progress'}
					<button
						onclick={() => (showResultsForm = !showResultsForm)}
						class="w-full min-h-[48px] bg-nascar-blue text-white font-medium rounded-lg
							   active:bg-blue-800 transition-colors"
					>
						{showResultsForm ? 'Cancel' : 'Enter Results'}
					</button>
				{/if}
			{/if}

			<!-- Worksheet button (visible to all when teams exist) -->
			{#if race.status !== 'upcoming'}
				<a
					href="/races/{race.id}/worksheet"
					class="block w-full min-h-[48px] mt-3 bg-gray-100 text-gray-700 font-medium rounded-lg
						   text-center leading-[48px] hover:bg-gray-200 transition-colors"
				>
					Pool Results / Print
				</a>
			{/if}
		</div>

		<!-- Edit race form (admin only) -->
		{#if isAdmin}
		{#if showEditForm}
			<div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 mb-4">
				<h2 class="text-lg font-bold text-gray-900 mb-4">Edit Race</h2>
				<form onsubmit={(e) => { e.preventDefault(); submitEdit(); }} class="space-y-4">
					<div>
						<label for="edit-name" class="block text-sm font-medium text-gray-700 mb-1">
							Race Name
						</label>
						<input
							type="text"
							id="edit-name"
							bind:value={editData.name}
							required
							class="w-full min-h-[44px] px-3 py-2 border border-gray-300 rounded-lg
								   focus:ring-2 focus:ring-nascar-blue focus:border-transparent"
						/>
					</div>
					<div>
						<label for="edit-date" class="block text-sm font-medium text-gray-700 mb-1">
							Date
						</label>
						<input
							type="date"
							id="edit-date"
							bind:value={editData.date}
							class="w-full min-h-[44px] px-3 py-2 border border-gray-300 rounded-lg
								   focus:ring-2 focus:ring-nascar-blue focus:border-transparent"
						/>
					</div>
					<div>
						<label for="edit-extra" class="block text-sm font-medium text-gray-700 mb-1">
							Extra Drivers <span class="text-gray-400">(not in pool)</span>
						</label>
						<input
							type="text"
							id="edit-extra"
							bind:value={editData.extra_drivers}
							placeholder="e.g., Allgaier, Johnson, Castroneves"
							class="w-full min-h-[44px] px-3 py-2 border border-gray-300 rounded-lg
								   focus:ring-2 focus:ring-nascar-blue focus:border-transparent"
						/>
					</div>
					<label class="flex items-center gap-3 min-h-[44px] cursor-pointer">
						<input
							type="checkbox"
							bind:checked={editData.is_special_race}
							class="w-5 h-5 text-nascar-red rounded focus:ring-nascar-red"
						/>
						<span class="text-gray-700">
							Special Race <span class="text-gray-500">(2x points)</span>
						</span>
					</label>
					<div class="flex gap-3">
						<button
							type="button"
							onclick={() => showEditForm = false}
							class="flex-1 min-h-[48px] bg-gray-200 text-gray-700 font-medium rounded-lg
								   active:bg-gray-300 transition-colors"
						>
							Cancel
						</button>
						<button
							type="submit"
							disabled={actionLoading}
							class="flex-1 min-h-[48px] bg-nascar-blue text-white font-medium rounded-lg
								   active:bg-blue-800 transition-colors disabled:opacity-50"
						>
							{actionLoading ? 'Saving...' : 'Save Changes'}
						</button>
					</div>
				</form>
			</div>
		{/if}
		{/if}

		<!-- Results entry form (admin only) -->
		{#if isAdmin}
		{#if showResultsForm}
			{@const multiplier = race.is_special_race ? 2 : 1}
			{@const fields = [
				{ id: 'first', label: '1st Place', base: 135, rollover: race.rollover_first, bind: () => results.first_place_car_number, set: (v) => results.first_place_car_number = v },
				{ id: 'second', label: '2nd Place', base: 25, rollover: race.rollover_second, bind: () => results.second_place_car_number, set: (v) => results.second_place_car_number = v },
				{ id: 'last', label: 'Last Place', base: 15, rollover: race.rollover_last, bind: () => results.last_place_car_number, set: (v) => results.last_place_car_number = v },
				{ id: 'stage1', label: 'Stage 1 Winner', base: 25, rollover: race.rollover_stage1, bind: () => results.stage1_winner_car_number, set: (v) => results.stage1_winner_car_number = v },
				{ id: 'stage2', label: 'Stage 2 Winner', base: 25, rollover: race.rollover_stage2, bind: () => results.stage2_winner_car_number, set: (v) => results.stage2_winner_car_number = v },
			]}
			<div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 mb-4">
				<h2 class="text-lg font-bold text-gray-900 mb-1">Enter Results</h2>
				<p class="text-sm text-gray-500 mb-4">Select by car number</p>

				<form onsubmit={(e) => { e.preventDefault(); submitResults(); }} class="space-y-4">
					{#each fields as field}
						{@const total = (field.base + field.rollover) * multiplier}
						<div>
							<label for={field.id} class="block text-sm font-medium text-gray-700 mb-1">
								{field.label}
								<span class="text-gray-400">
									({total} pts{field.rollover > 0 ? ` incl. ${field.rollover * multiplier} rollover` : ''})
								</span>
							</label>
							<select
								id={field.id}
								value={field.bind()}
								onchange={(e) => field.set(Number(e.currentTarget.value))}
								required
								class="w-full min-h-[44px] px-3 py-2 border border-gray-300 rounded-lg
									   focus:ring-2 focus:ring-nascar-blue focus:border-transparent bg-white"
							>
								<option value={0}>Select car...</option>
								{#each drivers as driver}
									<option value={driver.car_number}>#{driver.car_number} {driver.name}</option>
								{/each}
							</select>
						</div>
					{/each}

					<button
						type="submit"
						disabled={actionLoading}
						class="w-full min-h-[48px] bg-green-600 text-white font-medium rounded-lg
							   active:bg-green-700 transition-colors disabled:opacity-50"
					>
						{actionLoading ? 'Saving...' : 'Save Results'}
					</button>
				</form>
			</div>
		{/if}
		{/if}

		<!-- Race results display -->
		{#if race.status === 'completed' && race.results}
			{@const first = getResultByType('first')}
			{@const second = getResultByType('second')}
			{@const last = getResultByType('last')}
			{@const stage1 = getResultByType('stage1')}
			{@const stage2 = getResultByType('stage2')}

			<div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4 mb-4">
				<h2 class="text-lg font-bold text-gray-900 mb-3">Results</h2>

				<div class="space-y-2">
					{#if first}
						<div class="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
							<div>
								<div class="text-xs text-gray-500 uppercase">1st Place</div>
								<div class="font-semibold">#{first.car_number} {first.driver_name}</div>
							</div>
							<div class="text-lg font-bold text-nascar-blue">135</div>
						</div>
					{/if}
					{#if second}
						<div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
							<div>
								<div class="text-xs text-gray-500 uppercase">2nd Place</div>
								<div class="font-semibold">#{second.car_number} {second.driver_name}</div>
							</div>
							<div class="text-lg font-bold text-gray-600">25</div>
						</div>
					{/if}
					{#if last}
						<div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
							<div>
								<div class="text-xs text-gray-500 uppercase">Last Place</div>
								<div class="font-semibold">#{last.car_number} {last.driver_name}</div>
							</div>
							<div class="text-lg font-bold text-gray-600">15</div>
						</div>
					{/if}
					<div class="grid grid-cols-2 gap-2">
						{#if stage1}
							<div class="p-3 bg-blue-50 rounded-lg">
								<div class="text-xs text-gray-500 uppercase">Stage 1</div>
								<div class="font-medium text-sm">#{stage1.car_number} {stage1.driver_name}</div>
								<div class="text-sm font-bold text-nascar-blue">25 pts</div>
							</div>
						{/if}
						{#if stage2}
							<div class="p-3 bg-blue-50 rounded-lg">
								<div class="text-xs text-gray-500 uppercase">Stage 2</div>
								<div class="font-medium text-sm">#{stage2.car_number} {stage2.driver_name}</div>
								<div class="text-sm font-bold text-nascar-blue">25 pts</div>
							</div>
						{/if}
					</div>
				</div>

				{#if race.is_special_race}
					<div class="mt-3 pt-3 border-t border-gray-100 text-center">
						<span class="text-sm font-medium text-amber-600">All points doubled!</span>
					</div>
				{/if}
			</div>
		{/if}

		<!-- Teams -->
		{#if teams && teams.length > 0}
			<div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
				<div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
					<h2 class="text-lg font-bold text-gray-900">Teams</h2>
					<a
						href="/races/{race.id}/worksheet"
						class="text-sm text-nascar-blue hover:underline"
					>
						Pool Results
					</a>
				</div>

				<div class="divide-y divide-gray-100">
					{#each teams as team, i}
						<div class="p-4 {i === 0 && race.status === 'completed' ? 'bg-yellow-50' : ''}">
							<div class="flex items-center justify-between mb-2">
								<div class="font-semibold text-gray-900">{team.participant_name}</div>
								<div class="text-xl font-bold {team.points_earned > 0 ? 'text-nascar-blue' : 'text-gray-400'}">
									{team.points_earned}
								</div>
							</div>
							<div class="flex flex-wrap gap-2">
								<span class="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded">
									#{team.driver1_number} {team.driver1_name}
								</span>
								<span class="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded">
									#{team.driver2_number} {team.driver2_name}
								</span>
								<span class="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded">
									#{team.driver3_number} {team.driver3_name}
								</span>
								<span class="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded">
									#{team.driver4_number} {team.driver4_name}
								</span>
							</div>
						</div>
					{/each}
				</div>
			</div>

			<!-- Extra drivers (not in pool) -->
			{#if race.extra_drivers}
				<div class="mt-4 p-3 bg-gray-50 rounded-lg text-sm text-gray-600 italic">
					<strong>Driver(s) not used (money rolls to next week):</strong> {race.extra_drivers}
				</div>
			{/if}
		{/if}
	{/if}
</div>
