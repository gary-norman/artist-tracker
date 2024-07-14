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

    members.forEach(member => {
        member.addEventListener('mouseover', () => toggleMemberCard(member, true));
        member.addEventListener('mouseleave', () => toggleMemberCard(member, false));
    });

    function toggleMemberCard(member, hover) {
        const memberNameElement = member.querySelector('.center');
        const memberPicElement = member.querySelector('.pic');
        const parent = member.parentElement;
        console.log("parent is:", parent)

        if (!memberNameElement) {
            console.log("No element");
            return;
        }

        if (hover) {
            // // Add placeholder to replace absolute member-item missing from the flow
            // const placeholder = document.createElement('div');
            // placeholder.classList.add('placeholder');
            // placeholder.style.width = `${member.offsetWidth}px`;
            // placeholder.style.height = `${member.offsetHeight}px`;
            // member.parentElement.insertBefore(placeholder, member);

            // Adjust the member item
            console.log("Mouse over member");
            memberNameElement.classList.remove('cut');
            memberPicElement.classList.remove('pic--sm');
            memberNameElement.style.whiteSpace = 'normal';
            // member.parentElement.classList.remove('scroll');

            console.log("member.parentElement is:", member.parentElement)
        } else {
            // Remove the placeholder
            // const placeholder = document.querySelector('.placeholder');
            // if (placeholder) {
            //     placeholder.parentElement.removeChild(placeholder);
            // }

            // Reset the member item
            console.log("Mouse leave member");
            memberNameElement.classList.add('cut');
            memberPicElement.classList.add('pic--sm');
            memberNameElement.style.whiteSpace = 'nowrap';
            // member.parentElement.classList.add('scroll');
        }
    }
});
