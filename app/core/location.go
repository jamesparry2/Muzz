package core

import "context"

type LocationRequest struct {
	Lat            int
	Long           int
	DistanceFromMe int

	UserID uint
}

func (c *Client) Location(ctx context.Context, request *LocationRequest) error {
	return nil
}
