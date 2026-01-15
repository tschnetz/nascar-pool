const API_BASE = 'http://localhost:8080/api';

export interface Participant {
	id: number;
	name: string;
	created_at: string;
}

export interface Driver {
	id: number;
	name: string;
	car_number: number;
	team_name?: string;
	manufacturer?: string;
	is_chartered: boolean;
	created_at: string;
}

export interface RaceResult {
	id: number;
	race_id: number;
	car_number: number;
	position?: number;
	is_first_place: boolean;
	is_second_place: boolean;
	is_last_place: boolean;
	is_stage1_winner: boolean;
	is_stage2_winner: boolean;
	driver_name: string;
	created_at: string;
}

export interface Race {
	id: number;
	name: string;
	race_number: number;
	date?: string;
	is_special_race: boolean;
	status: 'upcoming' | 'in_progress' | 'completed';
	rollover_first: number;
	rollover_second: number;
	rollover_last: number;
	rollover_stage1: number;
	rollover_stage2: number;
	extra_drivers?: string;
	results?: RaceResult[];
	created_at: string;
}

export interface Team {
	id: number;
	race_id: number;
	participant_id: number;
	driver1_id: number;
	driver2_id: number;
	driver3_id: number;
	driver4_id: number;
	points_earned: number;
	participant_name: string;
	driver1_name: string;
	driver1_number: number;
	driver2_name: string;
	driver2_number: number;
	driver3_name: string;
	driver3_number: number;
	driver4_name: string;
	driver4_number: number;
	created_at: string;
}

export interface Standing {
	participant_id: number;
	participant_name: string;
	total_points: number;
	races_completed: number;
	rank: number;
}

// Results are entered BY CAR NUMBER
export interface RaceResultsPayload {
	first_place_car_number: number;
	second_place_car_number: number;
	last_place_car_number: number;
	stage1_winner_car_number: number;
	stage2_winner_car_number: number;
}

export interface CreateRacePayload {
	name: string;
	race_number: number;
	date?: string;
	is_special_race: boolean;
	extra_drivers?: string;
}

async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<T> {
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
	getParticipants: () => fetchApi<Participant[]>('/participants'),

	// Drivers (sorted by car number)
	getDrivers: () => fetchApi<Driver[]>('/drivers'),

	// Races
	getRaces: () => fetchApi<Race[]>('/races'),
	getRace: (id: number) => fetchApi<Race>(`/races/${id}`),
	createRace: (data: CreateRacePayload) =>
		fetchApi<{ id: number }>('/races', {
			method: 'POST',
			body: JSON.stringify(data)
		}),
	updateRace: (id: number, data: Partial<CreateRacePayload>) =>
		fetchApi<{ message: string }>(`/races/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		}),

	// Teams
	generateTeams: (raceId: number) =>
		fetchApi<{ message: string }>(`/races/${raceId}/generate-teams`, {
			method: 'POST'
		}),
	getRaceTeams: (raceId: number) => fetchApi<Team[]>(`/races/${raceId}/teams`),

	// Results (entered by car number)
	enterRaceResults: (raceId: number, results: RaceResultsPayload) =>
		fetchApi<{ message: string }>(`/races/${raceId}/results`, {
			method: 'POST',
			body: JSON.stringify(results)
		}),

	// Standings
	getStandings: () => fetchApi<Standing[]>('/standings')
};
