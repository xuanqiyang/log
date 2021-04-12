package mlog
var Log *Logger
func init() {
	var err error
	Log, err = NewLogger(Config{
		LogPath: "./",
		Level:   FATAL | ERROR | DEBUG |INFO,
		fileMaxSize: 1024*50,
		TypeMapFile: map[LogType]string{
			FATAL: "errors",
			ERROR: "errors",
			DEBUG: "info",
			INFO: "info",
		},
	})
	if err != nil {
		panic(err)
	}
}
