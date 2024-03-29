package exception

import "fmt"

var (
	ErrFileNotExists                       error = fmt.Errorf("file is not exits")
	ErrConfigNotValid                      error = fmt.Errorf("your config file is not valid json file, see log for more detail")
	ErrConfigSourceNotExist                error = fmt.Errorf("config 'sources' is not exist")
	ErrConfigDestinationNotExist           error = fmt.Errorf("config 'destination' is not exist")
	ErrConfigSourceDriverNotAvailable      error = fmt.Errorf("config sources driver not available, see log for more detail")
	ErrConfigDestinationDriverNotAvailable error = fmt.Errorf("config sources driver not available, see log for more detail")
	ErrConfigSourceError                   error = fmt.Errorf("config source has error, see log for more detail")
	ErrConfigDestinationError              error = fmt.Errorf("config destination has error, see log for more detail")
	ErrMysqldumpNotExist                   error = fmt.Errorf("error mysqldump is not exist on your system")
)
