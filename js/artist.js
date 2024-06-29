const firstalbum = document.getElementById('first-album')
firstalbum.addEventListener('mouseover', ()=> {
    document.getElementById('createSwitch').innerText = "original api gives: " + document.getElementById
    ("albumInfo").getAttribute("data-album");
})
firstalbum.addEventListener('mouseout', ()=> {
    document.getElementById('createSwitch').innerText = document.getElementById
    ("actualAlbumInfo").getAttribute("data-album");
})

