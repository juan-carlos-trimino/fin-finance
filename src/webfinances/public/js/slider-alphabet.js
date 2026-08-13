document.addEventListener('DOMContentLoaded', () => {
  console.log("Entering slider-alphabet.js...");
  const slider = document.getElementById('alphabet-range');
  const display = document.getElementById('range-display');
  //Protect code from crashing if a page includes the script but lacks the slider element
  if (!slider || !display) {
    console.log("Exiting slider-alphabet.js (no slider element)...");
    return;
  }
  const rangeLabels = {
    '1': 'A - G',
    '2': 'H - N',
    '3': 'O - T',
    '4': 'U - Z'
  };
  slider.addEventListener('input', (event) => {
    const selectedValue = event.target.value;
    display.textContent = rangeLabels[selectedValue];
  });
  console.log("Exiting slider-alphabet.js...");
});
