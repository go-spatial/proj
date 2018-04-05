package tables

// RawPrimeMeridian holds a constant prime meridian item
type RawPrimeMeridian struct {
	ID         string
	Definition string
}

// RawPrimeMeridians holds all the globally known prime meridians
var RawPrimeMeridians = []RawPrimeMeridian{
	{"greenwich", "0dE"},
	{"lisbon", "9d07'54.862\"W"},
	{"paris", "2d20'14.025\"E"},
	{"bogota", "74d04'51.3\"W"},
	{"madrid", "3d41'16.58\"W"},
	{"rome", "12d27'8.4\"E"},
	{"bern", "7d26'22.5\"E"},
	{"jakarta", "106d48'27.79\"E"},
	{"ferro", "17d40'W"},
	{"brussels", "4d22'4.71\"E"},
	{"stockholm", "18d3'29.8\"E"},
	{"athens", "23d42'58.815\"E"},
	{"oslo", "10d43'22.5\"E"},
	{"copenhagen", "12d34'40.35\"E"},
}
