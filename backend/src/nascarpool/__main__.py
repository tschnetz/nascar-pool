#!/usr/bin/env python3
"""
NASCAR Pool - Desktop App Launcher

Starts the FastAPI backend and opens it in a native Toga window.
"""
import threading
import time

import toga
from toga.style import Pack
from toga.style.pack import COLUMN
import uvicorn
from .app.main import app as fastapi_app
from .app.database import init_db


def start_server():
    """Start the FastAPI server in a background thread"""
    uvicorn.run(fastapi_app, host="127.0.0.1", port=8080, log_level="warning")


class NASCARPool(toga.App):
    def startup(self):
        # Initialize database
        init_db()

        # Start server in background thread
        server_thread = threading.Thread(target=start_server, daemon=True)
        server_thread.start()

        # Wait for server to be ready
        time.sleep(1.5)

        # Create WebView
        self.webview = toga.WebView(
            style=Pack(flex=1)
        )

        # Create main window
        self.main_window = toga.MainWindow(title="NASCAR Pool", size=(1024, 768))
        self.main_window.content = toga.Box(
            children=[self.webview],
            style=Pack(direction=COLUMN, flex=1)
        )
        self.main_window.show()

        # Load URL after window is shown
        self.webview.url = "http://127.0.0.1:8080"


def main():
    return NASCARPool("NASCAR Pool", "com.nascarpool")


if __name__ == "__main__":
    main().main_loop()
