package handlers

import (
	"github/Ely0rda/projectMic/product_api/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {

		//For Put requests we need to get the id of the element
		//we want to update this why we need to get the path of the URL
		//and we will be using regex to get the id from the URL

		//URL is a member in the Request struct
		// URL specifies either the URI being requested (for server
		// requests) or the URL to access (for client requests).
		//Path is a member in the URL struct
		// path (relative paths may omit leading slash)

		//MustCompile
		//Using the expression giving to it, it returns a regex that can be matched against text
		//and if it can not parsed that expression ,it panics.
		regex := regexp.MustCompile(`/([0-9]+)`)
		// FindAllStringSubmatch is the 'All' version of FindStringSubmatch
		// FindStringSubmatch returns a slice of strings holding the text of the
		// leftmost match of the regular expression in s and the Submatches
		//Submatches are matches of parenthesized subexpressions (also known as capturing groups) within the regular expression
		//If 'All' is present, the routine matches successive non-overlapping matches of the entire expression
		// if n>0 it return at most n matches
		//if not it return all the matches, this why I put n as -1
		//why use FindAllStringSubmatch to handle the case
		//where more than one id was passed as a BadRequest
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		//g should be something like this [["ab" ""]]
		//because it should be there one match
		//this why we should check the len(g) and it suppose
		//to be 1 something like the slice from before.
		if len(g) != 1 {
			http.Error(rw, "Invalid URI more than one id", http.StatusBadRequest)
			return
		}
		//and the lentgh of g[0] should be 2 one match and
		//one submatch [["ab" ""]]
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URL more than one captured group", http.StatusBadRequest)
		}
		//remeber that idString is a string
		idString := g[0][1]
		//converting id form string to int
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI unable to convert the idString into int", http.StatusBadRequest)
		}
		p.l.Println("got id", id)
		p.updateProduct(id, rw, r)

		return
	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

// getProducts returns the products from the data store
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handler PUT Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
}
