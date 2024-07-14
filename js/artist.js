const firstalbum = document.getElementById('first-album')
firstalbum.addEventListener('mouseover', ()=> {
    document.getElementById('createSwitch').innerText = "original api gives: " + document.getElementById
    ("albumInfo").getAttribute("data-album");
})
firstalbum.addEventListener('mouseout', ()=> {
    document.getElementById('createSwitch').innerText = document.getElementById
    ("actualAlbumInfo").getAttribute("data-album");
})


document.addEventListener('DOMContentLoaded', () => {
    const members = document.querySelectorAll('.member-item');

    if (!members.length) {
        console.log("No members present");
        return; // Exit if no members are found
    } else {
        console.log(members.length, " members present");
    }

    console.log("Members present__", members.length);
    console.error("Members present__", members.length);

    members.forEach(member => {
        member.addEventListener('mouseover', () => toggleMemberCard(member, true));
        member.addEventListener('mouseleave', () => toggleMemberCard(member, false));
    });

    function toggleMemberCard(member, hover) {
        const memberNameElement = member.querySelector('.center');
        const memberPicElement = member.querySelector('.pic');

        if (!memberNameElement) {
            console.log("No element");
            return;
        }

        if (hover) {
            console.log("Mouse over member");
            memberNameElement.classList.remove('cut');
            memberPicElement.classList.remove('pic--sm');
            memberNameElement.style.whiteSpace = 'normal';
            // member.parentElement.classList.remove('scroll');
        } else {
            console.log("Mouse leave member");
            memberNameElement.classList.add('cut');
            memberPicElement.classList.add('pic--sm');
            memberNameElement.style.whiteSpace = 'nowrap';
            // member.parentElement.classList.add('scroll');
        }
    }
});
