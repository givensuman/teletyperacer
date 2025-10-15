# teletyperacer TODO

## High Priority

- [ ] **setup_modules**: Set up basic Go modules and dependencies for client and server (including websocket libraries)
- [ ] **home_screen**: Implement home screen in client with options to host private/public game or join private/browse public games
- [ ] **lobby_private_host**: Create private lobby screen for hosting (generate room code, wait for players)
- [ ] **lobby_public_host**: Create public lobby screen for hosting (no code, lobby appears in public browse list)
- [ ] **lobby_private_join**: Create screen for joining private lobby (enter room code, connect)
- [ ] **lobby_public_browse**: Create screen for browsing and joining public lobbies (list available public games)
- [x] **server_websockets_setup**: Set up websocket server infrastructure
- [x] **server_websockets_room_create**: Implement room creation logic on server (private and public)
- [ ] **server_websockets_player_join**: Implement player joining logic on server (validate codes for private, add to public)
- [ ] **server_websockets_lobby_broadcast**: Implement lobby updates broadcasting (player list, ready status)

## Medium Priority

- [ ] **game_screen_prompt**: Implement game screen UI with typing prompt display
- [ ] **game_screen_progress_self**: Implement progress bar for player's own typing progress
- [ ] **game_screen_progress_opponents**: Implement progress bars for opponents' typing progress
- [ ] **typing_logic_input**: Add typing input logic on client (capture keystrokes, calculate progress)
- [ ] **typing_logic_updates**: Add progress update sending via websocket
- [ ] **server_game_state_tracking**: Implement server-side game state tracking (player positions, game status)
- [ ] **server_game_state_broadcast**: Implement progress update broadcasting to all players in room
- [ ] **sticky_connections**: Add cookie-based session management for sticky connections
- [ ] **game_start_logic**: Implement game start logic (host initiates, server broadcasts start signal)

## Low Priority

- [ ] **game_end_detection**: Add game completion detection and winner calculation
- [ ] **game_end_results**: Implement results screen with winner announcement and stats
- [ ] **docker_deployment**: Configure Docker and docker-compose for client and server deployment
- [ ] **makefile_targets**: Add Makefile targets for building, running, and deploying
