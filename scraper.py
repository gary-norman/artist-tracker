import requests
from yq import BeautifulSoup

def get_tadb_artist_id_via_scraping(artist_name):
    search_url = f"https://www.theaudiodb.com/search.php?s={artist_name}"
    response = requests.get(search_url)
    soup = BeautifulSoup(response.text, 'html.parser')

    artist_link = soup.find('a', href=True, text=artist_name)

    if artist_link:
        artist_id = artist_link['href'].split('/')[-1]
        return artist_id
    else:
        return None

artist_name = 'Radiohead'
artist_id = get_tadb_artist_id_via_scraping(artist_name)

if artist_id:
    print(f"The scraped Artist ID for {artist_name} is {artist_id}")
else:
    print(f"Artist {artist_name} not found")
