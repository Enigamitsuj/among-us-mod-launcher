package installer

import "errors"

var ErrAlreadyExists = errors.New("Town of Us already exists")

// ErrProtectedInstallFolder is returned when the destination beside Among Us is not writable.
var ErrProtectedInstallFolder = errors.New("Among Us is in a protected folder (often Program Files). Move the game to a user library folder, or run the launcher as administrator")
