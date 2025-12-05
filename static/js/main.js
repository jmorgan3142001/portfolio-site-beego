document.addEventListener('DOMContentLoaded', () => {
    // --- Dark Mode Logic ---
    const toggles = document.querySelectorAll('.theme-toggle');
    const body = document.body;

    // Helper to update both toggle button icons on click
    const updateIcons = (isDark) => {
        toggles.forEach(btn => {
            if (isDark) {
                btn.innerHTML = '<i class="bi bi-sun-fill"></i>';
            } else {
                btn.innerHTML = '<i class="bi bi-moon-fill"></i>';
            }
        });
    };

    // Check Local Storage for Preference
    const savedTheme = localStorage.getItem('jake-morgan-dev-theme');
    if (savedTheme === 'dark') {
        body.setAttribute('data-theme', 'dark');
        updateIcons(true); // Set initial state for all buttons
    }

    toggles.forEach(toggleBtn => {
        toggleBtn.addEventListener('click', (e) => {
            e.preventDefault();
            
            if (body.hasAttribute('data-theme')) {
                body.removeAttribute('data-theme');
                localStorage.setItem('jake-morgan-dev-theme', 'light');
                updateIcons(false); 
            } else {
                body.setAttribute('data-theme', 'dark');
                localStorage.setItem('jake-morgan-dev-theme', 'dark');
                updateIcons(true); 
            }
        });
    });

    /* --- Typewriter Effect --- */
    const textElement = document.getElementById('typewriter-text');
    const phrases = [
        "Impactful Software",
        "Reliable Backends",
        "Clean UI/UX",
        "Scalable Systems",
        "Cloud Architectures",
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