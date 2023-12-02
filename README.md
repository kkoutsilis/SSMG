
# SSMG

SSMG (Secret Santa Match Generator) is a simple CLI tool written in Go. It reads a JSON file containing Secret Santa participants, generates the matches, and sends an email to each participant with their assigned match.


### Usage

To use SSMG, you need to provide a JSON file containing the participant data. The JSON file should have the following structure:
```json
[
    {
      "name": "participant1",
      "email": "participant@test.com",
      "wishlist": ["foo","bar"] // Optional
    },
    {
      "name": "participant2",
      "email": "participant2@test.com"
    },
    ...
]

```
