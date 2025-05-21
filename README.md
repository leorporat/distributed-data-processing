# Reddit Summarization & Sentiment Analysis Platform

A comprehensive platform for Reddit content analysis that enables summarization of Reddit discussions and sentiment analysis of posts and comments, built with Go, React, gRPC, and Kafka.

## Overview

This project combines:
1. **Content Summarization**: Quickly get concise summaries of Reddit discussions based on search queries
2. **Sentiment Analysis**: Analyze the sentiment and tone of Reddit posts and comments
3. **Distributed Processing**: Scalable architecture for handling large volumes of Reddit data

The platform uses Go's powerful concurrency model for the backend, React for the frontend, gRPC for efficient communication, and is designed with a microservices architecture that can incorporate Kafka for stream processing.

## Features

- Search and summarize Reddit discussions
- Real-time data processing from Reddit API
- Sentiment analysis capabilities for posts and comments
- Modern React frontend with responsive design
- Efficient backend using Go and gRPC
- Scalable microservices architecture
- JSON and Protobuf data format support

## Setup

### Prerequisites
- Go 1.16+
- Node.js 16+
- npm 7+
- Reddit API credentials

### Backend Setup

1. Copy the environment template:
   ```bash
   cd backend
   cp .env.example .env
   ```

2. Get Reddit API credentials:
   - Go to https://www.reddit.com/prefs/apps
   - Create a new script-type application
   - Fill in the .env file with your credentials

### Frontend Setup

1. Copy the environment template:
   ```bash
   cd frontend
   cp .env.example .env.local
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

## Development

Start the backend server:
```bash
cd backend
source .env
cd cmd/reddit-server
go run main.go
```

Start the frontend development server:
```bash
cd frontend
npm run dev
```

Visit http://localhost:3000 to use the application.

## Architecture

The system is designed as a set of microservices:
- **Reddit API Client**: Interfaces with Reddit to fetch data
- **Summarization Service**: Processes and summarizes Reddit content
- **Sentiment Analysis**: Analyzes the emotional tone of content
- **Frontend**: Provides user interface and visualization

Services communicate using gRPC for efficient, typed interactions, with JSON fallback for broader compatibility.

## Security Notes

- Never commit .env files containing real credentials
- Keep your Reddit API credentials secure
- Use environment variables for all sensitive data
- Follow the principle of least privilege for API access

## Future Enhancements

- Full Kafka integration for stream processing
- Enhanced summarization using NLP techniques
- More detailed sentiment analysis visualization
- User authentication and saved searches
