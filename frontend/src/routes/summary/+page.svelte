<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';

	let races = $state([]);
	let participants = $state([]);
	let allTeams = $state(new Map());
	let loading = $state(true);
	let error = $state(null);

	// Season investment per person
	const SEASON_INVESTMENT = 1000;

	onMount(async () => {
		try {
			const [racesData, participantsData] = await Promise.all([
				api.getRaces(),
				api.getParticipants()
			]);
			races = racesData;
			participants = participantsData.sort((a, b) => a.name.localeCompare(b.name));

			// Fetch teams for completed races
			const completedRaces = races.filter(r => r.status === 'completed');
			for (const race of completedRaces) {
				const teams = await api.getRaceTeams(race.id);
				allTeams.set(race.id, teams);
			}
			allTeams = allTeams; // trigger reactivity
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load data';
		} finally {
			loading = false;
		}
	});

	// Get winnings for a participant in a race
	function getWinnings(raceId, participantId) {
		const teams = allTeams.get(raceId);
		if (!teams) return 0;
		const team = teams.find(t => t.participant_id === participantId);
		return team?.points_earned || 0;
	}

	// Get total winnings for a participant
	function getTotalWinnings(participantId) {
		let total = 0;
		for (const race of races) {
			if (race.status === 'completed') {
				total += getWinnings(race.id, participantId);
			}
		}
		return total;
	}

	// Count victories (races with winnings >= $135, indicating a race win)
	function getVictoryCount(participantId) {
		let count = 0;
		for (const race of races) {
			if (race.status === 'completed' && getWinnings(race.id, participantId) >= 135) {
				count++;
			}
		}
		return count;
	}

	// Count races in the money (any winnings > 0)
	function getRacesInMoney(participantId) {
		let count = 0;
		for (const race of races) {
			if (race.status === 'completed' && getWinnings(race.id, participantId) > 0) {
				count++;
			}
		}
		return count;
	}

	// Get race total
	function getRaceTotal(raceId) {
		const teams = allTeams.get(raceId);
		if (!teams) return 0;
		return teams.reduce((sum, t) => sum + t.points_earned, 0);
	}

	// Grand total of all winnings
	function getGrandTotal() {
		return participants.reduce((sum, p) => sum + getTotalWinnings(p.id), 0);
	}
</script>

<svelte:head>
	<title>Season Summary - NASCAR Pool</title>
	<style>
		@media print {
			body { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
			.no-print { display: none !important; }
		}
	</style>
</svelte:head>

<div class="max-w-full mx-auto">
	<!-- Print button -->
	<div class="no-print mb-4 flex justify-between items-center">
		<h1 class="text-2xl font-bold">Season Summary</h1>
		<button onclick={() => window.print()} class="px-4 py-2 bg-nascar-blue text-white rounded hover:bg-blue-700">
			Print / Save PDF
		</button>
	</div>

	{#if loading}
		<div class="text-center py-12 text-gray-500">Loading...</div>
	{:else if error}
		<div class="text-red-600">{error}</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="w-full border-collapse text-xs">
				<thead>
					<!-- Title row -->
					<tr>
						<td colspan={participants.length + 4} class="text-xl font-bold pb-2">2026 NASCAR POOL</td>
					</tr>
					<!-- Victories row -->
					<tr class="text-blue-600">
						<td></td>
						<td></td>
						<td class="px-1 text-right">Victories</td>
						{#each participants as p}
							<td class="px-2 py-1 text-center font-bold">{getVictoryCount(p.id)}</td>
						{/each}
						<td></td>
					</tr>
					<!-- Header row -->
					<tr class="border-b-2 border-gray-400">
						<td class="px-1 py-1 font-medium"></td>
						<td class="px-1 py-1 font-medium">Race Name (* 2x race)</td>
						<td class="px-1 py-1 font-medium">Date</td>
						{#each participants as p}
							<td class="px-2 py-1 text-center font-medium underline">{p.name}</td>
						{/each}
						<td class="px-2 py-1 text-center font-medium">TOTALS</td>
					</tr>
				</thead>
				<tbody>
					{#each races as race}
						{@const raceTotal = getRaceTotal(race.id)}
						<tr class="border-b border-gray-200 {race.status !== 'completed' ? 'text-gray-400' : ''}">
							<td class="px-1 py-1">{race.race_number}</td>
							<td class="px-1 py-1 whitespace-nowrap">
								{race.name}{race.is_special_race ? '*' : ''}
							</td>
							<td class="px-1 py-1 whitespace-nowrap">{race.date || ''}</td>
							{#each participants as p}
								{@const winnings = getWinnings(race.id, p.id)}
								<td class="px-2 py-1 text-center {winnings > 0 ? 'bg-green-200 font-medium' : ''}">
									{#if race.status === 'completed' && winnings > 0}
										${winnings}
									{/if}
								</td>
							{/each}
							<td class="px-2 py-1 text-center font-medium">
								{#if race.status === 'completed'}
									${raceTotal}
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
				<tfoot>
					<!-- Spacer -->
					<tr><td colspan={participants.length + 4} class="h-2"></td></tr>
					<!-- Winnings row -->
					<tr class="border-t-2 border-gray-400 font-bold">
						<td></td>
						<td class="px-1 py-2">WINNINGS</td>
						<td></td>
						{#each participants as p}
							<td class="px-2 py-2 text-center">${getTotalWinnings(p.id).toFixed(2)}</td>
						{/each}
						<td class="px-2 py-2 text-center">${getGrandTotal().toFixed(2)}</td>
					</tr>
					<!-- Investment row -->
					<tr>
						<td></td>
						<td class="px-1 py-1">Investment to Date</td>
						<td></td>
						{#each participants as p}
							<td class="px-2 py-1 text-center">(${SEASON_INVESTMENT.toFixed(2)})</td>
						{/each}
						<td class="px-2 py-1 text-center">(${(participants.length * SEASON_INVESTMENT).toFixed(2)})</td>
					</tr>
					<!-- NET row -->
					<tr class="font-bold">
						<td></td>
						<td class="px-1 py-2">NET</td>
						<td></td>
						{#each participants as p}
							{@const net = getTotalWinnings(p.id) - SEASON_INVESTMENT}
							<td class="px-2 py-2 text-center {net >= 0 ? 'text-green-700' : 'text-red-600'}">
								{net >= 0 ? '' : '('}${Math.abs(net).toFixed(2)}{net >= 0 ? '' : ')'}
							</td>
						{/each}
						<td class="px-2 py-2 text-center">$0.00</td>
					</tr>
					<!-- Races in the Money -->
					<tr>
						<td></td>
						<td class="px-1 py-1">Races in the Money</td>
						<td></td>
						{#each participants as p}
							<td class="px-2 py-1 text-center">{getRacesInMoney(p.id)}</td>
						{/each}
						<td></td>
					</tr>
					<!-- Season Amount -->
					<tr>
						<td></td>
						<td class="px-1 py-1">Season Amount</td>
						<td></td>
						{#each participants as p}
							<td class="px-2 py-1 text-center">${SEASON_INVESTMENT.toFixed(2)}</td>
						{/each}
						<td class="px-2 py-1 text-center">${(participants.length * SEASON_INVESTMENT).toFixed(2)}</td>
					</tr>
				</tfoot>
			</table>
		</div>
	{/if}
</div>
