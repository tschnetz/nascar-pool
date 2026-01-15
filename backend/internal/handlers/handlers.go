package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tgschnetzer/nascar-pool/internal/database"
	"github.com/tgschnetzer/nascar-pool/internal/models"
)

// Helper function to write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function to write error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// GetParticipants returns all participants
func GetParticipants(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(context.Background(),
		"SELECT id, name, created_at FROM participants ORDER BY name")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch participants")
		return
	}
	defer rows.Close()

	var participants []models.Participant
	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan participant")
			return
		}
		participants = append(participants, p)
	}

	writeJSON(w, http.StatusOK, participants)
}

// GetDrivers returns all drivers with their car numbers
func GetDrivers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(context.Background(),
		"SELECT id, name, car_number, team_name, manufacturer, is_chartered, created_at FROM drivers ORDER BY car_number")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch drivers")
		return
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		if err := rows.Scan(&d.ID, &d.Name, &d.CarNumber, &d.TeamName, &d.Manufacturer, &d.IsChartered, &d.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan driver")
			return
		}
		drivers = append(drivers, d)
	}

	writeJSON(w, http.StatusOK, drivers)
}

// GetRaces returns all races with their results
func GetRaces(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(context.Background(), `
		SELECT id, name, race_number, date::text, is_special_race, status,
			   rollover_first, rollover_second, rollover_last, rollover_stage1, rollover_stage2,
			   extra_drivers, created_at
		FROM races ORDER BY race_number`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch races")
		return
	}
	defer rows.Close()

	var races []models.RaceWithResults
	for rows.Next() {
		var race models.RaceWithResults
		if err := rows.Scan(
			&race.ID, &race.Name, &race.RaceNumber, &race.Date,
			&race.IsSpecialRace, &race.Status,
			&race.RolloverFirst, &race.RolloverSecond, &race.RolloverLast,
			&race.RolloverStage1, &race.RolloverStage2,
			&race.ExtraDrivers, &race.CreatedAt,
		); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan race")
			return
		}
		races = append(races, race)
	}

	// Fetch results for completed races
	for i := range races {
		if races[i].Status == "completed" {
			results, err := getRaceResults(races[i].ID)
			if err == nil {
				races[i].Results = results
			}
		}
	}

	writeJSON(w, http.StatusOK, races)
}

// getRaceResults fetches results for a race with driver names
func getRaceResults(raceID int) ([]models.RaceResultWithDriver, error) {
	rows, err := database.Pool.Query(context.Background(), `
		SELECT rr.id, rr.race_id, rr.car_number, rr.position,
			   rr.is_first_place, rr.is_second_place, rr.is_last_place,
			   rr.is_stage1_winner, rr.is_stage2_winner, rr.created_at,
			   COALESCE(d.name, 'Unknown') as driver_name
		FROM race_results rr
		LEFT JOIN drivers d ON rr.car_number = d.car_number
		WHERE rr.race_id = $1
		ORDER BY rr.car_number`, raceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.RaceResultWithDriver
	for rows.Next() {
		var r models.RaceResultWithDriver
		if err := rows.Scan(
			&r.ID, &r.RaceID, &r.CarNumber, &r.Position,
			&r.IsFirstPlace, &r.IsSecondPlace, &r.IsLastPlace,
			&r.IsStage1Winner, &r.IsStage2Winner, &r.CreatedAt,
			&r.DriverName,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// GetRace returns a single race by ID with results
func GetRace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid race ID")
		return
	}

	var race models.RaceWithResults
	err = database.Pool.QueryRow(context.Background(), `
		SELECT id, name, race_number, date::text, is_special_race, status,
			   rollover_first, rollover_second, rollover_last, rollover_stage1, rollover_stage2,
			   extra_drivers, created_at
		FROM races WHERE id = $1`, id).Scan(
		&race.ID, &race.Name, &race.RaceNumber, &race.Date,
		&race.IsSpecialRace, &race.Status,
		&race.RolloverFirst, &race.RolloverSecond, &race.RolloverLast,
		&race.RolloverStage1, &race.RolloverStage2,
		&race.ExtraDrivers, &race.CreatedAt,
	)
	if err != nil {
		writeError(w, http.StatusNotFound, "Race not found")
		return
	}

	// Fetch results if completed
	if race.Status == "completed" {
		results, err := getRaceResults(race.ID)
		if err == nil {
			race.Results = results
		}
	}

	writeJSON(w, http.StatusOK, race)
}

// CreateRace creates a new race
func CreateRace(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var raceID int
	err := database.Pool.QueryRow(context.Background(),
		`INSERT INTO races (name, race_number, date, is_special_race, status)
		 VALUES ($1, $2, $3, $4, 'upcoming')
		 RETURNING id`,
		req.Name, req.RaceNumber, req.Date, req.IsSpecialRace).Scan(&raceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create race")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int{"id": raceID})
}

// UpdateRace updates an existing race's name, date, or special race status
func UpdateRace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid race ID")
		return
	}

	var req models.CreateRaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := database.Pool.Exec(context.Background(),
		`UPDATE races SET name = $1, date = $2, is_special_race = $3, extra_drivers = $4 WHERE id = $5`,
		req.Name, req.Date, req.IsSpecialRace, req.ExtraDrivers, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update race")
		return
	}

	if result.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "Race not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Race updated"})
}

