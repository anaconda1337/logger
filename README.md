## GO Logger - developed by [Anaconda1337](https://github.com/anaconda1337)
_The Logger Module is a simple Go package for logging messages with color support to both the console and a log file. It provides the flexibility to configure log levels, log file settings, and log message colors through a YAML configuration file._

<hr>

### Features

Log messages to the console with color-coding based on log levels.
Optionally log messages to a file.
Customize log colors and formatting using a configuration file.
Configurable log levels: INFO, WARNING, ERROR, DANGER, and DEFAULT.

### Getting Started
- Installation

_To use the Logger Module in your Go project, simply import it:_

```go
import "github.com/anaconda1337/logger"
```

- Usage
  
_Initialize Logger_
    
```go
l, err := logger.InitializeLogger("config.yaml")
if err != nil {
  log.Fatalf("Error creating logger: %v", err)
}
defer l.Close()
```

_Log Messages_

```go
l.LogMessage(logger.LogLevelInfo, "This is an informational message")
l.LogMessage(logger.LogLevelWarning, "This is a warning message")
l.LogMessage(logger.LogLevelError, "This is an error message")
l.LogMessage(logger.LogLevelDanger, "This is a danger message")
```

### Configuration
The Logger Module can be configured using a YAML configuration file. Here's a sample config.yaml:

```yaml
log_settings:
  log_level: log
  log_file_bool: true
  log_file_name: logs
  log_file_extension: .log

log_colors:
  info_color: "blue"
  warning_color: "yellow"
  error_color: "red"
  danger_color: "magenta"
  default_color: "cyan"
```
_You can have multiple configuration files and pass the file name as a parameter to the `InitializeLogger` function._

Note: The config file must be in the `/config` directory.

```go
InitializeLogger("{config_file.yaml}")
```

### License
This module is open-source and available under the [MIT License](https://opensource.org/license/mit/).

<hr>

### TODO:

- [x] Add log levels
- [x] Configure log levels
- [x] Optional log file
- [x] Configurable log file (name, extension)
- [x] Configurable log colors
- [ ] Configurable log message format
- [ ] Configurable log date format
- [ ] Configurable log file max size


<hr>

- Please feel free to contribute to this project. I am happy about every contribution. :smiley:
- Give me a star if you like the project. :star:
- Give me a follow if you want to see more projects from me. :heart:
- Provide feedback if you have any suggestions. :speech_balloon:
- Provide ideas if you have any. :bulb:
