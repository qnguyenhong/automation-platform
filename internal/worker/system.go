package worker

import "os"

func getHostnameFromOS() (string, error) {
	return os.Hostname()
}
