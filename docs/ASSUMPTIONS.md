# Assumptions

1. Dates are read and stored in UTC format.
2. we test public functions/methods only. If we feel we need to test a private function then it means that the function needs its own package or library ~ Kent Beck.
3. All player fields are mandatory. ID, first name, last name, email, nickname, password and country.
4. Before creating a Player We need to validate that repository does not already contain any player with the given email or nickname values.
5. Deleting players can be done physically in the database and no logical state handling is required.
6. In the player search function, if the client does not provide any search criteria, the service will return an empty result.
7. In the player search function, if the filter criteria does not match any player data, the service will return an empty result.
8. We need a RDBMS (Relational Database Management system) repository to save player data. I used Postgres.
9. I didn't use any GORM framework, it could be helpful to speed up development, but just wanted to keep things simple.
10. I tried to follow this thought "easy to understand rather than easy to do".
11. The notifier infrastructure is still under discussion so the application has the logic to notify any eventbus asynchronously but no actual eventbus is configured or called from the service for this release.
12. Extensibility, maintainability, flexible coupling and high cohesion are important for this project.
13. You have go 1.23 installed.