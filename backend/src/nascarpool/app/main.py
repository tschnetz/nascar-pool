from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse
import os

from .database import init_db
from .routes import router

app = FastAPI(title="NASCAR Pool API")

# CORS for development
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include API routes
app.include_router(router)

# Serve static files (Svelte build) - bundled in resources folder
# Path: app -> nascarpool -> resources
STATIC_DIR = os.path.join(os.path.dirname(__file__), '..', 'resources')
if os.path.exists(STATIC_DIR) and os.path.exists(os.path.join(STATIC_DIR, 'index.html')):
    # SvelteKit uses _app for assets
    app_dir = os.path.join(STATIC_DIR, '_app')
    if os.path.exists(app_dir):
        app.mount("/_app", StaticFiles(directory=app_dir), name="app_assets")

    @app.get("/")
    async def serve_index():
        return FileResponse(os.path.join(STATIC_DIR, 'index.html'))

    @app.get("/{full_path:path}")
    async def serve_spa(full_path: str):
        # Don't intercept API routes
        if full_path.startswith("api/"):
            return None
        # Try to serve the file directly
        file_path = os.path.join(STATIC_DIR, full_path)
        if os.path.isfile(file_path):
            return FileResponse(file_path)
        # Fall back to index.html for SPA routing
        return FileResponse(os.path.join(STATIC_DIR, 'index.html'))


@app.on_event("startup")
def on_startup():
    init_db()


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8080)
