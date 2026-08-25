document.addEventListener('DOMContentLoaded', () => {
  console.log("Entering slider-alphabet.js (1)...");
  const slider = document.getElementById('alphabet-range');
  const display = document.getElementById('range-display');
  const ticksContainer = document.getElementById('slider-ticks-container');
  //Protect code from crashing if a page includes the script but lacks the slider element.
  if (!slider || !display) {
    console.log("Missing required slider elements. Exiting script.");
    console.log("Exiting slider-alphabet.js (2)...");
    return;
  }
  //READ FROM GO BACKEND: Extract the comma-separated labels and splits them into an array.
  const rawLabels = slider.getAttribute('data-labels') ? slider.getAttribute('data-labels').split(',') : [];
  // Keep the empty spacer index 0 logic intact so slider 1 aligns with index 1.
  const rangeLabels = ["", ...rawLabels];
  //Generate tick marks based on the array size.
  const generateTicks = () => {
    if (!ticksContainer) return;
    ticksContainer.innerHTML = '';
    for (let i = 1; i < rangeLabels.length; i++) {
      const tick = document.createElement('span');
      tick.className = 'slider-tick-mark';
      tick.textContent = '|';
      ticksContainer.appendChild(tick);
    }
  };
  //Helper function to update the DOM text uniformly
  const updateDisplay = (value) => {
    //Fallback directly to the first array item sent by Go.
    display.textContent = rangeLabels[value] || rawLabels[0];
  };
  //Dynamic Update: Update text instantly while dragging the slider.
  slider.addEventListener('input', (event) => {
    updateDisplay(event.target.value);
  });
  //INITIALIZATION: Run both setup steps on page load.
  generateTicks();
  updateDisplay(slider.value);
  console.log("Exiting slider-alphabet.js (1)...");
});
