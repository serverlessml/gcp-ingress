package handlers

// OutputDistribution defines the output distribution config.
type OutputDistribution struct {
	// Topic is the message broker topic to push payload to.
	Topic string
}

// OutputPayload defines output payload.
type OutputPayload struct {
	Errors      []string `json:"errors"`
	SubmittedID []string `json:"submitted_id"`
}
