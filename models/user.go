package models

type MGMUser struct {
	ID        string `bson:"_id,omitempty"`
	Index     int    `bson:"index"`
	Username  string `bson:"username"`
	Name      string `bson:"name"`
	Clearance int    `bson:"clearance"`
	Banned    bool   `bson:"banned"`
	RollNo    string `bson:"rollno"`
	KratosID  string `bson:"kid"`
}
