package main

import (
	"encoding/json"
	"net/http"

	"github.com/cenkalti/log"
)

func handleAdminGetPayment(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if username != "admin" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if password != config.AdminPassword {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	account := r.FormValue("account")
	if account == "" {
		http.Error(w, "invalid account", http.StatusBadRequest)
		return
	}
	payment, err := LoadPayment([]byte(account))
	if err == errPaymentNotFound {
		log.Debugln("account not found:", account)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := NewResponse(payment, "")
	b, err := json.MarshalIndent(&response, "", "  ")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Debug(err)
	}
}

func handleAdminCheckPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	if username != "admin" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if password != config.AdminPassword {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	account := r.FormValue("account")
	if account == "" {
		http.Error(w, "invalid account", http.StatusBadRequest)
		return
	}
	payment, err := LoadPayment([]byte(account))
	if err == errPaymentNotFound {
		log.Debugln("account not found:", account)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = payment.check()
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	response := NewResponse(payment, "")
	b, err := json.MarshalIndent(&response, "", "  ")
	if err != nil {
		log.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Debug(err)
	}
}
