# NASCAR Pool Manager

A desktop app to manage a NASCAR fantasy pool for 9 friends. Replaces manual Excel tracking with a mobile-first dashboard.

## How It Works

- **9 participants** compete each race week
- **36 chartered drivers** (permanent car numbers) are randomly assigned into teams of 4
- Points awarded for race results: 1st, 2nd, Last, Stage 1 winner, Stage 2 winner
- **Car number-based scoring** - results entered by car number, not driver name (handles substitute drivers)
- **Rollover system** - if a non-chartered driver (extra in field) wins points, those points roll to next race

## Scoring System

| Position | Base Points |
|----------|-------------|
| 1st Place | 135 |
| 2nd Place | 25 |
| Last Place | 15 |
| Stage 1 Winner | 25 |
| Stage 2 Winner | 25 |

- **Special races** (Daytona 500, championships, etc.) = 2x points
- Points displayed as dollars in the UI

## Development vs Production

### Local Development (default):
- Uses local PostgreSQL database (`postgres://localhost:5432/nascar_pool`)
- Safe to experiment with data without affecting production

### Production:
- PostgreSQL on Render (Ohio region)

### Environment Variables (.env):
```
# Database - PROD_DATABASE_URL takes precedence if set
DATABASE_URL=postgres://localhost:5432/nascar_pool
# PROD_DATABASE_URL=postgres://user:pass@host:5432/db
```

### Startup Output:
Backend shows database environment on startup:
```
Connecting to PostgreSQL (LOCAL)...
PostgreSQL connected successfully (LOCAL)
```
or `(PROD)` when using production database.

## Quick Start

### Desktop App (pywebview):
```bash
cd frontend && npm run build    # Build frontend (only when changed)
cd python_backend && python3 run.py   # Launch desktop app
```

### Development (separate servers):
```bash
cd python_backend && python3 -m uvicorn app.main:app --reload --port 8080
cd frontend && npm run dev      # Frontend on :5173
```

## Architecture

- **Frontend**: Svelte 5 + Tailwind CSS (mobile-first)
- **Backend**: Python FastAPI + SQLAlchemy
- **Database**: PostgreSQL
- **Desktop**: pywebview (native window wrapper)

## Admin Mode

Admin controls are hidden by default. Add `?admin=true` to any URL to enable:
- Create/edit races
- Generate random teams
- Enter race results

The admin indicator appears in the header when active. Navigation preserves the admin param.

## Key Concepts

### Car Numbers are Permanent
The `car_number` field identifies the car, not the driver. If a substitute driver races, results still credit the car number. This is why results are entered by car number.

### Chartered vs Non-Chartered Drivers
- 36 chartered drivers form the pool teams (4 drivers per participant × 9 participants)
- Extra drivers in some races (40+ car fields) are non-chartered
- If a non-chartered driver scores, those points **roll over** to the next race's same position

### Team Generation
Teams are randomly generated each race:
1. Shuffle the 36 chartered drivers
2. Assign 4 drivers to each of the 9 participants
3. Teams are fixed for that race

### Weekly Workflow
1. Thursday: Entry list comes out
2. Admin generates random teams
3. Admin shares worksheet (printable page)
4. Sunday: Race happens
5. Admin enters results (by car number)
6. Standings update automatically

## Database Schema

### Core Tables
- `participants` - The 9 pool members
- `drivers` - 36 chartered NASCAR drivers (car_number, team_name, manufacturer)
- `races` - 36 race schedule with rollover tracking
- `race_teams` - Random team assignments per race
- `race_results` - Scoring results per race

### Key Fields
- `races.is_special_race` - 2x points multiplier
- `races.rollover_*` - Accumulated rollover points per position
- `race_results.car_number` - Links to driver by car, not by ID

## API Endpoints

```
GET  /api/participants     - List all participants
GET  /api/drivers          - List all drivers (sorted by car number)
GET  /api/races            - List all races
POST /api/races            - Create new race
GET  /api/races/:id        - Get race details with results
PUT  /api/races/:id        - Update race
POST /api/races/:id/generate-teams  - Generate random teams
GET  /api/races/:id/teams  - Get teams for a race
POST /api/races/:id/results - Enter race results
GET  /api/standings        - Get current standings
```

## Frontend Routes

- `/` - Standings (leaderboard)
- `/races` - Race list
- `/races/:id` - Race detail with teams and results
- `/races/:id/worksheet` - Printable worksheet
- `/schedule` - Full season schedule with TV/track info
- `/drivers` - Driver info with manufacturer filter

## File Structure

```
python_backend/
  app/
    __init__.py
    main.py           - FastAPI app, static file serving
    database.py       - PostgreSQL connection (LOCAL/PROD)
    models.py         - SQLAlchemy models
    routes.py         - API endpoints
  run.py              - pywebview desktop launcher
  requirements.txt

frontend/
  src/
    lib/api.js        - API client
    routes/           - SvelteKit pages
    app.css           - Tailwind + custom styles
  build/              - Static build for desktop app
```

## Constraints

- Single admin user (the developer) - no auth system needed
- Mobile-first design with 44px minimum touch targets
- Car numbers are the source of truth for scoring
- Never delete race results - only add/update
