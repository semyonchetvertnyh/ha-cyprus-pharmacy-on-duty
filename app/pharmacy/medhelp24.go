package pharmacy

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/semyonchetvertnyh/ha-cyprus-pharmacy-on-duty/app/geo"
)

type MedHelp24Parser struct {
	url     string
	scraper *colly.Collector
}

func NewMedHelp24Parser() *MedHelp24Parser {
	return &MedHelp24Parser{
		url:     "https://medhelp24.com/EN/pharmacies/today",
		scraper: colly.NewCollector(),
	}
}

func (p *MedHelp24Parser) Parse() ([]Pharmacy, error) {
	urls, err := p.parseMainPage()
	if err != nil {
		return nil, err
	}

	var pharmacies []Pharmacy
	for _, url := range urls {
		pharmacy, err := p.parseDetailsPage(url)
		if err != nil {
			return nil, err
		}

		pharmacies = append(pharmacies, pharmacy)
	}
	return pharmacies, nil
}

func (p *MedHelp24Parser) parseMainPage() ([]string, error) {
	var urls []string

	p.scraper.OnHTML("div.view-night-pharmacies-today > div.view-content", func(e *colly.HTMLElement) {
		e.ForEach("div.item-list ul > li a", func(_ int, e *colly.HTMLElement) {
			urls = append(urls, e.Request.AbsoluteURL(e.Attr("href")))
		})
	})
	if err := p.scraper.Visit(p.url); err != nil {
		return nil, fmt.Errorf("failed to visit main page: %w", err)
	}

	return urls, nil
}

func (p *MedHelp24Parser) parseDetailsPage(url string) (Pharmacy, error) {
	pharmacy := &Pharmacy{}

	p.scraper.OnHTML("div.field-name-field-municipality div.field-item", func(e *colly.HTMLElement) {
		pharmacy.Municipality = sanitize(e.Text)
	})
	p.scraper.OnHTML("div.field-name-field-pharmacy-region div.field-item", func(e *colly.HTMLElement) {
		pharmacy.City = sanitize(e.Text)
	})
	p.scraper.OnHTML("div.field-name-field-pharmacy-address div.field-item", func(e *colly.HTMLElement) {
		pharmacy.Address = sanitize(e.Text)
	})
	p.scraper.OnHTML("div.field-name-field-pharmacy-address-info div.field-item", func(e *colly.HTMLElement) {
		pharmacy.Instructions = sanitize(e.Text)
	})
	p.scraper.OnHTML("div.field-name-field-pharmacy-description div.field-item", func(e *colly.HTMLElement) {
		pharmacy.Instructions += "\n" + sanitize(e.Text)
	})

	p.scraper.OnHTML("div.pane-content > a", func(e *colly.HTMLElement) {
		pharmacy.GoogleMapsLink = e.Attr("href")
	})

	// p.scraper.OnHTML("div.field-name-field-pharmacy-name div.field-item", func(e *colly.HTMLElement) {
	// 	pharmacy.Contacts[0].Name = sanitize(e.Text)
	// })
	p.scraper.OnHTML("div.field-name-field-pharmacy-phone div.field-item", func(e *colly.HTMLElement) {
		pharmacy.Phone = sanitizePhone(e.Text)
	})
	p.scraper.OnHTML("div.field-name-field-pharmacy-phone-2 div.field-item", func(e *colly.HTMLElement) {
		pharmacy.SecondPhone = sanitizePhone(e.Text)
	})

	if err := p.scraper.Visit(url); err != nil {
		return Pharmacy{}, fmt.Errorf("failed to visit main page: %w", err)
	}

	if pharmacy.GoogleMapsLink != "" {
		pharmacy.GeoPoint = geo.NewPointFromGoogleMapsLink(pharmacy.GoogleMapsLink)
	}

	// Hardcode.
	if url == "https://medhelp24.com/EN/content/ioannoy-alexandros" {
		pharmacy.City = "Limassol"
	}

	return *pharmacy, nil
}

func sanitizePhone(p string) string {
	p = sanitize(p)
	p = strings.ReplaceAll(p, " ", "")

	if !strings.HasPrefix("+357", p) {
		p = "+357 " + p
	}
	return p
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimSpace(s)
	return s
}
