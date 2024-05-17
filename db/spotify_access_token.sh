# this is required to get an access token that expires after 1 hour while the spotify app is still in developer mode
# it needs to be run from the command line (chmod +x spotify_access_token.sh if it is not executable from your system)
curl -X POST "https://accounts.spotify.com/api/token" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=client_credentials&client_id=bbbba5fd8a62451fb7192f71f596ab88&client_secret=ca85187ed2ec4f7bb601039343d23777"

