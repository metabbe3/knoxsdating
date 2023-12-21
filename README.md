# Dating App Backend

## Introduction

### Purpose
This document outlines the features and technical specifications for the development of a Dating App Backend.

### Scope
The backend covers user authentication, profile management, swiping, premium features, and additional enhancements.

### Improvement
Consider adding RabbitMQ, increasing unit testing, comprehensive logging, and a direct messaging feature.

## Features

### User Authentication
- Sign up securely.
- User registration with a unique username and valid email.
- Password hashing for security.
- JWT-based authentication for sessions.
- Register for new users.

### Profile Management
- Manage dating profiles.
- Upload 1-9 photos.
- Add a description (max 500 characters).
- Specify interests, goals, height, language, Zodiac sign, education, and more.

### Swiping Profiles
- Discover and express interest.
- View a limited number of profiles daily.
- Swipe left (pass) or right (like).
- Avoid showing profiles twice daily.

### Premium Features
- Enhance experience with premium packages.
- Example Premium Features: Unlimited Likes, See Who Likes You, Unlimited Rewinds, Passport (Travel Mode), Direct Message Matched User.

### Direct Messaging
- Send messages to matched profiles.
- Real-time chat with text, emojis, and multimedia.

### Notification Features
- Receive notifications for app activities.
- NotificationType (e.g., "New Message," "Matched," "Premium Feature Unlocked").
- Real-time delivery with timestamps.
- View Notification History.
- Mark Notifications as Read manually or in bulk.
- Option to mark all as read.
- Notification Triggers: New messages, matches, premium features, and activity reminders.

## Non-functional Requirements

### Performance
- Handle a minimum of [X] requests per second.
- Use Redis for caching frequently accessed data.

### Security
- Securely store and transmit user data.
- Implement measures against abusive behavior.

### Scalability
- Use RabbitMQ for handling asynchronous tasks.
- Implement Nginx as a load balancer.

## Tech Stacks and Justification

- Backend Language: Golang (High performance, simplicity, efficiency).
- Database: PostgreSQL (Robust, open-source).
- ORM: GORM (Simplifies database interactions, PostgreSQL support).
- Authentication: JWT (Stateless, secure).
- Cache: Redis (Caching data).
- Load Balancer: Nginx (Distribute incoming traffic).

## Installation with Docker

1. Clone the repository:
   ```bash
   git clone [repository_url]
   cd [repository_directory]

2. Build the Docker containers:
   ```bash
   docker-compose build

3. Run the application:
   ```bash
    docker-compose up -d

4. Shut down the application:
   ```bash
   docker-compose down

5. Restart the application (in case of errors or changes):
   ```bash
   docker-compose restart

Replace `[repository_url]` and `[repository_directory]` with your actual repository URL and directory. Users can follow these instructions to clone, build, run, shut down, and restart the Dating App Backend using Docker.