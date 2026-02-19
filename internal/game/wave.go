package game

// WaveManager handles wave spawning
type WaveManager struct {
	Stage       *StageDef
	Board       *Board
	CurrentWave int
	GroupTimers []float64 // timer for each group in current wave
	GroupCounts []int     // how many spawned per group
	WaveActive  bool
	AllDone     bool
	SpawnedEnemies []*Enemy
}

func NewWaveManager(stage *StageDef, board *Board) *WaveManager {
	return &WaveManager{
		Stage:       stage,
		Board:       board,
		CurrentWave: 0,
		WaveActive:  false,
	}
}

func (wm *WaveManager) TotalWaves() int {
	return len(wm.Stage.Waves)
}

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

func (wm *WaveManager) Update() []*Enemy {
	if !wm.WaveActive || wm.CurrentWave >= len(wm.Stage.Waves) {
		return nil
	}

	wave := wm.Stage.Waves[wm.CurrentWave]
	var newEnemies []*Enemy

	allSpawned := true
	for i, group := range wave.Groups {
		if wm.GroupCounts[i] >= group.Count {
			continue
		}
		allSpawned = false

		wm.GroupTimers[i] -= 1.0 / 60.0
		if wm.GroupTimers[i] <= 0 {
			// Spawn one enemy
			def := EnemyDefs[group.Enemy]
			if def != nil {
				e := NewEnemy(def, group.PathID, wm.Stage.EnemyHPMul, wm.Stage.EnemySpdMul, wm.Board)
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

func (wm *WaveManager) IsWaveActive() bool {
	return wm.WaveActive
}

func (wm *WaveManager) WaveCompleteCheck(enemies []*Enemy) bool {
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
