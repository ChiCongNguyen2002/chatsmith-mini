package catalog

import "chatsmith-mini/constants"

var providerOrder = []constants.Provider{
	constants.ProviderOpenAI,
	constants.ProviderAnthropic,
}

var providerNames = map[constants.Provider]string{
	constants.ProviderOpenAI:    "GPT",
	constants.ProviderAnthropic: "Claude",
}
