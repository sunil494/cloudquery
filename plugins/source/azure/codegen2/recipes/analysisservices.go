// Code generated by codegen1; DO NOT EDIT.
package recipes

import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/analysisservices/armanalysisservices"

func init() {
	tables := []Table{
		{
			Service:        "armanalysisservices",
			Name:           "servers",
			Struct:         &armanalysisservices.Server{},
			ResponseStruct: &armanalysisservices.ServersClientListResponse{},
			Client:         &armanalysisservices.ServersClient{},
			ListFunc:       (&armanalysisservices.ServersClient{}).NewListPager,
			NewFunc:        armanalysisservices.NewServersClient,
			URL:            "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/servers",
			Multiplex:      `client.SubscriptionMultiplexRegisteredNamespace(client.NamespaceMicrosoft_AnalysisServices)`,
			ExtraColumns:   DefaultExtraColumns,
		},
	}
	Tables = append(Tables, tables...)
}
