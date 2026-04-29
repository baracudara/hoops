package models

import "time"

type Player struct {
    ID       string
    Name     string
    Nickname string
    Position *string
    Age      *int32
    Stats    PlayerStats
    CreatedAt time.Time
    UpdatedAt time.Time
}

type PlayerStats struct {
    GamesPlayed int32
    Wins        int32
    Losses      int32
}