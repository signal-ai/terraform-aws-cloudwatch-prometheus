docs:
	terraform-docs .

build:
	cd lambda_code && env GOOS=linux GOARCH=amd64 go build -o main
	cd lambda_code && zip -r payload.zip main && rm main
