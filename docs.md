# Websocket Packet Documentation

## Initial
**Client**
Creating a new game
```json
{
	"type": "create_game",
	"data": {
		"name": "<name-of-game>",
		// "map_size": [],
		"self_starting": true // or false -> is creator starting 
	}
}
```
Joining a game
```json
{
	"type": "join_game",
	"data": {
		"game_id": "<game-id>",
		"username": "<name-of-user>"
	}
}
``` 
**Server**
```json
{
	"type": "greeting",
	"data": {
		"game_id": "<id-of-game>"
	}
}
```

## 0. Game Start
**Server -> Clients**
```json
{
	"type": "game_start",
	"data": {
		"map_size": ["<rows>", "<columns>"],
		"name": "<name of game>",
		"starting": true, // or false - this client starts,
		"enemy": {
			"name": "<name-of-enemy>"
		}
	}

}

```

## 1. Ship placement
**Client A -> Server**
```json
{
	"type": "place",
	"data": {
		"ships": [
			{
				"ship_type": "<ship-type>",
				"from": ["x1", "y1"],
				"to": ["x2", "y2"]
			},
			{
				"ship_type": "<ship-type>",
				"from": ["x1", "y1"],
				"to": ["x2", "y2"]
			},
		]
		
	}
	
}
```
**Server -> Client B**
```json
{
	"type": "enemy_ships",
	"data": 

	
}
```

## 2. Shooting a cell
```json
???
```

## 3. Game finished message
```json
???
```