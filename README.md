# Valsea Technical Test

This repository contains the solution for the technical test provided by Valsea for a Golang developer position. The project involves building a RESTful API for managing bank accounts and transactions.

## Project Structure

The project is divided into the following components:


### `bank-demo-app`
This folder contains the implementation of the RESTful API server. It has the following sub-packages:

- **restServer**: Implements the REST server, initializing the routes and handling the server's methods.
- **mongodb**: Provides a wrapper around the MongoDB client, managing connections, configuration, and collections.
- **inputParams**: A utility package for reading input parameters during the application's startup.
- **bank**: Core bank logic:
  - **memoryBank**: Implements an in-memory database to store accounts and transactions.
  - **dbBank**: Similar to `memoryBank` but stores data in MongoDB for persistence.

### `bank-test-client`
This folder contains a simple Go CLI application designed to test the functionality of the RESTful API. It’s a straightforward tool to verify server functionality but is not highly refactored.

- **tester**: Contains the code responsible for interacting with the API and performing tests.

### `scripts`
This folder contains two scripts to manage MongoDB:

1. **setup-mongodb.sh**: Downloads the MongoDB Docker image and initializes a database named `BankStore` with two collections: `accounts` and `transactions`. **Docker should be installed in your computer and you have to be loged into your docker account to run this script**
2. **empty_db_collections.sh**: Clears all data from the MongoDB collections, providing a way to reset the database before testing.

## How to Test the Application

To test the application, there are three scripts available: `run_server.sh` and `run_test_client.sh`.

1. **run_server.sh**: Builds and runs the `bank-demo-app`, which listens on `localhost:8080` by default. 
   - To test the MongoDB-backed bank instead of the in-memory version, set the `IN_MEMORY` variable to `false`.

2. **run_test_client.sh**: Builds and runs the `bank-test-client` application. Follow the instructions in the terminal to test the API's functionality.

   - *Important*: All errors covered in the application are defined in `bank-demo-app/internal/bank/errors.go`. You can trigger these errors by performing invalid actions in the test client.

2. **run_unit_tests.sh**: Some Unit tests as a proof of concepts have been implemented, not all the code is covered.

Additionally, a Postman collection file (`ValseaBankTestRequests.postman_collection.json`) is available, containing all the API requests described in the test. You can import it into Postman to test the API directly.

## Conclusion

I hope you find this application useful and consider it for the open position! Feel free to reach out with any questions or feedback.
