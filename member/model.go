package member

type Member struct {
    ID          string `json:"id" bson:"_id"`
    Name        string `json:"name" bson:"name"`
    Face        string `json:"face" bson:"face"`
    Attendance  int    `json:"attendance" bson:"attendance"`
    AddedAt     string `json:"addedAt" bson:"addedAt"`
}

type CreateGroupRequest struct {
    Name string `json:"name"`
}

type CreateGroupResponse struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt string    `json:"createdAt"`
}