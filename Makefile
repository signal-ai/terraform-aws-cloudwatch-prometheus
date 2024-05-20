docs:
	terraform-docs .

build:
	cd lambda_code && rm payload.zip || true
	cd lambda_code && env GOOS=linux GOARCH=amd64 go build -o bootstrap
	cd lambda_code && zip -r payload.zip bootstrap && rm bootstrap
