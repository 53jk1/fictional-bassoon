package main //main.go
// Language: go
import ( //importing packages
	"encoding/json" // for json
	"fmt"           // for printing
	"io/ioutil"     // for reading the website
	"log"           //logging
	"net/http"      //importing http
	"strconv"       //for converting string to int
	"strings"       //for strings.Join
	"text/template" // for the html
)

// Don't need dates/locations because relations already exist.

// StringArtists - struct for the main page. Save everything in one place for later.
type StringArtists struct { // struct for the main page
	ID       string              //ID
	Image    string              //Image
	Name     string              //Name of the band
	Members  string              //members
	Creation string              //creation date
	First    string              //first album
	Relation map[string][]string // map[band name] = [relation]
} // struct for the individual page. Save everything in one place for later.

// Artists - struct for the main page.
type Artists struct { //struct for the main page
	ID       int      `json:"id"`           //ID
	Image    string   `json:"image"`        //image
	Name     string   `json:"name"`         //list of names
	Members  []string `json:"members"`      //list of members
	Creation int      `json:"creationDate"` //int because it's a number
	First    string   `json:"firstAlbum"`   //first album
} //end of struct

// Relations - struct for the main page.
type Relations struct {
	Relation map[string][]string `json:"datesLocations"` // map[artistID][]string
} // saves it to "Relations"

// AegKoht - struct for the main page.
type AegKoht struct { //struct for the main page
	Index []Relations `json:"index"` //index is the name of the JSON file
} //Relations is the name of the struct

var artistid []Artists // saves the info from website
var ajadKohad AegKoht  // saves the relations

var nbr int    //nbr == band's ID
var m []string // for auditors convinience displaying the names

// Unmarshal - unmarshal the JSON file.
func Unmarshal() error { //Take the info from website and save it to struct

	names, err := http.Get("https://groupietrackers.herokuapp.com/api/artists") //get the info from website
	if err != nil {                                                             //if there is an error
		return fmt.Errorf("Error happened in http.Get. Err: %s", err) //log the error
	} //log the error
	bytes, err := ioutil.ReadAll(names.Body) //read the info from website
	if err != nil {                          //if there is an error
		return fmt.Errorf("Error happened in ioutil.ReadAll. Err: %s", err)
	} //end of if

	err = json.Unmarshal(bytes, &artistid) //saves it to "Artists"
	if err != nil {                        //if there is an error
		return fmt.Errorf("Error happened in json.Unmarshal. Err: %s", err) //log the error
	} //end of if

	relation, err := http.Get("https://groupietrackers.herokuapp.com/api/relation") //get the info from website
	if err != nil {                                                                 //if there is an error
		return fmt.Errorf("Error happened in http.Get. Err: %s", err)
	}
	Rbytes, err := ioutil.ReadAll(relation.Body) //read the info from website
	if err != nil {                              //if there is an error
		return fmt.Errorf("Error happened in ioutil.ReadAll. Err: %s", err) //log the error
	} //end of if
	err = json.Unmarshal(Rbytes, &ajadKohad) // saves it to "Relations"
	if err != nil {                          //if there is an error
		return fmt.Errorf("Error happened in json.Unmarshal. Err: %s", err) //log the error
	}
	return nil //if there is no error
} //end of Unmarshal

// FindByName - find the artist by name.
func FindByName(input string) { //input == what user entered on the website

	x := 0 // x == band's ID

	for x <= len(artistid) { //for every band

		if input == artistid[x].Name { // if the input is the same as the name in the list
			nbr = x // nbr == band's ID
			break   // break the loop
		}
		x++ // x == band's ID

	}

	nbr = x //nbr == band's ID
}

// InduvidualPage - individual page.
func InduvidualPage(w http.ResponseWriter, r *http.Request) { //individual page

	switch r.Method { //switch for the method
	case "GET": //if the method is GET
		http.ServeFile(w, r, "templates/individual.html") //use the info listed in "uwu" for the main page (templates/main.html)
	case "POST": //if the method is POST
		if r.FormValue("NAME") != "" { //Go to the "FindByName" function and use "NAME" (from html site) as input
			FindByName(r.FormValue("NAME")) //Go to the "FindByName" function and use "NAME" (from html site) as input

		}

		uwu := StringArtists{ // converting everything to string because the struct wants it that way
			ID:       ("ID: " + strconv.Itoa(artistid[nbr].ID)),                  //ID
			Image:    (artistid[nbr].Image),                                      //Image
			Name:     ("Name: " + artistid[nbr].Name),                            //Name
			Members:  ("Members: " + strings.Join(artistid[nbr].Members, ", ")),  //Members
			Creation: ("Creation date: " + strconv.Itoa(artistid[nbr].Creation)), //Creation date
			First:    ("First album: " + artistid[nbr].First),                    //First album
			Relation: (ajadKohad.Index[nbr].Relation),                            //Relation
		}

		t, err := template.ParseFiles("templates/individual.html") //parse the html file
		if err != nil {                                            //if there is an error
			log.Fatalf("Error happened in template.ParseFiles. Err: %s", err) //log the error
		}
		if uwu.Relation == nil { // if there is no relation, then the map is empty
			uwu.Relation = make(map[string][]string) // make the map empty
		} //end of if
		err = t.Execute(w, uwu) // use the info listed in "uwu" for the individual page (templates/individual.html)
		if err != nil {         //if there is an error
			log.Fatalf("Error happened in t.Execute. Err: %s", err) //log the error
		}
	} //end of switch

} //end of InduvidualPage

// Page - main page.
func Page(w http.ResponseWriter, r *http.Request) { //main page

	if r.URL.Path != "/" { //if the path is not "/"
		http.Error(w, "404 not found.", http.StatusNotFound) //error 404
		return                                               //return
	} //end of if

	for x := range artistid { // add all the names to "m"
		m = append(m, artistid[x].Name) //add all the names to "m"
	}
	uwu := Artists{ // from list to single string seperated by comma
		Name: (strings.Join(m, ", ")), //Name
	} //end of struct
	t, err := template.ParseFiles("templates/main.html") //parse the html file
	if err != nil {                                      //if there is an error
		log.Fatalf("Error happened in template.ParseFiles. Err: %s", err) //log the error
	} //end of if
	err = t.Execute(w, uwu) // use the info listed in "uwu" for the main page (templates/main.html)
	if err != nil {         //if there is an error
		log.Fatalf("Error happened in t.Execute. Err: %s", err) //log the error
	} //end of if

} //end of Page

func main() { //main function
	//Order counts
	err := Unmarshal() //Take the info from website and save it to struct

	if err != nil { //if there is an error
		log.Fatalf("Error happened in Unmarshal. Err: %s", err) //log the error
	} //end of if

	port := ":8099" //port

	http.HandleFunc("/", Page)                     //main page
	http.HandleFunc("/Individual", InduvidualPage) //individual page

	fmt.Printf("Server listen to - localhost%s\n", port) //print the port

	err = http.ListenAndServe(port, nil) //start the server
	if err != nil {                      //if there is an error
		log.Fatalf("Error happened in http.ListenAndServe. Err: %s", err) //log the error
	} //end of if
} //end of main
