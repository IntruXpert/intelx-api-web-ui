// Main JavaScript file for IntelX Web Interface

document.addEventListener('DOMContentLoaded', function() {
    // Show loading indicator when submitting search
    const searchForm = document.querySelector('form[action="/search"]');
    if (searchForm) {
        searchForm.addEventListener('submit', function() {
            const submitButton = this.querySelector('button[type="submit"]');
            if (submitButton) {
                submitButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Searching...';
                submitButton.disabled = true;
            }
        });
    }

    // Format dates in search results
    const formatDates = function() {
        document.querySelectorAll('td:nth-child(2)').forEach(function(cell) {
            try {
                const date = new Date(cell.textContent);
                if (!isNaN(date)) {
                    cell.textContent = date.toLocaleString();
                }
            } catch (e) {
                // Keep original format if parsing fails
            }
        });
    };

    // Call formatDates if we're on the results page
    if (window.location.pathname === '/results') {
        formatDates();
    }

    // Highlight code in preview
    const previewContainer = document.querySelector('.preview-container pre');
    if (previewContainer) {
        // Simple syntax highlighting for common formats
        const content = previewContainer.textContent;
        
        // Check if it looks like JSON
        if (content.trim().startsWith('{') || content.trim().startsWith('[')) {
            try {
                const jsonObj = JSON.parse(content);
                previewContainer.textContent = JSON.stringify(jsonObj, null, 2);
                previewContainer.classList.add('json');
            } catch (e) {
                // Not valid JSON, keep as is
            }
        }
        
        // Add more format detection and highlighting as needed
    }
}); 
