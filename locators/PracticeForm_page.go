package locators

import "autogo/driver"

type PracticeFormType struct {
	Practice   driver.Locator
	Gender     driver.Locator
	Phone      driver.Locator
	Subjects   driver.Locator
	File       driver.Locator
	State      driver.Locator
	City       driver.Locator
	Submit     driver.Locator
	Submitting driver.Locator
}

// Динамический локатор
func (d PracticeFormType) UserFormField(name string) driver.Locator {
	return driver.Locator{
		Name:  "Поле: " + name,
		XPath: "//div[contains(@class, 'practice-form-wrapper')]//input[@id='" + name + "']",
	}
}

var PracticeForm = PracticeFormType{
	Practice: driver.Locator{
		Name:  "Practice Form",
		XPath: "//div[contains(@class, 'accordion')]//span[text()='Practice Form']",
	},

	Gender: driver.Locator{
		Name:  "Gender",
		XPath: "//input[@name='gender'][@value='Other']",
	},
	Phone: driver.Locator{
		Name:  "userNumber",
		XPath: "//div[@id='userNumber-wrapper']//input[@id='userNumber']",
	},

	Subjects: driver.Locator{
		Name:  "subjects",
		XPath: "//div[@id='subjectsContainer']//input[@id='subjectsInput']",
	},
	File: driver.Locator{
		Name:  "input file",
		XPath: "//input[@id='uploadPicture']",
	},
	State: driver.Locator{
		Name:  "input state",
		XPath: "//input[@id='react-select-3-input']",
	},
	City: driver.Locator{
		Name:  "input state",
		XPath: "//input[@id='react-select-4-input']",
	},
	Submit: driver.Locator{
		Name:  "submit",
		XPath: "//div[contains(@class,'content-end')]//button[@id='submit']",
	},
	Submitting: driver.Locator{
		Name:  "Submitting",
		XPath: "//div[text()='Thanks for submitting the form']",
	},
}
