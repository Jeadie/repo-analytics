##
# Get email of a Github user.
#
#  ./user_email.sh <USERNAME>
#
# Returns: most common email from the user's most recent, non-forked, repository.
##

go run . user-emails $1