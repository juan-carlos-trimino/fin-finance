
const buttonIds = ["lhs-button1", "lhs-button2", "lhs-button3", "lhs-button4", "lhs-button5", "lhs-button6", "lhs-button7",
                   "lhs-button8", "lhs-button9"];

document.addEventListener('DOMContentLoaded', (event) => {
  console.log("Entering document.addEventListener...");
  //Pass Go variables into JavaScript via HTML5 data-* attributes. Grab the element by its exact ID.
  const scriptTag = document.getElementById('page-ui-script');
  //Read the Go data attribute natively.
  const cb = scriptTag ? scriptTag.dataset.cb : null;
  console.log(`Current Button from Go is: ${cb}`);
  setUI(cb);
  console.log("Exiting document.addEventListener...");
});

function setUI(eid) {
  console.log("Entering setUI...");
  for (const id of buttonIds) {
    const btn = document.getElementById(id);
    if (!btn) {
      console.log(`${id} does not exist, breaking loop...`);
      break;  //Exit the loop.
    } else if (btn.hasAttribute('disabled')) {
      console.log(`${id}: The HTML tag contains the 'disabled' attribute.`);
      continue;
    }
    console.log(`Button Id = ${id}`);
    btn.disabled = (id === eid);
  }
  console.log(`Disable button: ${eid}.`)
  //Find the absolute first input field on the page.
  const firstField = document.querySelector('input:not([type="hidden"]):not([disabled]), select, textarea');
  if (firstField) {
    console.log(`Setting focus and selecting text for: ${firstField.id || firstField.name}`);
    firstField.focus();  //Put the cursor in the field.
    firstField.select();  //Highlight the entire text string.
  } else {
    console.log("No interactive input fields found on this page.");
  }
  console.log("Exiting setUI...");
}
