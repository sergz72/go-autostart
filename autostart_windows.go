package autostart

import (
	"os"
	"path/filepath"
        ole "github.com/go-ole/go-ole"
        "github.com/go-ole/go-ole/oleutil"
        "strings"
)

func (a *App) Init() {
    if a.AllUsers {
        a.startupDir = filepath.Join(os.Getenv("ProgramData"), "Microsoft", "Windows", "Start Menu", "Programs", "StartUp")
    } else {
        a.startupDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
    }
}

func (a *App) path() string {
	return filepath.Join(a.startupDir, a.Name+".lnk")
}

func (a *App) IsEnabled() bool {
	_, err := os.Stat(a.path())
	return err == nil
}

func (a *App) Enable() error {
	path := a.Exec[0]
	args := strings.Join(a.Exec[1:], " ")

	if err := os.MkdirAll(a.startupDir, 0777); err != nil {
		return err
	}

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", a.path())
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", path)
	oleutil.PutProperty(idispatch, "Arguments", args)
	oleutil.PutProperty(idispatch, "WorkingDirectory", a.WorkDir)
	_, err = oleutil.CallMethod(idispatch, "Save")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Disable() error {
	return os.Remove(a.path())
}
