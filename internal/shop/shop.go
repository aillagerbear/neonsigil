package shop

import (
	"math/rand"

	"neonsigil/internal/config"
	"neonsigil/internal/data"
	"neonsigil/internal/entity"
)

// Shop manages the unit shop and economy
type Shop struct {
	Slots     [config.ShopSlots]*data.UnitDef // nil = empty slot
	Gold      int
	Level     int
	XP        int
	DeployCap int
	Rules     data.ShopRules
	Rng       *rand.Rand
}

// NewShop creates a new shop for the given stage
func NewShop(stage *data.StageDef, rng *rand.Rand) *Shop {
	s := &Shop{
		Gold:      stage.StartingGold,
		Level:     stage.StartingLv,
		DeployCap: stage.DeployCapBase,
		Rules:     stage.ShopRules,
		Rng:       rng,
	}
	s.Refresh()
	return s
}

// Refresh fills all shop slots with random units
func (s *Shop) Refresh() {
	for i := 0; i < config.ShopSlots; i++ {
		s.Slots[i] = s.rollUnit()
	}
}

func (s *Shop) rollUnit() *data.UnitDef {
	// Determine cost based on level weights
	weights, ok := data.ShopWeights[s.Level]
	if !ok {
		weights = data.ShopWeights[1]
	}

	// Filter by allowed costs
	allowedWeights := make([]float64, len(weights))
	for _, c := range s.Rules.AllowedCosts {
		if c-1 < len(weights) {
			allowedWeights[c-1] = weights[c-1]
		}
	}

	// Normalize
	total := 0.0
	for _, w := range allowedWeights {
		total += w
	}
	if total == 0 {
		// Fallback: all 1-cost
		units := data.GetUnitsForCost(1)
		if len(units) > 0 {
			return units[s.Rng.Intn(len(units))]
		}
		return nil
	}

	roll := s.Rng.Float64() * total
	cumulative := 0.0
	selectedCost := 1
	for i, w := range allowedWeights {
		cumulative += w
		if roll <= cumulative {
			selectedCost = i + 1
			break
		}
	}

	units := data.GetUnitsForCost(selectedCost)
	if len(units) == 0 {
		units = data.GetUnitsForCost(1)
	}
	if len(units) == 0 {
		return nil
	}
	return units[s.Rng.Intn(len(units))]
}

// CanBuy checks if the player can afford the unit in the given slot
func (s *Shop) CanBuy(slot int) bool {
	if slot < 0 || slot >= config.ShopSlots {
		return false
	}
	if s.Slots[slot] == nil {
		return false
	}
	return s.Gold >= s.Slots[slot].Cost
}

// Buy purchases the unit in the given slot
func (s *Shop) Buy(slot int) *entity.Unit {
	if !s.CanBuy(slot) {
		return nil
	}
	def := s.Slots[slot]
	s.Gold -= def.Cost
	s.Slots[slot] = nil
	return entity.NewUnit(def)
}

// CanReroll checks if the player can afford a reroll
func (s *Shop) CanReroll() bool {
	return s.Rules.RerollEnabled && s.Gold >= config.RerollCost
}

// Reroll refreshes all shop slots
func (s *Shop) Reroll() bool {
	if !s.CanReroll() {
		return false
	}
	s.Gold -= config.RerollCost
	s.Refresh()
	return true
}

// LevelUpCost returns the cost to level up
func (s *Shop) LevelUpCost() int {
	return 4 + s.Level*2
}

// CanLevelUp checks if the player can level up
func (s *Shop) CanLevelUp() bool {
	return s.Rules.LevelUpEnabled && s.Gold >= s.LevelUpCost() && s.Level < 6
}

// LevelUp increases the shop level
func (s *Shop) LevelUp() bool {
	if !s.CanLevelUp() {
		return false
	}
	s.Gold -= s.LevelUpCost()
	s.Level++
	s.DeployCap++
	return true
}

// SellUnit sells a unit and refunds gold
func (s *Shop) SellUnit(u *entity.Unit) int {
	refund := max(1, u.Def.Cost/2)
	if u.Star >= 2 {
		refund = u.Def.Cost * u.Star
	}
	s.Gold += refund
	return refund
}

// AddGold adds gold (from kills, wave bonus, etc.)
func (s *Shop) AddGold(amount int) {
	s.Gold += amount
}
