output "my_accounts" {
  value = data.pureport_accounts.main
}

output "my_networks" {
  value = data.pureport_networks.name_filter
}
