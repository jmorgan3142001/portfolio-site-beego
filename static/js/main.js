document.addEventListener('DOMContentLoaded', () => {
    /* --- Typewriter Effect --- */
    const textElement = document.getElementById('typewriter-text');
    const phrases = [
        "Scalable Systems",
        "Cloud Architectures",
        "Efficient Algorithms",
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

    /* --- Hardware Parallax Effect --- */
    // Move shapes slowly based on scroll position
    const shapes = document.querySelectorAll('.molded-shape');
    
    window.addEventListener('scroll', () => {
        const scrolled = window.scrollY;
        
        shapes.forEach(shape => {
            const speed = shape.getAttribute('data-speed');
            const yPos = -(scrolled * speed);
            // Apply subtle movement
            shape.style.transform = `translateY(${yPos}px)`;
        });
    });
});