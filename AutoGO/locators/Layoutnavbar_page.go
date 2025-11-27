package locators

import "autogo/driver"

type LayoutNavbarPageType struct {
	Elements     driver.Locator
	Forms        driver.Locator
	Practice     driver.Locator
	Alerts       driver.Locator
	Widgets      driver.Locator
	Interactions driver.Locator
	BookStore    driver.Locator
}

var LayoutNavbarPage = LayoutNavbarPageType{
	Elements: driver.Locator{
		Name:  "Elements",
		XPath: "//div[contains(@class, 'accordion')]//*[text()='Elements']",
	},
	Forms: driver.Locator{
		Name:  "Forms",
		XPath: "//div[contains(@class, 'accordion')]//*[text()='Forms']/..",
	},
	Alerts: driver.Locator{
		Name:  "Alerts, Frame & Windows",
		XPath: "//div[contains(@class, 'accordion')]//*[text()='Alerts, Frame & Windows']",
	},
	Widgets: driver.Locator{
		Name:  "Widgets",
		XPath: "//div[contains(@class, 'accordion')]//*[text()='Widgets']",
	},
	Interactions: driver.Locator{
		Name:  "Interactions",
		XPath: "//div[contains(@class, 'accordion')]//*[text()='Interactions']",
	},
	BookStore: driver.Locator{
		Name:  "Book Store Application",
		XPath: "//div[contains(@class, 'accordion')]//*[text()='Book Store Application']",
	},
}
