document.addEventListener("DOMContentLoaded", function () {
  console.log("Entering tabSplitPage.js...");
  const leftSide = document.querySelector(".left-side");
  const rightSide = document.querySelector(".right-side");
  const scrollContainer = document.getElementById("form-scroll-cnt");
  //Selection string for tracking active structural elements (inputs, select dropdowns, buttons)
  const targetSelector = 'button:not([disabled]), input:not([type="hidden"]):not([disabled]), select:not([disabled]), .clickable-row';
  //Safely builds our combined collection of active focusable objects
  function getWorkspaceElements() {
    let elements = [];
    //Safely pull active nodes from the left side pane first
    if (leftSide) {
      elements = elements.concat(Array.from(leftSide.querySelectorAll(targetSelector)));
    }
    //Append active nodes from the right side panel next
    if (rightSide) {
      elements = elements.concat(Array.from(rightSide.querySelectorAll(targetSelector)));
    }
    //Strictly filter out layout components that are disabled or hidden
    return elements.filter(el => {
      if (el.disabled === true || el.getAttribute('disabled') !== null) {
        return false;
      }
      const style = window.getComputedStyle(el);
      return el.offsetWidth > 0 && el.offsetHeight > 0 && style.display !== 'none' && style.visibility !== 'hidden';
    });
  }
  const activeElements = getWorkspaceElements();
  if (activeElements.length > 0) {
    let focusTarget;
    if (activeElements.length === 1) {
      //Single element protection (If only one active button is present)
      focusTarget = activeElements[0];
    } else {
      //Multi-field form rule. Priority to the first text input field
      const explicitInputField = document.querySelector('input:not([type="hidden"]):not([disabled])');
      focusTarget = explicitInputField ? explicitInputField : activeElements[0];
    }
    //
    if (focusTarget) {
      setTimeout(() => {
        focusTarget.focus();
        focusTarget.scrollIntoView({ behavior: "instant", block: "nearest" });
      }, 10);
    }
  }
  //Intercept Tab navigation globally to maintain a strict workspace loop trap
  document.addEventListener("keydown", function (e) {
    //Since table rows inside a <form> are now focusable, pressing Enter while selecting a row might accidentally submit your form prematurely.
    //Treat "Enter" on a table row like a click
    if (e.key === "Enter" && document.activeElement.classList.contains("clickable-row")) {
      e.preventDefault();
      document.activeElement.click();  //Fire your selection/highlight logic
      return;
    }
    //
    if (e.key !== "Tab") {
      return;
    }
    const elements = getWorkspaceElements();
    if (elements.length === 0) {
      return;
    }
    const firstItem = elements[0];
    const lastItem = elements[elements.length - 1];
    const activeItem = document.activeElement;
    if (elements.length === 1) {
      e.preventDefault();
      firstItem.focus();
      return;
    }
    //Strict Border Control: Edge Wrapping
    if (e.shiftKey && activeItem === firstItem) {
      console.log("Wrapping focus backward.");
      lastItem.focus();
      e.preventDefault();  //Block focus from shifting up to the browser address bar
      return;
    }
    //
    if (!e.shiftKey && activeItem === lastItem) {
      console.log("Wrapping focus forward.");
      firstItem.focus();
      e.preventDefault();  //Block focus from dropping down into footer layouts or unexpected page links
      return;
    }
    //Safe Focus Drift Detection
    //We must wait for the browser to actually finish moving focus before we accuse it of drifting out of the layout.
    setTimeout(() => {
      const currentElements = getWorkspaceElements();
      const freshActiveItem = document.activeElement;
      //If the element the browser just landed on isn't part of our workspace arrays
      if (!currentElements.includes(freshActiveItem)) {
        console.log("Focus drifted outside specified panels. Restoring continuous loop.");
        if (e.shiftKey) {
          lastItem.focus();
        } else {
          firstItem.focus();
        }
      }
    }, 10); //A tiny 10ms delay gives the DOM time to catch up safely
  });
  //Pull right-side scrolling fields forward gracefully during focus tracking sequences
  if (scrollContainer) {
    scrollContainer.addEventListener("focusin", function (e) {
      if (e.target.tagName === "INPUT" || e.target.tagName === "SELECT") {
        e.target.scrollIntoView({ behavior: "smooth", block: "nearest" });
      }
    });
  }
  console.log("Exiting tabSplitPage...");
});
