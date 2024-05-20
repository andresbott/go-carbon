## list of TODOS for this project



## problems to solve

* improve config
  * use a config file with env overrides (?)
* session login store
  * check suspicius logins, based on location / device ?
  * MFA
  * logout
  * wrong login attempts, to greylist or blacklist
  
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
* 