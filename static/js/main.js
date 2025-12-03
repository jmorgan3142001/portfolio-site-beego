document.addEventListener('DOMContentLoaded', () => {
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