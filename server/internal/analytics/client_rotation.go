/*
UTC-day-skipped client rotation keeps a pseudonymous client alive across adjacent
UTC days without adding a second lookup field.

For each request, analytics computes a daily hash for today and yesterday from:
site ID, truncated IP prefix, browser, and device.

- If today's hash exists, that client is reused.
- If only yesterday's hash exists, that same row is rewritten to today's hash.
- If neither hash exists, a new client row is created with today's hash.

A client therefore rotates only after at least one full UTC day was skipped.
Sessions remain separate and continue to use the 30-minute inactivity window.
*/
package analytics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	"github.com/uptrace/bun"
)

var errExitClientNotFound = errors.New("exit client not found")

type clientRotationHashes struct {
	Today     string
	Yesterday string
}

func (s *Service) buildClientRotationHashes(siteID int64, ip string, browser analyticspersistence.ClientBrowser, device analyticspersistence.ClientDevice, now time.Time) clientRotationHashes {
	return clientRotationHashes{
		Today:     s.generateVisitorID(siteID, ip, browser, device, now),
		Yesterday: s.generateVisitorID(siteID, ip, browser, device, now.AddDate(0, 0, -1)),
	}
}

func (s *Service) resolveClientWithRotation(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	ip string,
	device analyticspersistence.ClientDevice,
	browser analyticspersistence.ClientBrowser,
	os analyticspersistence.ClientOS,
	screenSize analyticspersistence.ClientScreenSize,
	country string,
	now time.Time,
) (*analyticspersistence.Client, error) {
	hashes := s.buildClientRotationHashes(siteID, ip, browser, device, now)
	client, err := s.analyticsRepo.FindClientByHashesTx(ctx, tx, siteID, hashes.Today, hashes.Yesterday)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("find client by rotation hashes: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		client = &analyticspersistence.Client{
			SiteID:     siteID,
			Hash:       hashes.Today,
			Country:    country,
			Device:     device,
			Browser:    browser,
			OS:         os,
			ScreenSize: screenSize,
		}
		if err := s.analyticsRepo.CreateClientTx(ctx, tx, client); err != nil {
			existing, findErr := s.analyticsRepo.FindClientByHashTx(ctx, tx, siteID, hashes.Today)
			if findErr == nil {
				return existing, nil
			}
			return nil, fmt.Errorf("create rotated client: %w", err)
		}
		return client, nil
	}

	changed := false
	if client.Hash == hashes.Yesterday {
		client.Hash = hashes.Today
		changed = true
	}
	changed = backfillClientAnalyticsDimensions(client, device, browser, os, screenSize, country) || changed
	if changed {
		if err := s.analyticsRepo.UpdateClientTx(ctx, tx, client); err != nil {
			existing, findErr := s.analyticsRepo.FindClientByHashTx(ctx, tx, siteID, hashes.Today)
			if findErr == nil {
				return existing, nil
			}
			return nil, fmt.Errorf("update rotated client: %w", err)
		}
	}

	return client, nil
}

func (s *Service) findClientForExit(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	ip string,
	browser analyticspersistence.ClientBrowser,
	device analyticspersistence.ClientDevice,
	now time.Time,
) (*analyticspersistence.Client, error) {
	hashes := s.buildClientRotationHashes(siteID, ip, browser, device, now)
	client, err := s.analyticsRepo.FindClientByHashesTx(ctx, tx, siteID, hashes.Today, hashes.Yesterday)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errExitClientNotFound
		}
		return nil, fmt.Errorf("find client by rotation hashes: %w", err)
	}

	return client, nil
}

func backfillClientAnalyticsDimensions(
	client *analyticspersistence.Client,
	device analyticspersistence.ClientDevice,
	browser analyticspersistence.ClientBrowser,
	os analyticspersistence.ClientOS,
	screenSize analyticspersistence.ClientScreenSize,
	country string,
) bool {
	if client == nil {
		return false
	}

	changed := false
	if client.Device == analyticspersistence.ClientDeviceUnknown && device != analyticspersistence.ClientDeviceUnknown {
		client.Device = device
		changed = true
	}
	if client.Browser == analyticspersistence.ClientBrowserUnknown && browser != analyticspersistence.ClientBrowserUnknown {
		client.Browser = browser
		changed = true
	}
	if client.OS == analyticspersistence.ClientOSUnknown && os != analyticspersistence.ClientOSUnknown {
		client.OS = os
		changed = true
	}
	if client.ScreenSize == analyticspersistence.ClientScreenSizeUnknown && screenSize != analyticspersistence.ClientScreenSizeUnknown {
		client.ScreenSize = screenSize
		changed = true
	}
	if (strings.TrimSpace(client.Country) == "" || client.Country == UnknownCountry.ISOCode || client.Country == LocalNetworkCountry.ISOCode) &&
		strings.TrimSpace(country) != "" &&
		country != UnknownCountry.ISOCode &&
		country != LocalNetworkCountry.ISOCode {
		client.Country = country
		changed = true
	}
	return changed
}
