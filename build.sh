# Define your go file
go_file="main.go" # change to the name of your Go file

# Build the Go file
go build -o bin/api $go_file

# Run the built file
./bin/api