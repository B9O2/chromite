var clickableElements = document.querySelectorAll('a, button, input[type="button"], input[type="submit"], input[type="image"], [role="button"], [role="link"], [onclick], [tabindex], [contenteditable="true"],[contenteditable=""], [contenteditable="inherit"]');
clickableElements.forEach(function(element) {
    element.click()
});