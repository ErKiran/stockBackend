package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"stockwatch/api/auth"
	"stockwatch/api/responses"
	"stockwatch/models"
	"strconv"

	"github.com/gorilla/mux"
)

type StockData []struct {
	ID           int    `json:"id"`
	Symbol       string `json:"symbol"`
	SecurityName string `json:"securityName"`
	Name         string `json:"name"`
	ActiveStatus string `json:"activeStatus"`
}

type SecurityDetails struct {
	SecurityDailyTradeDto struct {
		SecurityID          string      `json:"securityId"`
		OpenPrice           float64     `json:"openPrice"`
		HighPrice           float64     `json:"highPrice"`
		LowPrice            float64     `json:"lowPrice"`
		TotalTradeQuantity  int         `json:"totalTradeQuantity"`
		TotalTrades         int         `json:"totalTrades"`
		LastTradedPrice     float64     `json:"lastTradedPrice"`
		PreviousClose       float64     `json:"previousClose"`
		BusinessDate        string      `json:"businessDate"`
		ClosePrice          interface{} `json:"closePrice"`
		FiftyTwoWeekHigh    float64     `json:"fiftyTwoWeekHigh"`
		FiftyTwoWeekLow     float64     `json:"fiftyTwoWeekLow"`
		LastUpdatedDateTime string      `json:"lastUpdatedDateTime"`
	} `json:"securityDailyTradeDto"`
	Security struct {
		ID               int         `json:"id"`
		Symbol           string      `json:"symbol"`
		Isin             string      `json:"isin"`
		PermittedToTrade string      `json:"permittedToTrade"`
		ListingDate      string      `json:"listingDate"`
		CreditRating     interface{} `json:"creditRating"`
		TickSize         float64     `json:"tickSize"`
		InstrumentType   struct {
			ID           int    `json:"id"`
			Code         string `json:"code"`
			Description  string `json:"description"`
			ActiveStatus string `json:"activeStatus"`
		} `json:"instrumentType"`
		CapitalGainBaseDate string      `json:"capitalGainBaseDate"`
		FaceValue           float64     `json:"faceValue"`
		HighRangeDPR        interface{} `json:"highRangeDPR"`
		IssuerName          interface{} `json:"issuerName"`
		MeInstanceNumber    int         `json:"meInstanceNumber"`
		ParentID            interface{} `json:"parentId"`
		RecordType          int         `json:"recordType"`
		SchemeDescription   interface{} `json:"schemeDescription"`
		SchemeName          interface{} `json:"schemeName"`
		Secured             interface{} `json:"secured"`
		Series              interface{} `json:"series"`
		ShareGroupID        struct {
			ID              int         `json:"id"`
			Name            string      `json:"name"`
			Description     string      `json:"description"`
			CapitalRangeMin int         `json:"capitalRangeMin"`
			ModifiedBy      interface{} `json:"modifiedBy"`
			ModifiedDate    interface{} `json:"modifiedDate"`
			ActiveStatus    string      `json:"activeStatus"`
			IsDefault       string      `json:"isDefault"`
		} `json:"shareGroupId"`
		ActiveStatus       string  `json:"activeStatus"`
		Divisor            int     `json:"divisor"`
		CdsStockRefID      int     `json:"cdsStockRefId"`
		SecurityName       string  `json:"securityName"`
		TradingStartDate   string  `json:"tradingStartDate"`
		NetworthBasePrice  float64 `json:"networthBasePrice"`
		SecurityTradeCycle int     `json:"securityTradeCycle"`
		IsPromoter         string  `json:"isPromoter"`
		CompanyID          struct {
			ID                   int    `json:"id"`
			CompanyShortName     string `json:"companyShortName"`
			CompanyName          string `json:"companyName"`
			Email                string `json:"email"`
			CompanyWebsite       string `json:"companyWebsite"`
			CompanyContactPerson string `json:"companyContactPerson"`
			SectorMaster         struct {
				ID                int    `json:"id"`
				SectorDescription string `json:"sectorDescription"`
				ActiveStatus      string `json:"activeStatus"`
				RegulatoryBody    string `json:"regulatoryBody"`
			} `json:"sectorMaster"`
			CompanyRegistrationNumber string `json:"companyRegistrationNumber"`
			ActiveStatus              string `json:"activeStatus"`
		} `json:"companyId"`
	} `json:"security"`
	StockListedShares    float64 `json:"stockListedShares"`
	PaidUpCapital        int     `json:"paidUpCapital"`
	IssuedCapital        int     `json:"issuedCapital"`
	MarketCapitalization int64   `json:"marketCapitalization"`
	PublicShares         int     `json:"publicShares"`
	PublicPercentage     float64 `json:"publicPercentage"`
	PromoterShares       float64 `json:"promoterShares"`
	PromoterPercentage   float64 `json:"promoterPercentage"`
	UpdatedDate          string  `json:"updatedDate"`
	SecurityID           int     `json:"securityId"`
}

const STOCK_API = "https://newweb.nepalstock.com/api/nots/security"

func (server *Server) GetStockScript(w http.ResponseWriter, r *http.Request) {
	stockAPI := "https://newweb.nepalstock.com/api/nots/security?nonDelisted=true"
	resp, err := http.Get(stockAPI)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var stockData StockData
	err = json.Unmarshal(body, &stockData)
	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    stockData,
	})
}

func (server *Server) GetStockDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	stockID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	stockAPI := fmt.Sprintf("%v/%d", STOCK_API, stockID)
	resp, err := http.Get(stockAPI)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var stockData SecurityDetails
	err = json.Unmarshal(body, &stockData)
	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    stockData,
	})
}

func (server *Server) CreateStockWatch(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	stockWatch := models.StockWatch{}
	err = json.Unmarshal(body, &stockWatch)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	stockWatch.Prepare()
	err = stockWatch.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	stockWatch.UserID = int64(userID)

	savedStockWatch, err := stockWatch.Save(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, savedStockWatch)
}

func (server *Server) GetStockWatchOfUser(w http.ResponseWriter, r *http.Request) {
	stockWatch := models.StockWatch{}

	userID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	stockInfo, err := stockWatch.GetStockWatchByUserID(server.DB, int64(userID))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"list":    stockInfo,
	})
}
