##
# $1: org/repository of which to find stars in. e.g. Jeadie/stargazers
##
# user_repos <username> 
# Returns: git url, new line separated, sorted by last updated
user_repos() {
  curl \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $GITHUB_TOKEN" \
    "https://api.github.com/users/$1/repos?sort=updated&direction=desc&per_page=100" | jq '.[].git_url'
}


# repo_details <org/repo>
# Returns: JSON payload of details about repo.
repo_details() {
    curl \
     -H "Accept: application/vnd.github+json" \
     -H "Authorization: Bearer $GITHUB_TOKEN" \
    "https://api.github.com/repos/$1" 
}

# get_stars_from_page <org/repo> <page> 
get_stars_from_page(){
    curl -H "Accept: application/vnd.github+json" \
       -H "Authorization: Bearer $GITHUB_TOKEN" \
      "https://api.github.com/repos/$1/stargazers?per_page=100&page=$2" | jq '.[].login'
}
export -f get_stars_from_page

STARS=$(repo_details $1 | jq '.stargazers_count')
ITERS=$(echo "($STARS/100)+1" | bc)
  
seq 1 $ITERS | xargs -I% -n1 -P1 bash -c "get_stars_from_page $1 %" > users.json

