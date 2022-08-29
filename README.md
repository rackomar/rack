# IRCCloud Keep Alive Utility

A simple GitHub Actions application to keep your IRCCloud connection active!

This script is based on the version made by [osm](https://github.com/osm/icka/). Unlike that version, this fork can be completely deployed to and configured for GitHub Actions from the browser. GitHub Actions Secrets are used to store your IRCCloud credentials.

This code uses IRCCloud's [publicly-documented RPC API](https://github.com/irccloud/irccloud-tools/wiki).
While this script generally prevents IRCCloud disconnecting from IRC servers after 120 minutes of inactivity, it does not provide access to any other of the numerous features that are available in the [Pro version of IRCCloud](https://www.irccloud.com/pricing). Please support them if you can!

**IMPORTANT: This utility should not be used for critically-important connections. IRCCloud may still occasionally disconnect users (such as during a GitHub Actions outage or when GitHub Actions queues are long), causing messages to be missed, even if this tool is used.**

Requirements
============
* A free IRCCloud account
* A free GitHub account

Detailed Setup Instructions
===========================
1. Click the "Fork" button above.
2. If necessary, sign up for or log into GitHub.
3. Go to the Settings tab of your repo, then navigate to Secrets > Actions.
4. Click "New Repository Secret", enter `IRCCLOUD_EMAIL` as the name, and your IRCCloud account email address as the value. Then click "Add Secret". Click "New Repository Secret" again, and this time enter `IRCCLOUD_PASSWORD` as the name, and your IRCCloud account password as the value. Then click "Add Secret". These values are required to authenticate with IRCCloud and send a signal indicating an active client. While I understand some people might be reluctant to enter their credentials, [Secrets are encrypted and GitHub does actually recommend using them for storing authentication keys](https://docs.github.com/en/rest/actions/secrets).
5. Go to the Actions tab of your repo, then click "I understand my workflows, go ahead and enable them". You're all set!

To update to later versions of this tool simply use the "Sync fork" feature available in the Code tab of your repo.

Note: if your repo is public, logs from runs of this tool will be public. These logs will not include your IRCCloud email and passwords, but if you want to make a private fork of this repo, use [GitHub's import tool](https://github.com/new/import) with the URL `https://github.com/tech234a/irccloud` to create the repo. You will have to manually enable Actions by going to the Settings tab of your repo, then going to Actions>General and selecting the option to enable Actions. The "Sync fork" feature won't be available for repos created this way, so the easiest way to update your repo is to delete it and recreate it following the steps avove.


Migrating From Heroku
=====================
Simply delete the app from your Heroku account and delete your scheduled job on cron-job.org, then follow the steps above. (You should also be able to just leave everything as is and your Heroku app and cron-job.org job should automatically be disabled at the end of November/early December.) 
