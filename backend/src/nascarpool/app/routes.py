import random
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from sqlalchemy import func
from pydantic import BaseModel
from typing import Optional, List

from .database import get_db
from .models import (
    Participant, Driver, Race, RaceResult, Team,
    POINTS_FIRST, POINTS_SECOND, POINTS_LAST, POINTS_STAGE1, POINTS_STAGE2, SPECIAL_MULTIPLIER
)

router = APIRouter(prefix="/api")


# Pydantic models for request/response
class ParticipantResponse(BaseModel):
    id: int
    name: str
    created_at: str

    class Config:
        from_attributes = True


class DriverResponse(BaseModel):
    id: int
    name: str
    car_number: str
    team_name: Optional[str]
    manufacturer: Optional[str]
    is_chartered: bool
    created_at: str

    class Config:
        from_attributes = True


class RaceResultResponse(BaseModel):
    id: int
    race_id: int
    car_number: str
    position: Optional[int]
    is_first_place: bool
    is_second_place: bool
    is_last_place: bool
    is_stage1_winner: bool
    is_stage2_winner: bool
    driver_name: str
    created_at: str


class RaceResponse(BaseModel):
    id: int
    name: str
    race_number: int
    date: Optional[str]
    is_special_race: bool
    status: str
    rollover_first: int
    rollover_second: int
    rollover_last: int
    rollover_stage1: int
    rollover_stage2: int
    extra_drivers: Optional[str]
    results: Optional[List[RaceResultResponse]]
    created_at: str

    class Config:
        from_attributes = True


class TeamResponse(BaseModel):
    id: int
    race_id: int
    participant_id: int
    driver1_id: int
    driver2_id: int
    driver3_id: int
    driver4_id: int
    points_earned: int
    participant_name: str
    driver1_name: str
    driver1_number: str
    driver2_name: str
    driver2_number: str
    driver3_name: str
    driver3_number: str
    driver4_name: str
    driver4_number: str
    created_at: str


class StandingResponse(BaseModel):
    participant_id: int
    participant_name: str
    total_points: int
    races_completed: int
    rank: int


class CreateRaceRequest(BaseModel):
    name: str
    race_number: int
    date: Optional[str] = None
    is_special_race: bool = False
    extra_drivers: Optional[str] = None


class RaceResultsRequest(BaseModel):
    first_place_car_number: str
    second_place_car_number: str
    last_place_car_number: str
    stage1_winner_car_number: str
    stage2_winner_car_number: str


# Routes
@router.get("/participants", response_model=List[ParticipantResponse])
def get_participants(db: Session = Depends(get_db)):
    participants = db.query(Participant).order_by(Participant.name).all()
    return [
        ParticipantResponse(
            id=p.id,
            name=p.name,
            created_at=str(p.created_at)
        )
        for p in participants
    ]


@router.get("/drivers", response_model=List[DriverResponse])
def get_drivers(db: Session = Depends(get_db)):
    drivers = db.query(Driver).order_by(Driver.car_number).all()
    return [
        DriverResponse(
            id=d.id,
            name=d.name,
            car_number=d.car_number,
            team_name=d.team_name,
            manufacturer=d.manufacturer,
            is_chartered=d.is_chartered,
            created_at=str(d.created_at)
        )
        for d in drivers
    ]


def get_race_results_with_driver(db: Session, race_id: int) -> List[RaceResultResponse]:
    results = db.query(RaceResult).filter(RaceResult.race_id == race_id).all()
    response = []
    for r in results:
        driver = db.query(Driver).filter(Driver.car_number == r.car_number).first()
        driver_name = driver.name if driver else 'Unknown'
        response.append(RaceResultResponse(
            id=r.id,
            race_id=r.race_id,
            car_number=r.car_number,
            position=r.position,
            is_first_place=r.is_first_place,
            is_second_place=r.is_second_place,
            is_last_place=r.is_last_place,
            is_stage1_winner=r.is_stage1_winner,
            is_stage2_winner=r.is_stage2_winner,
            driver_name=driver_name,
            created_at=str(r.created_at)
        ))
    return response


@router.get("/races", response_model=List[RaceResponse])
def get_races(db: Session = Depends(get_db)):
    races = db.query(Race).order_by(Race.race_number).all()
    response = []
    for race in races:
        results = None
        if race.status == 'completed':
            results = get_race_results_with_driver(db, race.id)
        response.append(RaceResponse(
            id=race.id,
            name=race.name,
            race_number=race.race_number,
            date=race.date,
            is_special_race=race.is_special_race,
            status=race.status,
            rollover_first=race.rollover_first,
            rollover_second=race.rollover_second,
            rollover_last=race.rollover_last,
            rollover_stage1=race.rollover_stage1,
            rollover_stage2=race.rollover_stage2,
            extra_drivers=race.extra_drivers,
            results=results,
            created_at=str(race.created_at)
        ))
    return response


