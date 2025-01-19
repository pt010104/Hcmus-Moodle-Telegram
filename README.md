# Moodle HCmus Event Calendar and Notifications - Golang

This repository is designed to fetch event calendars and send notifications on Moodle HCMUS using Golang.

## Main Features:
- **Fetch Calendar Events**: Retrieve and notify users about all calendar events for the entire month.
- **Daily Notifications**: Pop-up notifications every day about upcoming events.
- **Deadline Alerts**: Alarm notifications 24, 12, 6, 3, 1 hours before any deadline.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/pt010104/Hcmus-Moodle-Telegram.git
   ```

2. Navigate to the project directory:
   ```bash
   cd Hcmus-Moodle-Telegram
   ```

3. Initialize environment variables:
   Create a `.env` file in the root directory and add the following configurations:
   ```bash
   APP_VERSION=1.0.0
   APP_PORT=8080
   LOGGER_LEVEL=debug
   LOGGER_MODE=development
   LOGGER_ENCODING=console
   HCMUS_URL=courses.ctda.hcmus.edu.vn
   HCMUS_USERNAME
   HCMUS_PASSWORD
   HCMUS_SESSKEY
   HCMUS_COOKIES
   MONGODB_DATABASE
   MONGODB_URI
   TELEGRAM_CHAT_ID
   TELEGRAM_BOT_TOKEN
   ```

4. Install dependencies and run the application:
   ```bash
   go mod tidy
   go run cmd/main.go
   ```
