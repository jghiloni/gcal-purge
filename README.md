# gcal-purge

This simple tool will delete all google calendar events within a time range.
This is currently not designed for public use, and will not be supported
outside of my Pivotal colleagues.

```
Options:
    --calendar-id, -c  string  The ID of the Google Calendar you wish to purge (default: primary)
    --dry-run          bool    Will print what will happen to stderr, but won't actually do it
    --end-date         string  The date when the tool should end purging at 11:59 PM local time in YYYY-MM-DD format (default: today)
    --help, -h         bool    Print this message and quit
    --start-date       string  The date when the tool should start purging at 12:01 AM local time in YYYY-MM-DD format (default: 1970-01-01)
    --verbose, -v      bool    Enable more verbose logging
```

Steps to use:

1. Download the latest compiled binary from the Releases page or build from
   source.
1. Enable the Google Calendar API for your Google Account (Some accounts, like
   corporate managed accounts, will not allow you to do this. You can enable it
   on another personal Google Account and log into any Google Account). Follow
   step 1 [here](https://developers.google.com/calendar/quickstart/go#step_1_turn_on_the).
   Make sure to configure your OAuth client as a desktop app. Save the downloaded
   `credentials.json` file to your home directory.

When you run the app, it will open a browser and let you log in to the Google
Account of your choice. When you copy the token and paste into the CLI, it will
save a token in `$HOME/.gcal-purge`. If you have authorization issues, try
deleting the file and trying again.

If the run quits unexpectedly in the middle of a run, look for a Rate Limit Exceeded
error. If that is the case, wait a few seconds and run again, and it will pick up
where it left off.

Any questions? Ping me on Slack.
