package external

import (
	"bufio"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
)

const usersFile = "/etc/passwd"

type User struct {
	ID    int    `json:"id"`
	GID   int    `json:"gid"`
	Name  string `json:"name"`
	Home  string `json:"home"`
	Shell string `json:"shell"`
}

func GetUser(username string) (User, error) {
	f, err := os.Open(usersFile)
	if err != nil {
		return User{}, err
	}

	scan := bufio.NewScanner(f)

	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line[0] == '#' {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}

		id, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}

		gid, err := strconv.Atoi(fields[3])
		if err != nil {
			continue
		}

		if fields[0] == username {
			return User{
				ID:    id,
				GID:   gid,
				Name:  fields[0],
				Home:  fields[5],
				Shell: fields[6],
			}, nil
		}
	}

	return User{}, errors.New("User not found")
}

func ListUsers() ([]User, error) {
	f, err := os.Open(usersFile)
	if err != nil {
		return []User{}, err
	}

	scan := bufio.NewScanner(f)

	users := []User{}
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line[0] == '#' {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}

		id, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}

		gid, err := strconv.Atoi(fields[3])
		if err != nil {
			continue
		}

		users = append(users, User{
			ID:    id,
			GID:   gid,
			Name:  fields[0],
			Home:  fields[5],
			Shell: fields[6],
		})
	}

	sort.Slice(users, func(i, j int) bool {
		forceOrderI := usersForcedOrder(users[i].ID)
		forceOrderJ := usersForcedOrder(users[j].ID)

		if forceOrderI == forceOrderJ {
			return users[i].Name < users[j].Name
		}

		return forceOrderI < forceOrderJ
	})

	return users, nil
}

func usersForcedOrder(id int) int {
	switch {
	case id == 65534:
		// Nobody comes last
		return 3
	case id == 0:
		// Root comes after regular users
		return 1
	case id < 1000:
		// System users come after root
		return 2
	default:
		// Regular users come first
		return 0
	}
}
