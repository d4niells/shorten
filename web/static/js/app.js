// Only handle copy button functionality
// Form submission is now handled server-side via Go templates

document.getElementById('copyBtn')?.addEventListener('click', () => {
    const shortUrl = document.getElementById('shortUrl');
    if (!shortUrl) return;

    // Copy to clipboard
    if (navigator.clipboard) {
        navigator.clipboard.writeText(shortUrl.value).then(() => {
            const btn = document.getElementById('copyBtn');
            const originalText = btn.textContent;
            btn.textContent = '✓ Copied!';
            setTimeout(() => {
                btn.textContent = originalText;
            }, 2000);
        });
    } else {
        // Fallback for older browsers
        shortUrl.select();
        document.execCommand('copy');
        
        const btn = document.getElementById('copyBtn');
        const originalText = btn.textContent;
        btn.textContent = '✓ Copied!';
        setTimeout(() => {
            btn.textContent = originalText;
        }, 2000);
    }
});
