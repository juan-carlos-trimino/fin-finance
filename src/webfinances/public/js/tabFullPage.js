
document.addEventListener("DOMContentLoaded", function () {
  console.log("Entering tabFullPage.js...");
  //Find form elements across the page layout (inputs and buttons)
  const targetSelector = 'button:not([disabled]), input:not([type="hidden"]):not([disabled]), select:not([disabled])';
  //Build the collection of active focusable objects.
  function getFocusableElements() {
    return Array.from(document.querySelectorAll(targetSelector)).filter(el => {
      if (el.disabled === true || el.getAttribute('disabled') !== null) {
        return false;
      }
      const style = window.getComputedStyle(el);
      return el.offsetWidth > 0 && el.offsetHeight > 0 && style.display !== 'none' && style.visibility !== 'hidden';
    });
  }
  const initialElements = getFocusableElements();
  if (initialElements.length > 0) {
    let focusTarget;
    if (initialElements.length === 1) {
      //If there is only ONE element on the screen (like a solitary button), target it directly.
      focusTarget = initialElements[0];
      console.log("Single element layout detected. Defaulting focus to the element:", focusTarget);
    } else {
      //For multi-field forms, search explicitly for an input field first to force focus upward.
      const explicitInputField = document.querySelector('input:not([type="hidden"]):not([disabled])');
      if (explicitInputField) {
        focusTarget = explicitInputField;
        console.log("Multi-field layout detected. Forcing focus to the first text field input:", focusTarget);
      } else {
        focusTarget = initialElements[0];
        console.log("Multi-field layout detected (No inputs found). Defaulting to index 0:", focusTarget);
      }
    }
    //
    if (focusTarget) {
      //Wrap in a tiny timeout to ensure the browser has fully rendered layout dimensions.
      setTimeout(() => {
        focusTarget.focus();
        focusTarget.scrollIntoView({ behavior: "instant", block: "nearest" });
      }, 10);
    }
  }
  //Intercept Tab navigation globally to maintain a strict form focus trap
  document.addEventListener("keydown", function (e) {
    if (e.key !== "Tab") {
      console.log("Exiting tabFullPage.js (not the TAB key)...");
      return;
    }
    const elements = getFocusableElements();
    if (elements.length === 0) {
      console.log("Exiting tabFullPage.js (no elements)...");
      return;
    }
    const firstItem = elements[0];
    const lastItem = elements[elements.length - 1];
    const activeItem = document.activeElement;
    //Single item protection.
    if (elements.length === 1) {
      e.preventDefault();
      firstItem.focus();
      console.log("Exiting tabFullPage.js (only one elements)...");
      return;
    }
    //Multiple element mode
    if (e.shiftKey && activeItem === firstItem) {
      console.log("Looping focus backward: Wrapping to the last element.");
      lastItem.focus();
      e.preventDefault();
    } else if (!e.shiftKey && activeItem === lastItem) {
      console.log("Looping focus forward: Wrapping back to the first element.");
      firstItem.focus();
      e.preventDefault();
    } else if (!elements.includes(activeItem)) {  //Escape safeguard.
      console.log("Focus drifted outside standard inputs. Restoring form loop.");
      if (e.shiftKey) {
        lastItem.focus();
      } else {
        firstItem.focus();
      }
      e.preventDefault();
    }
  });
  console.log("Exiting tabFullPage.js...");
});
