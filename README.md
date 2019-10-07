# Brightwheel email coding exercise
Very simple proxy email service to use either Mailgun or Sendgrid to send simple emails to a single recipient.

## How To Run
- Make sure you have an up-to-date go environment
- `go get github.com/austinbisharat/brightwheel/`
- `cd $GOPATH/src/github.com/austinbisharat/brightwheel/`
- optionally: `dep ensure`
- run tests: `go test ./...`
- Setup Sendgrid api key: `export SENDGRID_API_KEY=<your sendgrid api key>`
- Setup Mailgun api key and domain `export MAILGUN_API_KEY=<your mailgun api key>` and `export MAILGUN_DOMAIN=<a domain registered with your mailgun account>"`
- Build with `go build`
- Run the server with `./brightwheel --email_service=sendgrid` or `./brightwheel --email_service=mailgun`

## Design & Tradeoffs
This email service is a very simple go server. I chose to use go because I am familiar with it, it's fast, and it has good support for networking and concurrency primatives. Broadly, there are three layers to the server:

- The external api layer: this is mostly contained in main.go and defines the handler for accepting POST requests. It is responsible for parsing incoming requests and passing them through the next two layers
- The validation layer: this is mostly contained in validation.go and defines the application-level validation of incoming requests
- The email service layer: this is contained in emailservice/ and defines the abstraction above both sendgrid and mailgun, as well as the implementations for each of those providers.

I chose to use an external package, bluemonday, to do HTML sanitization of email bodys. I figured that using a well-tested package to do the sanitization not only would make it more robust, but also much more configurable. If we decide to allow certain types of of HTML tags, then we can easily use bluemonday's policy builder to construct whatever we want.

There were a number of directions that I would have liked to take the project, but did not have time for:
- Supporting some retry logic for failed requests to the email service
- Adding a layer above the email-service layer to allow automatically fallback from one service provider to another. We can wrap the EmailServices in an interface that knows about all implementations, and chooses the last one that was known to work
- Adding some queueing for batching, rate-limiting, etc. If we want to unlock the capability of batching email send requests, we need to add a queue to our implementation. This would also allowing us to control the rate at which we make requests to the email service to avoid being throttled. There are two (sane) options that I see for doing this:
  - Make an in-memory queue but still block the request from returning to the client until we've successfully used the email service to send the email. This can be accomplished using a go channel as the queue and make a pool of goroutines to consume from the queue. We would need to include a waitgroup in each element of the queue to prevent early returns to the client.
  - Make a persistent queue by using some external service like AWS Kineses. Or even just a simple sql database can easily support queue-like operations when we don't care that much about ordering. These might also unlock scheduling as a feature. This is my preferred direction, as I think adding a blocking waitgroup/callback to an in-memory queue is both awkward and doesn't take advantage of one of the main points of a queue: having some buffer space to smooth out spikey traffic
 - Lastly, if I were building this for production, I would want a much more thorough build/test/release process. This is pretty specific to the org and development practices of the org, so I don't know exactly where I would take this, but it would be nice to use kubernetes or something to provide a very stable build/release process
