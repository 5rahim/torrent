package torrent

import (
	"testing"
	"time"

	"github.com/go-quicktest/qt"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/tracker"
)

// A re-added Torrent used to inherit the dropped Torrent's announce history. It skipped Started and
// waited out the remaining interval, leaving DHT as its only peer source.
func TestReAddedTorrentAnnouncesStarted(t *testing.T) {
	for _, last := range []tracker.AnnounceEvent{tracker.Started, tracker.Stopped} {
		t.Run(last.String(), func(t *testing.T) {
			cfg := TestingConfig(t)
			cfg.DisableTrackers = false
			cl, err := NewClient(cfg)
			qt.Assert(t, qt.IsNil(err))
			t.Cleanup(func() { cl.Close() })

			ih := metainfo.HashBytes([]byte("re-added"))
			add := func() *Torrent {
				torr, _ := cl.AddTorrentOpt(AddTorrentOpts{InfoHash: ih})
				cl.lock()
				defer cl.unlock()
				torr.addTrackers([][]string{{"http://tracker.example.com:6969/announce"}})
				return torr
			}
			first := add()
			dispatcher := &cl.regularTrackerAnnounceDispatcher

			cl.lock()
			var key torrentTrackerAnnouncerKey
			for k := range first.regularTrackerAnnounceState {
				key = k
			}
			*dispatcher.announceStates[key] = announceState{
				lastOk: lastAnnounceOk{AnnouncedEvent: last, Interval: time.Hour, Completed: time.Now()},
			}
			cl.unlock()

			first.Drop()
			second := add()

			cl.lock()
			defer cl.unlock()
			event, _ := dispatcher.nextAnnounceEvent(key)
			qt.Assert(t, qt.Equals(event, tracker.Started))
			row, ok := dispatcher.announceData.Get(key)
			qt.Assert(t, qt.IsTrue(ok))
			qt.Assert(t, qt.Equals(row.AnnounceEvent, tracker.Started))
			qt.Assert(t, qt.IsFalse(row.When.After(time.Now())))
			qt.Assert(t, qt.Equals(dispatcher.getTorrentForAnnounceRequest(key.ShortInfohash), second))
		})
	}
}
