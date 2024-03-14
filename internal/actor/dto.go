package actor

type ActorInsert struct {
	Name      string `json:"name" valid:",required"`
	Gender    string `json:"gender" valid:",required"`
	BirthDate string `josn:"birth_date" valid:"rfc3339WithoutZone,required"`
}
