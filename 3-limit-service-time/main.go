//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

const maxTime time.Duration = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mu        sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium { // checks premium
		process()
		return true
	}
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.TimeUsed >= int64(maxTime) {
		return false
	}

	done := make(chan time.Time, 1)
	start := time.Now()

	go func() {
		process()
		done <- time.Now()
	}()

	select {
	case end := <-done: // a process completed
		duration := end.Sub(start)
		u.TimeUsed += int64(duration.Seconds())

		return true
	case <-time.After(time.Duration(10-u.TimeUsed) * time.Second): // a process aborted
		u.TimeUsed = 10
		return false
	}
}

func main() {
	RunMockServer()
}
