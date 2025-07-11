# Good-Blast-REST-API

Good Blast is a casual Match-3 game played by millions of users. This REST API powers the competitive features of the game, allowing users to join daily tournaments, track their progress, and see how they rank globally and locally.

## Features

### User Management
- `CreateUser`: Register a new user with a unique username.
- `UpdateProgress`: Update user level and coin progress.

### Tournament System
- Daily tournaments (00:00–23:59 UTC).
- Users can join if they:
  - Have completed at least 10 levels.
  - Spend 500 coins.
  - Join before 12:00 UTC.
- Grouped into groups of 35 players.
- Tournament score increases by 1 per level progressed.

#### Rewards
- 🥇 1st place: 5000 coins  
- 🥈 2nd place: 3000 coins  
- 🥉 3rd place: 2000 coins  
- 4th–10th: 1000 coins

### Leaderboards
- `GetGlobalLeaderboard`: Top 1000 players globally.
- `GetCountryLeaderboard`: Top 1000 players in your country.
- `GetTournamentLeaderboard`: Your tournament's leaderboard.
- `GetTournamentRank`: Your rank in the current tournament.

## Technologies Used

- Language: Go (Golang)
- Database: PostgreSQL
- Cache: Redis
- API Style: REST