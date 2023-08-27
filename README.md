# Simple monorepo app for tracking user wallet balance for a gaming platform.

# Getting started
Locally run the following services: 
- API (entry service for all request)
- User (handles new user creation and retrieval)
- Grpc (handles a single gRCP request from API to retrieve a user balance)
- Wallet (handles deposit and withdraw operations on user wallet)
- Websocket (handles creation of websocket chanel to retranslate live updates on user wallet)

* The above services need redis running locally on port 6379

### API Endpoints
| HTTP Verbs | Endpoints | Action |
| --- | --- | --- |
| POST | /api/users | To create a new user |
| POST | /api/wallet/deposit | To deposit user a specific amount (This action works both on empty and existing wallets) |
| POST | /api/wallet/withdraw | To withdraw a specific amount from user wallet (Works only on positive amount wallets, returns error otherwise) |
| GET  | /api/wallet/balance/:userId | To retrieve user wallet balance via gRPC call to Wallet service |

### Websocket communication
Websocket is hosted on  ws://localhost:5104/ws with test webpase on localhost:5104. The test webpage can accept two event types ("gameprogress-event" and "leaderboard-game-event") which allow users to "subscribe" to a specific game event. This is a two-way communication between backend and frontend to choose a game event and receive messages on that specific event.

Things to improve:
- Fix docker-compose and Dockerfiles to properly host and deploy the platform
- Migrate user storage from redis to persistent storage like postgres (easier to maintain indexed data)
- Expand logging using zero-logs to further store and search logs via elastic
- Add user authentication to ensure secure wallet transaction operations (not to easily manipulate wallet state for other users)
- Due to user db and wallet db using the same redis instance, and having the same key (userId), wallet operations overwrite user entries in redis
- Improve usage of environment variables, currently hardcoded as end-default values (good enought for this project state)
- We can further add security to websocket chanel to limit accessing websocket (configure handshaking via SHA encryption?)

Time spent on this project: around 15 hours.