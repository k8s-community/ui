package client_test

import (
	"fmt"
	"github.com/k8s-community/user-manager/client"
)

func ExampleSyncUser() {
	cl, err := client.NewClient(nil, "https://services.k8s.community/user-manager")

	if err != nil {
		fmt.Printf("Error during client creation: %s", err)
	}

	request := client.NewUser("rumyantseva")

	respCode, err := cl.User.Sync(request)
	if err != nil {
		fmt.Printf("Server error: %s", err)
	}

	fmt.Println(respCode)
}