@router.get("/races/{race_id}", response_model=RaceResponse)
def get_race(race_id: int, db: Session = Depends(get_db)):
    race = db.query(Race).filter(Race.id == race_id).first()
    if not race:
        raise HTTPException(status_code=404, detail="Race not found")

    results = None
    if race.status == 'completed':
        results = get_race_results_with_driver(db, race.id)

    return RaceResponse(
        id=race.id,
        name=race.name,
        race_number=race.race_number,
        date=race.date,
        is_special_race=race.is_special_race,
        status=race.status,
        rollover_first=race.rollover_first,
        rollover_second=race.rollover_second,
        rollover_last=race.rollover_last,
        rollover_stage1=race.rollover_stage1,
        rollover_stage2=race.rollover_stage2,
        extra_drivers=race.extra_drivers,
        results=results,
        created_at=str(race.created_at)
    )


@router.post("/races")
def create_race(req: CreateRaceRequest, db: Session = Depends(get_db)):
    race = Race(
        name=req.name,
        race_number=req.race_number,
        date=req.date,
        is_special_race=req.is_special_race,
        extra_drivers=req.extra_drivers,
        status='upcoming'
    )
    db.add(race)
    db.commit()
    db.refresh(race)
    return {"id": race.id}


@router.put("/races/{race_id}")
def update_race(race_id: int, req: CreateRaceRequest, db: Session = Depends(get_db)):
    race = db.query(Race).filter(Race.id == race_id).first()
    if not race:
        raise HTTPException(status_code=404, detail="Race not found")

    race.name = req.name
    race.date = req.date
    race.is_special_race = req.is_special_race
    race.extra_drivers = req.extra_drivers
    db.commit()
    return {"message": "Race updated"}


@router.post("/races/{race_id}/generate-teams")
def generate_teams(race_id: int, db: Session = Depends(get_db)):
    # Check if teams already exist
    existing = db.query(Team).filter(Team.race_id == race_id).count()
    if existing > 0:
        raise HTTPException(status_code=400, detail="Teams already generated for this race")

    # Get participants and drivers
    participants = db.query(Participant).order_by(Participant.id).all()
    drivers = db.query(Driver).filter(Driver.is_chartered == True).all()

    if len(drivers) < len(participants) * 4:
        raise HTTPException(status_code=400, detail="Not enough chartered drivers")

    # Shuffle drivers
    driver_ids = [d.id for d in drivers]
    random.shuffle(driver_ids)

    # Assign 4 drivers per participant
    for i, participant in enumerate(participants):
        start = i * 4
        team = Team(
            race_id=race_id,
            participant_id=participant.id,
            driver1_id=driver_ids[start],
            driver2_id=driver_ids[start + 1],
            driver3_id=driver_ids[start + 2],
            driver4_id=driver_ids[start + 3]
        )
        db.add(team)

    # Update race status
    race = db.query(Race).filter(Race.id == race_id).first()
    race.status = 'in_progress'
    db.commit()

    return {"message": "Teams generated successfully"}


@router.get("/races/{race_id}/teams", response_model=List[TeamResponse])
def get_race_teams(race_id: int, db: Session = Depends(get_db)):
    teams = db.query(Team).filter(Team.race_id == race_id).order_by(Team.points_earned.desc()).all()
    response = []
    for t in teams:
        response.append(TeamResponse(
            id=t.id,
            race_id=t.race_id,
            participant_id=t.participant_id,
            driver1_id=t.driver1_id,
            driver2_id=t.driver2_id,
            driver3_id=t.driver3_id,
            driver4_id=t.driver4_id,
            points_earned=t.points_earned,
            participant_name=t.participant.name,
            driver1_name=t.driver1.name,
            driver1_number=t.driver1.car_number,
            driver2_name=t.driver2.name,
            driver2_number=t.driver2.car_number,
            driver3_name=t.driver3.name,
            driver3_number=t.driver3.car_number,
            driver4_name=t.driver4.name,
            driver4_number=t.driver4.car_number,
            created_at=str(t.created_at)
        ))
    return response


