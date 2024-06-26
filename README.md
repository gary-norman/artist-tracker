
## Audit Info
### Functional

* Has the requirement for the allowed packages been respected? (Reminder for this project: only standard packages)
* Is the data from the artists being used?
* Is the data from the locations being used?
* Is the data from the dates being used?
* Is data from the relations being used?
### Try to see the "members" for the artist/band "Queen"


    "Freddie Mercury",
    "Brian May",
    "John Daecon",
    "Roger Meddows-Taylor",
    "Mike Grose",
    "Barry Mitchell",
    "Doug Fogie"

* Does it present the right "member", as above?
### Try to see the "firstAlbum" for the artist/band "Gorillaz"

    "26-03-2001"

* Does it present the right date for the "firstAlbum", as above?
### Try to see the "locations" for the artist/band "Travis Scott"


    "santiago-chile"
    "sao_paulo-brasil"
    "los_angeles-usa"
    "houston-usa"
    "atlanta-usa"
    "new_orleans-usa"
    "philadelphia-usa"
    "london-uk"
    "frauenfeld-switzerland"
    "turku-finland"

* Does it present the right "locations" as above?
### Try to see the ""members"" for the artist/band "Foo Fighters".


    "Dave Grohl"
    "Nate Mendel"
    "Taylor Hawkins"
    "Chris Shiflett"
    "Pat Smear"
    "Rami Jaffee"

* Does it present the right members as above?
#### Try to trigger an event/action using some kind of action (ex: Clicking the mouse over a certain element, pressing a key on the keyboard, resizing or closing the browser window, a form being submitted, an error occurring, etc).
* Does the event/action responds as expected?
* Did the server behaved as expected?(did not crashed)
* Does the server use the right HTTP method?
* Did the site run without crashing at any time?
* Are all the pages working? (Absence of 404 page?)
* Does the project handle HTTP status 500 - Internal Server Errors?
* Is the communication between server and client well established?
* Does the server present all the needed handlers and patterns for the http requests?
* As an auditor, is this project up to every standard? If not, why are you failing the project?(Empty Work, Incomplete Work, Invalid compilation, Cheating, Crashing, Leaks)

### General

* Does the event system run as asynchronous? (usage of go routines and channels)
* Is the site hosted/deployed? Can you access the website through a DNS (Domain Name System)?

### Basic

* Does the project runs quickly and effectively? (Favoring recursive, no unnecessary data requests, etc)
* Does the code obey the good practices?
* Is there a test file for this code?

##

## Semantic Commit Messages

### Intro

See how a minor change to your commit message style can make you a better programmer.

Format: `<type>(<scope>): <subject>`

`<scope>` is optional

### Example

    feat: add hat wobble
    ^--^  ^------------^
    |     |
    |     +-> Summary in present tense.
    |
    +-------> Type: chore, docs, feat, fix, refactor, style, or test.`

### More Examples:

`feat`: (new feature for the user, not a new feature for build script)

`fix`: (bug fix for the user, not a fix to a build script)

`docs`: (changes to the documentation)

`style`: (formatting, missing semi colons, etc; no production code change)

`refactor`: (refactoring production code, eg. renaming a variable)

`test`: (adding missing tests, refactoring tests; no production code change)

`chore`: (updating grunt tasks etc; no production code change)

`maint`: (updating current issue you're working on)

### References:

* https://www.conventionalcommits.org/
* https://seesparkbox.com/foundry/semantic_commit_messages
* http://karma-runner.github.io/1.0/dev/git-commit-msg.html
