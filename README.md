# atlas-buddies
Mushroom game buddies Service

## Overview

A RESTful resource which provides buddies services.

### Features

- **Buddy List Management**: Create and retrieve character buddy lists with configurable capacity
- **Dynamic Capacity Updates**: Modify buddy list capacity with validation and constraints
- **Real-time Buddy Status**: Track buddy online status, channel, and shop presence
- **Event-Driven Architecture**: Kafka-based messaging for real-time updates and commands

## Environment

- JAEGER_HOST_PORT - Jaeger [host]:[port]
- LOG_LEVEL - Logging level - Panic / Fatal / Error / Warn / Info / Debug / Trace
- DB_USER - Postgres user name
- DB_PASSWORD - Postgres user password
- DB_HOST - Postgres Database host
- DB_PORT - Postgres Database port
- DB_NAME - Postgres Database name
- BOOTSTRAP_SERVERS - Kafka [host]:[port]
- BASE_SERVICE_URL - [scheme]://[host]:[port]/api/
- COMMAND_TOPIC_BUDDY_LIST - Kafka Topic for transmitting buddy list commands.
- COMMAND_TOPIC_INVITE - Kafka Topic for transmitting invite commands.
- EVENT_TOPIC_BUDDY_LIST_STATUS - Kafka Topic for transmitting buddy list status events.
- EVENT_TOPIC_CASH_SHOP_STATUS - Kafka Topic for receiving cash shop status events.
- EVENT_TOPIC_CHARACTER_STATUS - Kafka Topic for receiving character status events.
- EVENT_TOPIC_INVITE_STATUS - Kafka Topic for receiving invite status events.

## Kafka Integration

The service uses Kafka for asynchronous command processing and event publishing:

### Commands Consumed
- **UPDATE_CAPACITY**: Updates buddy list capacity with validation

### Events Published
- **BUDDY_LIST_STATUS**: Status updates for buddy list operations

### Message Processing
All commands are processed asynchronously with proper error handling and validation. Status events are published to notify other services of state changes.

## Database Schema

### Buddy Lists
- **ID**: UUID primary key
- **Tenant ID**: Multi-tenancy identifier
- **Character ID**: Owner character identifier (uint32)
- **Capacity**: Maximum buddy count (1-255, default 50)
- **Created At**: Timestamp of list creation

### Capacity Constraints
- Minimum capacity: 1 buddy
- Maximum capacity: 255 buddies
- New capacity must be >= current buddy count
- Capacity updates are atomic operations

## API

### Header

All RESTful requests require the supplied header information to identify the server instance.

```
TENANT_ID:083839c6-c47c-42a6-9585-76492795d123
REGION:GMS
MAJOR_VERSION:83
MINOR_VERSION:1
```

### Requests

#### [GET] Get Characters Buddy List

```/api/characters/{characterId}/buddy-list```

Example Response:
```json
{
  "data": {
    "type": "buddy-list",
    "id": "1",
    "attributes": {
      "characterId": 12345,
      "capacity": 50,
      "buddies": [
        {
          "characterId": 67890,
          "group": "Friends",
          "characterName": "MapleHero",
          "channelId": 1,
          "inShop": false,
          "pending": false
        },
        {
          "characterId": 54321,
          "group": "Guild",
          "characterName": "MapleWarrior",
          "channelId": 2,
          "inShop": true,
          "pending": false
        }
      ]
    }
  }
}
```

#### [POST] Create Characters Buddy List

```/api/characters/{characterId}/buddy-list```

Example Request:
```json
{
  "data": {
    "type": "buddy-list",
    "attributes": {
      "capacity": 50
    }
  }
}
```

Response: 202 Accepted (No content)

#### [GET] Get Buddies in Character's Buddy List

```/api/characters/{characterId}/buddy-list/buddies```

Example Response:
```json
{
  "data": [
    {
      "type": "buddies",
      "id": "67890",
      "attributes": {
        "characterId": 67890,
        "group": "Friends",
        "characterName": "MapleHero",
        "channelId": 1,
        "inShop": false,
        "pending": false
      }
    },
    {
      "type": "buddies",
      "id": "54321",
      "attributes": {
        "characterId": 54321,
        "group": "Guild",
        "characterName": "MapleWarrior",
        "channelId": 2,
        "inShop": true,
        "pending": false
      }
    }
  ]
}
```

#### [POST] Add Buddy to Character's Buddy List

```/api/characters/{characterId}/buddy-list/buddies```

Example Request:
```json
{
  "data": {
    "type": "buddies",
    "attributes": {
      "characterId": 67890,
      "group": "Friends",
      "characterName": "MapleHero",
      "channelId": 1,
      "inShop": false,
      "pending": true
    }
  }
}
```

Response: 202 Accepted (No content)

#### [PUT] Update Character's Buddy List Capacity

```/api/characters/{characterId}/buddy-list/capacity```

Updates the maximum capacity of a character's buddy list. The new capacity must be between 1-255 and cannot be less than the current number of buddies in the list.

Example Request:
```json
{
  "data": {
    "type": "buddy-list",
    "attributes": {
      "capacity": 100
    }
  }
}
```

Response: 202 Accepted (No content)

**Error Responses:**
- **400 Bad Request**: Invalid capacity (zero or missing)
- **500 Internal Server Error**: Server error processing the request

**Validation Rules:**
- Capacity must be between 1-255
- New capacity must be greater than or equal to the current number of buddies in the list
- Request is processed asynchronously via Kafka messaging