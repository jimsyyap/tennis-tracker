# Tennis Tracker

A comprehensive web application for tennis players to track their unforced errors during matches and practice sessions.

## Overview

Tennis Tracker allows players to log, analyze, and visualize their unforced errors across different matches and practice sessions. The application provides insights through interactive charts and supports sharing data with coaches for feedback and improvement.

## Features

- **User Authentication**: Secure sign-up, login, and password reset functionality
- **Session Management**: Create and manage tennis matches and practice sessions
- **Error Tracking**: Log and analyze unforced errors for each session
- **Data Visualization**: Interactive charts to visualize progress and patterns over time
- **Mobile Support**: Responsive design that works on mobile, tablet, and desktop devices
- **Offline Functionality**: Continue tracking errors even without an internet connection
- **Data Sharing**: Share session data with coaches via secure links

## Tech Stack

### Backend
- **Language**: Go (Golang)
- **Database**: PostgreSQL
- **Authentication**: JWT-based authentication
- **API**: RESTful API

### Frontend
- **Framework**: React.js
- **Styling**: Tailwind CSS
- **Charts**: Recharts/Chart.js
- **State Management**: React Context API and hooks
- **Offline Support**: IndexedDB/localStorage

### Deployment
- **Backend Hosting**: Google Cloud (Cloud Run + Cloud SQL)
- **Frontend Hosting**: Vercel
- **CI/CD**: GitHub Actions

## Project Structure

### Backend
```
tennis-tracker-backend/
├── cmd/
│   └── server/
│       └── main.go        # Entry point for the application
├── internal/
│   ├── api/               # API handlers
│   ├── middleware/        # Middleware functions
│   ├── models/            # Database models
│   ├── database/          # Database connection and migration
│   └── services/          # Business logic
├── pkg/                   # Shared packages
├── config/                # Configuration files
├── go.mod                 # Go module file
└── go.sum                 # Go dependencies
```

### Frontend
```
tennis-tracker-frontend/
├── public/                # Static files
├── src/
│   ├── components/        # Reusable components
│   ├── pages/             # Page components
│   ├── services/          # API services
│   ├── utils/             # Utility functions
│   ├── context/           # React context providers
│   ├── hooks/             # Custom React hooks
│   ├── App.js             # Main application component
│   └── index.js           # Entry point
├── package.json           # NPM package configuration
└── tailwind.config.js     # Tailwind CSS configuration
```

## Development Roadmap

### Phase 1: Project Setup and Core Backend (Weeks 1-2)
- [x] Initialize project structure
- [x] Set up PostgreSQL database
- [ ] Implement user authentication
- [ ] Create core API endpoints

### Phase 2: Core Frontend and Session Management (Weeks 3-4)
- [ ] Set up React project
- [ ] Build authentication UI
- [ ] Implement session management
- [ ] Create basic UI components

### Phase 3: Error Tracking and Data Visualization (Weeks 5-6)
- [ ] Implement error logging
- [ ] Create data visualization components
- [ ] Build error summary features
- [ ] Add interactive filters

### Phase 4: Mobile Support and Offline Functionality (Weeks 7-8)
- [ ] Optimize for mobile devices
- [ ] Implement offline data storage
- [ ] Create synchronization logic
- [ ] Test offline-to-online workflow

### Phase 5: Data Sharing and Final Testing (Weeks 9-10)
- [ ] Implement shareable links
- [ ] Create shared view components
- [ ] Perform comprehensive testing
- [ ] Fix bugs and optimize performance

### Phase 6: Deployment and Documentation (Weeks 11-12)
- [ ] Set up Google Cloud environment
- [ ] Configure CI/CD pipeline
- [ ] Deploy application
- [ ] Create documentation

## Getting Started

### Prerequisites
- Go 1.18 or later
- PostgreSQL 13 or later
- Node.js 16 or later
- npm 8 or later

### Backend Setup
1. Clone the repository
   ```bash
   git clone https://github.com/yourusername/tennis-tracker.git
   cd tennis-tracker/backend
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Set up the database
   ```bash
   psql -U postgres -f database/migrations/init.sql
   ```

4. Configure environment variables
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. Run the server
   ```bash
   go run cmd/server/main.go
   ```

### Frontend Setup
1. Navigate to the frontend directory
   ```bash
   cd ../frontend
   ```

2. Install dependencies
   ```bash
   npm install
   ```

3. Configure environment variables
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. Run the development server
   ```bash
   npm start
   ```

## API Documentation

The API documentation is available at `/api/docs` when running the server locally.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Go](https://golang.org/)
- [React](https://reactjs.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [Tailwind CSS](https://tailwindcss.com/)