// GenerateTeams randomly assigns drivers to participants for a race
func GenerateTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	raceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid race ID")
		return
	}

	// Check if teams already exist for this race
	var existingCount int
	err = database.Pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM teams WHERE race_id = $1", raceID).Scan(&existingCount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check existing teams")
		return
	}
	if existingCount > 0 {
		writeError(w, http.StatusBadRequest, "Teams already generated for this race")
		return
	}

	// Get all participants
	participantRows, err := database.Pool.Query(context.Background(),
		"SELECT id FROM participants ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch participants")
		return
	}
	defer participantRows.Close()

	var participantIDs []int
	for participantRows.Next() {
		var id int
		if err := participantRows.Scan(&id); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan participant")
			return
		}
		participantIDs = append(participantIDs, id)
	}

	// Get all drivers
	driverRows, err := database.Pool.Query(context.Background(),
		"SELECT id FROM drivers ORDER BY id")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch drivers")
		return
	}
	defer driverRows.Close()

	var driverIDs []int
	for driverRows.Next() {
		var id int
		if err := driverRows.Scan(&id); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan driver")
			return
		}
		driverIDs = append(driverIDs, id)
	}

	// Shuffle drivers randomly
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(driverIDs), func(i, j int) {
		driverIDs[i], driverIDs[j] = driverIDs[j], driverIDs[i]
	})

	// Assign 4 drivers to each participant (9 participants x 4 drivers = 36 drivers)
	for i, participantID := range participantIDs {
		startIdx := i * 4
		driver1 := driverIDs[startIdx]
		driver2 := driverIDs[startIdx+1]
		driver3 := driverIDs[startIdx+2]
		driver4 := driverIDs[startIdx+3]

		_, err := database.Pool.Exec(context.Background(),
			`INSERT INTO teams (race_id, participant_id, driver1_id, driver2_id, driver3_id, driver4_id)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			raceID, participantID, driver1, driver2, driver3, driver4)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to create team")
			return
		}
	}

	// Update race status to in_progress
	_, err = database.Pool.Exec(context.Background(),
		"UPDATE races SET status = 'in_progress' WHERE id = $1", raceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update race status")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Teams generated successfully"})
}

// GetRaceTeams returns all teams for a specific race with driver car numbers
func GetRaceTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	raceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid race ID")
		return
	}

	rows, err := database.Pool.Query(context.Background(), `
		SELECT t.id, t.race_id, t.participant_id, t.driver1_id, t.driver2_id,
			   t.driver3_id, t.driver4_id, t.points_earned, t.created_at,
			   p.name,
			   d1.name, d1.car_number,
			   d2.name, d2.car_number,
			   d3.name, d3.car_number,
			   d4.name, d4.car_number
		FROM teams t
		JOIN participants p ON t.participant_id = p.id
		JOIN drivers d1 ON t.driver1_id = d1.id
		JOIN drivers d2 ON t.driver2_id = d2.id
		JOIN drivers d3 ON t.driver3_id = d3.id
		JOIN drivers d4 ON t.driver4_id = d4.id
		WHERE t.race_id = $1
		ORDER BY t.points_earned DESC, p.name`, raceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch teams")
		return
	}
	defer rows.Close()

	var teams []models.TeamWithDetails
	for rows.Next() {
		var t models.TeamWithDetails
		if err := rows.Scan(
			&t.ID, &t.RaceID, &t.ParticipantID, &t.Driver1ID, &t.Driver2ID,
			&t.Driver3ID, &t.Driver4ID, &t.PointsEarned, &t.CreatedAt,
			&t.ParticipantName,
			&t.Driver1Name, &t.Driver1Number,
			&t.Driver2Name, &t.Driver2Number,
			&t.Driver3Name, &t.Driver3Number,
			&t.Driver4Name, &t.Driver4Number,
		); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan team")
			return
		}
		teams = append(teams, t)
	}

	writeJSON(w, http.StatusOK, teams)
}

// EnterRaceResults enters the results for a race BY CAR NUMBER and calculates points
// Scoring flow: car_number -> check if chartered -> if yes, award points + rollover to team
// If car is not chartered, points roll over to next race
func EnterRaceResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	raceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid race ID")
		return
	}

	var req models.RaceResultsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get race info including rollover amounts coming into this race
	var isSpecialRace bool
	var raceNumber int
	var rolloverFirst, rolloverSecond, rolloverLast, rolloverStage1, rolloverStage2 int
	err = database.Pool.QueryRow(context.Background(),
		`SELECT is_special_race, race_number, rollover_first, rollover_second,
		        rollover_last, rollover_stage1, rollover_stage2
		 FROM races WHERE id = $1`, raceID).Scan(
		&isSpecialRace, &raceNumber, &rolloverFirst, &rolloverSecond,
		&rolloverLast, &rolloverStage1, &rolloverStage2)
	if err != nil {
		writeError(w, http.StatusNotFound, "Race not found")
		return
	}

	// Get all chartered car numbers (the 36 that are assigned to teams)
	charteredCars := make(map[int]bool)
	carRows, err := database.Pool.Query(context.Background(),
		"SELECT car_number FROM drivers WHERE is_chartered = true")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch chartered cars")
		return
	}
	for carRows.Next() {
		var carNum int
		carRows.Scan(&carNum)
		charteredCars[carNum] = true
	}
	carRows.Close()

	// Delete any existing results for this race (allows re-entering)
	_, err = database.Pool.Exec(context.Background(),
		"DELETE FROM race_results WHERE race_id = $1", raceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to clear existing results")
		return
	}

	// Insert race results by car number
	scoringResults := []struct {
		carNumber      int
		isFirst        bool
		isSecond       bool
		isLast         bool
		isStage1Winner bool
		isStage2Winner bool
	}{
		{req.FirstPlaceCarNumber, true, false, false, false, false},
		{req.SecondPlaceCarNumber, false, true, false, false, false},
		{req.LastPlaceCarNumber, false, false, true, false, false},
		{req.Stage1WinnerCarNumber, false, false, false, true, false},
		{req.Stage2WinnerCarNumber, false, false, false, false, true},
	}

	for _, result := range scoringResults {
		_, err := database.Pool.Exec(context.Background(), `
			INSERT INTO race_results (race_id, car_number, is_first_place, is_second_place,
				is_last_place, is_stage1_winner, is_stage2_winner)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (race_id, car_number) DO UPDATE SET
				is_first_place = race_results.is_first_place OR EXCLUDED.is_first_place,
				is_second_place = race_results.is_second_place OR EXCLUDED.is_second_place,
				is_last_place = race_results.is_last_place OR EXCLUDED.is_last_place,
				is_stage1_winner = race_results.is_stage1_winner OR EXCLUDED.is_stage1_winner,
				is_stage2_winner = race_results.is_stage2_winner OR EXCLUDED.is_stage2_winner`,
			raceID, result.carNumber, result.isFirst, result.isSecond, result.isLast,
			result.isStage1Winner, result.isStage2Winner)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to insert race result")
			return
		}
	}

	// Update race status to completed
	_, err = database.Pool.Exec(context.Background(),
		"UPDATE races SET status = 'completed' WHERE id = $1", raceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update race status")
		return
	}

	// Calculate multiplier
	multiplier := 1
	if isSpecialRace {
		multiplier = models.SpecialMultiplier
	}

	// Calculate total available points for each position (base + rollover) * multiplier
	availableFirst := (models.PointsFirst + rolloverFirst) * multiplier
	availableSecond := (models.PointsSecond + rolloverSecond) * multiplier
	availableLast := (models.PointsLast + rolloverLast) * multiplier
	availableStage1 := (models.PointsStage1 + rolloverStage1) * multiplier
	availableStage2 := (models.PointsStage2 + rolloverStage2) * multiplier

	// Track rollover for next race (for non-chartered scorers)
	nextRolloverFirst, nextRolloverSecond, nextRolloverLast := 0, 0, 0
	nextRolloverStage1, nextRolloverStage2 := 0, 0

	// Check if each scoring car is chartered
	if !charteredCars[req.FirstPlaceCarNumber] {
		nextRolloverFirst = models.PointsFirst + rolloverFirst
		availableFirst = 0 // No one gets these points
	}
	if !charteredCars[req.SecondPlaceCarNumber] {
		nextRolloverSecond = models.PointsSecond + rolloverSecond
		availableSecond = 0
	}
	if !charteredCars[req.LastPlaceCarNumber] {
		nextRolloverLast = models.PointsLast + rolloverLast
		availableLast = 0
	}
	if !charteredCars[req.Stage1WinnerCarNumber] {
		nextRolloverStage1 = models.PointsStage1 + rolloverStage1
		availableStage1 = 0
	}
	if !charteredCars[req.Stage2WinnerCarNumber] {
		nextRolloverStage2 = models.PointsStage2 + rolloverStage2
		availableStage2 = 0
	}

	// Update next race with rollover (if any)
	if nextRolloverFirst > 0 || nextRolloverSecond > 0 || nextRolloverLast > 0 ||
		nextRolloverStage1 > 0 || nextRolloverStage2 > 0 {
		_, err = database.Pool.Exec(context.Background(), `
			UPDATE races SET
				rollover_first = $1, rollover_second = $2, rollover_last = $3,
				rollover_stage1 = $4, rollover_stage2 = $5
			WHERE race_number = $6`,
			nextRolloverFirst, nextRolloverSecond, nextRolloverLast,
			nextRolloverStage1, nextRolloverStage2, raceNumber+1)
		// Ignore error if next race doesn't exist yet
	}

	// Get all teams with their drivers' car numbers and calculate points
	rows, err := database.Pool.Query(context.Background(), `
		SELECT t.id,
			   d1.car_number, d2.car_number, d3.car_number, d4.car_number
		FROM teams t
		JOIN drivers d1 ON t.driver1_id = d1.id
		JOIN drivers d2 ON t.driver2_id = d2.id
		JOIN drivers d3 ON t.driver3_id = d3.id
		JOIN drivers d4 ON t.driver4_id = d4.id
		WHERE t.race_id = $1`, raceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch teams")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var teamID, car1, car2, car3, car4 int
		if err := rows.Scan(&teamID, &car1, &car2, &car3, &car4); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan team")
			return
		}

		teamCarNumbers := []int{car1, car2, car3, car4}
		points := 0

		// Check each car number against scoring positions
		for _, carNum := range teamCarNumbers {
			if carNum == req.FirstPlaceCarNumber {
				points += availableFirst
			}
			if carNum == req.SecondPlaceCarNumber {
				points += availableSecond
			}
			if carNum == req.LastPlaceCarNumber {
				points += availableLast
			}
			if carNum == req.Stage1WinnerCarNumber {
				points += availableStage1
			}
			if carNum == req.Stage2WinnerCarNumber {
				points += availableStage2
			}
		}

		// Update team points
		_, err := database.Pool.Exec(context.Background(),
			"UPDATE teams SET points_earned = $1 WHERE id = $2", points, teamID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to update team points")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Race results entered and points calculated"})
}

// GetStandings returns the cumulative standings across all races
func GetStandings(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Pool.Query(context.Background(), `
		SELECT p.id, p.name,
			   COALESCE(SUM(t.points_earned), 0) as total_points,
			   COUNT(DISTINCT CASE WHEN r.status = 'completed' THEN t.race_id END) as races_completed
		FROM participants p
		LEFT JOIN teams t ON p.id = t.participant_id
		LEFT JOIN races r ON t.race_id = r.id
		GROUP BY p.id, p.name
		ORDER BY total_points DESC, p.name`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to fetch standings")
		return
	}
	defer rows.Close()

	var standings []models.Standing
	rank := 1
	for rows.Next() {
		var s models.Standing
		if err := rows.Scan(&s.ParticipantID, &s.ParticipantName, &s.TotalPoints, &s.RacesCompleted); err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to scan standing")
			return
		}
		s.Rank = rank
		standings = append(standings, s)
		rank++
	}

	writeJSON(w, http.StatusOK, standings)
}
