// Code generated by codegen1; DO NOT EDIT.
package recipes

import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql"

func init() {
	tables := []Table{
		{
			Service:        "armmysql",
			Name:           "servers",
			Struct:         &armmysql.Server{},
			ResponseStruct: &armmysql.ServersClientListResponse{},
			Client:         &armmysql.ServersClient{},
			ListFunc:       (&armmysql.ServersClient{}).NewListPager,
			NewFunc:        armmysql.NewServersClient,
			URL:            "/subscriptions/{subscriptionId}/providers/Microsoft.DBforMySQL/servers",
			Multiplex:      `client.SubscriptionMultiplexRegisteredNamespace(client.NamespaceMicrosoft_DBforMySQL)`,
			ExtraColumns:   DefaultExtraColumns,
		},
	}
	Tables = append(Tables, tables...)
}
