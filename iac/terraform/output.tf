output "api_url" {
  value       = "${aws_apigatewayv2_api.go_api.api_endpoint}/"
  description = "URL of the API"
}