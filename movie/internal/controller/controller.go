package controller

import (
	"context"
	"errors"
	metadata "movieexample.com/metadata/pkg"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/movie/pkg"
	model "movieexample.com/rating/pkg"
	rating "movieexample.com/rating/pkg"
)

// ErrNotFound is returned when the movie metadata is not
// found.
var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}
type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadata.Metadata, error)
}

// Controller defines a movie service controller.
type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service controller.
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway, metadataGateway}
}

// Get returns the movie details including the aggregated
// rating and movie metadata.
// Get returns the movie details including the aggregated rating and movie metadata.
func (c *Controller) Get(ctx context.Context, id string) (*pkg.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	details := &pkg.MovieDetails{Metadata: *metadata}

	rating, err := c.ratingGateway.GetAggregatedRating(ctx, rating.RecordID(id), rating.RecordTypeMovie)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// Just proceed in this case, it's ok not to have ratings yet.
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}

	return details, nil
}
