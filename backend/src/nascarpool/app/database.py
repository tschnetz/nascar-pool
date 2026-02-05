import os
import logging
from pathlib import Path

from dotenv import load_dotenv
from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker
from .models import Base

# Load .env - try bundled resources first, then dev path
# Bundled: app -> nascarpool -> resources
bundled_env = Path(__file__).parent.parent / 'resources' / '.env'
# Dev: app -> nascarpool -> src -> backend -> NASCAR
dev_env = Path(__file__).parent.parent.parent.parent.parent / '.env'

if bundled_env.exists():
    load_dotenv(bundled_env)
elif dev_env.exists():
    load_dotenv(dev_env)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Database environment tracking
current_env = None
engine = None
SessionLocal = None


def get_database_url():
    """Get database URL, preferring PROD over LOCAL"""
    global current_env

    # Check for production database first
    prod_url = os.getenv('PROD_DATABASE_URL')
    if prod_url:
        current_env = 'PROD'
        logger.info("Connecting to PostgreSQL (PROD)...")
        return prod_url

    # Fall back to local database
    local_url = os.getenv('DATABASE_URL')
    if local_url:
        current_env = 'LOCAL'
        logger.info("Connecting to PostgreSQL (LOCAL)...")
        return local_url

    raise ValueError("No database URL configured (set DATABASE_URL or PROD_DATABASE_URL)")


def _init_engine():
    """Initialize the database engine"""
    global engine, SessionLocal
    if engine is None:
        DATABASE_URL = get_database_url()
        # SQLAlchemy requires postgresql:// not postgres://
        if DATABASE_URL.startswith('postgres://'):
            DATABASE_URL = DATABASE_URL.replace('postgres://', 'postgresql://', 1)
        engine = create_engine(DATABASE_URL)
        SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def init_db():
    """Initialize database connection (tables already exist from Go migrations)"""
    _init_engine()
    # Verify connection works
    try:
        with engine.connect() as conn:
            conn.execute(text("SELECT 1"))
        logger.info(f"PostgreSQL connected successfully ({current_env})")
    except Exception as e:
        logger.error(f"Failed to connect to database: {e}")
        raise


def get_db():
    """Dependency to get database session"""
    _init_engine()
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
