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
package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

type clientRotationHashes struct {
	Today     string
	Yesterday string
}

func (s *AnalyticsService) buildClientRotationHashes(siteID int64, ip string, browser models.ClientBrowser, device models.ClientDevice, now time.Time) clientRotationHashes {
	return clientRotationHashes{
		Today:     s.generateVisitorID(siteID, ip, browser, device, now),
		Yesterday: s.generateVisitorID(siteID, ip, browser, device, now.AddDate(0, 0, -1)),
	}
}

func (s *AnalyticsService) resolveClientWithRotation(
	ctx context.Context,
	tx bun.Tx,
	siteID int64,
	ip string,
	device models.ClientDevice,
	browser models.ClientBrowser,
	os models.ClientOS,
	screenSize models.ClientScreenSize,
	country string,
	now time.Time,
) (*models.Client, error) {
	hashes := s.buildClientRotationHashes(siteID, ip, browser, device, now)
	client, err := s.analyticsRepo.FindClientByHashesTx(ctx, tx, siteID, hashes.Today, hashes.Yesterday)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("find client by rotation hashes: %w", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		client = &models.Client{
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

func backfillClientAnalyticsDimensions(
	client *models.Client,
	device models.ClientDevice,
	browser models.ClientBrowser,
	os models.ClientOS,
	screenSize models.ClientScreenSize,
	country string,
) bool {
	if client == nil {
		return false
	}

	changed := false
	if client.Device == models.ClientDeviceUnknown && device != models.ClientDeviceUnknown {
		client.Device = device
		changed = true
	}
	if client.Browser == models.ClientBrowserUnknown && browser != models.ClientBrowserUnknown {
		client.Browser = browser
		changed = true
	}
	if client.OS == models.ClientOSUnknown && os != models.ClientOSUnknown {
		client.OS = os
		changed = true
	}
	if client.ScreenSize == models.ClientScreenSizeUnknown && screenSize != models.ClientScreenSizeUnknown {
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
