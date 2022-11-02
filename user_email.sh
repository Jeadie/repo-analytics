"""
Get email of a Github user.

  ./user_email.sh <USERNAME>

Returns: most common email from the user's most recent, non-forked, repository.
"""

user_repos(){
  curl \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    "https://api.github.com/users/$1/repos?sort=updated&direction=desc&per_page=100" |  jq --raw-output ' sort_by(.updated_at) | . [] | select(.fork==false) | .html_url'
}

REPO=$(user_repos $1 | head -n 1)
DIR="/tmp/$(echo $1 | cut -d":" -f2 | sed -r 's/\//_/g')"
mkdir $DIR

git clone --quiet --bare -- $REPO $DIR

cd $DIR
EMAIL=$(git --no-pager log -s --format="%ae" | sort  | uniq -c | sort -r | awk '{ print $NF }' | head -n 1)
echo $EMAIL
cd - > /dev/null
