package member

type Member struct {
	ID         string `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	Face       string `json:"face" bson:"face"`
	Attendance int    `json:"attendance" bson:"attendance"`
	AddedAt    string `json:"addedAt" bson:"addedAt"`
}
