<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api.js';

	let race = $state(null);
	let teams = $state([]);
	let allDrivers = $state([]);
	let loading = $state(true);
	let error = $state(null);

	// Get list of chartered car numbers for this race's teams
	let assignedCarNumbers = $derived(
		teams.flatMap(t => [t.driver1_number, t.driver2_number, t.driver3_number, t.driver4_number])
	);

	// Non-chartered drivers (extras in the race field)
	let unusedDrivers = $derived(
		allDrivers.filter(d => !assignedCarNumbers.includes(d.car_number))
	);

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
			const [raceData, teamsData, driversData] = await Promise.all([
				api.getRace(id),
				api.getRaceTeams(id),
				api.getDrivers()
			]);
			race = raceData;
			teams = teamsData;
			allDrivers = driversData;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load data';
		} finally {
			loading = false;
		}
	}

	function getPointValue(base, rollover) {
		const mult = race?.is_special_race ? 2 : 1;
		return (base + rollover) * mult;
	}

	// Extract last name from full name (e.g., "William Byron" -> "Byron")
	function lastName(fullName) {
		const parts = fullName.trim().split(' ');
		return parts.length > 1 ? parts[parts.length - 1] : fullName;
	}

	// Find participant who owns a car number
	function findOwner(carNumber) {
		const team = teams.find(t =>
			t.driver1_number === carNumber ||
			t.driver2_number === carNumber ||
			t.driver3_number === carNumber ||
			t.driver4_number === carNumber
		);
		return team?.participant_name || 'N/A';
	}

	// Get highlight class for a driver based on race results
	function getDriverHighlight(carNumber) {
		if (!race?.results || race.status !== 'completed') return '';
		const result = race.results.find(r => r.car_number === carNumber);
		if (!result) return '';
		if (result.is_first_place) return 'bg-green-400';
		if (result.is_second_place) return 'bg-yellow-400';
		if (result.is_last_place) return 'bg-red-400';
		if (result.is_stage1_winner || result.is_stage2_winner) return 'bg-blue-200';
		return '';
	}
</script>

