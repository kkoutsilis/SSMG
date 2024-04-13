# SSMG 

![CI](https://github.com/kkoutsilis/SSMG/actions/workflows/ci.yml/badge.svg)

SSMG (Secret Santa Match Generator) is a simple CLI tool written in Go. It reads a JSON file containing Secret Santa participants, generates the matches, and sends an email to each participant with their assigned match.

## Motivation

Picture this: my friends and I, spread across the globe, were scheming up a festive get-together complete with a sprinkle of Secret Santa enchantment. Now, the catch? Our scattered geography made the classic match-making puzzle a bit tricky. Sure, there are online platforms for such scenarios, but most demanded user accounts—a hoop I'd rather not jump through. And then there's the whole sharing-emails-with-apps deal, not my cup of cocoa. So, cue my weekend project: a sleek CLI tool tailored to weave our Secret Santa magic, no strings attached. Here's to coding our way to holiday cheer! 🌐🎄✨

**Disclaimer**  
I am not so familiar with GO so the code probably needs cleaning and improvements, but hey, it works! 

### Usage


To use SSMG, you need to have go installed, clone the repo locally, provide a JSON file containing the participant data and some environmental variables to be used to send the emails.

#### Environmental Variables 
- `EMAIL_HOST`: The email host, for example `smtp.gmail.com` for gmail.

- `EMAIL_PORT`: The email port, tipically `583` for smtp.

- `EMAIL_FROM`: The email sender, the one pulling the strings.

- `EMAIL_USER`: The email account's username, usually your email itself.

- `EMAIL_PASSWORD`: The passwowrd for the given email account.

**Note**: Gmail does not allow plain password login and instead uses app passwords, which is more secure as well. You can learn how to create and use app passwords [here](https://support.google.com/accounts/answer/185833). 

#### Input File Format

```json
[
    {
      "name": "participant1",
      "email": "participant@test.com",
    },
    {
      "name": "participant2",
      "email": "participant2@test.com"
    }
]

```

#### Run
```bash
go run main.go --file my_data.json   
```
or by using the binary
```bash
 go build --ldflags "-s -w" ssmg
./ssmg --file my_data.json
```
**Note**: If `file` flag is not provided, the program will default to data.json.


### Testing Emails

There is a docker compose file that spins up a [MailHog](https://github.com/mailhog/MailHog) service that can be used for testing the emails.
Simply `docker compose up` to use the testing (MailHoghttps://github.com/mailhog/MailHog) service and set the environmentat variables to 
```bash
export EMAIL_HOST=localhost
export EMAIL_PORT=1025
export EMAIL_FROM=test@example.org
export EMAIL_USER=””
export EMAIL_PASSWORD=””
```
You can access the MailHog UI at `localhost:8025`