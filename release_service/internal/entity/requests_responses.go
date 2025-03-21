package entity

type DeleteResponse struct {
	IsDeleted bool `json:"isDeleted"`
}

type UpdateRequest struct {
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
}

type AddRequest struct {
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	ImageName   string `json:"image_name"`
	ImageData   []byte `json:"image_data"`
}

type AddResponse struct {
	ReleaseId int `json:"release_id"`
}
