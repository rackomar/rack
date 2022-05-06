# IRCCloud Keep Alive Utility

A simple Heroku application to keep your IRCCloud connection active!

This script is based on the version made by [osm](https://github.com/osm/icka/). Unlike that version, this fork can be completely deployed to and configured for Heroku from the browser. Heroku config vars are used to store your IRCCloud credentials.

This code uses IRCCloud's [publicly-documented RPC API] (https://github.com/irccloud/irccloud-tools/wiki).
While this script prevents IRCCloud disconnecting from IRC servers after 120 minutes of inactivity, it does not provide access to any other of the numerous features that are available in the [Pro version of IRCCloud](https://www.irccloud.com/pricing). Please support them if you can!

**IMPORTANT: This utility should not be used for critically-important connections. IRCCloud may still occasionally disconnect users, causing messages to be missed, even if this tool is used.**

Requirements
============
* A free IRCCloud account
* A free Heroku account
* A free cron-job.org account
  
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/tech234a/irccloud/)

Detailed Setup Instructions
===========================
1. Click the "Deploy to Heroku" button above.
2. Sign up for or log into Heroku.
3. Choose an app name on Heroku. It doesn't really matter what it is, but I wouldn't recommend making it a perfect name as you won't be entering it once you automate the application's execution. Enter this name into the "App name" field and the `heroku-app-name` config var.
4. Create a [new authorization token](https://dashboard.heroku.com/account/applications/authorizations/new). Leave the expiration field blank. Once again, it doesn't matter what you description you provide for it. Fill it into the `heroku-key` config var.
5. Fill out the `IRCCLOUD_PASSWORD` and `IRCCLOUD_EMAIL` config vars. This is required to authenticate with IRCCloud and send a signal indicating an active client. While I understand some people might be reluctant to enter their credentials, [Heroku does actually recommend using config vars for storing data encryption keys](https://devcenter.heroku.com/articles/getting-started-with-python#define-config-vars).
6. Go to [cron-job.org](https://cron-job.org/) and [sign up](https://console.cron-job.org/signup) for or [log into](https://console.cron-job.org/login) an account.
7. Go to the Cronjobs tab. [Create a new cron job](https://console.cron-job.org/jobs/create).
8. The title of the job is optional and does not matter.
9. For the URL, enter `https://api.heroku.com/apps/[Heroku app name]/formation`, where `[Heroku app name]` is the same app name you entered earlier.
10. For execution schedule, select "Custom". If not already selected, select all "Days of month", "Days of week", "Months" and "Hours". (Tip: Select the first entry in the list and, while holding the Shift key, select the last entry in the list to select the whole list.
11. Select any one specific minute. This configures the job to run once per hour (IRCCloud disconnects users after 2 hours of inactivity).
12. Optionally enable email notifications for when the cronjob fails.
13. Ensure the "Enable cronjob" is switch is turned on.
14. Under "Advanced", select `PATCH` for the request method.
15. Add the following request header keys and values:
  - `Content-Type`: `application/json`
  - `Accept`: `application/vnd.heroku+json; version=3`
  - `Authorization`: `Bearer [Heroku authorization token]` where `[Heroku authorization token]` is the authorization token you created earlier.
16. Enter the following for "Request body": ```{ "updates": [ { "quantity": 1, "type": "worker" } ] }```
17. Click "Create" to create the cronjob, and you're done!
