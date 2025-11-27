package locators

import (
	"autogo/driver"
)

type CommonLocatorsPageType struct {
}

func (d CommonLocatorsPageType) Div(div string) driver.Locator {
	return driver.Locator{
		Name:  div,
		XPath: "//div[text()='" + div + "']",
	}
}

func (d CommonLocatorsPageType) Span(span string) driver.Locator {
	return driver.Locator{
		Name:  span,
		XPath: "//Span[@placeholder='" + span + "']",
	}
}

func (d CommonLocatorsPageType) Text(text string) driver.Locator {
	return driver.Locator{
		Name:  text,
		XPath: "//input[@placeholder='" + text + "']",
	}
}

func (d CommonLocatorsPageType) Button(button string) driver.Locator {
	return driver.Locator{
		Name:  button,
		XPath: "//input[@placeholder='" + button + "']",
	}
}

var CommonPage = CommonLocatorsPageType{}
