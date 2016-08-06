package gopublica

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type Cosponsors struct {
	Result
	Results []struct {
		Cosponsors []struct {
			CosponsorId string `json:"cosponsor_id"`
			Name        string `json:"name"`
			Date        string `json:"date"`
		} `json:"cosponsors"`
		Congress              string `json:"congress"`
		Number                string `json:"number"`
		BillUri               string `json:"bill_uri"`
		Title                 string `json:"title"`
		SponsorId             string `json:"sponsor_id"`
		IntroducedDate        string `json:"introduced_date"`
		Committees            string `json:"committees"`
		LatestMajorActionDate string `json:"latest_major_action_date"`
		LatestMajorAction     string `json:"latest_major_action"`
	} `json:"results"`
}

type Amendments struct {
	Result
	Results []struct {
		Congress              string `json:"congress"`
		Bill                  string `json:"bill"`
		UrlNumber             string `json:"url_number"`
		Title                 string `json:"title"`
		Sponsor               string `json:"sponsor"`
		SponsorId             string `json:"sponsor_id"`
		IntroducedDate        string `json:"introduced_date"`
		Cosponsors            string `json:"cosponsors,string"`
		Committees            string `json:"committees"`
		LatestMajorActionDate string `json:"latest_major_action_date"`
		LatestMajorAction     string `json:"latest_major_action"`
		HousePassageVote      string `json:"house_passage_vote"`
		SenatePassageVote     string `json:"senate_passage_vote"`
		Amendments            []struct {
			AmendmentNumber       string `json:"amendment_number"`
			SponsorId             string `json:"sponsor_id"`
			IntroducedDate        string `json:"introduced_date"`
			Title                 string `json:"title"`
			LatestMajorActionDate string `json:"latest_major_action_date"`
			LatestMajorAction     string `json:"latest_major_action"`
		} `json:"amendments"`
	}
}


// :congress/bills/:bill-id/cosponsors.js
//
func GetBillCosponsors(congress, billId string) (*Cosponsors, error) {
	rsp, err := get(congress, "bills", billId, "cosponsors")
	if err != nil {
		return nil, err
	}
	bod, _ := ioutil.ReadAll(rsp.Body)
	res := &Cosponsors{}
	json.Unmarshal(bod, res)
	code, _ := strconv.Atoi(res.Result.Status)
	if !statusOk(code) && res.Status != "OK" {
		return nil, getError(code)
	}

	return res, nil
}

// :congress/bills/:bill-id/amendments
//
func GetBillAmendments(congress, billId string) (*Amendments, error) {
	rsp, err := get(congress, "bills", billId, "amendments")
	if err != nil {
		return nil, err
	}
	bod, _ := ioutil.ReadAll(rsp.Body)
	res := &Amendments{}
	json.Unmarshal(bod, res)
	code, _ := strconv.Atoi(res.Result.Status)
	if !statusOk(code) && res.Status != "OK" {
		return nil, getError(code)
	}

	return res, nil
}
