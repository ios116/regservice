package server

import (
	"context"
	"fmt"
	"github.com/ios116/regservice/session"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"sync"
)

const sessionKeyLen = 10

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[*session.SessionID]*session.Session
	logger   *zap.Logger
}

func NewSessionManager(logger *zap.Logger) *SessionManager {
	return &SessionManager{
		sync.RWMutex{},
		map[*session.SessionID]*session.Session{},
		logger,
	}
}

func (sm *SessionManager) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
	id := &session.SessionID{ID:RandStringRunes(sessionKeyLen)}
	sm.logger.Info("Create=== ", zap.String("id=", id.ID))
	sm.mu.Lock()
	sm.sessions[id] = in
	sm.mu.Unlock()
	return id, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *session.SessionID) (*session.Session, error) {
	sm.logger.Info("Check", zap.String("id=", in.ID))
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	fmt.Println(sm.sessions)
	if sess, ok := sm.sessions[in]; ok {
		return sess, nil
	}
	sm.logger.Info("id is false")
	return nil, status.Errorf(codes.NotFound, "Session not found")
}

func (sm *SessionManager) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	sm.logger.Info("Delete", zap.String("id=", in.ID))
	sm.mu.Lock()
	delete(sm.sessions, in)
	sm.mu.Unlock()
	return &session.Nothing{Dummy: true}, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
