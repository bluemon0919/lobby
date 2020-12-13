# lobby

This is a sample repository in the game lobby.  
Based on a one-on-one game, we will implement a mechanism to create a pair of opponents and WebSocket communication in the game.

# Try out

Run the following command to set up the server.

```bash
go run *.go
```

Then connect to localhost: 8080 to see how it works.

```bash
open http://localhost:8080
```

Enter the player name.  
It should lead to a lobby page waiting for your opponent.

Open another browser, access localhost, enter the player name, and the game will start.  
However, please note that it will be recognized as the same user as sharing the cookie. You just need to access it as another user.