<svelte:head>
	<title>{race?.name || 'Race'} Pool Results - NASCAR Pool</title>
	<style>
		@media print {
			body { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
			.no-print { display: none !important; }
		}
	</style>
</svelte:head>

<div class="max-w-4xl mx-auto p-4">
	<!-- Print button (hidden when printing) -->
	<div class="no-print mb-4 flex gap-2">
		<a href="/races/{$page.params.id}" class="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300">
			&larr; Back
		</a>
		<button onclick={() => window.print()} class="px-4 py-2 bg-nascar-blue text-white rounded hover:bg-blue-700">
			Print / Save PDF
		</button>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else if error}
		<div class="text-red-600">{error}</div>
	{:else if race && teams.length > 0}
		<!-- Worksheet header -->
		<div class="mb-4">
			<h1 class="text-2xl font-bold">2026 NASCAR POOL</h1>
			<div class="flex items-baseline gap-8 text-lg">
				<span class="font-bold">{race.name}</span>
				<span>{race.race_number}</span>
				{#if race.is_special_race}
					<span class="text-amber-600 font-bold text-sm">2X POINTS</span>
				{/if}
			</div>
		</div>

		<!-- Teams grid -->
		<div class="overflow-x-auto mb-6">
			<table class="w-full border-collapse text-sm">
				<thead>
					<tr class="bg-blue-100">
						<th class="border border-gray-300 p-2 w-20"></th>
						{#each teams as team}
							<th class="border border-gray-300 p-2 text-center underline min-w-[90px]">
								{team.participant_name}
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					<tr>
						<td class="border border-gray-300 p-2 font-bold bg-gray-50">Driver 1</td>
						{#each teams as team}
							<td class="border border-gray-300 p-2 text-center {getDriverHighlight(team.driver1_number)}">{lastName(team.driver1_name)}</td>
						{/each}
					</tr>
					<tr>
						<td class="border border-gray-300 p-2 font-bold bg-gray-50">Driver 2</td>
						{#each teams as team}
							<td class="border border-gray-300 p-2 text-center {getDriverHighlight(team.driver2_number)}">{lastName(team.driver2_name)}</td>
						{/each}
					</tr>
					<tr>
						<td class="border border-gray-300 p-2 font-bold bg-gray-50">Driver 3</td>
						{#each teams as team}
							<td class="border border-gray-300 p-2 text-center {getDriverHighlight(team.driver3_number)}">{lastName(team.driver3_name)}</td>
						{/each}
					</tr>
					<tr>
						<td class="border border-gray-300 p-2 font-bold bg-gray-50">Driver 4</td>
						{#each teams as team}
							<td class="border border-gray-300 p-2 text-center {getDriverHighlight(team.driver4_number)}">{lastName(team.driver4_name)}</td>
						{/each}
					</tr>
					<!-- Spacer row -->
					<tr class="bg-gray-300">
						<td colspan={teams.length + 1} class="h-2"></td>
					</tr>
					<!-- Total winnings row -->
					<tr>
						<td class="border border-gray-300 p-2 font-bold">TOTAL WINNINGS</td>
						{#each teams as team}
							<td class="border border-gray-300 p-2 text-center font-bold {team.points_earned > 0 ? 'text-green-700' : ''}">
								{#if race.status === 'completed' && team.points_earned > 0}
									${team.points_earned}
								{/if}
							</td>
						{/each}
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Scoring breakdown -->
		<div class="mb-6">
			<table class="text-sm">
				<tbody>
					<tr class="bg-green-400">
						<td class="py-1 font-bold text-green-800 px-2">First</td>
						<td class="px-4 py-1 min-w-[60px]">${getPointValue(135, race.rollover_first)}</td>
						{#if race.status === 'completed'}
							{@const result = race.results?.find(r => r.is_first_place)}
							{#if result}
								<td class="px-4 py-1">{lastName(result.driver_name)}/{findOwner(result.car_number)}</td>
							{/if}
						{/if}
					</tr>
					<tr class="bg-yellow-400">
						<td class="py-1 font-bold text-yellow-700 px-2">2nd</td>
						<td class="px-4 py-1">${getPointValue(25, race.rollover_second)}</td>
						{#if race.status === 'completed'}
							{@const result = race.results?.find(r => r.is_second_place)}
							{#if result}
								<td class="px-4 py-1">{lastName(result.driver_name)}/{findOwner(result.car_number)}</td>
							{/if}
						{/if}
					</tr>
					<tr class="bg-red-400">
						<td class="py-1 font-bold text-red-800 px-2">Last</td>
						<td class="px-4 py-1">${getPointValue(15, race.rollover_last)}</td>
						{#if race.status === 'completed'}
							{@const result = race.results?.find(r => r.is_last_place)}
							{#if result}
								<td class="px-4 py-1">{lastName(result.driver_name)}/{findOwner(result.car_number)}</td>
							{/if}
						{/if}
					</tr>
					<tr><td colspan="3" class="h-2"></td></tr>
					<tr>
						<td class="py-1 text-blue-700 px-2">Stage 1 Winner</td>
						<td class="px-4 py-1">${getPointValue(25, race.rollover_stage1)}</td>
						{#if race.status === 'completed'}
							{@const result = race.results?.find(r => r.is_stage1_winner)}
							{#if result}
								<td class="px-4 py-1">{lastName(result.driver_name)}/{findOwner(result.car_number)}</td>
							{/if}
						{/if}
					</tr>
					<tr>
						<td class="py-1 text-blue-700 px-2">Stage 2 Winner</td>
						<td class="px-4 py-1">${getPointValue(25, race.rollover_stage2)}</td>
						{#if race.status === 'completed'}
							{@const result = race.results?.find(r => r.is_stage2_winner)}
							{#if result}
								<td class="px-4 py-1">{lastName(result.driver_name)}/{findOwner(result.car_number)}</td>
							{/if}
						{/if}
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Extra drivers (not in pool, money rolls if they win) -->
		{#if race.extra_drivers}
			<div class="text-sm italic text-gray-600">
				<strong>Driver(s) not used (money rolls to next week):</strong>
				{race.extra_drivers}
			</div>
		{/if}
	{:else if race}
		<div class="text-amber-600">
			Teams have not been generated for this race yet.
			<a href="/races/{race.id}" class="underline">Go back</a> to generate teams first.
		</div>
	{/if}
</div>
