package recordingprotocol

import "errors"

var ErrFailedToStartCamera = errors.New("failed to start camera streaming")

var ErrRecordingFolderForCameraNotFound = errors.New("Folder where in current camera is set to save its recordings isn't available.")
