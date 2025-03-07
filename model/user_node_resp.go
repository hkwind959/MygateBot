package model

type UserNodeResp struct {
	Data struct {
		Item []struct {
			Id             string      `json:"_id"`
			DeletedAt      interface{} `json:"deletedAt"`
			CreatedBy      string      `json:"createdBy"`
			ID             string      `json:"id"`
			IP             string      `json:"ip"`
			Name           string      `json:"name"`
			Status         string      `json:"status"`
			IsHidden       bool        `json:"isHidden"`
			ActivationDate string      `json:"activationDate"`
			Date           string      `json:"date"`
			Uptime         string      `json:"uptime"`
			CreatedAt      string      `json:"createdAt"`
			UpdatedAt      string      `json:"updatedAt"`
			V              int64       `json:"__v"`
			TodayEarn      string      `json:"todayEarn"`
			SeasonEarn     string      `json:"seasonEarn"`
		} `json:"items"`
		Meta struct {
			CurrentPage int `json:"currentPage"`
			From        int `json:"from"`
			PerPage     int `json:"perPage"`
			LastPage    int `json:"lastPage"`
			To          int `json:"to"`
			Total       int `json:"total"`
		} `json:"meta"`
	} `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type RegisterNodeResp struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
