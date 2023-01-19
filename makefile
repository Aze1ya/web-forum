run-local:
		go run database-schema/main.go
		go run cmd/app/main.go

run-docker:
	    docker build --tag forum .
    
	    docker run -p 8181:8181 -it forum
     
