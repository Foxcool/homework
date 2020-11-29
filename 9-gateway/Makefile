run:
	helm install prom stable/prometheus-operator -f prometheus.yaml --atomic
	helm install nginx stable/nginx-ingress -f nginx-ingress.yaml
	helm repo add datawire https://getambassador.io
	helm install aes datawire/ambassador -f ambassador_values.yaml
	helm install profile profile/profile-chart/
	helm install identity identity/identity-chart/
	kubectl apply -f ambassador/