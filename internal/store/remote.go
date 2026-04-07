package store

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/mattermost/mattermost-marketplace/internal/model"
)

// RemoteStore fetches plugins from a remote URL and periodically refreshes.
type RemoteStore struct {
	url      string
	interval time.Duration
	logger   logrus.FieldLogger

	mu    sync.RWMutex
	store *StaticStore

	stopCh chan struct{}
}

// NewRemote creates a store that fetches plugins.json from a URL and refreshes on the given interval.
func NewRemote(url string, refreshInterval time.Duration, logger logrus.FieldLogger) (*RemoteStore, error) {
	r := &RemoteStore{
		url:      url,
		interval: refreshInterval,
		logger:   logger,
		stopCh:   make(chan struct{}),
	}

	if err := r.refresh(); err != nil {
		return nil, errors.Wrap(err, "failed initial fetch from remote URL")
	}

	go r.refreshLoop()

	return r, nil
}

func (r *RemoteStore) refresh() error {
	r.logger.WithField("url", r.url).Debug("Fetching plugins from remote URL")

	resp, err := http.Get(r.url) // #nosec G107 -- URL is operator-configured, not user input
	if err != nil {
		return errors.Wrap(err, "failed to fetch remote plugins")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote URL returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	plugins, err := model.PluginsFromReader(bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "failed to parse remote plugins")
	}

	store, err := NewStatic(plugins, r.logger)
	if err != nil {
		return errors.Wrap(err, "failed to create store from remote plugins")
	}

	r.mu.Lock()
	r.store = store
	r.mu.Unlock()

	r.logger.WithField("count", len(plugins)).Info("Loaded plugins from remote URL")
	return nil
}

func (r *RemoteStore) refreshLoop() {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.refresh(); err != nil {
				r.logger.WithError(err).Error("Failed to refresh plugins from remote URL")
			}
		case <-r.stopCh:
			return
		}
	}
}

// Stop stops the background refresh loop.
func (r *RemoteStore) Stop() {
	close(r.stopCh)
}

// GetPlugins returns plugins from the latest fetched data.
func (r *RemoteStore) GetPlugins(filter *model.PluginFilter) ([]*model.Plugin, error) {
	r.mu.RLock()
	store := r.store
	r.mu.RUnlock()

	return store.GetPlugins(filter)
}
