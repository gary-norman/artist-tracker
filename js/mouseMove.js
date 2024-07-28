export function setupGridMouseMoveListener() {
    document.addEventListener('DOMContentLoaded', () => {
        const grids = document.querySelectorAll('.grid');
            grids.forEach(grid => {
                document.addEventListener("mousemove", (e) => {
                    grid.style.setProperty('--x', e.x + 'px');
                    grid.style.setProperty('--y', e.y + 'px');
            });
        });
    });
}

setupGridMouseMoveListener();