@router.post("/races/{race_id}/results")
def enter_race_results(race_id: int, req: RaceResultsRequest, db: Session = Depends(get_db)):
    race = db.query(Race).filter(Race.id == race_id).first()
    if not race:
        raise HTTPException(status_code=404, detail="Race not found")

    # Get chartered car numbers
    chartered_cars = {d.car_number for d in db.query(Driver).filter(Driver.is_chartered == True).all()}

    # Delete existing results
    db.query(RaceResult).filter(RaceResult.race_id == race_id).delete()

    # Insert results
    scoring_results = [
        (req.first_place_car_number, True, False, False, False, False),
        (req.second_place_car_number, False, True, False, False, False),
        (req.last_place_car_number, False, False, True, False, False),
        (req.stage1_winner_car_number, False, False, False, True, False),
        (req.stage2_winner_car_number, False, False, False, False, True),
    ]

    for car_num, is_first, is_second, is_last, is_stage1, is_stage2 in scoring_results:
        existing = db.query(RaceResult).filter(
            RaceResult.race_id == race_id,
            RaceResult.car_number == car_num
        ).first()

        if existing:
            existing.is_first_place = existing.is_first_place or is_first
            existing.is_second_place = existing.is_second_place or is_second
            existing.is_last_place = existing.is_last_place or is_last
            existing.is_stage1_winner = existing.is_stage1_winner or is_stage1
            existing.is_stage2_winner = existing.is_stage2_winner or is_stage2
        else:
            result = RaceResult(
                race_id=race_id,
                car_number=car_num,
                is_first_place=is_first,
                is_second_place=is_second,
                is_last_place=is_last,
                is_stage1_winner=is_stage1,
                is_stage2_winner=is_stage2
            )
            db.add(result)

    race.status = 'completed'

    # Calculate multiplier
    multiplier = SPECIAL_MULTIPLIER if race.is_special_race else 1

    # Calculate available points (with rollover)
    available_first = (POINTS_FIRST + race.rollover_first) * multiplier
    available_second = (POINTS_SECOND + race.rollover_second) * multiplier
    available_last = (POINTS_LAST + race.rollover_last) * multiplier
    available_stage1 = (POINTS_STAGE1 + race.rollover_stage1) * multiplier
    available_stage2 = (POINTS_STAGE2 + race.rollover_stage2) * multiplier

    # Track rollover for non-chartered cars
    next_rollover_first = 0
    next_rollover_second = 0
    next_rollover_last = 0
    next_rollover_stage1 = 0
    next_rollover_stage2 = 0

    if req.first_place_car_number not in chartered_cars:
        next_rollover_first = POINTS_FIRST + race.rollover_first
        available_first = 0
    if req.second_place_car_number not in chartered_cars:
        next_rollover_second = POINTS_SECOND + race.rollover_second
        available_second = 0
    if req.last_place_car_number not in chartered_cars:
        next_rollover_last = POINTS_LAST + race.rollover_last
        available_last = 0
    if req.stage1_winner_car_number not in chartered_cars:
        next_rollover_stage1 = POINTS_STAGE1 + race.rollover_stage1
        available_stage1 = 0
    if req.stage2_winner_car_number not in chartered_cars:
        next_rollover_stage2 = POINTS_STAGE2 + race.rollover_stage2
        available_stage2 = 0

    # Update next race with rollover
    if any([next_rollover_first, next_rollover_second, next_rollover_last, next_rollover_stage1, next_rollover_stage2]):
        next_race = db.query(Race).filter(Race.race_number == race.race_number + 1).first()
        if next_race:
            next_race.rollover_first = next_rollover_first
            next_race.rollover_second = next_rollover_second
            next_race.rollover_last = next_rollover_last
            next_race.rollover_stage1 = next_rollover_stage1
            next_race.rollover_stage2 = next_rollover_stage2

    # Calculate team points
    teams = db.query(Team).filter(Team.race_id == race_id).all()
    for team in teams:
        car_numbers = [team.driver1.car_number, team.driver2.car_number,
                       team.driver3.car_number, team.driver4.car_number]
        points = 0

        if req.first_place_car_number in car_numbers:
            points += available_first
        if req.second_place_car_number in car_numbers:
            points += available_second
        if req.last_place_car_number in car_numbers:
            points += available_last
        if req.stage1_winner_car_number in car_numbers:
            points += available_stage1
        if req.stage2_winner_car_number in car_numbers:
            points += available_stage2

        team.points_earned = points

    db.commit()
    return {"message": "Race results entered and points calculated"}


@router.get("/standings", response_model=List[StandingResponse])
def get_standings(db: Session = Depends(get_db)):
    # Get all participants
    participants = db.query(Participant).order_by(Participant.name).all()

    # Build standings
    standings_data = []
    for p in participants:
        # Get total points from completed races
        total_points = db.query(func.coalesce(func.sum(Team.points_earned), 0))\
            .join(Race, Team.race_id == Race.id)\
            .filter(Team.participant_id == p.id, Race.status == 'completed')\
            .scalar() or 0

        # Count completed races
        races_completed = db.query(func.count(Team.id))\
            .join(Race, Team.race_id == Race.id)\
            .filter(Team.participant_id == p.id, Race.status == 'completed')\
            .scalar() or 0

        standings_data.append({
            'participant_id': p.id,
            'participant_name': p.name,
            'total_points': int(total_points),
            'races_completed': int(races_completed)
        })

    # Sort by points descending
    standings_data.sort(key=lambda x: (-x['total_points'], x['participant_name']))

    standings = []
    for rank, data in enumerate(standings_data, 1):
        standings.append(StandingResponse(
            participant_id=data['participant_id'],
            participant_name=data['participant_name'],
            total_points=data['total_points'],
            races_completed=data['races_completed'],
            rank=rank
        ))
    return standings
