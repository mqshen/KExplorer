package services

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"kafkaexplorer/backend/consts"
	"kafkaexplorer/backend/storage"
	"kafkaexplorer/backend/types"
	"strings"
	"sync"
)

type preferencesService struct {
	pref          *storage.PreferencesStorage
	clientVersion string
}

var preferences *preferencesService
var oncePreferences sync.Once

func Preferences() *preferencesService {
	if preferences == nil {
		oncePreferences.Do(func() {
			preferences = &preferencesService{
				pref:          storage.NewPreferences(),
				clientVersion: "",
			}
		})
	}
	return preferences
}

func (p *preferencesService) SetAppVersion(ver string) {
	if !strings.HasPrefix(ver, "v") {
		p.clientVersion = "v" + ver
	} else {
		p.clientVersion = ver
	}
}

func (p *preferencesService) GetWindowSize() (width, height int, maximised bool) {
	data := p.pref.GetPreferences()
	width, height, maximised = data.Behavior.WindowWidth, data.Behavior.WindowHeight, data.Behavior.WindowMaximised
	if width <= 0 {
		width = consts.DEFAULT_WINDOW_WIDTH
	}
	if height <= 0 {
		height = consts.DEFAULT_WINDOW_HEIGHT
	}
	return
}

func (p *preferencesService) GetWindowPosition(ctx context.Context) (x, y int) {
	data := p.pref.GetPreferences()
	x, y = data.Behavior.WindowPosX, data.Behavior.WindowPosY
	width, height := data.Behavior.WindowWidth, data.Behavior.WindowHeight
	var screenWidth, screenHeight int
	if screens, err := runtime.ScreenGetAll(ctx); err == nil {
		for _, screen := range screens {
			if screen.IsCurrent {
				screenWidth, screenHeight = screen.Size.Width, screen.Size.Height
				break
			}
		}
	}
	if screenWidth <= 0 || screenHeight <= 0 {
		screenWidth, screenHeight = consts.DEFAULT_WINDOW_WIDTH, consts.DEFAULT_WINDOW_HEIGHT
	}
	if x <= 0 || x+width > screenWidth || y <= 0 || y+height > screenHeight {
		// out of screen, reset to center
		x, y = (screenWidth-width)/2, (screenHeight-height)/2
	}
	return
}

func (p *preferencesService) SaveWindowPosition(x, y int) {
	if x > 0 || y > 0 {
		p.UpdatePreferences(map[string]any{
			"behavior.windowPosX": x,
			"behavior.windowPosY": y,
		})
	}
}

func (p *preferencesService) UpdatePreferences(value map[string]any) (resp types.JSResp) {
	err := p.pref.UpdatePreferences(value)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

func (p *preferencesService) SaveWindowSize(width, height int, maximised bool) {
	if maximised {
		// do not update window size if maximised state
		p.UpdatePreferences(map[string]any{
			"behavior.windowMaximised": true,
		})
	} else if width >= consts.MIN_WINDOW_WIDTH && height >= consts.MIN_WINDOW_HEIGHT {
		p.UpdatePreferences(map[string]any{
			"behavior.windowWidth":     width,
			"behavior.windowHeight":    height,
			"behavior.windowMaximised": false,
		})
	}
}
