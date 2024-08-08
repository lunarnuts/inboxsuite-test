package dto

type ProfileID int64

type ClassID uint8

type RoadmapID uint8

type Job struct {
	ProfileID ProfileID `json:"profile_id"`
	ClassID   ClassID   `json:"class_id"`
}

type Result struct {
	ProfileID ProfileID `json:"profile_id"`
	RoadmapID RoadmapID `json:"roadmap_id"`
}

type Statistics struct {
	Count int64 `json:"count"`
}
