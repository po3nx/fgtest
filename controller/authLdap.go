package controller

import (
    "fmt"
    "github.com/go-ldap/ldap"
	"github.com/po3nx/fgtest/config"
)

const (
    ldapServer   = config.Config("LDAP_SERVER")
    ldapPort     = config.Config("LDAP_PORT")
    ldapBindDN   = config.Config("LDAP_BINDN")
    ldapPassword = config.Config("LDAP_PASSWORD")
    ldapSearchDN = config.Config("LDAP_DN")
)
type UserLDAPData struct {
    ID       string
    Email    string
    Name     string
    FullName string
}

func AuthUsingLDAP(username, password string) (bool, *UserLDAPData, error) {
    l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
    if err != nil {
        return false, nil, err
    }
    defer l.Close()

    err = l.Bind(ldapBindDN, ldapPassword)
	if err != nil {
		return false, nil, err
	}
	searchRequest := ldap.NewSearchRequest(
		ldapSearchDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"dn", "cn", "sn", "mail"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, nil, err
	}

	if len(sr.Entries) == 0 {
		return false, nil, fmt.Errorf("User not found")
	}
	entry := sr.Entries[0]

	err = l.Bind(entry.DN, password)
	if err != nil {
		return false, nil, err
	}
	data := new(UserLDAPData)
	data.ID = username

	for _, attr := range entry.Attributes {
		switch attr.Name {
		case "sn":
			data.Name = attr.Values[0]
		case "mail":
			data.Email = attr.Values[0]
		case "cn":
			data.FullName = attr.Values[0]
		}
	}

	return true, data, nil
}