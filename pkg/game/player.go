package game

import (
	"time"
)

// Player represents a player in the game
type Player struct {
	ID        string
	PublicKey string
	Number    int
	Timestamp time.Time
	Signature []byte
}

// NewPlayer creates a new player with the given ID
func NewPlayer(id string) *Player {
	return &Player{
		ID: id,
	}
}

// SetPublicKey sets the player's public key
func (p *Player) SetPublicKey(publicKey string) {
	p.PublicKey = publicKey
}

// SetNumber sets the player's number
func (p *Player) SetNumber(number int) {
	p.Number = number
}

// SetTimestamp sets the player's timestamp
func (p *Player) SetTimestamp(timestamp time.Time) {
	p.Timestamp = timestamp
}

// SetSignature sets the player's signature
func (p *Player) SetSignature(signature []byte) {
	p.Signature = signature
}
