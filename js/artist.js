function debounce(func, wait) {
    let timeout;
    return function(...args) {
        const context = this;
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(context, args), wait);
    };
}

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
    const parent = members[0].parentElement;
    if (!parent) {
        console.log("parent not present");

    }
    let isHovering = false;

    if (!members.length) {
        console.log("No members present");
        return; // Exit if no members are found
    } else {
        console.log(members.length, " members present");
    }


    window.addEventListener('resize',  debounce( () => {

        if (window.innerWidth < 850) {
            parent.classList.add('scroll');
        } else {
            parent.classList.remove('scroll');
        }
    }, 300));


    members.forEach(member => {
        if (!isHovering) {
            member.addEventListener('mouseover',  debounce( () => {
                toggleMemberCard(member, true);
            }, 300));

        } else {
            member.addEventListener('mouseover',  debounce( () => {
                toggleMemberCard(member, false);
            }, 300));
        };
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
            isHovering = true;

            // if (member === member.parentNode.firstElementChild) {
            //     console.log("entering hover first child")
            //     const placeholder = document.createElement('div');
            //
            //     // Add placeholder to replace absolute member-item missing from the flow
            //     if (!placeholder){
            //         console.log("inserting a placeholder")
            //         placeholder.classList.add('placeholder');
            //         placeholder.style.width = `${member.offsetWidth}px`;
            //         placeholder.style.height = `${member.offsetHeight}px`;
            //         parent.insertBefore(placeholder, member.nextSibling);
            //     }
            //
            // }


            // Adjust the member item
            console.log("Mouse over member");
            memberNameElement.classList.remove('cut');
            // memberPicElement.classList.remove('pic--sm');
            memberNameElement.style.whiteSpace = 'normal';
            // member.parentElement.classList.remove('scroll');

            console.log("member.parentElement is:", member.parentElement)
        } else {
            isHovering = false;
            // if (member === member.parentNode.firstElementChild) {
            //     console.log("leaving hover first child")
            //
            //     // Remove the placeholder(s)
            //     const placeholders = document.querySelectorAll('.placeholder');
            //     placeholders.forEach(placeholder => {
            //         placeholder.remove();
            //     });
            //
            // }


            // Reset the member item
            console.log("Mouse leave member");
            memberNameElement.classList.add('cut');
        }
    }
});
