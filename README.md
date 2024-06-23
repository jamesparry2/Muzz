# Muzz

The Muzz API is a miniature API that is designed to handle the management of users and provide a public interface for a consuming application to use to coordinate the discovery of new users for a user to match with and provide the user the opportunity to match with other users.

Below will be subsections denoting the requested functionality, the assumptions that were made that led to technical decisions and how the assumptions impacted the final endpoints. Also, will be some setup steps on how to run the API locally, for manual verification of the endpoints.

To view the swagger documentation once the api is running, if you follow the URL: `http://localhost:5001/swagger/index.html` you will able to view each endpoint, the expected request and responses.

## Feature Requests
- [x] Allow the ability to create users
- [x] Allow users to login
- [x] Allow users to discover other users to swipe on
- [x] Allow users to swipe and potentially match with others
- [x] Allow users to filter discoveries based on age and gender
- [x] Order the results based on distance to the users
- [ ] Bonus: Apply a secondary order filter based on attractiveness

## Assumptions made
- For determining the distance, the assumption that the consumer would provided the lat and long of the users current location. This led to a new location resource which would be a child of the user resource. 
- Regarding the variety of potential data integration options, the API would be built with the design of Ports/Adapters to allow easy switching of underlying technical stores. 

## Running the application locally
For running the user will need to have docker installed. So, once installed I've built a docker compose file that will use the local image that's been built along with a DB container and also handle the networking between the two. There will be two ways to achieve this depending on the users setup:

- 1 Within the root there is a Makefile with a command which handles building the image and tagging the local image, it will then run the docker-compose in detached mode to free up the terminal. `make setup_and_run_local`
- 2 If the user does not have the 'Make' command installed, and do not wish to install it, they can follow the 3 commands to spin up the API: 
    `docker build --tag muzz-api .`
	`docker tag muzz-api:latest muzz-api:local`
	`docker-compose up -d`

## Technical Design & Decisions

### Project Structure

The project follows a structure that supports the implementation of the ports and adapters pattern outlined within hexagonal architecture. Leveraging Go interfaces allows the packages to be loosely coupled by contracts and allows multiple injectable packages. For example, the store contract could allow a MYSQL package or a Dynamo Package to be injected. This allow mocks to be injected allowing each unit to be tested in isolation conforming to positive unit testing practices. 

The project structure is defined into the following:

<b>app/auth</b> - This package is concerned about generating JWTs and any form of authentication within the application.

<b>app/core</b> This package houses the main business login and orchestrates the flows within these files.

<b>app/handlers</b> This package is concerned with managing HTTP inbound and outbound concerns and ensuring the input into core is expected.

<b>app/store</b> This package acts as the data layer between the project and the external data management tool that the project deems suitable to use.

### Database Decision and ORM

For the database its been built against SQL primarily to leverage the relational structure of the data that was being stored as the relationship can be mapped to one is many and vice vera such as a User can have many swipes but a swipe must belong to a single user.

When deciding what DB Tool to use with this, it was decided to sue MYSql to leverage the built in functionality that would support with the distancing query rather than trying to implement my own solution. The only drawback to using the built in functionality is, according to documentation, its not the most performant at scale but for this solution and use case it wouldn't be a problem.

For the ORM in Go, I went with GORM because it allowed for the following features to be used:
- automigration for the table structures and allowed a model first approach
- data parameterization to attempt to limit any potential SQL Injection
- Quicker startup due to its ease of setup

However, there are some drawback, its automigration is a powerful tool but it only supports additive processes so deletion and manual column changes are a challenge.