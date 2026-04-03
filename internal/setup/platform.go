package setup

import "runtime"

func DetectPlatform() string {
	return runtime.GOOS
}
