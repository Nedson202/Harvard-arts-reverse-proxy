package reverse_proxy

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v7"
)

type App struct {
	baseURL       string
	harvardAPIKey string
	redisClient   *redis.Client
	elasticClient *elasticsearch.Client
}

// Route defines a structure for routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes defines the list of routes of our API
type Routes []Route

// RootPayload structure for error responses
type RootPayload struct {
	Error   bool        `json:"error"`
	Payload interface{} `json:"payload"`
}

// DataPayload structure for error responses
type DataPayload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type RecordsPayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Records interface{} `json:"records"`
}

type RecordPayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Record  interface{} `json:"record"`
}

type CollectionsResponse struct {
	Records []CollectionsObject `json:"records"`
}

type PublicationsResponse struct {
	Records []interface{} `json:"records"`
}

type PublicationsPayload struct {
	Error        bool          `json:"error"`
	Message      string        `json:"message"`
	Publications []interface{} `json:"publications"`
}

type PlaceIdPayload struct {
	Error   bool      `json:"error"`
	Message string    `json:"message"`
	Places  []PlaceID `json:"places"`
}

type PlacesPayload struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Places  []Place `json:"places"`
}

type SearchResultsPayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Results interface{} `json:"results"`
}

type PlaceID struct {
	ParentPlaceId int64  `json:"parentPlaceID"`
	PathForward   string `json:"pathForward"`
}

type Place struct {
	Objectcount   int    `json:"objectcount"`
	ID            int    `json:"id"`
	LastUpdate    string `json:"lastupdate"`
	HasChildren   int    `json:"haschildren"`
	Level         int    `json:"level"`
	PlaceID       int    `json:"placeid"`
	PathForward   string `json:"pathforward"`
	ParentPlaceID int64  `json:"parentplaceid"`
	Name          string `json:"name"`
	TgnID         int    `json:"tgn_id,omitempty"`
	Geo           struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"geo,omitempty"`
}

type CollectionsObject struct {
	Accessionmethod     string      `json:"accessionmethod"`
	Accessionyear       int         `json:"accessionyear"`
	Accesslevel         int         `json:"accesslevel"`
	Century             string      `json:"century"`
	Classification      string      `json:"classification"`
	Classificationid    int         `json:"classificationid"`
	Colorcount          int         `json:"colorcount"`
	Commentary          interface{} `json:"commentary"`
	Contact             string      `json:"contact"`
	Contextualtextcount int         `json:"contextualtextcount"`
	Copyright           interface{} `json:"copyright"`
	Creditline          string      `json:"creditline"`
	Culture             string      `json:"culture"`
	Datebegin           int         `json:"datebegin"`
	Dated               string      `json:"dated"`
	Dateend             int         `json:"dateend"`
	Dateoffirstpageview string      `json:"dateoffirstpageview,omitempty"`
	Dateoflastpageview  string      `json:"dateoflastpageview,omitempty"`
	Department          string      `json:"department"`
	Description         interface{} `json:"description"`
	Details             struct {
		Coins struct {
			Dateonobject       interface{} `json:"dateonobject"`
			Denomination       interface{} `json:"denomination"`
			Dieaxis            interface{} `json:"dieaxis"`
			Metal              string      `json:"metal"`
			Obverseinscription interface{} `json:"obverseinscription"`
			Reverseinscription interface{} `json:"reverseinscription"`
		} `json:"coins"`
	} `json:"details"`
	Dimensions           string      `json:"dimensions"`
	Division             string      `json:"division"`
	Edition              interface{} `json:"edition"`
	Exhibitioncount      int         `json:"exhibitioncount"`
	Groupcount           int         `json:"groupcount"`
	ID                   int         `json:"id"`
	Imagecount           int         `json:"imagecount"`
	Imagepermissionlevel int         `json:"imagepermissionlevel"`
	Images               []struct {
		Baseimageurl    string      `json:"baseimageurl"`
		Copyright       string      `json:"copyright"`
		Displayorder    int         `json:"displayorder"`
		Format          string      `json:"format"`
		Height          int         `json:"height"`
		Idsid           int         `json:"idsid"`
		Iiifbaseuri     string      `json:"iiifbaseuri"`
		Imageid         int         `json:"imageid"`
		Publiccaption   interface{} `json:"publiccaption"`
		Renditionnumber string      `json:"renditionnumber"`
		Width           int         `json:"width"`
	} `json:"images"`
	Labeltext    interface{} `json:"labeltext"`
	Lastupdate   string      `json:"lastupdate"`
	Markscount   int         `json:"markscount"`
	Mediacount   int         `json:"mediacount"`
	Medium       string      `json:"medium"`
	Objectid     int         `json:"objectid"`
	Objectnumber string      `json:"objectnumber"`
	Peoplecount  int         `json:"peoplecount"`
	Period       interface{} `json:"period"`
	Periodid     interface{} `json:"periodid"`
	Places       []struct {
		Displayname string `json:"displayname"`
		Placeid     int    `json:"placeid"`
		Type        string `json:"type"`
	} `json:"places"`
	Primaryimageurl  string      `json:"primaryimageurl"`
	Provenance       interface{} `json:"provenance"`
	Publicationcount int         `json:"publicationcount"`
	Rank             int         `json:"rank"`
	Relatedcount     int         `json:"relatedcount"`
	SeeAlso          []struct {
		Format  string `json:"format"`
		ID      string `json:"id"`
		Profile string `json:"profile"`
		Type    string `json:"type"`
	} `json:"seeAlso"`
	Signed                  interface{} `json:"signed"`
	Standardreferencenumber interface{} `json:"standardreferencenumber"`
	State                   interface{} `json:"state"`
	Style                   interface{} `json:"style"`
	Technique               string      `json:"technique"`
	Techniqueid             int         `json:"techniqueid"`
	Terms                   struct {
		Century []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"century"`
		Culture []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"culture"`
		Medium []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"medium"`
		Place []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"place"`
	} `json:"terms"`
	Title  string `json:"title"`
	Titles []struct {
		Displayorder int    `json:"displayorder"`
		Title        string `json:"title"`
		Titleid      int    `json:"titleid"`
		Titletype    string `json:"titletype"`
	} `json:"titles"`
	Titlescount                  int    `json:"titlescount"`
	Totalpageviews               int    `json:"totalpageviews"`
	Totaluniquepageviews         int    `json:"totaluniquepageviews"`
	URL                          string `json:"url"`
	Verificationlevel            int    `json:"verificationlevel"`
	Verificationleveldescription string `json:"verificationleveldescription"`
	Worktypes                    []struct {
		Worktype   string `json:"worktype"`
		Worktypeid string `json:"worktypeid"`
	} `json:"worktypes"`
}
