package support

// Datum contains information about one datum ellipse
type Datum struct {
	ID               string
	DefinitionString string
	EllipseID        string
	Comments         string
	Definition       *PairList
}

// DatumTable is the list of all the known datums
type DatumTable []*Datum

func init() {
	for _, d := range *Datums {
		pl, err := NewPairListFromString(d.DefinitionString)
		if err != nil {
			panic(err)
		}
		d.Definition = pl
	}
}

// Lookup returns the given Datum
func (dt *DatumTable) Lookup(ellipseID string) *Datum {
	for _, d := range *dt {
		if d.EllipseID == ellipseID {
			return d
		}
	}
	return nil
}

//---------------------------------------------------------------------

// Datums is the global table of known datums
var Datums = &DatumTable{
	{"WGS84", "towgs84=0,0,0", "WGS84", "", nil},
	{"GGRS87", "towgs84=-199.87,74.79,246.62", "GRS80", "Greek_Geodetic_Reference_System_1987", nil},
	{"NAD83", "towgs84=0,0,0", "GRS80", "North_American_Datum_1983", nil},
	{"NAD27", "nadgrids=@conus,@alaska,@ntv2_0.gsb,@ntv1_can.dat", "clrk66", "North_American_Datum_1927", nil},
	{"potsdam" /*"towgs84=598.1,73.7,418.2,0.202,0.045,-2.455,6.7",*/, "nadgrids=@BETA2007.gsb", "bessel", "Potsdam Rauenberg 1950 DHDN", nil},
	{"carthage", "towgs84=-263.0,6.0,431.0", "clrk80ign", "Carthage 1934 Tunisia", nil},
	{"hermannskogel", "towgs84=577.326,90.129,463.919,5.137,1.474,5.297,2.4232", "bessel", "Hermannskogel", nil},
	{"ire65", "towgs84=482.530,-130.596,564.557,-1.042,-0.214,-0.631,8.15", "mod_airy", "Ireland 1965", nil},
	{"nzgd49", "towgs84=59.47,-5.04,187.44,0.47,-0.1,1.024,-4.5993", "intl", "New Zealand Geodetic Datum 1949", nil},
	{"OSGB36", "towgs84=446.448,-125.157,542.060,0.1502,0.2470,0.8421,-20.4894", "airy", "Airy 1830", nil},
}

/*type PrimeMeridian struct {
	ID         string
	Definition string
}

var primeMeridianTable = []PrimeMeridian{
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
*/
