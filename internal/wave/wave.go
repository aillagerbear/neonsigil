package wave

import (
	"neonsigil/internal/board"
	"neonsigil/internal/data"
	"neonsigil/internal/entity"
)

// WaveManager handles wave spawning
type WaveManager struct {
	Stage          *data.StageDef
	Board          *board.Board
	CurrentWave    int
	GroupTimers    []float64 // timer for each group in current wave
	GroupCounts    []int     // how many spawned per group
	WaveActive     bool
	AllDone        bool
	SpawnedEnemies []*entity.Enemy
}

// NewWaveManager creates a new wave manager
func NewWaveManager(stage *data.StageDef, b *board.Board) *WaveManager {
	return &WaveManager{
		Stage:       stage,
		Board:       b,
		CurrentWave: 0,
		WaveActive:  false,
	}
}

// TotalWaves returns the total number of waves
func (wm *WaveManager) TotalWaves() int {
	return len(wm.Stage.Waves)
}

// StartWave begins the next wave
func (wm *WaveManager) StartWave() {
	if wm.CurrentWave >= len(wm.Stage.Waves) {
		wm.AllDone = true
		return
	}

	wave := wm.Stage.Waves[wm.CurrentWave]
	wm.GroupTimers = make([]float64, len(wave.Groups))
	wm.GroupCounts = make([]int, len(wave.Groups))
	wm.WaveActive = true
	wm.SpawnedEnemies = nil
}

// Update spawns enemies and checks wave completion
func (wm *WaveManager) Update() []*entity.Enemy {
	if !wm.WaveActive || wm.CurrentWave >= len(wm.Stage.Waves) {
		return nil
	}

	wave := wm.Stage.Waves[wm.CurrentWave]
	var newEnemies []*entity.Enemy

	allSpawned := true
	for i, group := range wave.Groups {
		if wm.GroupCounts[i] >= group.Count {
			continue
		}
		allSpawned = false

		wm.GroupTimers[i] -= 1.0 / 60.0
		if wm.GroupTimers[i] <= 0 {
			// Spawn one enemy
			def := data.EnemyDefs[group.Enemy]
			if def != nil {
				e := entity.NewEnemy(def, group.PathID, wm.Stage.EnemyHPMul, wm.Stage.EnemySpdMul, wm.Board)
				newEnemies = append(newEnemies, e)
				wm.SpawnedEnemies = append(wm.SpawnedEnemies, e)
			}
			wm.GroupCounts[i]++
			wm.GroupTimers[i] = group.Interval
		}
	}

	// Check if wave is complete (all spawned and all dead/reached)
	if allSpawned {
		allDead := true
		for _, e := range wm.SpawnedEnemies {
			if e.Alive && !e.Reached {
				allDead = false
				break
			}
		}
		if allDead {
			wm.WaveActive = false
			wm.CurrentWave++
			if wm.CurrentWave >= len(wm.Stage.Waves) {
				wm.AllDone = true
			}
		}
	}

	return newEnemies
}

// IsWaveActive returns whether a wave is currently active
func (wm *WaveManager) IsWaveActive() bool {
	return wm.WaveActive
}

// WaveCompleteCheck checks if the current wave is complete
func (wm *WaveManager) WaveCompleteCheck(enemies []*entity.Enemy) bool {
	if !wm.WaveActive {
		return false
	}

	wave := wm.Stage.Waves[wm.CurrentWave]
	allSpawned := true
	for i, group := range wave.Groups {
		if wm.GroupCounts[i] < group.Count {
			allSpawned = false
			break
		}
	}

	if !allSpawned {
		return false
	}

	for _, e := range wm.SpawnedEnemies {
		if e.Alive && !e.Reached {
			return false
		}
	}

	return true
}
