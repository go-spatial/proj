package api

type pair struct {
	source, dest int
}

type pairTable map[pair]string

var knownProjections = pairTable{
	pair{4326, 3395}: "+proj=merc +lon_0=0 +k=1 +x_0=0 +y_0=0 +datum=WGS84", // TODO: support +units=m +no_defs
	pair{4326, 3857}: "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext +no_defs",
	// TODO: add (3857,3395)?
}
