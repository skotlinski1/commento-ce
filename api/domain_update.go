package main

import (
	"net/http"
)

func domainUpdate(d domain) error {
	statement := `
		UPDATE domains
    SET name=$2, state=$3, autoSpamFilter=$4, requireModeration=$5, requireIdentification=$6
		WHERE domain=$1;
	`

	_, err := db.Exec(statement, d.Domain, d.Name, d.State, d.AutoSpamFilter, d.RequireModeration, d.RequireIdentification)
	if err != nil {
		logger.Errorf("cannot update non-moderators: %v", err)
		return errorInternal
	}

	return nil
}

func domainUpdateHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Session *string `json:"session"`
		D       *domain `json:"domain"`
	}

	var x request
	if err := unmarshalBody(r, &x); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetBySession(*x.Session)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := stripDomain((*x.D).Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		writeBody(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	if err = domainUpdate(*x.D); err != nil {
		writeBody(w, response{"success": false, "message": err.Error()})
		return
	}

	writeBody(w, response{"success": true})
}