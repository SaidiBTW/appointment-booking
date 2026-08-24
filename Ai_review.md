# AI Review

## Section Assistance

### Section 1

I used AI to debate between the fixed-time slot model and the time blocking system as specified in the README.md. We decided to move on with the time-blocking system due to its flexibility. It could also realize the benefits of the fixed-time slot model using application level control.

### Section 2

During the desing of the application. I used the AI (Antigravity) to write some of the tests after providing several base cases for example when writing the handler tests. I also used the AI to generate the initial .gitgnore and .dockerignore files.

### Section 3

I used the AI to develop inline YAML that was used for the CloudBuild configurations. It needed some corrections in terms of versioning the images that were used during the build e.g., recommending `go 1.20` instead of `go 1.25`

## Sections where AI was wrong

### Error Handling

When designing the failure cases. The AI always directly returned from the function/method using `fmt.Error..` to return errors to the client. This is not correct as it does not provide a maintainable way for managing errors. Instead we should use custom error handling to strictly define our errors, making them more maintanable and easily consumed by downstream services.

### DTO Design (time Based Values)

During design of the dto's most AI discussion assumed that time was a `String` which meant that validation e.g., for start and end times, by packages such as `govalidator` had to fall to deeper entites such as handler instead of being handled upfront.

### Application Design (Seperation of Concerns)

The AI almost always assumed that validation of values provided could be carried out at any layer e.g., the repository layer and return client errors from there. This is an anti-pattern and should only be used for cases that require db_fetches etc. i.e., Validation of invariants of the system e.g., start_time and end_time being 30 minutes apart should happen as early as possible,e.g., in the handler class that is receiving the values.

## Decisions Made Without AI

- The selection of the database based on the schema that was extracted from the database. It seemed to have a relational aspect to it hence the selection of PostgreSQL database.
- The selection of the deployment infrastructure was also decided from the initial exerpt provided. As they needed a scalable, cost effective solution that is suitable for starting projects.
- The structure of the project. This was based on known best practices that help with isolated testing of components e.g., repositories, services, handlers etc
