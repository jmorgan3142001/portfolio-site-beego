document.addEventListener('DOMContentLoaded', () => {
    // --- Dark Mode Logic ---
    const toggleBtn = document.getElementById('theme-toggle');
    const body = document.body;
    
    // Check Local Storage for Preference
    const savedTheme = localStorage.getItem('jake-morgan-dev-theme');
    if (savedTheme === 'dark') {
        body.setAttribute('data-theme', 'dark');
        if (toggleBtn) toggleBtn.innerHTML = '<i class="bi bi-sun-fill"></i>'; // Switch icon
    }

    // 2. Handle Toggle Click
    if (toggleBtn) {
        toggleBtn.addEventListener('click', (e) => {
            e.preventDefault();
            
            if (body.hasAttribute('data-theme')) {
                // Switch to Light
                body.removeAttribute('data-theme');
                localStorage.setItem('jake-morgan-dev-theme', 'light');
                toggleBtn.innerHTML = '<i class="bi bi-moon-fill"></i>';
            } else {
                // Switch to Dark
                body.setAttribute('data-theme', 'dark');
                localStorage.setItem('jake-morgan-dev-theme', 'dark');
                toggleBtn.innerHTML = '<i class="bi bi-sun-fill"></i>';
            }
        });
    }

    /* --- Typewriter Effect --- */
    const textElement = document.getElementById('typewriter-text');
    const phrases = [
        "Scalable Systems",
        "Cloud Architectures",
        "Impactful Software",
        "Reliable Backends"
    ];
    let phraseIndex = 0;
    let charIndex = 0;
    let isDeleting = false;
    let typeSpeed = 100;

    function typeWriter() {
        const currentPhrase = phrases[phraseIndex];
        
        if (isDeleting) {
            textElement.textContent = currentPhrase.substring(0, charIndex - 1);
            charIndex--;
            typeSpeed = 50;
        } else {
            textElement.textContent = currentPhrase.substring(0, charIndex + 1);
            charIndex++;
            typeSpeed = 100;
        }

        if (!isDeleting && charIndex === currentPhrase.length) {
            isDeleting = true;
            typeSpeed = 2000; // Pause at end
        } else if (isDeleting && charIndex === 0) {
            isDeleting = false;
            phraseIndex = (phraseIndex + 1) % phrases.length;
            typeSpeed = 500;
        }

        setTimeout(typeWriter, typeSpeed);
    }

    typeWriter();

    /* --- Hardware Parallax Effect (Optimized) --- */
    const shapes = document.querySelectorAll('.molded-shape');
    let lastScrollY = window.scrollY;
    let ticking = false;

    function updateParallax() {
        shapes.forEach(shape => {
            const speed = shape.getAttribute('data-speed');
            const yPos = -(lastScrollY * speed);
            shape.style.transform = `translate3d(0, ${yPos}px, 0)`; // Use translate3d for hardware acceleration
        });
        ticking = false;
    }
    
    window.addEventListener('scroll', () => {
        lastScrollY = window.scrollY;
        if (!ticking) {
            window.requestAnimationFrame(updateParallax);
            ticking = true;
        }
    }, { passive: true });
});