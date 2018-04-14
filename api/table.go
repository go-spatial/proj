package api

type pair struct {
	source, dest int
}

type pairTable map[pair]string

var knownProjections = pairTable{
	pair{4326, 3395}: "+proj=merc +lon_0=0 +k=1 +x_0=0 +y_0=0 +ellps=WGS84 +datum=WGS84", // +units=m +no_defs",
}
