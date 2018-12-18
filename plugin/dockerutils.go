package plugin

import (
	"bufio"
	"errors"
	"os"
	"regexp"
)

// DockerUtils is an interface with docker utilities
type DockerUtils interface {
	GetCurrentContainerID() (string, error)
}

type dockerUtils struct{}

// CreateDockerUtils creates a new instance of docker utils
func CreateDockerUtils() DockerUtils {
	return &dockerUtils{}
}

func (wrapper *dockerUtils) GetCurrentContainerID() (string, error) {
	file, err := os.Open("/proc/self/cgroup")

	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		_, lines, err := bufio.ScanLines([]byte(scanner.Text()), true)
		if err != nil {
			return nil, err
		}
		strLines := string(lines)
		if id := matchDockerCurrentContainerID(strLines); id != "" {
			return id, nil
		} else if id := matchECSCurrentContainerID(strLines); id != "" {
			return id, nil
		}
	}
	return nil, errors.New("Cannot find container id")
}

func matchDockerCurrentContainerID(lines string) string {
	regex := "/docker[/-]([[:alnum:]]{64})(\\.scope)?$"
	re := regexp.MustCompilePOSIX(regex)

	if re.MatchString(lines) {
		submatches := re.FindStringSubmatch(string(lines))
		containerID := submatches[1]

		return containerID
	}
	return ""
}

func matchECSCurrentContainerID(lines string) string {
	regex := "/ecs\\/[^\\/]+\\/(.+)$"
	re := regexp.MustCompilePOSIX(regex)

	if re.MatchString(string(lines)) {
		submatches := re.FindStringSubmatch(string(lines))
		containerID := submatches[1]

		return containerID
	}

	return ""
}
