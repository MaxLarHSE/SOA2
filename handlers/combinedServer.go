package handlers

import (
	"SOA2/handlers/auth"
	"SOA2/handlers/products"
)

type CombinedServer struct {
	*products.ProductServer
	*auth.AuthServer
}
