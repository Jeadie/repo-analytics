##
# Get email of a Github user.
#
#  ./user_email.sh <USERNAME>
#
# Returns: most common email from the user's most recent, non-forked, repository.
##


user_repos(){
  curl -s \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    "https://api.github.com/users/$1/repos?sort=updated&direction=desc&per_page=100" |  jq --raw-output ' sort_by(.updated_at) | . [] | select(.fork==false) | .full_name'
}

author_emails() {
    curl -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/repos/$1/commits" -s  | jq --raw-output ' .[].commit.author.email' 2> /dev/null
}

export -f author_emails

user_repos $1 | xargs -I% -P0 bash -c "author_emails %" | sort | uniq -c | sort -r | head -n 1 | awk '{print $2}'
