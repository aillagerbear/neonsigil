package ui

import (
	"math/rand"
)

const (
	ShopSlots   = 5
	BenchSlots  = 8
	RerollCost  = 2
)

// Shop manages the unit shop and economy
type Shop struct {
	Slots      [ShopSlots]*UnitDef // nil = empty slot
	Gold       int
	Level      int
	XP         int
	DeployCap  int
	Rules      ShopRules
	Rng        *rand.Rand
}

func NewShop(stage *StageDef, rng *rand.Rand) *Shop {
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

func (s *Shop) Refresh() {
	for i := 0; i < ShopSlots; i++ {
		s.Slots[i] = s.rollUnit()
	}
}

func (s *Shop) rollUnit() *UnitDef {
	// Determine cost based on level weights
	weights, ok := ShopWeights[s.Level]
	if !ok {
		weights = ShopWeights[1]
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
		units := GetUnitsForCost(1)
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

	units := GetUnitsForCost(selectedCost)
	if len(units) == 0 {
		units = GetUnitsForCost(1)
	}
	if len(units) == 0 {
		return nil
	}
	return units[s.Rng.Intn(len(units))]
}

func (s *Shop) CanBuy(slot int) bool {
	if slot < 0 || slot >= ShopSlots {
		return false
	}
	if s.Slots[slot] == nil {
		return false
	}
	return s.Gold >= s.Slots[slot].Cost
}

func (s *Shop) Buy(slot int) *Unit {
	if !s.CanBuy(slot) {
		return nil
	}
	def := s.Slots[slot]
	s.Gold -= def.Cost
	s.Slots[slot] = nil
	return NewUnit(def)
}

func (s *Shop) CanReroll() bool {
	return s.Rules.RerollEnabled && s.Gold >= RerollCost
}

func (s *Shop) Reroll() bool {
	if !s.CanReroll() {
		return false
	}
	s.Gold -= RerollCost
	s.Refresh()
	return true
}

func (s *Shop) LevelUpCost() int {
	return 4 + s.Level*2
}

func (s *Shop) CanLevelUp() bool {
	return s.Rules.LevelUpEnabled && s.Gold >= s.LevelUpCost() && s.Level < 6
}

func (s *Shop) LevelUp() bool {
	if !s.CanLevelUp() {
		return false
	}
	s.Gold -= s.LevelUpCost()
	s.Level++
	s.DeployCap++
	return true
}

func (s *Shop) SellUnit(u *Unit) int {
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

