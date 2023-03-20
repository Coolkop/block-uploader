build:
	docker build -t file-storage-app:latest .

init-env:
	cp .env.dist .env

up:
	docker-compose -f docker-compose-test.yml up -d

down:
	docker-compose -f docker-compose-test.yml down


run-tests:
	docker run -it --rm --env-file .env --network file-storage_file-storage-test --name=file-storage-app-test file-storage-app:latest go test ./... -tags=integration

check: build init-env up run-tests down

start-cluster:
	minikube start

apply:
	kubectl apply -f ./kube/app.yaml
	kubectl apply -f ./kube/psql.yaml
	kubectl apply -f ./kube/minio-0.yaml
	kubectl apply -f ./kube/minio-1.yaml
	kubectl apply -f ./kube/minio-2.yaml
	kubectl apply -f ./kube/minio-3.yaml
	kubectl apply -f ./kube/minio-4.yaml
	kubectl apply -f ./kube/minio-5.yaml

delete:
	kubectl delete -f ./kube/app.yaml
	kubectl delete -f ./kube/psql.yaml
	kubectl delete -f ./kube/minio-0.yaml
	kubectl delete -f ./kube/minio-1.yaml
	kubectl delete -f ./kube/minio-2.yaml
	kubectl delete -f ./kube/minio-3.yaml
	kubectl delete -f ./kube/minio-4.yaml
	kubectl delete -f ./kube/minio-5.yaml

tunnel-app:
	minikube service block-storage-service --url

tunnel-postgres:
	minikube service postgresdb-service --url

tunnel-minio-0:
	minikube service minio-0-service --url

tunnel-minio-1:
	minikube service minio-1-service --url

tunnel-minio-2:
	minikube service minio-2-service --url

tunnel-minio-3:
	minikube service minio-3-service --url

tunnel-minio-4:
	minikube service minio-4-service --url

tunnel-minio-5:
	minikube service minio-5-service --url

tunnel-mino-0-admin:
	minikube service minio-0-admin-service --url
