# Tweedy

Twitter Delete Yourself is a Go application for deleting your Tweets, while preserving part of your history for n amount of time.  It also supports consuming a Twitter JSON backup once you convert the `tweet.js` to `tweet.json` by opening it up, removing the variable, and leaving only the array.


## Usage

```
Tweedy deletes your Tweets, Retweets, and Likes

Usage:
  tweedy [flags]

Flags:
  -f, --file string   Tweets JSON file
  -r, --from int      Delete Tweets from ID
  -h, --help          help for tweedy
  -k, --keep int      Days to keep (default 3)
  -0, --silent        Silent output (Warnings only) (default true)
  -u, --user string   The user
  -1, --verbose       Verbose output
```

### Running

You can either visit https://github.com/envygeeks/tweedy/releases/latest, and download the latest binary build (on macOS they are signed,) or you can install `go` and do `go run . --help` and go from there.  You will also need to signup for the Twitter API through https://developer.twitter.com, and then get generate a key for your account.  Once you have done that you can set `TWITTER_ACCESS_TOKEN`, `TWITTER_ACCESS_TOKEN_SECRET`, `TWITTER_API_KEY`, and `TWITTER_API_SECRET`, and just run the binary, and you'll be good to go.  If Twitter asks for clarification during the app creation process let them know that are "trying to delete your Tweets using https://github.com/envygeeks/tweedy so you need an API key", most of the time they will approve it.
