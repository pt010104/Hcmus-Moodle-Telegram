# Moodle HCmus Event Calendar and Notifications - Golang

This repository is designed to fetch event calendars and send notifications on Moodle HCmus using Golang.

## Main Features:
- **Fetch Calendar Events**: Retrieve and notify users about all calendar events for the entire month.
- **Daily Notifications**: Pop-up notifications every day about upcoming events.
- **Deadline Alerts**: Alarm notifications 2 hours before any deadline.

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/moodle-calendar-notifications.git
   ```
2. Navigate to the project directory:
   ```bash
   cd moodle-calendar-notifications
   ```
3. Install dependencies and run:
   ```bash
   go mod tidy
   go run main.go
   ```
