package collectors

import (
	"context"
	"os/user"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

// HostUsersLastLoginCollector derives last-login data from host user sessions.
type HostUsersLastLoginCollector struct{}

// NewLastLoginCollector constructs the default last-login collector.
func NewLastLoginCollector() LastLoginCollector {
	return HostUsersLastLoginCollector{}
}

// CollectLastLogin implements LastLoginCollector.
func (HostUsersLastLoginCollector) CollectLastLogin(ctx context.Context) (*LastLoginInfo, error) {
	current, err := user.Current()
	if err != nil {
		recordError("last_login", err)
	}
	users, err := host.UsersWithContext(ctx)
	if err != nil {
		recordError("last_login", err)
		return nil, nil
	}
	var latest *host.UserStat
	for i := range users {
		u := users[i]
		if current != nil && u.User != current.Username {
			continue
		}
		if latest == nil || u.Started > latest.Started {
			latest = &u
		}
	}
	if latest == nil {
		return nil, nil
	}
	return &LastLoginInfo{
		Timestamp: time.Unix(int64(latest.Started), 0),
		Source:    latest.Host,
	}, nil
}
