package data

// ShopWeights defines shop pool weights by level
var ShopWeights = map[int][]float64{
	1: {1.0, 0, 0, 0},     // level 1: only 1-cost
	2: {0.75, 0.25, 0, 0}, // level 2
	3: {0.55, 0.30, 0.15, 0},
	4: {0.40, 0.30, 0.20, 0.10},
	5: {0.30, 0.30, 0.25, 0.15},
	6: {0.20, 0.25, 0.30, 0.25},
}

// GetUnitsForCost returns all unit definitions with the given cost
func GetUnitsForCost(cost int) []*UnitDef {
	var result []*UnitDef
	for _, u := range UnitDefs {
		if u.Cost == cost {
			result = append(result, u)
		}
	}
	return result
}
