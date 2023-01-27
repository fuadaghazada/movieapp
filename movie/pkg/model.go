package pkg

import model "movieexample.com/metadata/pkg"

// MovieDetails includes movie metadata its aggregated
// rating.
type MovieDetails struct {
	Rating   *float64       `json:"rating,omitEmpty"`
	Metadata model.Metadata `json:"metadata`
}
