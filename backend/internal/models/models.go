package models

import "time"

// Participant represents a person in the fantasy pool
type Participant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Driver represents a NASCAR driver with permanent car number
type Driver struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	CarNumber    int       `json:"car_number"`
	TeamName     *string   `json:"team_name,omitempty"`
	Manufacturer *string   `json:"manufacturer,omitempty"`
	IsChartered  bool      `json:"is_chartered"`
	CreatedAt    time.Time `json:"created_at"`
}

// Race represents a NASCAR race
type Race struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	RaceNumber      int       `json:"race_number"`
	Date            *string   `json:"date,omitempty"`
	IsSpecialRace   bool      `json:"is_special_race"`
	Status          string    `json:"status"`
	RolloverFirst   int       `json:"rollover_first"`
	RolloverSecond  int       `json:"rollover_second"`
	RolloverLast    int       `json:"rollover_last"`
	RolloverStage1  int       `json:"rollover_stage1"`
	RolloverStage2  int       `json:"rollover_stage2"`
	ExtraDrivers    *string   `json:"extra_drivers,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// RaceResult represents a scoring result for a car in a race
type RaceResult struct {
	ID            int       `json:"id"`
	RaceID        int       `json:"race_id"`
	CarNumber     int       `json:"car_number"`
	Position      *int      `json:"position,omitempty"`
	IsFirstPlace  bool      `json:"is_first_place"`
	IsSecondPlace bool      `json:"is_second_place"`
	IsLastPlace   bool      `json:"is_last_place"`
	IsStage1Winner bool     `json:"is_stage1_winner"`
	IsStage2Winner bool     `json:"is_stage2_winner"`
	CreatedAt     time.Time `json:"created_at"`
}

// RaceResultWithDriver includes driver info for display
type RaceResultWithDriver struct {
	RaceResult
	DriverName string `json:"driver_name"`
}

// RaceWithResults includes the race and its results for display
type RaceWithResults struct {
	Race
	Results []RaceResultWithDriver `json:"results,omitempty"`
}

// Team represents a participant's driver assignments for a race
type Team struct {
	ID            int       `json:"id"`
	RaceID        int       `json:"race_id"`
	ParticipantID int       `json:"participant_id"`
	Driver1ID     int       `json:"driver1_id"`
	Driver2ID     int       `json:"driver2_id"`
	Driver3ID     int       `json:"driver3_id"`
	Driver4ID     int       `json:"driver4_id"`
	PointsEarned  int       `json:"points_earned"`
	CreatedAt     time.Time `json:"created_at"`
}

// TeamWithDetails includes participant and driver names/numbers
type TeamWithDetails struct {
	Team
	ParticipantName string `json:"participant_name"`
	Driver1Name     string `json:"driver1_name"`
	Driver1Number   int    `json:"driver1_number"`
	Driver2Name     string `json:"driver2_name"`
	Driver2Number   int    `json:"driver2_number"`
	Driver3Name     string `json:"driver3_name"`
	Driver3Number   int    `json:"driver3_number"`
	Driver4Name     string `json:"driver4_name"`
	Driver4Number   int    `json:"driver4_number"`
}

// Standing represents a participant's cumulative standings
type Standing struct {
	ParticipantID   int    `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	TotalPoints     int    `json:"total_points"`
	RacesCompleted  int    `json:"races_completed"`
	Rank            int    `json:"rank"`
}

// RaceResultsRequest is the payload for entering race results BY CAR NUMBER
type RaceResultsRequest struct {
	FirstPlaceCarNumber   int `json:"first_place_car_number"`
	SecondPlaceCarNumber  int `json:"second_place_car_number"`
	LastPlaceCarNumber    int `json:"last_place_car_number"`
	Stage1WinnerCarNumber int `json:"stage1_winner_car_number"`
	Stage2WinnerCarNumber int `json:"stage2_winner_car_number"`
}

// CreateRaceRequest is the payload for creating/updating a race
type CreateRaceRequest struct {
	Name          string `json:"name"`
	RaceNumber    int    `json:"race_number"`
	Date          string `json:"date,omitempty"`
	IsSpecialRace bool   `json:"is_special_race"`
	ExtraDrivers  string `json:"extra_drivers,omitempty"`
}

// Scoring constants
const (
	PointsFirst       = 135
	PointsSecond      = 25
	PointsLast        = 15
	PointsStage1      = 25
	PointsStage2      = 25
	SpecialMultiplier = 2
)
