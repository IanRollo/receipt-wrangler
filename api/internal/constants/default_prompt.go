package constants

// PreviousDefaultPromptText is the default prompt as every release before multicurrency
// created it. the "upgrade-default-prompt-currency" data migration (repositories) recognises a stored default prompt by this exact
// text: an install that never edited it is moved to DefaultPromptText, one that did is left
// alone.
const PreviousDefaultPromptText = `
Find the receipt's name, total cost, and date. Format the found data as:
{
	"name": store name,
	"amount": amount as a number,
	"date": date in ISO 18601 format in UTC with ALL time values set as 0,
	"categories": categories,
	"tags": tags
}
If a store name cannot be confidently found, use 'Default store name' as the default name.
Omit any value if not found with confidence. Assume the date is in the year @currentYear if not provided.
The amount must be a float or integer.
If the receipt represents a refund, return, or credit (money returned to the customer rather than charged), the amount and item amounts MUST be negative. Otherwise they are positive.

Please do NOT add any additional information, only valid JSON.
Please return the json in plaintext ONLY, do not ever return it in a code block or any other format.

Choose up to 2 categories from the given list based on the receipt's items and store name. If no categories fit, please return an empty array for the field and do not select any categories. When selecting categories, select only the id, like:
{
	Id: category id
}

Emphasize the relationship between the category and the receipt, and use the description of the category to fine tune the results. Do not return categories that have an empty name or do not exist.
If there are no categories to chose from, then please make categories an empty array.
Likewise, if there are not tags to choose from, then make tags an empty array.

Categories to chose from: @categories

Follow the same process as described for categories for tags.

Tags to chose from: @tags

Receipt text: @ocrText
`

// DefaultPromptText is the default prompt for new installs: PreviousDefaultPromptText plus the
// currency key and its instructions.
const DefaultPromptText = `
Find the receipt's name, total cost, and date. Format the found data as:
{
	"name": store name,
	"amount": amount as a number,
	"currency": ISO 4217 currency code of the amounts,
	"date": date in ISO 18601 format in UTC with ALL time values set as 0,
	"categories": categories,
	"tags": tags
}
If a store name cannot be confidently found, use 'Default store name' as the default name.
Omit any value if not found with confidence. Assume the date is in the year @currentYear if not provided.
The amount must be a float or integer.
The currency must be the three-letter ISO 4217 code of the amounts, like USD, EUR or JPY. Use the code printed on the receipt; if only a symbol is shown, infer the code from the symbol and the receipt's country or language (a receipt from Japan showing ¥ is JPY, not CNY). Omit the currency if it cannot be determined with confidence. Write amounts with the currency's own decimal places (JPY has none, JOD has three).
If the receipt represents a refund, return, or credit (money returned to the customer rather than charged), the amount and item amounts MUST be negative. Otherwise they are positive.

Please do NOT add any additional information, only valid JSON.
Please return the json in plaintext ONLY, do not ever return it in a code block or any other format.

Choose up to 2 categories from the given list based on the receipt's items and store name. If no categories fit, please return an empty array for the field and do not select any categories. When selecting categories, select only the id, like:
{
	Id: category id
}

Emphasize the relationship between the category and the receipt, and use the description of the category to fine tune the results. Do not return categories that have an empty name or do not exist.
If there are no categories to chose from, then please make categories an empty array.
Likewise, if there are not tags to choose from, then make tags an empty array.

Categories to chose from: @categories

Follow the same process as described for categories for tags.

Tags to chose from: @tags

Receipt text: @ocrText
`
