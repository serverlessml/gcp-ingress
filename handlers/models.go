package handlers

// OutputDistribution defines the output distribution config.
type OutputDistribution struct {
	// Topic is the message broker topic to push payload to.
	Topic string
}

// ErrorOutput defines output error.
type ErrorOutput struct {
	// contains error message
	Message string `json:"message"`
	// contains pipeline config
	PipelineConfig interface{} `json:"pipeline_config"`
}

// OutputPayload defines output payload.
type OutputPayload struct {
	Errors      []ErrorOutput `json:"errors"`
	SubmittedID []string      `json:"submitted_id"`
}
