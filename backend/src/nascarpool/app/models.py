from sqlalchemy import Column, Integer, String, Boolean, DateTime, ForeignKey, Text
from sqlalchemy.orm import relationship, declarative_base
from sqlalchemy.sql import func

Base = declarative_base()

# Scoring constants
POINTS_FIRST = 135
POINTS_SECOND = 25
POINTS_LAST = 15
POINTS_STAGE1 = 25
POINTS_STAGE2 = 25
SPECIAL_MULTIPLIER = 2


class Participant(Base):
    __tablename__ = 'participants'

    id = Column(Integer, primary_key=True)
    name = Column(String(100), nullable=False)
    created_at = Column(DateTime, server_default=func.now())

    teams = relationship("Team", back_populates="participant")


class Driver(Base):
    __tablename__ = 'drivers'

    id = Column(Integer, primary_key=True)
    name = Column(String(100), nullable=False)
    car_number = Column(String(10), unique=True, nullable=False)
    team_name = Column(String(100))
    manufacturer = Column(String(50))
    is_chartered = Column(Boolean, default=True)
    created_at = Column(DateTime, server_default=func.now())


class Race(Base):
    __tablename__ = 'races'

    id = Column(Integer, primary_key=True)
    name = Column(String(200), nullable=False)
    race_number = Column(Integer, nullable=False)
    date = Column(String(50))
    is_special_race = Column(Boolean, default=False)
    status = Column(String(20), default='upcoming')  # upcoming, in_progress, completed
    rollover_first = Column(Integer, default=0)
    rollover_second = Column(Integer, default=0)
    rollover_last = Column(Integer, default=0)
    rollover_stage1 = Column(Integer, default=0)
    rollover_stage2 = Column(Integer, default=0)
    extra_drivers = Column(Text)
    nascar_race_id = Column(Integer)
    created_at = Column(DateTime, server_default=func.now())

    results = relationship("RaceResult", back_populates="race")
    teams = relationship("Team", back_populates="race")


class RaceResult(Base):
    __tablename__ = 'race_results'

    id = Column(Integer, primary_key=True)
    race_id = Column(Integer, ForeignKey('races.id'), nullable=False)
    car_number = Column(String(10), nullable=False)
    position = Column(Integer)
    is_first_place = Column(Boolean, default=False)
    is_second_place = Column(Boolean, default=False)
    is_last_place = Column(Boolean, default=False)
    is_stage1_winner = Column(Boolean, default=False)
    is_stage2_winner = Column(Boolean, default=False)
    created_at = Column(DateTime, server_default=func.now())

    race = relationship("Race", back_populates="results")


class Team(Base):
    __tablename__ = 'teams'

    id = Column(Integer, primary_key=True)
    race_id = Column(Integer, ForeignKey('races.id'), nullable=False)
    participant_id = Column(Integer, ForeignKey('participants.id'), nullable=False)
    driver1_id = Column(Integer, ForeignKey('drivers.id'), nullable=False)
    driver2_id = Column(Integer, ForeignKey('drivers.id'), nullable=False)
    driver3_id = Column(Integer, ForeignKey('drivers.id'), nullable=False)
    driver4_id = Column(Integer, ForeignKey('drivers.id'), nullable=False)
    points_earned = Column(Integer, default=0)
    created_at = Column(DateTime, server_default=func.now())

    race = relationship("Race", back_populates="teams")
    participant = relationship("Participant", back_populates="teams")
    driver1 = relationship("Driver", foreign_keys=[driver1_id])
    driver2 = relationship("Driver", foreign_keys=[driver2_id])
    driver3 = relationship("Driver", foreign_keys=[driver3_id])
    driver4 = relationship("Driver", foreign_keys=[driver4_id])
