[![Code Climate](https://img.shields.io/codeclimate/maintainability/envygeeks/tweedy.svg?style=for-the-badge)](https://codeclimate.com/github/envygeeks/tweedy/maintainability)
[![GitHub release](https://img.shields.io/github/release/envygeeks/tweedy.svg?style=for-the-badge)](http://github.com/envygeeks/tweedy/releases/latest)

# Tweedy

Twitter Delete Yourself is a Go application for deleting your Tweets, while preserving part of your history for n amount of time.  It also supports consuming a Twitter JSON backup once you convert the `tweet.js` to `tweet.json` by opening it up, removing the variable, and leaving only the array.

## Usage

### 1. API Key
#### Get

Visit https://developer.twitter.com, and then sign in, and go to `Username` > Apps | Create an App, and fill in the form in as much detail as possible. Be prepared to wait a day or two, for me it took an hour (as I already had an app,) for some people it might take longer, so don't get upset if you don't get approved instantly.

#### Set

```
export TWITTER_TOKEN=""
export TWITTER_TOKEN_SECRET=""
export TWITTER_SECRET_KEY=""
export TWITTER_KEY=""
```

### 3. Download
#### Binary

* Visit https://github.com/envygeeks/tweedy/releases/latest
* Download a binary
  * Make sure it matches your operating system
  * Make sure it matches your arch

```
mv ~/Downloads/<Binary> /usr/local/bin/tweedy
```

### 4. Use

```
tweedy --user=<Handle> \
  --keep=<Days>
```
