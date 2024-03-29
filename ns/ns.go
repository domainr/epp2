package ns

import (
	"github.com/domainr/epp2/schema/contact"
	"github.com/domainr/epp2/schema/domain"
	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/eppcom"
	"github.com/domainr/epp2/schema/host"
)

// TODO: get rid of this file and package.

const (
	// EPP is the IETF URN for the EPP namespace.
	// See https://www.iana.org/assignments/xml-registry/ns/epp-1.0.txt.
	EPP = epp.NS

	// Common is the IETF URN for the EPP common namespace.
	// See https://www.iana.org/assignments/xml-registry/ns/eppcom-1.0.txt.
	Common = eppcom.NS

	// Host is the IETF URN for the EPP contact namespace.
	// See https://www.iana.org/assignments/xml-registry/ns/contact-1.0.txt.
	Contact = contact.NS

	// Domain is the IETF URN for the EPP domain namespace.
	// See https://www.iana.org/assignments/xml-registry/ns/domain-1.0.txt
	// and https://datatracker.ietf.org/doc/html/rfc5731.
	Domain = domain.NS

	// Host is the IETF URN for the EPP host namespace.
	// See https://www.iana.org/assignments/xml-registry/ns/host-1.0.txt.
	Host = host.NS

	// SecDNS is the IETF URN for the EPP DNSSEC namespace.
	// See https://datatracker.ietf.org/doc/html/rfc5910.
	SecDNS = "urn:ietf:params:xml:ns:secDNS-1.1"

	Fee05      = "urn:ietf:params:xml:ns:fee-0.5"
	Fee06      = "urn:ietf:params:xml:ns:fee-0.6"
	Fee07      = "urn:ietf:params:xml:ns:fee-0.7"
	Fee08      = "urn:ietf:params:xml:ns:fee-0.8"
	Fee09      = "urn:ietf:params:xml:ns:fee-0.9"
	Fee10      = "urn:ietf:params:xml:ns:epp:fee-1.0"
	Fee11      = "urn:ietf:params:xml:ns:fee-0.11"
	Fee21      = "urn:ietf:params:xml:ns:fee-0.21"
	IDN        = "urn:ietf:params:xml:ns:idn-1.0"
	Launch     = "urn:ietf:params:xml:ns:launch-1.0"
	Neulevel   = "urn:ietf:params:xml:ns:neulevel"
	Neulevel10 = "urn:ietf:params:xml:ns:neulevel-1.0"
	Price      = "urn:ar:params:xml:ns:price-1.1"
	RGP        = "urn:ietf:params:xml:ns:rgp-1.0"

	Finance   = "http://www.unitedtld.com/epp/finance-1.0"
	Charge    = "http://www.unitedtld.com/epp/charge-1.0"
	Namestore = "http://www.verisign-grs.com/epp/namestoreExt-1.1"
)
