package storage

import (
	"gopkg.in/yaml.v3"
	"kafkaexplorer/backend/consts"
	"kafkaexplorer/backend/types"
	"sync"
)

type PreferencesStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewPreferences() *PreferencesStorage {
	return &PreferencesStorage{
		storage: NewLocalStore("preferences.yaml"),
	}
}

func (p *PreferencesStorage) getPreferences() (ret types.Preferences) {
	ret = p.DefaultPreferences()
	b, err := p.storage.Load()
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(b, &ret); err != nil {
		ret = p.DefaultPreferences()
		return
	}
	return
}

func (p *PreferencesStorage) DefaultPreferences() types.Preferences {
	return types.NewPreferences()
}

// GetPreferences Get preferences from local
func (p *PreferencesStorage) GetPreferences() (ret types.Preferences) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ret = p.getPreferences()
	if ret.General.ScanSize <= 0 {
		ret.General.ScanSize = consts.DEFAULT_SCAN_SIZE
	}
	ret.Behavior.AsideWidth = max(ret.Behavior.AsideWidth, consts.DEFAULT_ASIDE_WIDTH)
	ret.Behavior.WindowWidth = max(ret.Behavior.WindowWidth, consts.MIN_WINDOW_WIDTH)
	ret.Behavior.WindowHeight = max(ret.Behavior.WindowHeight, consts.MIN_WINDOW_HEIGHT)
	return
}
