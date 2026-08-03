package applock

import (
	"sync"
	"time"
)

const (
	maxUnlockFailures = 5                // 잠금까지 허용하는 연속 실패 횟수
	unlockLockWindow  = 30 * time.Second // 실패 초과 시 잠금 유지 시간
)

// unlockAttempt는 사용자별 PIN 실패 상태를 보관한다.
type unlockAttempt struct {
	failures    int
	lockedUntil time.Time
}

var (
	unlockAttempts   = make(map[string]*unlockAttempt)
	unlockAttemptsMu sync.Mutex
)

// unlockLockRemaining는 현재 잠금 중이면 남은 초를 반환한다.
func unlockLockRemaining(uid string) int {
	unlockAttemptsMu.Lock()
	defer unlockAttemptsMu.Unlock()

	attempt, ok := unlockAttempts[uid]
	if !ok {
		return 0
	}
	remaining := time.Until(attempt.lockedUntil)
	if remaining <= 0 {
		return 0
	}
	return int(remaining.Seconds()) + 1
}

// recordUnlockFailure는 실패를 누적하고, 임계치 초과 시 잠금을 건다.
func recordUnlockFailure(uid string) {
	unlockAttemptsMu.Lock()
	defer unlockAttemptsMu.Unlock()

	attempt, ok := unlockAttempts[uid]
	if !ok {
		attempt = &unlockAttempt{}
		unlockAttempts[uid] = attempt
	}
	attempt.failures++
	if attempt.failures >= maxUnlockFailures {
		attempt.lockedUntil = time.Now().Add(unlockLockWindow)
		attempt.failures = 0
	}
}

// resetUnlockFailure는 성공 시 실패 상태를 초기화한다.
func resetUnlockFailure(uid string) {
	unlockAttemptsMu.Lock()
	defer unlockAttemptsMu.Unlock()

	delete(unlockAttempts, uid)
}
