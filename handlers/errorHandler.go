package handler

type HasError interface {
	GetHandler() Handler
}

type Handler struct {
	Status       int    `json:"status"`
	DetailStatus int    `json:"detail_status"`
	MessageID    string `json:"message_id"`
	MessageEN    string `json:"message_en"`
	Error        string `json:"error"`
}

type ErrorHandlers struct {
	Error *Handler
}

func (h Handler) GetHandler() Handler {
	return Handler{}
}

func ErrorHandler(status int, detailStatus int, errDesc string) Handler {

	var handler Handler
	switch status {
	case 400:
		if detailStatus == 401 {
			handler.Status = status
			handler.DetailStatus = detailStatus
			handler.MessageID = "Unauthorized"
			handler.MessageEN = "Unauthorized"
			handler.Error = errDesc
			// e.Error.DetailStatus = detailStatus,
		} else if detailStatus == 402 {
			handler.Status = status
			handler.DetailStatus = detailStatus
			handler.MessageID = "User tidak ditemukan"
			handler.MessageEN = "User not found"
			handler.Error = errDesc

		} else {
			handler.Status = status
			handler.DetailStatus = detailStatus
			handler.MessageID = "Bad Request"
			handler.MessageEN = "Bad Request"
			handler.Error = errDesc
		}

	case 500:
		handler.Status = status
		handler.DetailStatus = detailStatus
		handler.MessageID = "Internal server error"
		handler.MessageEN = "Internal server error"
		handler.Error = errDesc
	}

	return handler
}
