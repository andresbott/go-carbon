## list of TODOS for this project



## problems to solve
* add misddleware to write repsonse headers. e.g. cors during dev
* keep me logged in is part of session management
* add logout into session management
* json error middlware to strip tailing \n
* improve request loggger to print the error message
* improve config
  * use a config file with env overrides (?)
* session login store
  * check suspicius logins, based on location / device ?
  * MFA
  * logout
  * wrong login attempts, to greylist or blacklist
* Add UI testing using ROD on the spa package
  
+ change the log middleware to intercept reponse error messages and replace them with generic errors in prod envs
  + this way all the handlers can write specific errors into the response, but the middleware makes sure taht they are not leaked
* totp
* template handler
  * replace user creation 
  * login form
    * create form and json absed login hander
* separate the auth handling from the login from rendenring
* add an SPA handler
* multiple log output stdout /  file / systemd journal
* prometheus metrics
* tracing?


## to be verified
* make sure the user manager uses a salt when storing passwords


##  DOC
### Router
notice that when instantiating subrouters e.g. r.PathPrefix("/basic").Subrouter() then the handler in the subrouter
needs to be like 	"r.Path("").Handler(demoPage)" to handler http://localhost:8080/basic 
iff you add "r.Path("/").Handler(demoPage)" the request needs to be http://localhost:8080/basic/