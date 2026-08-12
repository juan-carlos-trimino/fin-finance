document.addEventListener("DOMContentLoaded", function () {
  console.log("Entering document.addEventListener (scrollbarTabFields.js)...");

  // 1. SELECT WORKSPACE CONTAINERS
  const leftSide = document.querySelector(".left-side");
  const rightSide = document.querySelector(".right-side");
  const scrollContainer = document.getElementById("form-scroll-cnt");

  // Selection string for tracking active structural elements (inputs, select dropdowns, buttons)
  const targetSelector = 'button:not([disabled]), input:not([type="hidden"]):not([disabled]), select:not([disabled])';

  // HELPER FUNCTION: Safely builds our combined collection of active focusable objects
  function getWorkspaceElements() {
    let elements = [];

    // Safely pull active nodes from the left side pane first
    if (leftSide) {
      elements = elements.concat(Array.from(leftSide.querySelectorAll(targetSelector)));
    }
    // Append active nodes from the right side panel next
    if (rightSide) {
      elements = elements.concat(Array.from(rightSide.querySelectorAll(targetSelector)));
    }

    // Strictly filter out layout components that are disabled or hidden
    return elements.filter(el => {
      if (el.disabled === true || el.getAttribute('disabled') !== null) {
        return false;
      }
      const style = window.getComputedStyle(el);
      return el.offsetWidth > 0 && el.offsetHeight > 0 && style.display !== 'none' && style.visibility !== 'hidden';
    });
  }

  // SMART INITIAL FOCUS RULE:
  const activeElements = getWorkspaceElements();
  if (activeElements.length > 0) {
    let focusTarget;

    if (activeElements.length === 1) {
      // Rule 1: Single element protection (If only one active button is present)
      focusTarget = activeElements[0];
    } else {
      // Rule 2: Multi-field form rule. Priority to the first text input field (Username input)
      const explicitInputField = document.querySelector('input:not([type="hidden"]):not([disabled])');
      focusTarget = explicitInputField ? explicitInputField : activeElements[0];
    }

    if (focusTarget) {
      setTimeout(() => {
        focusTarget.focus();
        focusTarget.scrollIntoView({ behavior: "instant", block: "nearest" });
      }, 10);
    }
  }

  // Intercept Tab navigation globally to maintain a strict workspace loop trap
  document.addEventListener("keydown", function (e) {
    if (e.key !== "Tab") {
      return;
    }

    const elements = getWorkspaceElements();
    if (elements.length === 0) {
      return;
    }

    const firstItem = elements[0];                // Usually the 'Register' button on the left
    const lastItem = elements[elements.length - 1];  // The 'Save' button on the right
    const activeItem = document.activeElement;

    // RULE 1: SINGLE ITEM PROTECTION
    if (elements.length === 1) {
      e.preventDefault();
      firstItem.focus();
      return;
    }

    // RULE 2: MULTIPLE ELEMENTS MODE
    if (e.shiftKey && activeItem === firstItem) {
      console.log("Wrapping focus backward: From first left element down to right panel Save button.");
      lastItem.focus();
      e.preventDefault(); // Blocks focus from shifting up to the browser address bar
    }
    else if (!e.shiftKey && activeItem === lastItem) {
      console.log("Wrapping focus forward: From right panel Save button up to first left side element.");
      firstItem.focus();
      e.preventDefault(); // Blocks focus from dropping down into footer layouts or unexpected page links
    }
    // RULE 3: ESCAPE SAFEGUARD (If the user clicks on raw text background or focus slips)
    else if (!elements.includes(activeItem)) {
      console.log("Focus drifted outside specified panels. Restoring continuous loop.");
      if (e.shiftKey) {
        lastItem.focus();
      } else {
        firstItem.focus();
      }
      e.preventDefault();
    }
  });

  // Pulls right-side scrolling fields forward gracefully during focus tracking sequences
  if (scrollContainer) {
    scrollContainer.addEventListener("focusin", function (e) {
      if (e.target.tagName === "INPUT" || e.target.tagName === "SELECT") {
        e.target.scrollIntoView({ behavior: "smooth", block: "nearest" });
      }
    });
  }

  console.log("Exiting document.addEventListener (scrollbarTabFields.js)...");
});
