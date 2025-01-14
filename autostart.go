package autostart

// An application that will be started when the user logs in.
type App struct {
	// Unique identifier for the app.
	Name string
	// The command to execute, followed by its arguments.
	Exec []string
        //working directory
        WorkDir string
	// The app name.
	DisplayName string
	// The app icon.
	Icon string
        // for all users
        AllUsers bool
        // startup file folder
        startupDir string
}
