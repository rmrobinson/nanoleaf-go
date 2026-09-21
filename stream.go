package nanoleaf

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/r3labs/sse"
	backoff "gopkg.in/cenkalti/backoff.v1"
)

const (
	streamReconnectMinDelay = 1 * time.Second
	streamReconnectMaxDelay = 30 * time.Second
)

// SubscribeEvents opens the panel's /events stream for the given event type IDs (1 = state,
// 2 = layout, 3 = effects, 4 = touch gestures) and returns a channel of decoded updates plus a
// channel of non-fatal connection errors (buffered by 1; a full channel drops the error rather
// than blocking the read loop - it's advisory only, the retry loop below is what actually keeps
// the stream alive).
//
// Both channels are closed once ctx is canceled - that's the only way this ever stops on its
// own. A dropped connection (including the panel simply being unreachable) is retried with
// exponential backoff, capped at streamReconnectMaxDelay, indefinitely; there's no
// give-up-after-N-attempts behaviour, since a panel that's temporarily off the network should
// still be reconnected to whenever it comes back.
func (c *Client) SubscribeEvents(ctx context.Context, ids ...int) (<-chan *PanelUpdate, <-chan error) {
	updates := make(chan *PanelUpdate)
	errs := make(chan error, 1)

	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	url := c.getURLBase() + "events?id=" + strings.Join(idStrs, ",")

	go func() {
		defer close(updates)
		defer close(errs)

		delay := streamReconnectMinDelay
		for {
			sc := sse.NewClient(url)
			sc.Connection = c.httpClient
			// A single attempt per loop iteration - this function owns the retry/backoff
			// decision (and its interaction with ctx cancellation) itself, rather than
			// leaning on sse.Client's own internal retry, which doesn't stop promptly once
			// ctx is canceled.
			sc.ReconnectStrategy = &backoff.StopBackOff{}

			attemptErr := sc.SubscribeWithContext(ctx, "", func(msg *sse.Event) {
				id, err := strconv.Atoi(string(msg.ID))
				if err != nil {
					// Not a typed panel event (e.g. a keep-alive comment) - ignore.
					return
				}

				update := &PanelUpdate{TypeID: id}
				if err := json.Unmarshal(msg.Data, update); err != nil {
					return
				}

				select {
				case updates <- update:
				case <-ctx.Done():
				}
			})

			if ctx.Err() != nil {
				return
			}

			if attemptErr != nil {
				select {
				case errs <- attemptErr:
				default:
				}
			} else {
				// A clean EOF (e.g. the panel closing the connection) counts as success -
				// reset the backoff so the next reconnect attempt starts fast again.
				delay = streamReconnectMinDelay
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}

			delay *= 2
			if delay > streamReconnectMaxDelay {
				delay = streamReconnectMaxDelay
			}
		}
	}()

	return updates, errs
}
