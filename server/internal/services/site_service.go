package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
)

var (
	ErrSiteNotFound    = errors.New("site not found")
	ErrSiteExists      = errors.New("site with this domain already exists")
	ErrNotAuthorized   = errors.New("not authorized to access this site")
)

type SiteService struct {
	siteRepo *repository.SiteRepository
}

func NewSiteService(siteRepo *repository.SiteRepository) *SiteService {
	return &SiteService{siteRepo: siteRepo}
}

type CreateSiteInput struct {
	Domain string
	Name   string
	UserID int64
}

func (s *SiteService) Create(ctx context.Context, input CreateSiteInput) (*models.Site, error) {
	// Check if domain already exists
	existing, _ := s.siteRepo.GetByDomain(ctx, input.Domain)
	if existing != nil {
		return nil, ErrSiteExists
	}

	publicKey, err := generatePublicKey()
	if err != nil {
		return nil, err
	}

	site := &models.Site{
		UserID:    input.UserID,
		Domain:    input.Domain,
		Name:      input.Name,
		PublicKey: publicKey,
	}

	if err := s.siteRepo.Create(ctx, site); err != nil {
		return nil, err
	}

	return site, nil
}

func (s *SiteService) GetByID(ctx context.Context, id, userID int64) (*models.Site, error) {
	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSiteNotFound
	}

	if site.UserID != userID {
		return nil, ErrNotAuthorized
	}

	return site, nil
}

func (s *SiteService) GetByPublicKey(ctx context.Context, publicKey string) (*models.Site, error) {
	return s.siteRepo.GetByPublicKey(ctx, publicKey)
}

func (s *SiteService) GetUserSites(ctx context.Context, userID int64) ([]*models.Site, error) {
	return s.siteRepo.GetByUserID(ctx, userID)
}

func (s *SiteService) Update(ctx context.Context, id, userID int64, name string) (*models.Site, error) {
	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSiteNotFound
	}

	if site.UserID != userID {
		return nil, ErrNotAuthorized
	}

	site.Name = name
	if err := s.siteRepo.Update(ctx, site); err != nil {
		return nil, err
	}

	return site, nil
}

func (s *SiteService) Delete(ctx context.Context, id, userID int64) error {
	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return ErrSiteNotFound
	}

	if site.UserID != userID {
		return ErrNotAuthorized
	}

	return s.siteRepo.Delete(ctx, id)
}

func (s *SiteService) RegeneratePublicKey(ctx context.Context, id, userID int64) (*models.Site, error) {
	site, err := s.siteRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSiteNotFound
	}

	if site.UserID != userID {
		return nil, ErrNotAuthorized
	}

	publicKey, err := generatePublicKey()
	if err != nil {
		return nil, err
	}

	site.PublicKey = publicKey
	if err := s.siteRepo.Update(ctx, site); err != nil {
		return nil, err
	}

	return site, nil
}

func generatePublicKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
