package catalog

import "chatsmith-mini/constants"

const (
	tableModels      = "models"
	whereEnabled     = "enabled = ?"
	whereEnabledByID = "enabled = ? AND id = ?"
	orderCatalog     = "sort_order, id"
)

var providerOrder = []constants.Provider{
	constants.ProviderOpenAI,
	constants.ProviderAnthropic,
}

var providerNames = map[constants.Provider]string{
	constants.ProviderOpenAI:    "GPT",
	constants.ProviderAnthropic: "Claude",
}
