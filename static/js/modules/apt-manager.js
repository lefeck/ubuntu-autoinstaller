/**
 * APT Configuration Management Module
 * Handles disable components and disable suites functionality
 */

/**
 * APT Configuration Management Functions
 */

/**
 * Add suite from dropdown selection
 */
function addSuiteFromSelect(selectElement) {
    const value = selectElement.value;
    if (value) {
        addSuiteTag(value);
        selectElement.value = ''; // Reset selection
    }
}

/**
 * Add custom suite from text input
 */
function addCustomSuite(event) {
    if (event.key === 'Enter') {
        const input = event.target;
        const value = input.value.trim();
        if (value && !value.includes(' ')) { // No spaces allowed in suite names
            addSuiteTag(value);
            input.value = ''; // Clear input
        }
    }
}

/**
 * Add a suite tag to the display
 */
function addSuiteTag(suite) {
    const display = document.getElementById('aptDisableSuitesDisplay');

    // Check if tag already exists
    const existingTags = Array.from(display.querySelectorAll('.tag-item')).map(tag => tag.dataset.value);
    if (existingTags.includes(suite)) {
        return; // Don't add duplicates
    }

    const tagElement = document.createElement('span');
    tagElement.className = `tag-item ${suite}`;
    tagElement.dataset.value = suite;
    tagElement.innerHTML = `
        ${suite}
        <button type="button" class="tag-remove" onclick="removeSuiteTag('${suite}')">&times;</button>
    `;

    display.appendChild(tagElement);
}

/**
 * Remove a suite tag
 */
function removeSuiteTag(suite) {
    const display = document.getElementById('aptDisableSuitesDisplay');
    const tagElement = display.querySelector(`[data-value="${suite}"]`);
    if (tagElement) {
        tagElement.remove();
    }
}

/**
 * Get selected disable components
 */
function getDisableComponents() {
    const checkboxes = document.querySelectorAll('#aptDisableComponents input[type="checkbox"]:checked');
    return Array.from(checkboxes).map(cb => cb.value);
}

/**
 * Get selected disable suites
 */
function getDisableSuites() {
    const tags = document.querySelectorAll('#aptDisableSuitesDisplay .tag-item');
    return Array.from(tags).map(tag => tag.dataset.value);
}

/**
 * Set disable components from configuration
 */
function setDisableComponents(components) {
    const checkboxes = document.querySelectorAll('#aptDisableComponents input[type="checkbox"]');
    checkboxes.forEach(cb => {
        cb.checked = components.includes(cb.value);
    });
}

/**
 * Set disable suites from configuration
 */
function setDisableSuites(suites) {
    const display = document.getElementById('aptDisableSuitesDisplay');
    display.innerHTML = ''; // Clear existing tags
    suites.forEach(suite => addSuiteTag(suite));
}

/**
 * Add APT repository configuration
 */
function addAptRepo() {
    const container = document.getElementById('aptPrimaryRepos');
    const repoDiv = document.createElement('div');
    repoDiv.className = 'apt-repo-config';
    repoDiv.innerHTML = `
        <div class="apt-repo-header">
            <span class="apt-repo-type">Repository</span>
            <button type="button" class="remove-btn" onclick="removeAptRepo(this)">Remove</button>
        </div>
        <div class="form-row">
            <div class="form-group">
                <label class="optional">Architectures <span class="hint-icon" data-tooltip="CPU architectures this mirror serves (e.g. amd64,i386)">?</span></label>
                <input type="text" class="apt-arches">
            </div>
            <div class="form-group">
                <label class="optional">URI <span class="hint-icon" data-tooltip="Repository base URL (e.g. http://archive.ubuntu.com/ubuntu)">?</span></label>
                <input type="text" class="apt-uri">
            </div>
        </div>
    `;
    container.appendChild(repoDiv);
}

/**
 * Remove APT repository configuration
 */
function removeAptRepo(button) {
    button.closest('.apt-repo-config').remove();
}

/**
 * Initialize APT configuration when tab is loaded
 */
function initAptConfiguration() {
    console.log('Initializing APT configuration...');

    // Ensure all Disable Components checkboxes are unchecked by default
    const disableComponentsCheckboxes = document.querySelectorAll('#aptDisableComponents input[type="checkbox"]');
    disableComponentsCheckboxes.forEach(checkbox => {
        checkbox.checked = false;
    });

    // Ensure all Disable Suites checkboxes are unchecked by default
    const disableSuitesCheckboxes = document.querySelectorAll('#aptDisableSuitesOptions input[type="checkbox"]');
    disableSuitesCheckboxes.forEach(checkbox => {
        checkbox.checked = false;
    });

    console.log('APT configuration initialized - all checkboxes unchecked by default');
}

// Export functions for global use
window.AptManager = {
    addSuiteFromSelect,
    addCustomSuite,
    addSuiteTag,
    removeSuiteTag,
    getDisableComponents,
    getDisableSuites,
    setDisableComponents,
    setDisableSuites,
    addAptRepo,
    removeAptRepo,
    initAptConfiguration
};

// Make functions globally available for HTML onclick handlers
window.addSuiteFromSelect = addSuiteFromSelect;
window.addCustomSuite = addCustomSuite;
window.removeSuiteTag = removeSuiteTag;
window.addAptRepo = addAptRepo;
window.removeAptRepo = removeAptRepo;
