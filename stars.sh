##
# Get all stargazers from a repository.
#
#  ./stars.sh <org/repository>
#
# Returns: Username of all stargazers on a repository. 
##

# repo_details <org/repo>
# Returns: JSON payload of details about repo.
repo_details() {
    curl -s \
     -H "Accept: application/vnd.github+json" \
     -H "Authorization: Bearer $GITHUB_TOKEN" \
    "https://api.github.com/repos/$1" 
}

# get_stars_from_page <org/repo> <page> 
get_stars_from_page(){
    curl -s -H "Accept: application/vnd.github+json" \
       -H "Authorization: Bearer $GITHUB_TOKEN" \
      "https://api.github.com/repos/$1/stargazers?per_page=100&page=$2" | jq --raw-output '.[].login'
}
export -f get_stars_from_page

STARS=$(repo_details $1 | jq '.stargazers_count')
ITERS=$(echo "($STARS/100)+1" | bc)
  
seq 1 $ITERS | xargs -I% -P0 bash -c "get_stars_from_page $1 %"

