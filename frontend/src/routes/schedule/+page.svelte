<script>
	import { schedule } from '$lib/schedule.js';

	function getTvStyle(tv) {
		switch (tv) {
			case 'FOX':
			case 'FS1':
				return 'bg-blue-100 text-blue-700';
			case 'NBC':
			case 'USA':
				return 'bg-purple-100 text-purple-700';
			case 'Prime':
				return 'bg-cyan-100 text-cyan-700';
			case 'TNT':
				return 'bg-red-100 text-red-700';
			default:
				return 'bg-gray-100 text-gray-600';
		}
	}
</script>

<svelte:head>
	<title>Schedule - NASCAR Pool</title>
</svelte:head>

<div>
	<h1 class="text-2xl font-bold text-gray-900 mb-4">2026 Schedule</h1>

	<div class="space-y-2">
		{#each schedule as entry}
			{#if entry.isOffWeek}
				<!-- Off week -->
				<div class="bg-gray-50 rounded-lg p-3 text-center text-gray-500 italic">
					{entry.date} - {entry.name}
				</div>
			{:else}
				<div
					class="bg-white rounded-xl shadow-sm border border-gray-100 p-3
						   {entry.isSpecial ? 'ring-2 ring-yellow-400' : ''}
						   {entry.isNonPoints ? 'opacity-60' : ''}"
				>
					<div class="flex items-start gap-3">
						<!-- Race number -->
						<div
							class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm flex-shrink-0
								   {entry.isSpecial
								? 'bg-yellow-400 text-yellow-900'
								: entry.isNonPoints
									? 'bg-gray-200 text-gray-500'
									: 'bg-nascar-blue text-white'}"
						>
							{entry.raceNumber ?? '*'}
						</div>

						<!-- Race info -->
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2 flex-wrap">
								<span class="font-semibold text-gray-900 truncate">{entry.name}</span>
								{#if entry.isSpecial}
									<span class="text-xs font-bold bg-yellow-400 text-yellow-900 px-2 py-0.5 rounded">
										2X
									</span>
								{/if}
								{#if entry.isNonPoints}
									<span class="text-xs text-gray-500">(Non-points)</span>
								{/if}
							</div>
							<div class="text-sm text-gray-500 mt-0.5">{entry.track}</div>
							<div class="flex items-center gap-3 mt-1 text-xs text-gray-400">
								<span>{entry.date}</span>
								{#if entry.startTime}
									<span>{entry.startTime} ET</span>
								{/if}
								{#if entry.tv}
									<span class="px-1.5 py-0.5 rounded {getTvStyle(entry.tv)}">{entry.tv}</span>
								{/if}
							</div>
						</div>
					</div>
				</div>
			{/if}
		{/each}
	</div>

	<!-- Legend -->
	<div class="mt-6 p-4 bg-gray-50 rounded-lg">
		<h3 class="text-sm font-medium text-gray-700 mb-2">Legend</h3>
		<div class="flex flex-wrap gap-3 text-sm">
			<div class="flex items-center gap-1">
				<span class="w-4 h-4 rounded-full bg-yellow-400"></span>
				<span class="text-gray-600">2x Points</span>
			</div>
			<div class="flex items-center gap-1">
				<span class="w-4 h-4 rounded-full bg-nascar-blue"></span>
				<span class="text-gray-600">Regular Race</span>
			</div>
			<div class="flex items-center gap-1">
				<span class="w-4 h-4 rounded-full bg-gray-200"></span>
				<span class="text-gray-600">Non-Points</span>
			</div>
		</div>
	</div>
</div>
