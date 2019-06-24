output "express_route_service_key" {
  description = "Expressroute Service Key"
  value       = [azurerm_express_route_circuit.main.service_key]
}

