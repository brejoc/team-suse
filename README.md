# team-suse
`team-suse` is a very cheesy applications that lists all PR where either you or team-suse was requested to review in the [saltstack/salt](https://github.com/saltstack/salt). So the target group is very small! ;)

## Installation

If you are running a 64 bit Linux, you can grab the binary from the [releases](https://github.com/brejoc/team-suse/releases) and just copy it to `$PATH`. Otherwise you'd have to build it from source. It should be enough to do  a `go build` inside of *team-suse* folder. You should get a binary named `team-suse`, that can be copied somewhere to `$PATH`. E.g. `~/bin` or `/usr/bin`.

You also need a ['Personal access token'](https://github.com/settings/tokens) from GitHub with at least read permissions for the repositories. `team-suse` expects the token in the environment variable `GITHUB_TOKEN`. You can export this in bash like this: `export GITHUB_TOKEN <insert_your_token_here>`.