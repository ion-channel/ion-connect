# Usage

```
./ion-connect project create-project-spdx ./ext/examples/ex.spdx --spdx-version 2.1 --team-id <team-id> --ruleset-id <ruleset-id> --package-name hello-bin
```

### Example response
```
{"id":"1234-5678-910","team_id":"1234-5678-910","ruleset_id":"1234-5678-910","name":"hello-bin","type":"git","source":"git@github.com:swinslow/spdx-examples.git#example2/content/build","branch":"master","description":"","active":true,"chat_channel":"","created_at":"2021-01-27T21:43:26.265708Z","updated_at":"2021-01-27T21:43:26.265708Z","deploy_key":"","should_monitor":true,"monitor_frequency":"","poc_name":"","poc_email":"","username":"","password":"","key_fingerprint":"","private":false,"aliases":[],"tags":[],"ruleset_history":null}
```