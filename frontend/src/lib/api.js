// Use relative URL for bundled desktop app, absolute for dev
const API_BASE = '/api';

async function fetchApi(endpoint, options) {
	const response = await fetch(`${API_BASE}${endpoint}`, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			...options?.headers
		}
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ error: 'Unknown error' }));
		throw new Error(error.error || 'API request failed');
	}

	return response.json();
}

export const api = {
	// Participants
	getParticipants: () => fetchApi('/participants'),

	// Drivers (sorted by car number)
	getDrivers: () => fetchApi('/drivers'),

	// Races
	getRaces: () => fetchApi('/races'),
	getRace: (id) => fetchApi(`/races/${id}`),
	createRace: (data) =>
		fetchApi('/races', {
			method: 'POST',
			body: JSON.stringify(data)
		}),
	updateRace: (id, data) =>
		fetchApi(`/races/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	// Teams
	generateTeams: (raceId) =>
		fetchApi(`/races/${raceId}/generate-teams`, {
			method: 'POST'
		}),
	getRaceTeams: (raceId) => fetchApi(`/races/${raceId}/teams`),

	// Results (entered by car number)
	enterRaceResults: (raceId, results) =>
		fetchApi(`/races/${raceId}/results`, {
			method: 'POST',
			body: JSON.stringify(results)
		}),

	// Standings
	getStandings: () => fetchApi('/standings')
};