package stream

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/markmorcos/naharda/api/internal/store"
)

// Channel is the Postgres NOTIFY channel ingest publishes on.
const Channel = "naharda_updates"

// Listen holds a dedicated Postgres LISTEN connection and broadcasts a fresh
// snapshot to all stream clients whenever ingest NOTIFYs a change. It reconnects
// on error until ctx is cancelled. Decoupled from ingest via Postgres (§5).
func Listen(ctx context.Context, st *store.Store, hub *Hub, log *slog.Logger) {
	if st == nil || st.Pool == nil {
		log.Info("stream listener disabled (no database)")
		return
	}
	for ctx.Err() == nil {
		if err := listenOnce(ctx, st, hub, log); err != nil && ctx.Err() == nil {
			log.Warn("stream listener error; retrying", "err", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}
}

func listenOnce(ctx context.Context, st *store.Store, hub *Hub, log *slog.Logger) error {
	conn, err := st.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "LISTEN "+Channel); err != nil {
		return err
	}
	log.Info("stream listener connected", "channel", Channel)

	for {
		if _, err := conn.Conn().WaitForNotification(ctx); err != nil {
			return err
		}
		snap, err := BuildSnapshot(ctx, st)
		if err != nil {
			log.Warn("snapshot build failed", "err", err)
			continue
		}
		hub.Broadcast([]byte(fmt.Sprintf("event: update\ndata: %s\n\n", snap)))
	}
}
