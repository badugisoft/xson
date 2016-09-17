package data

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

type SocialNetwork struct {
	Name string
	URL  string
}

type Variables struct {
	StrValue  string
	IntValue  int
	BoolValue bool
}

type Data struct {
	Persons        []Person
	SocialNetworks []SocialNetwork
	Variables      Variables
}

var SampleData = Data{
	Persons: []Person{
		Person{
			FirstName: "Liam",
			LastName:  "Neeson",
			Age:       64,
		},
		Person{
			FirstName: "Tom",
			LastName:  "Cruise",
			Age:       53,
		},
		Person{
			FirstName: "Conan",
			LastName:  "O'Brien",
			Age:       52,
		},
	},
	SocialNetworks: []SocialNetwork{
		SocialNetwork{
			Name: "Twitter",
			URL:  "https://www.twitter.com",
		},
		SocialNetwork{
			Name: "Facebook",
			URL:  "https://www.facebook.com/",
		},
	},
	Variables: Variables{
		StrValue:  "Some text",
		IntValue:  1997,
		BoolValue: false,
	},
}
