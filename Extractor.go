package main

import "fmt"
import (
	//"encoding/xml"
	"github.com/vanng822/go-solr/solr"
	//"text/template/parse"
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"time"
)

type feed struct {
	Entry []Product `xml:"entry"`
}

//availability   condition brand
type Product struct {
	Avalilability    string    `xml:"g:availability"`
	brand            string    `xml:"g:brand"`
	Title            string    `xml:"g:title"`
	OP               string    `xml:"g:condition"`
	Price            float64   `xml:"g:price"`
	Id               string    `xml:"g:id"`
	Desc             string    `xml:"g:description"`
	B                string    `xml:"g:brand"`
	Link             string    `xml:"g:link"`
	ImageLink        string    `xml:"g:image_link"`
	DL               string    `xml:"g:deeplink"`
	Visits           float64   `xml:"g:visits"`
	Inventory        string    `xml:"g:inventory"`
	ProductType      string    `xml:"g:product_type"`
	Score            float64   `xml:"g:score"`
	CustomLabel      string    `xml:"g:custom_label_2"`
	CustomLabel3     float64   `xml:"g:custom_label_3"`
	ExpirationDate   string    `xml:"g:expiration_date"`
	Extras           []Applink `xml:"applink,omitempty"`
	IdentifierExists string    `xml:"g:identifier_exists"`
}
type Applink struct {
	Key   string `xml:"property,attr"`
	Value string `xml:"content,attr"`
}

func getParams(w http.ResponseWriter, r *http.Request) {

	keyword := r.URL.Query().Get("keyword")
	priceFrom := r.URL.Query().Get("priceFrom")
	category := r.URL.Query().Get("category")
	priceTo := r.URL.Query().Get("price")
	full := r.URL.Query().Get("full")
	if priceFrom != "" || priceTo != "" {
		priceFrom = "0"
		priceTo = "999999"
	}

	ipf, _ := strconv.Atoi(priceFrom)
	ipt, _ := strconv.Atoi(priceTo)
	makeXML(keyword, ipf, ipt, category, full)
}

func makeXML(keyword string, priceFrom int, priceTo int, category string, full string) {

	si, _ := solr.NewSolrInterface("http://solr-server/solr", "products-core")
	query := solr.NewQuery()
	query.Q("*:*")

	query.Rows(900000)


	s := si.Search(query)
	r, _ := s.Result(nil)
	fmt.Println(r.Results.Docs)
	docs := r.Results.Docs

	//gerando xml
	pArray := []Product{}

	for _, doc := range docs {
		//fmt.Println(i,doc.Get("name").(string))
		//fmt.Println(doc.Get("product_images"))

		deepLink := "protocol://product/" + doc.Get("object_id").(string)
		extra_deeplink_ios := Applink{"ios_url", deepLink}
		extra_store := Applink{"ios_app_store_id", "1031983829"}
		extra_deeplink_android := Applink{"android_url", deepLink}

		//android_app_name
		tagArray := []Applink{
			extra_deeplink_ios,
			extra_store,
			extra_ios_app_name,
			extra_deeplink_android,
			extra_packge_android,
			extra_android_appname}

		cid := doc.Get("category_id")
		if cid == nil {
			cid = 1
		}
		categoryName := categoryList[cid.(float64)]
		productLink := ""
		productLink = "https://yourstore.com/products/" + slug.Make(categoryName) + "/" + slug.Make(doc.Get("name").(string)) + "/" + doc.Get("object_id").(string)



		score2 := float64(3)

		score := doc.Get("score_prod")
		if score != nil {
			score2 = score.(float64)
		}
		}
	
		
		old := doc.Get("product_images").([]interface{})
		new := make([]interface{}, len(old))
		for i, v := range old {
    		new[i] = v
		}
		
		fmt.Println(new[0])
		productImage := new[0]

		
		namep := doc.Get("name")
		namep2 := ""
		if namep != nil {
			namep2 = namep.(string)
		}

		descText := namep2
		descText2 := doc.Get("description")
		if descText2 != nil {
			descText = descText2.(string)
		}

		if len(descText) > 0 {
			descText = strings.ToLower(descText)
		} else {
			descText = namep2
		}


		p := Product{
			Title:            strings.ToLower(namep2),
			Price:            doc.Get("price").(float64),
			Id:               doc.Get("object_id").(string),
			Desc:             descText,
			Inventory:        "1",
			Visits:           visits2,
			Score:            score2,
			CustomLabel:      usePayment,
			ProductType:      categoryName,
			OP:               "used",
			B:                "store",
			brand:            "store",
			Avalilability:    "in stock",
			CustomLabel3:     cid.(float64),
			ExpirationDate:   t.UTC().Format("2006-01-02T15:04:05-0700"),
			ImageLink:        productImage.(string),
			Link:             productLink,
			IdentifierExists: "false",
			Extras:           tagArray,
		}

		pArray = append(pArray, p)
	}

	prodList := &feed{
		Entry: pArray}

	out, err := xml.MarshalIndent(prodList, "", "   ")

	if err != nil {
		panic(err)
	}
	//fmt.Println(string(out))sudo
	if category == "todas" {
		keyword = keyword + category
	}
	file, err := os.Create("/var/www/xml/" + keyword + ".xml")
	defer file.Close()

	header := "<?xml version=\"1.0\" encoding=\"UTF-8\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\" xmlns:g=\"http://base.google.com/ns/1.0\"><title>Store</title><link rel=\"self\" href=\"https://store.com\"/>"
	xml := strings.Replace(string(out), "<feed>", header, 2)
	xml2 := xml
	xml3 := strings.Replace(string(xml2), "></applink>", "/>", -1)

	n3, err := file.WriteString(xml3)
	fmt.Println(xml3)
	fmt.Println(n3)
}

func main() {

	http.HandleFunc("/", getParams)          // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
