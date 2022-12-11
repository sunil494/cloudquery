// Code generated by codegen1; DO NOT EDIT.
package recipes

import "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/reservations/armreservations"

func init() {
	tables := []Table{
		{
			Service:        "armreservations",
			Name:           "reservation",
			Struct:         &armreservations.ReservationResponse{},
			ResponseStruct: &armreservations.ReservationClientListAllResponse{},
			Client:         &armreservations.ReservationClient{},
			ListFunc:       (&armreservations.ReservationClient{}).NewListAllPager,
			NewFunc:        armreservations.NewReservationClient,
			URL:            "/providers/Microsoft.Capacity/reservations",
			Multiplex:      `client.SubscriptionMultiplexRegisteredNamespace(client.NamespaceMicrosoft_Capacity)`,
			ExtraColumns:   DefaultExtraColumns,
		},
		{
			Service:        "armreservations",
			Name:           "reservation_order",
			Struct:         &armreservations.ReservationOrderResponse{},
			ResponseStruct: &armreservations.ReservationOrderClientListResponse{},
			Client:         &armreservations.ReservationOrderClient{},
			ListFunc:       (&armreservations.ReservationOrderClient{}).NewListPager,
			NewFunc:        armreservations.NewReservationOrderClient,
			URL:            "/providers/Microsoft.Capacity/reservationOrders",
			Multiplex:      `client.SubscriptionMultiplexRegisteredNamespace(client.NamespaceMicrosoft_Capacity)`,
			ExtraColumns:   DefaultExtraColumns,
		},
	}
	Tables = append(Tables, tables...)
}
